package repositories

import (
	"teammate/server/seedwork/domain"
	"testing"
)

// TestRepositoryModelCompliance verifies that our repository models implement RepositoryModel interface
func TestRepositoryModelCompliance(t *testing.T) {
	// Test Meeting model
	var meeting domain.RepositoryModel = &Meeting{}
	meeting.SetID("test-meeting-id")

	if meeting.GetID() != "test-meeting-id" {
		t.Errorf("Meeting GetID failed: expected 'test-meeting-id', got '%s'", meeting.GetID())
	}

	if meeting.TableName() != "meetings" {
		t.Errorf("Meeting TableName failed: expected 'meetings', got '%s'", meeting.TableName())
	}

	// Test Participant model
	var participant domain.RepositoryModel = &Participant{}
	participant.SetID("test-participant-id")

	if participant.GetID() != "test-participant-id" {
		t.Errorf("Participant GetID failed: expected 'test-participant-id', got '%s'", participant.GetID())
	}

	if participant.TableName() != "participants" {
		t.Errorf("Participant TableName failed: expected 'participants', got '%s'", participant.TableName())
	}

	// Test BotSession model
	var botSession domain.RepositoryModel = &BotSession{}
	botSession.SetID("test-bot-session-id")

	if botSession.GetID() != "test-bot-session-id" {
		t.Errorf("BotSession GetID failed: expected 'test-bot-session-id', got '%s'", botSession.GetID())
	}

	if botSession.TableName() != "bot_sessions" {
		t.Errorf("BotSession TableName failed: expected 'bot_sessions', got '%s'", botSession.TableName())
	}
}
