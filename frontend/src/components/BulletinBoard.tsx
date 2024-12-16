import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

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
    <div className="bulletin-board">
      <h1>Bulletin Boards</h1>
      <input
        type="text"
        placeholder="Enter board name"
        value={boardName}
        onChange={(e) => setBoardName(e.target.value)}
      />
      <input
        type="text"
        placeholder="Enter username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={handleJoinBoard}>Join Board</button>
    </div>
  );
};

export default BulletinBoard;
