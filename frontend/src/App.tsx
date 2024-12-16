import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Box } from '@mui/material';
import BulletinBoard from './components/BulletinBoard';
import ChatRoom from './components/ChatRoom';

const App: React.FC = () => {
  return (
    <Router>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100vh',
          backgroundColor: '#f5f5f5', // Optional background color
          textAlign: 'center', // Ensures text and child components align centrally
        }}
      >
        <Routes>
          <Route path="/" element={<BulletinBoard />} />
          <Route path="/chat/:boardName" element={<ChatRoom />} />
        </Routes>
      </Box>
    </Router>
  );
};

export default App;
