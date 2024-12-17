# BulletinHub

## Instructions

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
