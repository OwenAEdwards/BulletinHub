# BulletinHub

## Instructions

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