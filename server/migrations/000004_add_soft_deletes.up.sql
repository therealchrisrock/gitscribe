-- Add soft delete support to tables missing deleted_at column
-- Migration: 000004_add_soft_deletes

-- Add deleted_at to transcriptions table
ALTER TABLE transcriptions 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to transcript_segments table  
ALTER TABLE transcript_segments 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to participants table
ALTER TABLE participants 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to bot_sessions table
ALTER TABLE bot_sessions 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to action_items table (already has updated_at, just need deleted_at)
ALTER TABLE action_items 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to ticket_references table
ALTER TABLE ticket_references 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to integration_configs table  
ALTER TABLE integration_configs 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add deleted_at to processing_jobs table
ALTER TABLE processing_jobs 
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Create indexes for soft delete queries
CREATE INDEX idx_transcriptions_deleted_at ON transcriptions(deleted_at);
CREATE INDEX idx_transcript_segments_deleted_at ON transcript_segments(deleted_at);
CREATE INDEX idx_participants_deleted_at ON participants(deleted_at);
CREATE INDEX idx_bot_sessions_deleted_at ON bot_sessions(deleted_at);
CREATE INDEX idx_action_items_deleted_at ON action_items(deleted_at);
CREATE INDEX idx_ticket_references_deleted_at ON ticket_references(deleted_at);
CREATE INDEX idx_integration_configs_deleted_at ON integration_configs(deleted_at);
CREATE INDEX idx_processing_jobs_deleted_at ON processing_jobs(deleted_at);

-- Add comments
COMMENT ON COLUMN transcriptions.deleted_at IS 'Soft delete timestamp for transcriptions';
COMMENT ON COLUMN transcript_segments.deleted_at IS 'Soft delete timestamp for transcript segments';
COMMENT ON COLUMN participants.deleted_at IS 'Soft delete timestamp for participants';
COMMENT ON COLUMN bot_sessions.deleted_at IS 'Soft delete timestamp for bot sessions';
COMMENT ON COLUMN action_items.deleted_at IS 'Soft delete timestamp for action items';
COMMENT ON COLUMN ticket_references.deleted_at IS 'Soft delete timestamp for ticket references';
COMMENT ON COLUMN integration_configs.deleted_at IS 'Soft delete timestamp for integration configs';
COMMENT ON COLUMN processing_jobs.deleted_at IS 'Soft delete timestamp for processing jobs'; 