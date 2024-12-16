import React, { useState, useEffect } from 'react';

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
    return <div>Loading users...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="user-list">
      <h2>Users in {boardName}</h2>
      {users.length === 0 ? (
        <p>No users currently in this board.</p>
      ) : (
        <ul>
          {users.map((user, index) => (
            <li key={index}>{user}</li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default UserList;
