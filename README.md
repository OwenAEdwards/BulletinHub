# BulletinHub

## Docker Instructions (Preferred)

1. Navigate to the root of the project directory.
1. Ensure you have [Docker](https://docs.docker.com/desktop/setup/install/windows-install/) installed on your machine and have Docker Desktop running in the background. Run `docker-compose up`. This should trigger the `docker-compose.yml` file which runs the `Dockerfile`.
1. Go to http://localhost:3000 to see the frontend.
1. Happy chatting (:

## Special Commands

Joins a new board and leaves old board (if you're in one):
> /join <board_name>

Leaves the current board (if you're in one):
> /leave

Returns a list of open chatrooms:
> /list

Returns a list of all users in the chatroom:
> /users

## Manual Instructions (if Docker not working)

### Running the backend
Open the terminal on the root and navigate to the `backend/` directory:
```bash
cd backend/
```

If you do not already have Golang installed, you can download it here: https://go.dev/doc/install.

Install dependencies for Golang with the following command:
```bash
go mod tidy
```

Run `main.go`, the entry point for the Golang program:
```bash
go run main.go
```

### Running the frontend
Open the terminal on the root and navigate to the `frontend/` directory:
```bash
cd frontend/
```

If you do not already have npm installed, you can download it here: https://nodejs.org/en/download/package-manager.

Ensure all dependencies are installed using `npm`:
```bash
npm install
```

Run the frontend using `npm`:
```bash
npm run dev
```

You should now be able to connect over http://localhost:5173/ (if you choose port number 5173 to serve React on).
