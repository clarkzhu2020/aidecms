# AideCMS Project

A comprehensive enterprise-level content management system consisting of a Go backend, User Frontend, and Admin Dashboard.

## Project Structure

- **admin**: Management frontend (Vue 3 + Element Plus + Vite)
- **web**: User frontend (Vue 3 + Vite)
- **backend**: Backend API (Golang + Hertz Framework)
- **doc**: Documentation

## Tech Stack

- **Backend**: Golang 1.24, Hertz, GORM, Redis client
- **Frontend**: Vue 3, TypeScript, Vite, Element Plus (Admin)
- **Database**: PostgreSQL 15 (Exclusive support)
- **Analytics**: ClickHouse (for DEX quotes/market data)
- **Cache**: Redis 7
- **Containerization**: Docker & Docker Compose

## Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your machine.
- [Docker Compose](https://docs.docker.com/compose/install/) (included with Docker Desktop).

## Quick Start

1. **Clone the repository** (if not already done).

2. **Start the application**:
   Run the following command in the project root to build and start all services:
   ```bash
   docker-compose up -d --build
   ```
   This will start PostgreSQL, Redis, ClickHouse, Backend, Web, and Admin services.

3. **Verify Status**:
   ```bash
   docker-compose ps
   ```

## Accessing Services

| Service | URL | Description |
|---------|-----|-------------|
| **Web Frontend** | [http://localhost:8080](http://localhost:8080) | Public facing user website |
| **Admin Dashboard** | [http://localhost:8081](http://localhost:8081) | Administration interface |
| **Backend API** | [http://localhost:8888](http://localhost:8888) | REST API endpoints |
| **ClickHouse HTTP** | [http://localhost:8123](http://localhost:8123) | ClickHouse HTTP Interface |
| **Swagger Docs** | [http://localhost:8888/swagger/index.html](http://localhost:8888/swagger/index.html) | API Documentation |

## Default Credentials

### Admin Dashboard
- **Username**: `admin`
- **Password**: `admin`
*(Note: These are mock credentials defined in the frontend login component)*

### Database (PostgreSQL)
- **Host**: `localhost` (port 5432)
- **User**: `aidecms`
- **Password**: `aidecms_secret`
- **Database**: `aidecms`

### Redis
- **Host**: `localhost` (port 6379)
- **Password**: (none)

## Configuration

### Backend
The backend configuration is handled via environment variables passed in `docker-compose.yml`. Key variables:
- `DB_CONNECTION`: `postgres`
- `DB_HOST`: `pgsql` (service name)
- `REDIS_HOST`: `redis` (service name)

### Frontend
The frontend applications are built as static assets and served via Nginx.
- To change the API endpoint in the frontend, you would typically update the `.env` file in the respective frontend directory (e.g., `VITE_API_BASE_URL`) and rebuild. Currently both frontends are running in decoupled mode.

## Stops and Cleanup

To stop the services:
```bash
docker-compose down
```

To stop and remove volumes (reset database):
```bash
docker-compose down -v
```

## Development

If you want to run services locally without Docker:
- **Backend**: `cd backend && go run .`
- **Web**: `cd web && pnpm install && pnpm dev`
- **Admin**: `cd admin && pnpm install && pnpm dev`
