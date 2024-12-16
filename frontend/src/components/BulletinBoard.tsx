import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, TextField, Button, Typography } from '@mui/material';

const BulletinBoard: React.FC = () => {
  const [boardName, setBoardName] = useState('');
  const [username, setUsername] = useState('');
  const navigate = useNavigate();

  const handleJoinBoard = () => {
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
      }}
    >
      <Typography variant="h4" component="h1" gutterBottom>
        BulletinHub
      </Typography>
      <TextField
        label="Enter board name"
        variant="outlined"
        value={boardName}
        onChange={(e) => setBoardName(e.target.value)}
      />
      <TextField
        label="Enter username"
        variant="outlined"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <Button
        variant="contained"
        color="primary"
        onClick={handleJoinBoard}
      >
        Join Board
      </Button>
    </Box>
  );
};

export default BulletinBoard;
