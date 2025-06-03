-- Remove soft delete support (rollback for 000004_add_soft_deletes)
-- Migration: 000004_add_soft_deletes DOWN

-- Drop indexes first
DROP INDEX IF EXISTS idx_transcriptions_deleted_at;
DROP INDEX IF EXISTS idx_transcript_segments_deleted_at;
DROP INDEX IF EXISTS idx_participants_deleted_at;
DROP INDEX IF EXISTS idx_bot_sessions_deleted_at;
DROP INDEX IF EXISTS idx_action_items_deleted_at;
DROP INDEX IF EXISTS idx_ticket_references_deleted_at;
DROP INDEX IF EXISTS idx_integration_configs_deleted_at;
DROP INDEX IF EXISTS idx_processing_jobs_deleted_at;

-- Remove deleted_at columns
ALTER TABLE transcriptions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE transcript_segments DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE participants DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE bot_sessions DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE action_items DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE ticket_references DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE integration_configs DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE processing_jobs DROP COLUMN IF EXISTS deleted_at; 