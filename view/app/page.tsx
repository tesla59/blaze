"use client";
import React, { useState, useEffect } from 'react';

export default function Home() {
  const [clientID, setClientID] = useState('');
  const [targetID, setTargetID] = useState('');
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<any[]>([]);
  const [message, setMessage] = useState('');

  const handleConnect = () => {
    if (!clientID || !targetID) {
      alert('Please enter both Client ID and Target ID.');
      return;
    }
    try {
    const ws = new WebSocket('ws://localhost:8080/ws');
      ws.onopen = () => {
        console.log('Connected to WebSocket');
        ws.send(JSON.stringify({ "id": clientID, "target": targetID }));
      };

      ws.onmessage = (event) => {
        const data = event.data;
        setMessages((prevMessages) => [...prevMessages, data]);
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
      socket.send(JSON.stringify({ message, clientID, targetID }));
      setMessage('');
    }
  };

  useEffect(() => {
    const id = getOrCreateClientId();
    setClientID(id);
    return () => {
      if (socket) socket.close();
    };
  }, [socket]);

  return (
      <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
        <h1>Your Client ID</h1>
        <p>{clientID}</p>
        <div>
          <input
              type="text"
              placeholder="Target ID"
              value={targetID}
              onChange={(e) => setTargetID(e.target.value)}
          />
          <button onClick={handleConnect}>Connect</button>
        </div>
        <div style={{ marginTop: '20px' }}>
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
    // if (!clientId) { # Temporary, will change when we add authentication
      clientId = crypto.randomUUID();
      localStorage.setItem('clientId', clientId);
    // }
    return clientId;
  }
  // Fallback for server-side rendering
  return `client-${Date.now()}`;
};