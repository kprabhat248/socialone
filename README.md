# SocialOne

SocialOne is a Backend web application  built with a Go (Golang) backend and a React frontend. The application provides CRUD operations, a PostgreSQL database, JWT authentication, and bearer token authorization for secure access to various endpoints.

## Features
- **Golang Backend**: RESTful API built with Go to handle all CRUD operations.
- **PostgreSQL Database**: Connected to a PostgreSQL database for data persistence.
- **Microservice Architecture**: Backend is designed as a microservice running in a Docker container for scalability.
- **JWT Authentication**: Secure authentication mechanism using JWT (JSON Web Tokens).
- **Bearer Token Authorization**: Protects routes via bearer token authentication.
- **React Frontend**: User interface developed with React, interacting with the backend via API calls.

## Technologies Used
- **Golang** (Backend)
- **PostgreSQL** (Database)
- **Docker** (Containerization)
- **React** (Frontend)
- **JWT** (Authentication)
- **Microservices Architecture**
- **Bearer Token** (Authorization)


## Getting Started

### Prerequisites
To run this project locally, ensure that you have the following installed:
- **Go** (for backend)
- **Node.js** and **npm** (for frontend)
- **Docker** (for containerization)
- **PostgreSQL** (database)

### Setting up the Backend

1. Clone the repository:
   ```bash
   git clone https://github.com/kprabhat248/socialone.git
   cd socialone/backend


Set up the PostgreSQL database:
Update the DATABASE_URL in your environment file with the correct PostgreSQL credentials.
Create a .env file in the backend folder with the following content:
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=socialone
JWT_SECRET=yourjwtsecret
Build and run the Go backend:
go run main.go
Alternatively, you can run the backend using Docker:
Docker Setup
To run both the backend and frontend with Docker, use Docker Compose. The docker-compose.yml file is set up to create containers for both the backend and the frontend.

In the root project directory, run:
docker-compose up --build
This will start the services defined in docker-compose.yml, which includes:
A Golang backend service.
A React frontend service.
A PostgreSQL service.
The application will be available at http://localhost:3000 (frontend) and the backend API at http://localhost:8080.

###Troubleshooting

Ensure that Docker is running before starting the containers.
Verify that the PostgreSQL service is up and running. If there are connection issues, check the database credentials in the .env file.
Contributing

If you want to contribute to SocialOne, feel free to fork the repository and submit a pull request. Make sure to follow the standard Go and JavaScript coding conventions.

