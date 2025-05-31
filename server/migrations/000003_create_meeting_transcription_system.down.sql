-- Drop meeting transcription and action item system tables
-- Migration: 000003_create_meeting_transcription_system (DOWN)

-- Drop tables in reverse order to handle foreign key constraints
DROP TABLE IF EXISTS processing_jobs;
DROP TABLE IF EXISTS integration_configs;
DROP TABLE IF EXISTS ticket_references;
DROP TABLE IF EXISTS action_items;
DROP TABLE IF EXISTS transcript_segments;
DROP TABLE IF EXISTS transcriptions;
DROP TABLE IF EXISTS bot_sessions;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS meetings; 