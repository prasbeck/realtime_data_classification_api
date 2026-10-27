# Real-time Data Classification API

A Go HTTP service for managing classification rules, with user signup and login, JWT-protected rule endpoints, and a WebSocket hub that broadcasts messages to every connected client.

> **Status: work in progress.** The rules package (`dsl/`) is still being refactored, and the service does not build yet. See [Roadmap](#roadmap).

## Features

- User signup and login backed by PostgreSQL, with bcrypt password hashing
- CRUD endpoints for rules, protected by JWT bearer-token middleware
- Configurable cap on the number of rules (default 10)
- In-memory rule store guarded by a read/write mutex
- WebSocket hub (`/ws`) with register, unregister, and broadcast channels

## Stack

Go 1.22 · `net/http` · PostgreSQL (`lib/pq`) · gorilla/websocket · JWT

## API

| Intended method | Route | Auth | Description |
|---|---|---|---|
| `POST` | `/signup` | none | Create a user |
| `POST` | `/login` | none | Log in |
| `POST` | `/rules/create` | JWT | Create a rule |
| `GET` | `/rules` | JWT | List rules |
| `GET` | `/rules/get?name=` | JWT | Get one rule |
| `PUT` | `/rules/update?name=` | JWT | Update a rule |
| `DELETE` | `/rules/delete?name=` | JWT | Delete a rule |
| `GET` | `/ws` | none | Open a WebSocket connection |

## Getting started

Requirements: Go 1.22+ and a PostgreSQL database.

Once the build is fixed, configuration will come from environment variables:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/rtmc?sslmode=disable"
export JWT_SECRET="change-me"
go run .
```

The server listens on port `8080` and creates the `users` table on startup.

## Roadmap

- [ ] Define the `Rule` type and fix the rule store so the project builds
- [ ] Read the database URL and JWT secret from environment variables
- [ ] Issue a JWT from `/login`, and stop returning the password hash
- [ ] Start the WebSocket write loop for each client
- [ ] Classify incoming data against rules and broadcast results over `/ws`
- [ ] Add tests
