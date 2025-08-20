import React, { useEffect, useState, useRef } from "react";
import styled from "styled-components";

const Container = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

const LeftPanel = styled.div`
  flex: 1;
  border-right: 1px solid #ddd;
`;

const RightPanel = styled.div`
  width: 300px;
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

const Message = styled.div`
  background: white;
  padding: 8px 12px;
  margin-bottom: 8px;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
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

export default function App() {
  const [messages, setMessages] = useState<string[]>([]);
  const [input, setInput] = useState("");
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Replace with your Go backend WebSocket URL
    ws.current = new WebSocket("ws://localhost:8000/chat/ws");

    ws.current.onmessage = (event) => {
      console.log(`data: ${event.data}`);
      setMessages((prev) => [...prev, event.data]);
    };

    ws.current.onclose = (event) => {
      if (!event.wasClean) {
        console.log("❌ Unexpected WebSocket disconnect", event);
      } else {
        console.log("✅ WebSocket closed cleanly");
      }
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  const sendMessage = () => {
    if (ws.current && input.trim()) {
      ws.current.send(input);
      setInput("");
    }
  };

  return (
    <Container>
      <LeftPanel>{/* Empty for now */}</LeftPanel>
      <RightPanel>
        <Messages>
          {messages.map((msg, idx) => (
            <Message key={idx}>{msg}</Message>
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
}