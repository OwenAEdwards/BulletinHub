import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Alert,
} from '@mui/material';

interface UserListProps {
  boardName: string;
}

const UserList: React.FC<UserListProps> = ({ boardName }) => {
  const [users, setUsers] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        console.log(`Fetching users for board: ${boardName}`);
        const response = await fetch(`http://localhost:8080/boards/${boardName}/users`);
        if (!response.ok) {
          throw new Error('Failed to fetch user list');
        }
        const userList = await response.json();
        console.log(`Received users: ${JSON.stringify(userList)}`);
        setUsers(userList);
        setLoading(false);
      } catch (err) {
        console.error('Error fetching user list:', err);
        setError('Error fetching user list');
        setLoading(false);
      }
    };

    fetchUsers();
  }, [boardName]);

  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100px',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ mt: 2 }}>
        <Alert severity="error">{error}</Alert>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        mt: 3,
        p: 2,
        bgcolor: '#ffffff',
        border: '1px solid #ddd',
        borderRadius: '8px',
        boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Typography variant="h6" gutterBottom>
        Users in <strong>{boardName}</strong>
      </Typography>
      {users.length === 0 ? (
        <Typography variant="body1" color="textSecondary">
          No users currently in this board.
        </Typography>
      ) : (
        <List>
          {users.map((user, index) => (
            <ListItem key={index} sx={{ pl: 0 }}>
              <ListItemText primary={user} />
            </ListItem>
          ))}
        </List>
      )}
    </Box>
  );
};

export default UserList;
