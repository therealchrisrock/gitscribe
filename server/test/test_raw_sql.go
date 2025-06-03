package main

import (
	"context"
	"log"
	"time"

	"teammate/server/seedwork/infrastructure/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying SQL DB
	db, err := database.GetDB().DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL DB: %v", err)
	}

	log.Println("Testing raw SQL operations...")

	// Test 1: Simple query
	var count int
	err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query users count: %v", err)
	}
	log.Printf("✅ Users count: %d", count)

	// Test 2: Check meetings table
	err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM meetings").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query meetings count: %v", err)
	}
	log.Printf("✅ Meetings count: %d", count)

	// Test 3: Try to insert a test meeting
	testMeetingID := "test_meeting_" + time.Now().Format("20060102150405")
	testUserID := "test_user_123"

	// First insert a test user if not exists
	_, err = db.ExecContext(context.Background(), `
		INSERT INTO users (id, email, name, auth_provider, firebase_uid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`, testUserID, "test@example.com", "Test User", "firebase", "firebase_123", time.Now(), time.Now())
	if err != nil {
		log.Fatalf("Failed to insert test user: %v", err)
	}
	log.Printf("✅ Test user inserted/exists")

	// Now insert a test meeting
	_, err = db.ExecContext(context.Background(), `
		INSERT INTO meetings (id, user_id, title, type, status, start_time, meeting_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, testMeetingID, testUserID, "Test Meeting", "generic", "scheduled", time.Now(), "ws://test", time.Now(), time.Now())
	if err != nil {
		log.Fatalf("Failed to insert test meeting: %v", err)
	}
	log.Printf("✅ Test meeting inserted: %s", testMeetingID)

	// Cleanup
	_, err = db.ExecContext(context.Background(), "DELETE FROM meetings WHERE id = $1", testMeetingID)
	if err != nil {
		log.Printf("Warning: Failed to cleanup test meeting: %v", err)
	} else {
		log.Printf("✅ Test meeting cleaned up")
	}

	log.Println("✅ All raw SQL tests passed!")
}
