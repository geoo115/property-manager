#!/bin/bash

# Start PostgreSQL Database Script for Property Management System

echo "🚀 Starting PostgreSQL Database..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install Docker Compose."
    exit 1
fi

# Start PostgreSQL service
echo "📦 Starting PostgreSQL container..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
until docker exec postgres-pms pg_isready -U postgres; do
    echo "⏳ PostgreSQL is unavailable - sleeping"
    sleep 2
done

echo "✅ PostgreSQL is ready!"
echo "📊 Database connection details:"
echo "   Host: localhost"
echo "   Port: 5432"
echo "   Database: property_management"
echo "   User: postgres"
echo "   Password: your_password"
echo ""
echo "🔗 You can connect to the database using:"
echo "   psql -h localhost -U postgres -d property_management"
echo ""
echo "🏃 To start the application:"
echo "   go run cmd/main.go"
echo ""
echo "🛑 To stop the database:"
echo "   docker-compose down"
