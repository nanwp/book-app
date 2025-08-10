# Book App

A full-stack application for managing books, built with Go (backend) and Next.js (frontend).

## Project Structure

```
byfood-interview/
├── backend/         # Go backend API
├── frontend/        # Next.js frontend
├── docker-compose.yml
├── README.md        # Project overview and instructions
```

- **backend/**: Contains the Go REST API, business logic, migrations, and tests.
- **frontend/**: Contains the Next.js web application.
- **docker-compose.yml**: For running backend and frontend together using Docker.

## Setup Instructions

### Prerequisites
- Docker & Docker Compose
- Go (for backend development)
- Node.js & npm (for frontend development)

### Running with Docker Compose

```zsh
docker-compose up --build
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000

### Running Backend Locally

```zsh
cd backend
make run
```

### Running Frontend Locally

```zsh
cd frontend
npm install
npm run dev
```

## API Endpoints

See backend API documentation:
- Swagger: [`backend/docs/swagger.json`](backend/docs/swagger.json) or [`backend/docs/swagger.yaml`](backend/docs/swagger.yaml)
- URL Swagger (example): http://localhost:8080/swagger
- Example endpoints:
  - `GET /books` - List all books
  - `POST /books` - Create a new book
  - `GET /books/{id}` - Get book details
  - `PUT /books/{id}` - Update a book
  - `DELETE /books/{id}` - Delete a book

## Testing

- **Backend**: Unit and integration test instructions are provided in [`backend/TEST_README.md`](backend/TEST_README.md). Please refer to that file for details on running and understanding backend tests.
- **Frontend**: Standard Next.js testing setup (add your preferred testing library).

## Contribution

Feel free to open issues or submit pull requests for improvements.

## License

MIT
