# Plantie API

Backend for a plant-watering reminder app. Manage plants, schedule reminders, and receive push notifications to keep your plants healthy.

## Features

- JWT auth (signup, login, refresh)
- Plant CRUD
- Reminders: create, update, list, delete
- Push notifications via Firebase Cloud Messaging
- Scheduler to dispatch reminders

## Tech

- Go + Gin
- PostgreSQL with GORM
- Firebase Cloud Messaging
- GoCron scheduler

## Getting started

1) Install deps

```bash
go mod download
```

2) Configure environment

- DB_URL: Postgres connection string
- JWT_KEY: secret for JWT signing
- FIREBASE_PATH: path to Firebase service account JSON

3) Run

```bash
go run .
```

## Testing

Run the test suite:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -race -coverprofile=coverage.out -covermode=atomic
```

## CI/CD

This project uses GitHub Actions for continuous integration. The CI pipeline:

- Runs on all pull requests to `main` and `development` branches
- Runs on pushes to `main` and `development` branches  
- Tests the code with `go test ./...`
- Generates coverage reports
- Prevents merging PRs if tests fail (requires branch protection setup)

**Note**: To complete the setup and enforce test passing before merging, see [`.github/BRANCH_PROTECTION_SETUP.md`](.github/BRANCH_PROTECTION_SETUP.md) for branch protection configuration instructions.

## API overview

All endpoints (except /ping, /login, /signup, /refresh) require Authorization: Bearer <access_token>.

### Health
- GET /ping
  - Response:
    ```json
    "pong"
    ```

### Auth
- POST /signup
  - Response:
    ```json
    { "access_token": "...", "refresh_token": "...", "user": { /* ... */ } }
    ```
- POST /login
  - Response:
    ```json
    { "access_token": "...", "refresh_token": "...", "user": { /* ... */ } }
    ```
- POST /refresh
  - Response:
    ```json
    { "access_token": "...", "refresh_token": "..." }
    ```

### Users
- GET /user/me
  - Response:
    ```json
    { "user": { /* ... */ } }
    ```
- POST /user/push_token
  - Body:
    ```json
    { "token": "..." }
    ```
  - Response:
    ```json
    { "message": "push token set successfully" }
    ```
- DELETE /user
  - Response:
    ```json
    { "message": "user and all associated data deleted successfully" }
    ```

### Plants
- POST /plant
  - Body:
    ```json
    { "name": "...", "note": "...", "tagColor": "...", "plantIcon": "..." }
    ```
  - Response:
    ```json
    { "plant": { /* ... */ } }
    ```
- GET /plant/:id
  - Response:
    ```json
    { "plant": { /* ... */ } }
    ```
- GET /plants
  - Response:
    ```json
    { "plants": [ /* ... */ ] }
    ```
- PUT /plant/:id
  - Body:
    ```json
    { "name": "...", "note": "...", "tagColor": "...", "plantIcon": "..." }
    ```
  - Response:
    ```json
    { "plant": { /* ... */ } }
    ```
- DELETE /plant/:id
  - Response:
    ```
    204 No Content
    ```

### Reminders

- POST /plant/:id/reminder
  - Body:
    ```json
    { "plantId": 1, "repeatType": "daily|weekly|monthly", "timeOfDay": "HH:MM", "dayOfWeek": 2, "dayOfMonth": 15 }
    ```
  - Response:
    ```json
    { "reminder": { /* ... */ } }
    ```

- PUT /plant/:id/reminder
  - Body:
    ```json
    { "id": 10, "plantId": 1, "repeatType": "weekly", "timeOfDay": "08:00", "dayOfWeek": 3 }
    ```
  - Response:
    ```json
    { "reminder": { /* ... */ } }
    ```

- GET /plant/:id/reminders
  - Response:
    ```json
    { "reminders": [ /* ... */ ] }
    ```
- GET /plant/reminders
  - Response:
    ```json
    { "reminders": [ /* ... */ ] }
    ```
- DELETE /plant/:id/reminder/:reminderId
  - Response:
    ```
    204 No Content
    ```
- POST /reminders/test
  - Response:
    ```json
    { "staus": "ok" }
    ```
  - Note: Intentional spelling as returned by the current implementation

Notes
- dayOfWeek is required for weekly reminders (0-6, Sunday=0)
- dayOfMonth is required for monthly reminders (1-31)
- For daily reminders, dayOfWeek/dayOfMonth must be omitted

## Project layout

- config/: configuration and DB setup
- constants/: enums and validation helpers
- container/: DI wiring
- controllers/: HTTP handlers
- dto/: request/response DTOs
- middleware/: auth middleware
- models/: GORM models
- routes/: router setup
- service/: business logic
- utils/: helpers (jwt, notifier, etc.)