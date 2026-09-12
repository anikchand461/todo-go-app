<p align="center">
  <img src="mobile/assets/logo.png" width="120" alt="Go Todo logo" />
</p>

<h1 align="center">Go Todo</h1>

<p align="center">
  <i>A minimal, native Todo app written entirely in Go</i>
</p>

<p align="center">
  <b>One language, front to back. No JS, no Electron, no bloat.</b>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&logoColor=white" alt="Go version" />
  <img src="https://img.shields.io/badge/UI-Gio-5849BE" alt="Gio UI" />
  <img src="https://img.shields.io/badge/Backend-Gin-00ACD7" alt="Gin" />
  <img src="https://img.shields.io/badge/DB-PostgreSQL%20(Neon)-336791?logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Auth-JWT%20%2B%20bcrypt-000000" alt="JWT + bcrypt" />
</p>

<p align="center">
  <a href="#-features">Features</a> ·
  <a href="#-screenshots">Screenshots</a> ·
  <a href="#-tech-stack">Tech Stack</a> ·
  <a href="#-project-structure">Project Structure</a> ·
  <a href="#-getting-started">Getting Started</a> ·
  <a href="#-environment">Environment</a>
</p>

---

## 🔖 About

**Go Todo** is a small full-stack task manager: a **Gio**-based native client talking to a **Gin** REST API, backed by **PostgreSQL**. Every layer — frontend, backend, and glue — is written in Go.

## 📱 Screenshots

<p align="center">
  <img src="mobile/assets/login.png" width="260" alt="Login screen" />
  &nbsp;&nbsp;&nbsp;
  <img src="mobile/assets/home.png" width="260" alt="Home screen" />
</p>

## ✨ Features

- User registration
- User login
- JWT authentication
- Create todos
- Complete todos
- Delete todos
- User-specific todos
- Logout

## 🧱 Tech Stack

| Layer          | Technology         |
|----------------|--------------------|
| Frontend       | Go + Gio           |
| Backend        | Go + Gin           |
| Database       | PostgreSQL (Neon)  |
| Authentication | JWT + bcrypt        |

## 📂 Project Structure

```text
todo-go/
├── backend/    # REST API and database
└── mobile/     # Gio application
    └── assets/ # logo.png · login.png · home.png
```

## 🚀 Getting Started

### Run Backend

```bash
cd backend
go run .
```

### Run Mobile App

```bash
cd mobile
go run .
```

The mobile app uses `http://localhost:8000` for the backend by default.

> For Android, use your computer's local network IP instead of `localhost`.

## ⚙️ Environment

Create `backend/.env` using `backend/.env.example` and configure your PostgreSQL connection and JWT secret.