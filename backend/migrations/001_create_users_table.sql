-- Migration: 001_create_users_table.sql
-- Description: Create users table for storing user information

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    full_name VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    organization_name VARCHAR(255),
    consent_pd BOOLEAN NOT NULL DEFAULT false,
    location_id INTEGER,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_location_id ON users(location_id);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data for testing (optional)
INSERT INTO users (telegram_id, full_name, phone, email, organization_name, consent_pd, location_id, role)
VALUES
    (123456789, 'Иван Иванов', '+79991234567', 'ivan@example.com', 'ООО Ромашка', true, 1, 'user'),
    (987654321, 'Петр Петров', '+79998765432', 'petr@example.com', 'ИП Петров', true, 2, 'clinic_admin')
ON CONFLICT (telegram_id) DO NOTHING;
