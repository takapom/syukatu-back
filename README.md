# Go Shop Backend

## Local Development Setup

1. Copy the environment example file:
```bash
cp .env.example .env
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run .
```

The server will start on http://localhost:8080

## Environment Variables

- `DATABASE_URL`: Database connection string (sqlite:// for local, postgres:// for production)
- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT tokens
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins
- `ENVIRONMENT`: Environment mode (development/production)

## Deployment on Render

1. Fork or push this repository to GitHub
2. Connect your GitHub account to Render
3. Create a new Web Service on Render
4. Select this repository
5. Render will automatically detect the `render.yaml` configuration
6. The PostgreSQL database will be created automatically
7. Update `CORS_ALLOWED_ORIGINS` in the environment variables to match your frontend URL

## Database

- Local: SQLite (./example.db)
- Production: PostgreSQL (automatically provisioned by Render)

## API Endpoints

### Public Routes
- `POST /register` - User registration
- `POST /login` - User login

### Protected Routes (requires JWT token)
- Company Lists: `/company_lists` (GET, POST, PUT, DELETE)
- Internships: `/internships` (GET, POST, PUT, DELETE)
- Posts: `/posts` (GET, POST, DELETE)
- Comments: `/posts/:id/comments` (POST)
- Likes: `/posts/:id/like` (POST, DELETE)