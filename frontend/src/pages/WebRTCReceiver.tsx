import React, { useEffect, useRef } from "react";

const WebRTCReceiver: React.FC = () => {
  const videoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    let pc: RTCPeerConnection | null = null;
    (async () => {
      pc = new RTCPeerConnection({
        iceServers: [{ urls: ["stun:stun.l.google.com:19302"] }]
      });

      pc.ontrack = (ev) => {
        if (videoRef.current) {
          videoRef.current.srcObject = ev.streams[0];
        }
      };

      // We only receive
      pc.addTransceiver("video", { direction: "recvonly" });

      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);

      const res = await fetch("http://localhost:8000/chat/webrtc/offer", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ sdp: offer.sdp, type: "offer" }),
      });
      const answer = await res.json();
      await pc.setRemoteDescription({ type: "answer", sdp: answer.sdp });

      // optional: handle ICE candidates (pion will include them on answer side)
    })();

    return () => {
      if (pc) pc.close();
    };
  }, []);

  return <video ref={videoRef} autoPlay playsInline muted style={{ width: "100%", height: "100%", background: "black" }} />;
};

export default WebRTCReceiver;