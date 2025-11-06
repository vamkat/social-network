-- Create follow_requests table
CREATE TABLE IF NOT EXISTS follow_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_follow UNIQUE (follower_id, followed_id),
    CONSTRAINT no_self_follow CHECK (follower_id <> followed_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_follow_requests_follower_id ON follow_requests (follower_id);
CREATE INDEX IF NOT EXISTS idx_follow_requests_followed_id ON follow_requests (followed_id);
CREATE INDEX IF NOT EXISTS idx_follow_requests_status ON follow_requests (status);

-- Trigger to auto-update updated_at column
CREATE OR REPLACE FUNCTION update_follow_request_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_follow_request_timestamp
BEFORE UPDATE ON follow_requests
FOR EACH ROW
EXECUTE FUNCTION update_follow_request_timestamp();
