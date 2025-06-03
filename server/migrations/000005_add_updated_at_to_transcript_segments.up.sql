-- Add updated_at column to transcript_segments table
-- Migration: 000005_add_updated_at_to_transcript_segments

-- Add updated_at column to transcript_segments table
ALTER TABLE transcript_segments 
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Create trigger to automatically update updated_at on row updates
CREATE OR REPLACE FUNCTION update_transcript_segments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER transcript_segments_updated_at_trigger 
    BEFORE UPDATE ON transcript_segments 
    FOR EACH ROW 
    EXECUTE FUNCTION update_transcript_segments_updated_at();

-- Set existing rows to have created_at as updated_at initially
UPDATE transcript_segments SET updated_at = created_at WHERE updated_at IS NULL;

-- Make updated_at NOT NULL after setting default values
ALTER TABLE transcript_segments ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN transcript_segments.updated_at IS 'Timestamp when the transcript segment was last updated'; 