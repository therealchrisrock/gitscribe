-- Remove updated_at column from transcript_segments table
-- Down migration: 000005_add_updated_at_to_transcript_segments

-- Drop the trigger first
DROP TRIGGER IF EXISTS transcript_segments_updated_at_trigger ON transcript_segments;

-- Drop the function
DROP FUNCTION IF EXISTS update_transcript_segments_updated_at();

-- Remove the updated_at column
ALTER TABLE transcript_segments DROP COLUMN IF EXISTS updated_at; 