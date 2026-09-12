# Go Todo

A simple Todo application built with Go.

## Tech Stack

-   **Frontend:** Go + Gio
-   **Backend:** Go + Gin
-   **Database:** PostgreSQL
-   **Authentication:** JWT + bcrypt

## Features

-   User registration
-   User login
-   JWT authentication
-   Create todos
-   Complete todos
-   Delete todos
-   User-specific todos
-   Logout

## Project Structure

``` text
todo-go/
├── backend/    # REST API and database
└── mobile/     # Gio application
```

## Run Backend

``` bash
cd backend
go run .
```

## Run Mobile App

``` bash
cd mobile
go run .
```

The mobile app uses:

``` text
http://localhost:8000
```

for the backend by default.

> For Android, use your computer's local network IP instead of
> `localhost`.

## Environment

Create `backend/.env` using `backend/.env.example` and configure your
PostgreSQL connection and JWT secret.
