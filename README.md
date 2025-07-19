# Plantie – Plant Watering Reminder API

Plantie is a backend API for a plant watering reminder application. It helps users manage their plants, set up watering reminders, and receive push notifications so they never forget to care for their green friends.

---

## Features

- **User Authentication**: Sign up and log in securely.
- **Plant Management**: Add, update, view, and delete your plants.
- **Reminders**: Set, update, view, and delete watering reminders for each plant.
- **Push Notifications**: Receive notifications when it’s time to water your plants (via Firebase).
- **Healthcheck**: Simple endpoint to check if the server is running.

---

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL (via GORM ORM)
- **Authentication**: JWT (JSON Web Tokens)
- **Notifications**: Firebase Cloud Messaging
- **Other**: godotenv, cron jobs for reminders

---

## Architecture

The project follows the **MVC (Model-View-Controller)** pattern:

- **Models**: Data structures and database logic (`models/`)
- **Controllers**: Business logic and request handling (`controllers/`)
- **Routes**: API endpoint definitions (`routes/`)
- **Middleware**: Authentication and other request pre-processing (`middleware/`)
- **Utils**: Utility functions (JWT, password hashing, notifications, etc.)

---

## API Endpoints

All endpoints (except `/ping`, `/login`, `/signup`) require a valid JWT in the `Authorization: Bearer <token>` header.

### Healthcheck

- **GET `/ping`**
  - **Response**: `"pong"`

### Authentication

- **POST `/signup`**
  - **Body**: `{ "email": string, "password": string }`
  - **Response**: `{ "token": string, "user": { ... } }`

- **POST `/login`**
  - **Body**: `{ "email": string, "password": string }`
  - **Response**: `{ "token": string, "user": { ... } }`

### User

- **POST `/user/push_token`**
  - **Body**: `{ "token": string }`
  - **Response**: `{ "message": "push token set successfully" }`

### Plants

- **POST `/plant`**
  - **Body**: `{ "name": string, ... }`
  - **Response**: `{ "plant": { ... } }`

- **GET `/plant/:id`**
  - **Response**: `{ "plant": { ... } }`

- **GET `/plants`**
  - **Response**: `{ "plants": [ ... ] }`

- **PUT `/plant`**
  - **Body**: `{ "id": number, ... }`
  - **Response**: `{ "plant": { ... } }`

- **DELETE `/plant/:id`**
  - **Response**: HTTP 204 No Content

### Reminders

- **POST `/plant/:id/reminder`**
  - **Body**: `{ "repeatType": string, "time": string, ... }`
  - **Response**: `{ "reminder": { ... } }`

- **GET `/plant/:id/reminders`**
  - **Response**: `{ "reminders": [ ... ] }`

- **PUT `/plant/:id/reminder`**
  - **Body**: `{ "id": number, ... }`
  - **Response**: `{ "reminder": { ... } }`

- **DELETE `/plant/:id/reminder/:reminderId`**
  - **Response**: HTTP 204 No Content

---

## Environment Variables

Create a `.env` file in the project root with the following variables:

- `DB_URL` – PostgreSQL connection string (e.g., `postgres://user:password@localhost:5432/plantie`)
- `JWT_KEY` – Secret key for signing JWT tokens
- `PORT` – (Optional) Port for the server (default: `:8080`)

Also, ensure you have a valid `google_services.json` file for Firebase Cloud Messaging in the project root.

---

## Getting Started

1. **Clone the repository**
2. **Install dependencies**:  
   ```
   go mod download
   ```
3. **Set up your `.env` file** (see above)
4. **Run the server**:  
   ```
   go run main.go
   ```
5. **API is now available at** `http://localhost:8080` (or your specified port)

---

## License

MIT 