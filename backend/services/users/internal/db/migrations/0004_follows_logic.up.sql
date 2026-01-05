
-----------------------------------------
-- Function to follow a user regardless of privacy setting
-----------------------------------------

CREATE OR REPLACE FUNCTION follow_user(p_follower BIGINT, p_target BIGINT)
RETURNS TEXT AS $$
DECLARE
    v_is_public BOOLEAN;
    v_already_following BOOLEAN;
    v_existing_request_status follow_request_status;
BEGIN
    -- Validate: cannot follow yourself
    IF p_follower = p_target THEN
        RAISE EXCEPTION 'Cannot follow yourself'
        USING ERRCODE = '22023';
    END IF;

    -- Check if target user exists and get privacy setting
    SELECT profile_public INTO v_is_public 
    FROM users 
    WHERE id = p_target 
      AND deleted_at IS NULL
      AND current_status = 'active';

    IF v_is_public IS NULL THEN
        RAISE EXCEPTION 'User not found, deleted, or not active'
        USING ERRCODE = 'P0002';
    END IF;

    -- Check if already following
    SELECT EXISTS(
        SELECT 1 FROM follows 
        WHERE follower_id = p_follower AND following_id = p_target
    ) INTO v_already_following;

    IF v_already_following THEN
        RETURN 'already_following';
    END IF;

    -- Public profile: add follow directly
    IF v_is_public THEN
        INSERT INTO follows (follower_id, following_id)
        VALUES (p_follower, p_target)
        ON CONFLICT DO NOTHING;
        
        -- Clean up any existing request
        DELETE FROM follow_requests
        WHERE requester_id = p_follower AND target_id = p_target;

        RETURN 'followed';
    
    -- Private profile: create or update follow request
    ELSE
        -- Check for existing request
        SELECT status INTO v_existing_request_status
        FROM follow_requests
        WHERE requester_id = p_follower AND target_id = p_target;

        IF v_existing_request_status = 'pending' THEN
            RETURN 'request_already_pending';
        ELSIF v_existing_request_status IN ('rejected', 'accepted') THEN
            -- Reset to pending if previously rejected/accepted
            UPDATE follow_requests
            SET status = 'pending',
                updated_at = CURRENT_TIMESTAMP
            WHERE requester_id = p_follower AND target_id = p_target;
            RETURN 'request_resent';
        ELSE
            -- Create new request
            INSERT INTO follow_requests (requester_id, target_id, status)
            VALUES (p_follower, p_target, 'pending')
            ON CONFLICT (requester_id, target_id) DO UPDATE
            SET status = 'pending',
                updated_at = CURRENT_TIMESTAMP;
            RETURN 'requested';
        END IF;
    END IF;
END;
$$ LANGUAGE plpgsql;




-----------------------------------------
-- Trigger to add follower when follow request is accepted
-----------------------------------------
   CREATE OR REPLACE FUNCTION add_follower_on_accept()
   RETURNS TRIGGER AS $$
   BEGIN
       INSERT INTO follows (follower_id, following_id, created_at)
       VALUES (NEW.requester_id, NEW.target_id, CURRENT_TIMESTAMP)
       ON CONFLICT DO NOTHING;
       
       RETURN NEW;
   END;
   $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_add_follower_on_accept
AFTER UPDATE ON follow_requests
FOR EACH ROW
WHEN (NEW.status = 'accepted' AND OLD.status IS DISTINCT FROM 'accepted')
EXECUTE FUNCTION add_follower_on_accept();

-----------------------------------------
-- Trigger to accept pending follow requests when a profile changes to public
-----------------------------------------
CREATE OR REPLACE FUNCTION accept_pending_requests_on_public()
RETURNS TRIGGER AS $$
BEGIN
    -- Only act if profile switches from private to public
    IF OLD.profile_public = FALSE AND NEW.profile_public = TRUE THEN
        
        -- Insert all pending requests directly as follows
        INSERT INTO follows (follower_id, following_id, created_at)
        SELECT requester_id, target_id, CURRENT_TIMESTAMP
        FROM follow_requests
        WHERE target_id = NEW.id
          AND status = 'pending'
          AND deleted_at IS NULL
        ON CONFLICT (follower_id, following_id) DO NOTHING;

        -- Mark all pending requests as accepted
        UPDATE follow_requests
        SET status = 'accepted', 
            updated_at = CURRENT_TIMESTAMP
        WHERE target_id = NEW.id
          AND status = 'pending'
          AND deleted_at IS NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_accept_pending_requests_on_public
AFTER UPDATE ON users
FOR EACH ROW
WHEN (OLD.profile_public = FALSE AND NEW.profile_public = TRUE)
EXECUTE FUNCTION accept_pending_requests_on_public();
