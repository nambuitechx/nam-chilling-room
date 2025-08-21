import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import styled from "styled-components";

const Container = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

const LeftPanel = styled.div`
  flex: 1;
  border-right: 1px solid #ddd;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const VideoPlayer = styled.video`
  width: 100%;
  height: 100%;
  background: black;
`;

const RightPanel = styled.div`
  width: 400px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
`;

const Messages = styled.div`
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: #f8f9fa;
`;

const MessageRow = styled.div<{ isSelf: boolean }>`
  display: flex;
  justify-content: ${(props) => (props.isSelf ? "flex-end" : "flex-start")};
  margin-bottom: 8px;
`;

const MessageBubble = styled.div<{ isSelf: boolean }>`
  max-width: 80%;
  padding: 8px 12px;
  border-radius: 12px;
  word-wrap: break-word;
  background: ${(props) => (props.isSelf ? "#007bff" : "white")};
  color: ${(props) => (props.isSelf ? "white" : "black")};
`;

const MessageUsername = styled.p`
  margin: 0px 0px 5px 0px;
  padding: 0px;
`;

const MessageContent = styled.p`
  margin: 0px;
  padding: 0px;
`;

const InputArea = styled.div`
  display: flex;
  padding: 12px;
  border-top: 1px solid #ddd;
`;

const Input = styled.input`
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 6px;
  margin-right: 8px;
  font-size: 14px;
`;

const Button = styled.button`
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-size: 14px;

  &:hover {
    background: #0056b3;
  }
`;

const TopBar = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #ddd;
`;

const LogoutButton = styled.button`
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  background: #dc3545;
  color: white;
  cursor: pointer;
  font-size: 13px;

  &:hover {
    background: #b02a37;
  }
`;

// ---------------- Message Component ----------------

type MessageItemProps = {
  isSelf: boolean;
  username: string;
  content: string;
};

const MessageItem: React.FC<MessageItemProps> = ({ isSelf, username, content }) => {
  return (
    <MessageRow isSelf={isSelf}>
      <MessageBubble isSelf={isSelf}>
        <MessageUsername>{username}</MessageUsername>
        <MessageContent>{content}</MessageContent>
      </MessageBubble>
    </MessageRow>
  );
};

// ---------------- ChatPage Component ----------------

type ChatPageProps = {
  token: string | null;
  setToken: (token: string | null) => void;
};

type ServerMessage = {
  tokenString: string;
  username: string;
  content: string;
};

const ChatPage: React.FC<ChatPageProps> = ({ token, setToken }) => {
  const [messages, setMessages] = useState<ServerMessage[]>([]);
  const [input, setInput] = useState("");
  const ws = useRef<WebSocket | null>(null);
  const navigate = useNavigate();

  // Media refs
  const mediaSourceRef = useRef<MediaSource | null>(null);
  const sourceBufferRef = useRef<SourceBuffer | null>(null);
  const videoRef = useRef<HTMLVideoElement | null>(null);

  useEffect(() => {
    ws.current = new WebSocket("ws://localhost:8000/chat/ws");

    // Set binary type for media chunks
    ws.current.binaryType = "arraybuffer";

    ws.current.onmessage = (event) => {
      // Check if it's JSON (chat) or binary (media)
      if (typeof event.data === "string") {
        const data = JSON.parse(event.data) as ServerMessage;
        setMessages((prev) => [...prev, data]);
      } else {
        // Binary data -> feed into SourceBuffer
        if (sourceBufferRef.current && !sourceBufferRef.current.updating) {
          sourceBufferRef.current.appendBuffer(new Uint8Array(event.data));
        }
      }
    };

    ws.current.onclose = (event) => {
      if (!event.wasClean) {
        console.log("❌ Unexpected WebSocket disconnect", event);
      } else {
        console.log("✅ WebSocket closed cleanly");
      }
    };

    // Setup MediaSource
    if (videoRef.current) {
      mediaSourceRef.current = new MediaSource();
      videoRef.current.src = URL.createObjectURL(mediaSourceRef.current);

      mediaSourceRef.current.addEventListener("sourceopen", () => {
        // Adjust codec to your file type
        sourceBufferRef.current = mediaSourceRef.current!.addSourceBuffer(
          'video/mp4; codecs="avc1.64001f, mp4a.40.2"'
        );
      });
    }

    return () => {
      ws.current?.close();
    };
  }, []);

  const sendMessage = () => {
    if (ws.current && input.trim()) {
      const message = {
        tokenString: token,
        content: input,
      };
      ws.current.send(JSON.stringify(message));
      setInput("");
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    setToken(null);
    navigate("/login", { replace: true });
  };

  return (
    <Container>
      <LeftPanel>
        {/* Live video */}
        <VideoPlayer
          ref={videoRef}
          autoPlay
          playsInline
          muted
        />
      </LeftPanel>

      <RightPanel>
        <TopBar>
          <LogoutButton onClick={handleLogout}>Logout</LogoutButton>
        </TopBar>

        <Messages>
          {messages.map((msg, idx) => (
            <MessageItem
              key={idx}
              isSelf={msg.tokenString === token}
              username={msg.username}
              content={msg.content}
            />
          ))}
        </Messages>

        <InputArea>
          <Input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            onKeyDown={(e) => e.key === "Enter" && sendMessage()}
          />
          <Button onClick={sendMessage}>Send</Button>
        </InputArea>
      </RightPanel>
    </Container>
  );
};

export default ChatPage;