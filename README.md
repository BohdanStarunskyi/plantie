# Plantie – Plant Watering Reminder API

Plantie is a backend API for a plant watering reminder application. It helps users manage their plants, set up watering reminders, and receive push notifications so they never forget to care for their green friends.

---

## Features

- **User Authentication**: Secure sign up, login, and token refresh with JWT
- **Plant Management**: Add, update, view, and delete your plants
- **Reminders**: Set, update, view, and delete watering reminders for each plant
- **Push Notifications**: Receive notifications when it's time to water your plants (via Firebase Cloud Messaging)
- **Cron Jobs**: Automated reminder scheduling and notification delivery
- **Automatic Data Cleanup**: Related data is automatically removed when users or plants are deleted
- **Healthcheck**: Simple endpoint to check if the server is running
- **Graceful Shutdown**: Proper server shutdown handling

---

## Tech Stack

- **Language**: Go 1.24.0
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL (via GORM ORM)
- **Authentication**: JWT (JSON Web Tokens) with refresh tokens
- **Notifications**: Firebase Cloud Messaging
- **Scheduling**: GoCron for reminder management
- **Validation**: Go Playground Validator
- **CORS**: Built-in CORS support
- **Environment**: godotenv for configuration

---

## Architecture

The project follows the **MVC (Model-View-Controller)** pattern:

- **Models**: Data structures and database logic (`models/`)
- **Controllers**: Business logic and request handling (`controllers/`)
- **Routes**: API endpoint definitions (`routes/`)
- **Middleware**: Authentication and other request pre-processing (`middleware/`)
- **Utils**: Utility functions (JWT, password hashing, notifications, etc.)
- **Config**: Database and application configuration (`config/`)
- **Constants**: Application constants (`constants/`)

---

## API Endpoints

All endpoints (except `/ping`, `/login`, `/signup`, `/refresh`) require a valid JWT in the `Authorization: Bearer <token>` header.

### Healthcheck

- **GET `/ping`**
  - **Response**: `"pong"`

### Authentication

- **POST `/signup`**
  - **Body**: `{ "email": string, "password": string, "name": string }`
  - **Response**: `{ "access_token": string, "refresh_token": string, "user": { ... } }`

- **POST `/login`**
  - **Body**: `{ "email": string, "password": string }`
  - **Response**: `{ "access_token": string, "refresh_token": string, "user": { ... } }`

- **POST `/refresh`**
  - **Body**: `{ "refresh_token": string }`
  - **Response**: `{ "access_token": string, "refresh_token": string }`

### User

- **POST `/user/push_token`**
  - **Body**: `{ "token": string }`
  - **Response**: `{ "message": "push token set successfully" }`

- **DELETE `/user`**
  - **Response**: `{ "message": "user and all associated data deleted successfully" }`
  - **Note**: This will delete the user, all their plants, and all associated reminders

### Plants

- **POST `/plant`**
  - **Body**: `{ "name": string, "note": string, "tagColor": string, "plantIcon": string }`
  - **Response**: `{ "plant": { ... } }`
  - **Note**: `plantIcon` must be one of the predefined icon types (e.g., "bananaPlant", "bigCactus", "bigPlant", "bigRose", etc.)

- **GET `/plant/:id`**
  - **Response**: `{ "plant": { ... } }`

- **GET `/plants`**
  - **Response**: `{ "plants": [ ... ] }`

- **PUT `/plant`**
  - **Body**: `{ "id": number, "name": string, "note": string, "tagColor": string, "plantIcon": string }`
  - **Response**: `{ "plant": { ... } }`
  - **Note**: `plantIcon` must be one of the predefined icon types

- **DELETE `/plant/:id`**
  - **Response**: HTTP 204 No Content
  - **Note**: This will delete the plant and all its associated reminders

### Reminders

- **POST `/plant/:id/reminder`**
  - **Body**: `{ "repeatType": "daily"|"weekly"|"monthly", "timeOfDay": "HH:MM" }`
  - **Response**: `{ "reminder": { ... } }`

- **GET `/plant/:id/reminders`**
  - **Response**: `{ "reminders": [ ... ] }`

- **GET `/plant/reminders`**
  - **Response**: `{ "reminders": [ ... ] }`
  - **Note**: Returns all reminders for the authenticated user

- **PUT `/plant/:id/reminder`**
  - **Body**: `{ "id": number, "repeatType": "daily"|"weekly"|"monthly", "timeOfDay": "HH:MM" }`
  - **Response**: `{ "reminder": { ... } }`

- **DELETE `/plant/:id/reminder/:reminderId`**
  - **Response**: HTTP 204 No Content

---

## Data Relationships

The application automatically handles data cleanup to maintain integrity:

- **Delete User**: Removes the user and all their plants and reminders
- **Delete Plant**: Removes the plant and all its associated reminders
- **Delete Reminder**: Removes only the specific reminder

This ensures that orphaned data is automatically cleaned up when parent records are deleted.

---

## Environment Variables

Create a `.env` file in the project root with the following variables:

### Required Variables
- `DB_URL` – PostgreSQL connection string (e.g., `postgres://user:password@localhost:5432/plantie`)
- `JWT_KEY` – Secret key for signing JWT tokens (should be at least 32 characters)

### Optional Variables
- `PORT` – Port for the server (default: `8080`)
- `ENV` – Environment mode (`production` or development, default: development)

### Firebase Configuration
Ensure you have a valid `firebase.json` file for Firebase Cloud Messaging in the project root.

---

## Getting Started

### Prerequisites
- Go 1.24.0 or higher
- PostgreSQL database
- Firebase project with Cloud Messaging enabled

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd plantie
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up your environment**
   - Create a `.env` file with required variables
   - Add your `firebase.json` file for Firebase

4. **Run the server**
   ```bash
   go run main.go
   ```

5. **API is now available at** `http://localhost:8080` (or your specified port)

### Development

The server includes:
- **Auto-migration**: Database tables are automatically created on startup
- **CORS support**: Configured for cross-origin requests
- **Graceful shutdown**: Proper cleanup on server termination
- **Cron jobs**: Automatic reminder scheduling
- **Automatic data cleanup**: Related data is removed when parent records are deleted

---

## Project Structure

```
plantie/
├── config/          # Database and app configuration
├── constants/       # Application constants
├── controllers/     # Request handlers
├── interfaces/      # Interface definitions
├── middleware/      # Authentication and other middleware
├── models/          # Database models
├── routes/          # API route definitions
├── services/        # Business logic services
├── utils/           # Utility functions
├── main.go          # Application entry point
└── README.md        # This file
```

---

## Database Schema

### Users
- `id` (Primary Key)
- `email` (Unique)
- `password` (Hashed)
- `name`
- `creation_date`
- `push_token`

### Plants
- `id` (Primary Key)
- `user_id` (Foreign Key with automatic cleanup)
- `name` (Required)
- `note`
- `tag_color` (Required)
- `plant_icon` (Required - one of predefined plant icon types)
- GORM timestamps

### Reminders
- `id` (Primary Key)
- `plant_id` (Foreign Key with automatic cleanup)
- `repeat` (daily/weekly/monthly)
- `time_of_day` (HH:MM format)
- `next_trigger_time`
- `user_id` (Foreign Key)
- GORM timestamps

---

## License

MIT License - see [LICENSE](LICENSE) file for details. 