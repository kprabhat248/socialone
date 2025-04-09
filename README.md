# SocialOne

SocialOne is a Backend web application  built with a Go (Golang) backend and a React frontend. The application provides CRUD operations, a PostgreSQL database, JWT authentication, and bearer token authorization for secure access to various endpoints.
<img width="1369" alt="Screenshot 2025-04-09 at 5 01 15 PM" src="https://github.com/user-attachments/assets/6fbfeb03-f3b1-45f4-8263-81db02af6557" />


## Features
- **Golang Backend**: RESTful API built with Go to handle all CRUD operations.
- **PostgreSQL Database**: Connected to a PostgreSQL database for data persistence.
- **Microservice Architecture**: Backend is designed as a microservice running in a Docker container for scalability.
- **JWT Authentication**: Secure authentication mechanism using JWT (JSON Web Tokens).
- **Bearer Token Authorization**: Protects routes via bearer token authentication.
- **React Frontend**: User interface developed with React, interacting with the backend via API calls.
<img width="1710" alt="Screenshot 2025-04-09 at 5 04 39 PM" src="https://github.com/user-attachments/assets/770a3bcb-e00d-48fc-8171-52c436a0114e" />

## Technologies Used
- **Golang** (Backend)
- **PostgreSQL** (Database)
- **Docker** (Containerization)
- **React** (Frontend)
- **JWT** (Authentication)
- **Microservices Architecture**
- **Bearer Token** (Authorization)
<img width="1710" alt="Screenshot 2025-04-09 at 5 05 51 PM" src="https://github.com/user-attachments/assets/724a59f3-66d6-4b58-8e2a-d7271a1cc5d9" />
<img width="1710" alt="Screenshot 2025-04-09 at 5 06 21 PM" src="https://github.com/user-attachments/assets/1cdc7fca-c0a1-4a26-8e40-758a123f8d28" />


## Getting Started
<img width="1710" alt="Screenshot 2025-04-09 at 5 03 14 PM" src="https://github.com/user-attachments/assets/534376b7-707b-40f2-b6ad-ce5b94137f0c" />


### Prerequisites
To run this project locally, ensure that you have the following installed:
- **Go** (for backend)
- **Node.js** and **npm** (for frontend)
- **Docker** (for containerization)
- **PostgreSQL** (database)


### Setting up the Backend
<img width="1710" alt="Screenshot 2025-04-09 at 5 03 20 PM" src="https://github.com/user-attachments/assets/bf861460-8e2f-423d-8218-2cc53da3bc2c" />
<img width="1710" alt="Screenshot 2025-04-09 at 5 04 12 PM" src="https://github.com/user-attachments/assets/9d01c69e-def7-4422-82e4-495a185f88b9" />

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

