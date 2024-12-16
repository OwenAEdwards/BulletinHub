import React, { useState, useEffect, useRef } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import UserList from './UserList';

const ChatRoom: React.FC = () => {
  const { boardName } = useParams<{ boardName: string }>();
  const location = useLocation();
  const [messages, setMessages] = useState<string[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const socket = useRef<WebSocket | null>(null);

  // Get the username from the router state
  const username = location.state?.username;

  useEffect(() => {
    if (!boardName || !username) {
      console.error("Board name is missing");
      return;
    }

    // Connect to WebSocket server
    console.log(`Connecting to WebSocket for board: ${boardName}`);
    socket.current = new WebSocket(`ws://localhost:8080/ws`);

    socket.current.onopen = () => {
      console.log('WebSocket connection established');

      // Send the username and join message to the server
      if (boardName) {
        const joinMessage = `/join ${boardName}`;
        console.log(`Sending username: ${username}`);
        socket.current?.send(username);
        console.log(`Sending message to join board: ${joinMessage}`);
        socket.current?.send(joinMessage);
      }
    };

    socket.current.onmessage = (event) => {
      console.log('Message received from WebSocket:', event.data);
      setMessages((prev) => [...prev, event.data]);
    };

    socket.current.onclose = (event) => {
      console.log('WebSocket connection closed');
      if (event.wasClean) {
        console.log(`Closed cleanly, code=${event.code}, reason=${event.reason}`);
      } else {
        console.error(`WebSocket closed unexpectedly, code=${event.code}, reason=${event.reason}`);
      }
    };

    socket.current.onerror = (error) => {
      console.error('WebSocket error:', error);
      if (error instanceof Event) {
        console.error('Error event details:', error);
      } else {
        console.error('Error details:', error);
      }
    };

    return () => {
      if (socket.current) {
        console.log('Closing WebSocket connection');
        socket.current.close();
      }
    };
  }, [boardName, username]);

  const sendMessage = () => {
    if (!newMessage.trim()) return;
    if (socket.current) {
      console.log('Sending message:', newMessage);
      socket.current.send(newMessage);
      setNewMessage('');
    } else {
      console.error('WebSocket is not connected');
      alert('Unable to send message: WebSocket is disconnected');
    }
  };

  return (
    <div className="chat-room">
      <h1>Chat Room: {boardName}</h1>
      <div className="chat-messages">
        {messages.map((msg, index) => (
          <div key={index} className="message">
            {msg}
          </div>
        ))}
      </div>
      <div className="message-input">
        <input
          type="text"
          placeholder="Type a message..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
        />
        <button onClick={sendMessage}>Send</button>
      </div>
      <UserList boardName={boardName || ''} />
    </div>
  );
};

export default ChatRoom;
