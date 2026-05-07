# YourNotes

YourNotes is a backend API for a notes application built with Go and PostgreSQL. It allows users to register, log in, manage authentication tokens, and create, read, update, and delete their own notes.

## Features

- User registration
- User login
- JWT access tokens
- Refresh token support
- Revoke refresh tokens
- Create notes
- Get all notes for the logged-in user
- Get a single note
- Update a note
- Delete a note
- Docker support

## Tech Stack

- Go
- PostgreSQL
- sqlc
- Goose
- JWT
- Docker

## Clone the Repository

```bash
git clone https://github.com/ananyabhardwaj10/yournotes.git
cd yournotes
``` 

## Install Dependencies
```bash
go mod tidy
```

## Environment Variables

Create a new `.env` file in the project root and add:
```env
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/yournotes?sslmode=disable
JWT_SECRET=your-secret-key
```

you can generate the secret key using the command : 
```bash 
openssl rand -hex 32
```

## Run the project locally

Make sure PostgreSQL is running and the database exists.

Run the server: 
```bash 
go run . 
```

## Run with Docker

Build the Docker image: 
```bash 
docker build -t yournotes . 
```

Run the container: 
```bash 
docker run --env-file .env -p 8085:8085 yournotes 
```

## Docker Image

Docker Hub: https://hub.docker.com/r/ananyabhardwaj10/yournotes

### How to Run
```bash
docker pull ananyabhardwaj10/yournotes:latest
docker run --env-file .env -p 8080:8080 ananyabhardwaj10/yournotes:latest
```

## API Endpoints

### Auth
- POST /api/register
- POST /api/login
- POST /api/refresh
- POST /api/revoke

### Notes
- GET /api/notes
- GET /api/notes/{noteID}
- POST /api/createnote
- PATCH /api/notes/{noteID}
- DELETE /api/notes/{noteID}

## Example Requests

### Register 
```bash 
curl -X POST http://localhost:8085/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ana","email":"ana@example.com","password":"secret123"}'
```

### Login 
```bash 
curl -X POST http://localhost:8085/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ana@example.com","password":"secret123"}'
```

