import React, { useState, useEffect, useRef } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';

const ChatRoom: React.FC = () => {
  const { boardName: initialBoardName } = useParams<{ boardName: string }>();
  const location = useLocation();
  const [boardName, setBoardName] = useState<string>(initialBoardName || ''); // Set initial board name from route param
  const [messages, setMessages] = useState<string[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const socket = useRef<WebSocket | null>(null);

  // Track whether the username and join message have been sent
  const messagesSentRef = useRef(false);

  // Get the username from the router state
  const username = location.state?.username;

  useEffect(() => {
    if (!boardName || !username) {
      console.error('Board name or username is missing');
      return;
    }

    // Connect to WebSocket server only if socket is not already connected
    if (!socket.current) {
    // Connect to WebSocket server
    console.log(`Connecting to WebSocket for board: ${boardName}`);
    socket.current = new WebSocket(`ws://localhost:8080/ws`);

    socket.current.onopen = () => {
      console.log('WebSocket connection established');
    
      let retryCount = 0;
      const maxRetries = 10;

      const sendMessageWithRetry = () => {
        // Wait until WebSocket is in OPEN state
        if (socket.current?.readyState === WebSocket.OPEN) {
          // Check if message has already been sent
          if (!messagesSentRef.current) {
            // Send the username to the server
            console.log(`Sending username: ${username}`);
            socket.current.send(username);
      
            // Send the join message to the server
            const joinMessage = `/join ${boardName}`;
            console.log(`Sending message to join board: ${joinMessage}`);
            socket.current.send(joinMessage);

            // Mark messages as sent
            messagesSentRef.current = true;
          }
        } else if (retryCount < maxRetries) {
          console.warn(`WebSocket not ready. Retrying in 100ms... (${retryCount + 1}/${maxRetries})`);
          retryCount++;
          setTimeout(sendMessageWithRetry, 100); // Retry after 100ms
        } else {
          console.error('WebSocket failed to become ready after multiple retries.');
        }
      };
    
      // Start the retry logic
      sendMessageWithRetry();
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
  }
    // Cleanup function
    return () => {
      if (socket.current && socket.current.readyState === WebSocket.OPEN) {
        console.log('Closing WebSocket connection');
        socket.current.close();
      }
    };
  }, []);

  // Handle sending messages
  const sendMessage = (e?: React.FormEvent) => {
    if (e) e.preventDefault(); // Prevent default form submission behavior
    if (!newMessage.trim()) return;

    if (newMessage.startsWith('/join ')) {
      // Extract the board name from the /join command
      const newBoardName = newMessage.split(' ')[1];
      if (newBoardName) {
        setBoardName(newBoardName); // Update boardName
        socket.current?.send(`/join ${newBoardName}`); // Send the /join message to the server
      }
    } else if (newMessage === '/leave') {
      setBoardName(''); // Reset boardName when leaving
      socket.current?.send('/leave'); // Send the /leave message to the server
    }
    else {
      if (socket.current) {
        console.log('Sending message:', newMessage);
        socket.current.send(newMessage);
        setNewMessage('');
      } else {
        console.error('WebSocket is not connected');
        alert('Unable to send message: WebSocket is disconnected');
      }
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        bgcolor: '#f5f5f5',
        padding: 2,
      }}
    >
      <Typography variant="h4" component="h1" gutterBottom>
        Chat Room: {boardName || 'No board'}
      </Typography>
      <Paper
        elevation={3}
        sx={{
          width: '100%',
          maxWidth: 600,
          height: 400,
          overflowY: 'auto',
          mb: 2,
          p: 2,
          bgcolor: '#ffffff',
        }}
      >
        <List>
          {messages.map((msg, index) => (
            <ListItem key={index}>
              <ListItemText primary={msg} />
            </ListItem>
          ))}
        </List>
      </Paper>
      <Box
        component="form"
        onSubmit={sendMessage}
        sx={{
          display: 'flex',
          width: '100%',
          maxWidth: 600,
          alignItems: 'center',
          gap: 1,
        }}
      >
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type a message..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          sx={{ flexGrow: 1 }}
        />
        <Button type="submit" variant="contained" color="primary">
          Send
        </Button>
      </Box>
    </Box>
  );
};

export default ChatRoom;
