# Plantie – Plant Watering Reminder API

Plantie is a backend API for a plant watering reminder application built with clean architecture principles. It helps users manage their plants, set up watering reminders, and receive push notifications so they never forget to care for their green friends.

**Latest Update**: Refactored to clean architecture with dependency injection, comprehensive testing, and interface-based design.

---

## Recent Architectural Improvements

This project was recently refactored from a basic MVC pattern to a clean architecture with the following improvements:

### What Changed
- **Service Layer**: Extracted business logic from controllers into dedicated service layer
- **Dependency Injection**: Controllers now depend on service interfaces rather than concrete implementations
- **Interface Design**: Created service interfaces (`PlantServiceInterface`, `UserServiceInterface`, `ReminderServiceInterface`)
- **DTO Pattern**: Implemented Data Transfer Objects for all API requests and responses
- **Testing**: Added comprehensive unit tests with mock implementations
- **Container**: Added dependency injection container for clean service management

### Benefits
- **Testability**: Interface-based design allows for easy mocking and unit testing
- **Maintainability**: Clear separation of concerns between HTTP, business logic, and data layers
- **Extensibility**: Easy to add new features or modify existing behavior
- **Type Safety**: Strong typing with DTO validation ensures API contract integrity
- **Code Quality**: Improved code organization and reduced coupling between components

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
- **Clean Architecture**: Dependency injection with service interfaces for maintainability
- **Comprehensive Testing**: Full unit test coverage with mock implementations
- **Type Safety**: Strong typing with DTO patterns for API contracts

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

## Testing

The project includes comprehensive unit tests for all controllers and services:

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./controllers/
go test ./service/
```

### Test Structure
- **Controller Tests**: Mock service dependencies and test HTTP handling
- **Service Tests**: Test business logic with mock database interactions
- **Mock Implementations**: Comprehensive mocks for all service interfaces
- **Test Coverage**: Full coverage of happy paths and error scenarios

### Test Features
- Interface-based mocking for clean test isolation
- HTTP request/response testing with Gin test mode
- Validation testing for all API endpoints
- Error handling verification

---

## Architecture

The project follows the **Clean Architecture** pattern with clear separation of concerns:

- **Models**: Data structures and database logic (`models/`)
- **Controllers**: HTTP request handling and response formatting (`controllers/`)
- **Services**: Business logic implementation with dependency injection via interfaces (`service/`)
- **DTOs**: Data Transfer Objects for API requests and responses (`dto/`)
- **Routes**: API endpoint definitions (`routes/`)
- **Middleware**: Authentication and other request pre-processing (`middleware/`)
- **Utils**: Utility functions (JWT, password hashing, notifications, etc.)
- **Config**: Database and application configuration (`config/`)
- **Constants**: Application constants (`constants/`)

### Key Architectural Improvements

- **Dependency Injection**: Controllers use service interfaces for better testability and maintainability
- **Interface-based Design**: Services implement interfaces allowing for easy mocking and testing
- **Comprehensive Testing**: Full test coverage for controllers and services with mock implementations
- **Clean Separation**: Clear boundaries between HTTP layer, business logic, and data layer
- **DTO Pattern**: Data Transfer Objects ensure type safety and clear API contracts
- **Dependency Injection Container**: Centralized dependency management with interfaces

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

- **GET `/user/me`**
  - **Response**: `{ "user": { "id": number, "email": string, "name": string, "createdAt": string } }`
  - **Note**: Returns the profile information of the authenticated user

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
  - **Note**: `plantIcon` must be one of: "bananaPlant", "bigCactus", "bigPlant", "bigRose", "chilliPlant", "daisy", "flowerBed", "flower", "leafyPlant", "mediumPlant", "redTulip", "seaweedPlant", "shortPlant", "skinnyPlant", "smallCactus", "smallPlant", "smallRose", "spikyPlant", "tallPlant", "threeFlowers", "twoFlowers", "twoPlants", "whiteFlower", "yellowTulip"

- **GET `/plant/:id`**
  - **Response**: `{ "plant": { ... } }`

- **GET `/plants`**
  - **Response**: `{ "plants": [ ... ] }`

- **PUT `/plant/:id`**
  - **Body**: `{ "name": string, "note": string, "tagColor": string, "plantIcon": string }`
  - **Response**: `{ "plant": { ... } }`
  - **Note**: `plantIcon` must be one of the predefined icon types listed above. The plant ID is specified in the URL path.

- **DELETE `/plant/:id`**
  - **Response**: HTTP 204 No Content
  - **Note**: This will delete the plant and all its associated reminders

### Reminders

- **POST `/plant/:id/reminder`**
  - **Body**: `{ "plantId": number, "repeatType": "daily"|"weekly"|"monthly", "timeOfDay": "HH:MM" }`
  - **Response**: `{ "reminder": { ... } }`
  - **Note**: The `plantId` in the body should match the `:id` in the URL path

- **GET `/plant/:id/reminders`**
  - **Response**: `{ "reminders": [ ... ] }`

- **GET `/plant/reminders`**
  - **Response**: `{ "reminders": [ ... ] }`
  - **Note**: Returns all reminders for the authenticated user

- **PUT `/plant/:id/reminder`**
  - **Body**: `{ "id": number, "plantId": number, "repeatType": "daily"|"weekly"|"monthly", "timeOfDay": "HH:MM" }`
  - **Response**: `{ "reminder": { ... } }`
  - **Note**: The `plantId` in the body should match the `:id` in the URL path

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
   go run .
   ```

5. **Run tests** (optional)
   ```bash
   go test ./...
   ```

6. **API is now available at** `http://localhost:8080` (or your specified port)

### Development

The server includes:
- **Auto-migration**: Database tables are automatically created on startup
- **CORS support**: Configured for cross-origin requests
- **Graceful shutdown**: Proper cleanup on server termination
- **Cron jobs**: Automatic reminder scheduling
- **Automatic data cleanup**: Related data is removed when parent records are deleted
- **Dependency injection**: Services are injected via interfaces for better testability
- **Unit testing**: Run tests with `go test ./...` for full test suite coverage
- **Clean architecture**: Separation of concerns between HTTP, business logic, and data layers

---

## Project Structure

```
plantie/
├── config/          # Database and app configuration
├── constants/       # Application constants
├── container/       # Dependency injection container
├── controllers/     # HTTP request handlers
│   ├── *_test.go   # Controller unit tests with mocks
├── dto/             # Data Transfer Objects
├── middleware/      # Authentication and other middleware
├── models/          # Database models
├── routes/          # API route definitions
├── service/         # Business logic layer
│   ├── interfaces.go # Service interface definitions
│   ├── *_test.go   # Service unit tests
├── utils/           # Utility functions
├── firebase.json    # Firebase configuration
├── main.go          # Application entry point
├── go.mod           # Go module file
├── go.sum           # Go dependencies
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