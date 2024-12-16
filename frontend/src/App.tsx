import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import BulletinBoard from './components/BulletinBoard';
import ChatRoom from './components/ChatRoom';
import './App.css'; // Optional: For global styles

const App: React.FC = () => {
  return (
    <Router>
      <div className="app">
        <Routes>
          <Route path="/" element={<BulletinBoard />} />
          <Route path="/chat/:boardName" element={<ChatRoom />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
