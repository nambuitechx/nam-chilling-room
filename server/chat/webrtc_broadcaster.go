package chat

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

// ---------- Global broadcaster state ----------
var (
	videoTrackMu   sync.RWMutex
	videoTrack     *webrtc.TrackLocalStaticSample
	peersMu        sync.Mutex
	peers          = map[*webrtc.PeerConnection]struct{}{}
	broadcasterMu  sync.Mutex
	isBroadcasting bool
	cancelBroadcaster context.CancelFunc
)

// readIVFHeader reads 32 bytes header and returns nil on success.
func readIVFHeader(r io.Reader) error {
	h := make([]byte, 32)
	if _, err := io.ReadFull(r, h); err != nil {
		return err
	}
	if string(h[0:4]) != "DKIF" || string(h[8:12]) != "VP80" {
		return fmt.Errorf("not VP8 IVF stream")
	}
	return nil
}

// readIVFFrame reads one IVF frame (size: uint32 little-endian, then pts uint64 little-endian, then payload)
func readIVFFrame(r io.Reader) ([]byte, error) {
	var sz uint32
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return nil, err
	}
	// read 8-byte PTS (we ignore value)
	var _pts uint64
	if err := binary.Read(r, binary.LittleEndian, &_pts); err != nil {
		return nil, err
	}
	buf := make([]byte, sz)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// ---------- startBroadcaster
// Run ffmpeg -> IVF on stdout, parse IVF frames, write to global videoTrack.
// This runs until EOF. Only one broadcaster runs at a time.
func startBroadcaster(ctx context.Context, mp4Path string) error {
	// ensure only one broadcaster at a time
	broadcasterMu.Lock()
	if isBroadcasting {
		broadcasterMu.Unlock()
		return fmt.Errorf("a broadcast is already running")
	}
	// create cancellable ctx so we can stop broadcaster later
	ctx, cancel := context.WithCancel(ctx)
	isBroadcasting = true
	cancelBroadcaster = cancel
	broadcasterMu.Unlock()

	// ffmpeg command: -re to read in realtime; transcode to VP8 IVF
	// We transcode to VP8 for simplicity; change as needed for H264
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-re", "-i", mp4Path,
		"-an", // drop audio for now (add audio handling later if needed)
		"-c:v", "libvpx", "-deadline", "realtime",
		"-f", "ivf", "pipe:1",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancelBroadcastState()
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		cancelBroadcastState()
		return fmt.Errorf("ffmpeg start: %w", err)
	}

	// log ffmpeg stderr asynchronously
	go func() {
		b, _ := io.ReadAll(stderr)
		if len(b) > 0 {
			log.Printf("ffmpeg: %s", bytes.TrimSpace(b))
		}
	}()

	// parse IVF header
	if err := readIVFHeader(stdout); err != nil {
		_ = cmd.Process.Kill()
		cancelBroadcastState()
		return fmt.Errorf("ivf header: %w", err)
	}

	// ensure global track exists (create if nil)
	videoTrackMu.Lock()
	if videoTrack == nil {
		t, err := webrtc.NewTrackLocalStaticSample(
			webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8},
			"video", "pion",
		)
		if err != nil {
			videoTrackMu.Unlock()
			_ = cmd.Process.Kill()
			cancelBroadcastState()
			return fmt.Errorf("create track: %w", err)
		}
		videoTrack = t
	}
	vt := videoTrack
	videoTrackMu.Unlock()

	// Read frames and write to track
	go func() {
		defer func() {
			// when done: cleanup
			_ = cmd.Wait()
			cancelBroadcastState()
		}()

		for {
			select {
			case <-ctx.Done():
				// stop ffmpeg process if still running
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
				return
			default:
			}

			frame, err := readIVFFrame(stdout)
			if err != nil {
				if err == io.EOF {
					log.Println("broadcaster finished (EOF)")
				} else {
					log.Println("ivf frame read error:", err)
				}
				// close peers after finishing
				peersMu.Lock()
				for pc := range peers {
					_ = pc.Close()
					delete(peers, pc)
				}
				peersMu.Unlock()
				return
			}

			// Write sample to video track; duration is approximate since IVF doesn't deliver timebase here.
			// Using 33ms per frame (~30fps) as reasonable default.
			err = vt.WriteSample(media.Sample{Data: frame, Duration: 33 * time.Millisecond})
			if err != nil {
				// WriteSample can fail if peer disconnected; log and continue
				log.Printf("WriteSample error: %v", err)
			}
			// small sleep to avoid busy-loop — ffmpeg -re already paces, but we keep a tiny sleep
			time.Sleep(33 * time.Millisecond)
		}
	}()

	return nil
}

func cancelBroadcastState() {
	broadcasterMu.Lock()
	isBroadcasting = false
	if cancelBroadcaster != nil {
		cancelBroadcaster()
		cancelBroadcaster = nil
	}
	broadcasterMu.Unlock()
}

// ---------- HTTP handler: /webrtc/offer ----------
type sdpPayload struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

func webrtcOfferHandler(w http.ResponseWriter, r *http.Request) {
	// parse offer
	var in sdpPayload
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid offer", http.StatusBadRequest)
		return
	}

	// create peer connection
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		http.Error(w, "pc create failed", http.StatusInternalServerError)
		return
	}

	// add global video track to this peer
	videoTrackMu.RLock()
	vt := videoTrack
	videoTrackMu.RUnlock()

	// If no broadcaster started yet, create a placeholder track so viewer doesn't fail.
	if vt == nil {
		tmp, _ := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
		vt = tmp
		videoTrackMu.Lock()
		videoTrack = vt
		videoTrackMu.Unlock()
	}

	if _, err := pc.AddTrack(vt); err != nil {
		log.Println("AddTrack error:", err)
	}

	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Println("peer state:", s.String())
		if s == webrtc.PeerConnectionStateFailed ||
			s == webrtc.PeerConnectionStateClosed ||
			s == webrtc.PeerConnectionStateDisconnected {
			pc.Close()
			peersMu.Lock()
			delete(peers, pc)
			peersMu.Unlock()
		}
	})

	// set remote offer
	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: in.SDP}
	if err := pc.SetRemoteDescription(offer); err != nil {
		http.Error(w, "SetRemoteDescription failed", http.StatusInternalServerError)
		return
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		http.Error(w, "CreateAnswer failed", http.StatusInternalServerError)
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		http.Error(w, "SetLocalDescription failed", http.StatusInternalServerError)
		return
	}
	<-gatherComplete

	peersMu.Lock()
	peers[pc] = struct{}{}
	peersMu.Unlock()

	out := sdpPayload{SDP: pc.LocalDescription().SDP, Type: "answer"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
