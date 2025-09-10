-- Enable UUID extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- gen_random_uuid() is provided by pgcrypto on many PG versions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,  -- Limited to 100 chars for performance
    email VARCHAR(255) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),  -- Basic email validation
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),  -- Must be positive
    description VARCHAR(500),  -- Limited to 500 chars
    note TEXT,
    category VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Existing indexes
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_id ON users(id);

-- New indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);  -- For login queries
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);  -- For category filters
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);  -- For type filters
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);  -- For date queries
CREATE INDEX IF NOT EXISTS idx_transactions_user_created ON transactions(user_id, created_at);  -- Composite for user + date
