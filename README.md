# Chirphy API

Chirphy is a backend server for managing users, authentication, and chirps (posts). This document provides an overview of the available endpoints and how to set up the project.

## Table of Contents

- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [Authentication](#authentication)
  - [Users](#users)
  - [Chirps](#chirps)
  - [Polka Webhooks](#polka-webhooks)
  - [Admin](#admin)

---

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/Tran-Duc-Hoa/chirphy.git
   cd chirphy
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory and configure the required environment variables ([see below](#environment-variables)).

4. Run the application:

   ```bash
   go run main.go
   ```

---

## Environment Variables

The following environment variables must be set in the `.env` file:

- `DB_URL`: The PostgreSQL database connection string.
- `PLATFORM`: The platform name (e.g., `production`, `dev`).
- `JWT_SECRET`: The secret key for signing JWT tokens.
- `POLKA_KEY`: The API key for Polka integration.

Example `.env` file:

```bash
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy"
PLATFORM="dev"
JWT_SECRET="n6K2+LcKatcCONgBb6wuI/66hz/qbOJj0chzkNbVqFDmDBqoe3LmbRMax2xd47VMzb+ZD3ViGyW1cK+fq6U7FA=="
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```
