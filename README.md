# Expense Tracker Server

A Go-based backend for the expense tracker application using Gin framework and Supabase as the database.

## Features

- User authentication and registration
- Transaction management (income and expenses)
- RESTful API design
- PostgreSQL database via Supabase
- JWT-based authentication
- Password hashing and validation

## Prerequisites

- Go 1.25 or higher
- Supabase account and project

## Setup

### 1. Clone the repository and install dependencies

```bash
go mod download
```

### 2. Environment Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env.local
```

Configure your `.env.local` file with your Supabase credentials:

```env
# Server Configuration
GIN_MODE=release
PORT=8080
JWT_SECRET=your_jwt_secret_key_here

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_supabase_anon_key
SUPABASE_DB_PASSWORD=your_database_password

# Alternative: Direct database connection
# DATABASE_URL=postgresql://postgres:your_password@db.your-project.supabase.co:5432/postgres
```

### 3. Database Connection

The application supports two ways to connect to Supabase:

#### Option 1: Individual Environment Variables

Set `SUPABASE_URL` and `SUPABASE_DB_PASSWORD` - the application will automatically construct the connection string.

#### Option 2: Direct Connection String (Recommended for Production)

Set `DATABASE_URL` with the complete PostgreSQL connection string from Supabase.

### 4. Getting Your Supabase Credentials

1. Go to [Supabase Dashboard](https://app.supabase.com)
2. Select your project
3. Go to **Settings** → **Database**
4. Copy the connection string and extract the required values
5. Go to **Settings** → **API** to get your API keys

### 5. Run the Application

```bash
# Development
go run cmd/api/main.go

# Production build
go build -o main cmd/api/main.go
./main
```

The server will start on the port specified in your environment (default: 8080).

## Database Schema

The application uses PostgreSQL with the following tables:

- `users` - User authentication and profile information
- `transactions` - Income and expense records

Database migrations are automatically applied on startup from the `internal/migrations/` directory.

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Users

- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update user profile
- `DELETE /api/user/account` - Delete user account

### Transactions

- `GET /api/transactions` - Get user transactions
- `POST /api/transactions` - Create new transaction
- `PUT /api/transactions/:id` - Update transaction
- `DELETE /api/transactions/:id` - Delete transaction

## Project Structure

```
server/
├── cmd/api/           # Application entry point
├── internal/
│   ├── auth/          # Authentication middleware
│   ├── connection/    # Database connection setup
│   ├── constants/     # Error constants
│   ├── handlers/      # HTTP handlers
│   ├── helpers/       # Utility functions
│   ├── migrations/    # Database migrations
│   ├── models/        # Data models and repository
│   └── router/        # Route definitions
└── .env.local         # Environment configuration
```

## Development

### Adding New Migrations

Place SQL files in `internal/migrations/` directory. They will be executed in alphabetical order on application startup.

### Testing

```bash
go test ./...
```

## Deployment

The application can be deployed to any platform that supports Go applications. Make sure to:

1. Set the `DATABASE_URL` environment variable with your production database connection string
2. Set a strong `JWT_SECRET`
3. Set `GIN_MODE=release` for production
4. Configure appropriate `PORT` if required by your platform

## Troubleshooting

### Connection Issues

- Verify your Supabase project is active
- Check that your database password is correct
- Ensure your IP is allowlisted in Supabase (or disable IP restrictions for development)
- Verify the project reference in your SUPABASE_URL matches your actual project

### Migration Issues

- Check that SQL files in `internal/migrations/` are valid
- Ensure proper permissions on your database
- Check application logs for specific error messages
