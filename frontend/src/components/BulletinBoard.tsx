import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, TextField, Button, Typography } from '@mui/material';

const BulletinBoard: React.FC = () => {
  const [boardName, setBoardName] = useState('');
  const [username, setUsername] = useState('');
  const navigate = useNavigate();

  const handleJoinBoard = (e?: React.FormEvent) => {
    if (e) e.preventDefault(); // Prevent default form submission behavior
    if (boardName.trim() && username.trim()) {
      navigate(`/chat/${boardName}`, { state: { username } });
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        gap: 2,
        maxWidth: '400px',
        width: '100%',
        padding: 2,
        margin: '0 auto', // Center horizontally
        mt: 4, // Add some vertical spacing
      }}
    >
      <Typography variant="h4" component="h1" gutterBottom>
        BulletinHub
      </Typography>
      <form onSubmit={handleJoinBoard}>
        <TextField
          label="Enter board name"
          variant="outlined"
          value={boardName}
          onChange={(e) => setBoardName(e.target.value)}
          fullWidth
          margin="normal"
        />
        <TextField
          label="Enter username"
          variant="outlined"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          fullWidth
          margin="normal"
        />
        <Button
          type="submit" // Enables Enter key submission
          variant="contained"
          color="primary"
          fullWidth
          sx={{ mt: 2 }}
        >
          Join Board
        </Button>
      </form>
    </Box>
  );
};

export default BulletinBoard;
