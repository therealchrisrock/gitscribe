-- Create meeting transcription and action item system tables
-- Migration: 000003_create_meeting_transcription_system

-- Meetings table
CREATE TABLE meetings (
    id VARCHAR(128) PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('zoom', 'google_meet', 'microsoft_teams', 'generic')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('scheduled', 'in_progress', 'completed', 'failed')),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NULL,
    meeting_url TEXT NOT NULL,
    bot_join_url TEXT NULL,
    recording_path TEXT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- Participants table
CREATE TABLE participants (
    id VARCHAR(128) PRIMARY KEY,
    meeting_id VARCHAR(128) NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NULL,
    role VARCHAR(100) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Bot sessions table
CREATE TABLE bot_sessions (
    id VARCHAR(128) PRIMARY KEY,
    meeting_id VARCHAR(128) NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('joining', 'active', 'recording', 'completed', 'failed')),
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL,
    left_at TIMESTAMP WITH TIME ZONE NULL,
    bot_user_id VARCHAR(255) NULL,
    metadata JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Transcriptions table
CREATE TABLE transcriptions (
    id VARCHAR(128) PRIMARY KEY,
    meeting_id VARCHAR(128) NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    audio_file_path TEXT NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    content TEXT NULL,
    confidence DECIMAL(5,4) NULL,
    provider VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Transcript segments table
CREATE TABLE transcript_segments (
    id VARCHAR(128) PRIMARY KEY,
    transcription_id VARCHAR(128) NOT NULL REFERENCES transcriptions(id) ON DELETE CASCADE,
    speaker VARCHAR(255) NULL,
    text TEXT NOT NULL,
    start_time DECIMAL(10,3) NOT NULL,
    end_time DECIMAL(10,3) NOT NULL,
    confidence DECIMAL(5,4) NULL,
    sequence_number INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Action items table
CREATE TABLE action_items (
    id VARCHAR(128) PRIMARY KEY,
    meeting_id VARCHAR(128) NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
    transcription_id VARCHAR(128) NOT NULL REFERENCES transcriptions(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    assignee VARCHAR(255) NULL,
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    due_date TIMESTAMP WITH TIME ZONE NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('extracted', 'pending', 'approved', 'created', 'rejected')),
    context TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ticket references table
CREATE TABLE ticket_references (
    id VARCHAR(128) PRIMARY KEY,
    action_item_id VARCHAR(128) NOT NULL REFERENCES action_items(id) ON DELETE CASCADE,
    system VARCHAR(50) NOT NULL,
    ticket_id VARCHAR(255) NOT NULL,
    ticket_url TEXT NOT NULL,
    project_key VARCHAR(100) NULL,
    reference_type VARCHAR(20) NOT NULL CHECK (reference_type IN ('existing', 'created')),
    metadata JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Integration configs table
CREATE TABLE integration_configs (
    id VARCHAR(128) PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_type VARCHAR(50) NOT NULL,
    provider_name VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, provider_type, provider_name)
);

-- Processing jobs table
CREATE TABLE processing_jobs (
    id VARCHAR(128) PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(128) NOT NULL,
    job_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    payload JSONB NOT NULL,
    error_message TEXT NULL,
    retry_count INTEGER DEFAULT 0,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NULL,
    completed_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_meetings_user_id ON meetings(user_id);
CREATE INDEX idx_meetings_status ON meetings(status);
CREATE INDEX idx_meetings_start_time ON meetings(start_time);
CREATE INDEX idx_meetings_type ON meetings(type);

CREATE INDEX idx_participants_meeting_id ON participants(meeting_id);

CREATE INDEX idx_bot_sessions_meeting_id ON bot_sessions(meeting_id);
CREATE INDEX idx_bot_sessions_session_id ON bot_sessions(session_id);
CREATE INDEX idx_bot_sessions_status ON bot_sessions(status);

CREATE INDEX idx_transcriptions_meeting_id ON transcriptions(meeting_id);
CREATE INDEX idx_transcriptions_status ON transcriptions(status);
CREATE INDEX idx_transcriptions_provider ON transcriptions(provider);

CREATE INDEX idx_transcript_segments_transcription_id ON transcript_segments(transcription_id);
CREATE INDEX idx_transcript_segments_sequence ON transcript_segments(transcription_id, sequence_number);

CREATE INDEX idx_action_items_meeting_id ON action_items(meeting_id);
CREATE INDEX idx_action_items_transcription_id ON action_items(transcription_id);
CREATE INDEX idx_action_items_status ON action_items(status);
CREATE INDEX idx_action_items_assignee ON action_items(assignee);
CREATE INDEX idx_action_items_priority ON action_items(priority);

CREATE INDEX idx_ticket_references_action_item_id ON ticket_references(action_item_id);
CREATE INDEX idx_ticket_references_system_ticket ON ticket_references(system, ticket_id);

CREATE INDEX idx_integration_configs_user_provider ON integration_configs(user_id, provider_type);
CREATE INDEX idx_integration_configs_active ON integration_configs(is_active);

CREATE INDEX idx_processing_jobs_status ON processing_jobs(status);
CREATE INDEX idx_processing_jobs_scheduled ON processing_jobs(scheduled_at);
CREATE INDEX idx_processing_jobs_entity ON processing_jobs(entity_type, entity_id);
CREATE INDEX idx_processing_jobs_job_type ON processing_jobs(job_type);

-- Add table comments
COMMENT ON TABLE meetings IS 'Stores meeting information and metadata';
COMMENT ON TABLE participants IS 'Stores meeting participants information';
COMMENT ON TABLE bot_sessions IS 'Tracks bot participation in meetings';
COMMENT ON TABLE transcriptions IS 'Stores meeting transcription data';
COMMENT ON TABLE transcript_segments IS 'Stores individual segments of transcriptions with speaker attribution';
COMMENT ON TABLE action_items IS 'Stores extracted action items from meeting transcriptions';
COMMENT ON TABLE ticket_references IS 'Links action items to external ticketing systems';
COMMENT ON TABLE integration_configs IS 'Stores user configuration for external service integrations';
COMMENT ON TABLE processing_jobs IS 'Tracks asynchronous processing tasks'; 