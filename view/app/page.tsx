"use client";
import React, { useState, useEffect } from 'react';

export default function Home() {
  const [clientID, setClientID] = useState('');
  const [targetID, setTargetID] = useState('');
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [message, setMessage] = useState('');

  const handleConnect = () => {
    try {
        if (socket && socket.readyState === WebSocket.OPEN) {
            return
        }
        const ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
            console.log('Connected to WebSocket');
            ws.send(JSON.stringify({ "type": "identity", "id": clientID }));
        };

        ws.onmessage = (event) => {
            try {
                const parsedData = JSON.parse(event.data);
                console.log("Received JSON Parsed as:", parsedData);

                switch (parsedData.type) {
                    case "identity":
                        setTargetID(parsedData.target);
                        break;
                    case "message":
                        setMessages((prevMessages) => [...prevMessages, parsedData.message ]);
                        break;
                    case "matched":
                        setTargetID(parsedData.value);
                        break;
                    case "rematch":
                        setMessages([]);
                        setTargetID(parsedData.value);
                        break;
                    default:
                        console.warn("Unknown message type:", parsedData.type);
                }
            } catch (error) {
                console.error("Error parsing JSON:", error);
            }

        };

        ws.onclose = () => {
            console.log('WebSocket connection closed');
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        setSocket(ws);
    } catch (e) {
        console.error('WebSocket error:', e);
    }
  };

  const handleSendMessage = () => {
    if (socket && message) {
        console.log("Sending message:", JSON.stringify({ "type": "message" , "message": message }));
      socket.send(JSON.stringify({ "type": "message" , "message": message}));
      setMessage('');
    }
  };

  const handleShuffle = () => {
    if (socket) {
        socket.send(JSON.stringify({ "type": "join" }));
    }
  };
    
  useEffect(() => {
    const id = getOrCreateClientId();
    setClientID(id);
    return () => {
      if (socket) {
          socket.send(JSON.stringify({ "type": "disconnect" }));
          socket.close();
      }
    };
  }, [socket]);

  return (
      <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
        <h1>Your Client ID</h1>
        <p>{clientID}</p>
        <h2>Connect to the Target</h2>
        <p>{targetID}</p>
          <div>
              <button onClick={handleConnect}>Connect</button>
              <button onClick={handleShuffle}>Shuffle</button>
          </div>
          <div style={{marginTop: '20px'}}>
        <textarea
            placeholder="Type a message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
        />
          <button onClick={handleSendMessage}>Send</button>
        </div>
        <div style={{ marginTop: '20px' }}>
          <h2>Messages</h2>
          <ul>
            {messages.map((msg, index) => (
                <li key={index}>{JSON.stringify(msg)}</li>
            ))}
          </ul>
        </div>
      </div>
  );
}

const getOrCreateClientId = () => {
  // Check if we're in a browser environment
  if (typeof window !== 'undefined') {
    let clientId = localStorage.getItem('clientId');
    if (!clientId) {
      clientId = crypto.randomUUID();
      localStorage.setItem('clientId', clientId);
    }
    return clientId;
  }
  // Fallback for server-side rendering
  return `client-${Date.now()}`;
};