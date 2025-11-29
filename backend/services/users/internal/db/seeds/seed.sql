-------------------------------------------------------
-- Development Seed Data
-- Safe for non-production environments only
-- Idempotent: safe to run multiple times
-- Uses ON CONFLICT DO NOTHING to skip existing records
-------------------------------------------------------

-------------------------------------------------------
-- Users
-------------------------------------------------------
INSERT INTO users (id, username, first_name, last_name, date_of_birth, avatar, about_me, profile_public, current_status)
OVERRIDING SYSTEM VALUE
VALUES
(1, 'alice', 'Alice', 'Wonder', '1990-01-01', 'https://example.com/avatars/alice.jpg', 'Love nature and outdoor activities', TRUE, 'active'),
(2, 'bob', 'Bob', 'Builder', '1992-02-02', 'https://example.com/avatars/bob.jpg', 'Professional builder and contractor', FALSE, 'active'),
(3, 'charlie', 'Charlie', 'Day', '1991-03-03', 'https://example.com/avatars/charlie.jpg', 'Gamer and tech enthusiast', TRUE, 'active'),
(4, 'diana', 'Diana', 'Prince', '1988-04-04', 'https://example.com/avatars/diana.jpg', 'Entrepreneur and business owner', TRUE, 'active'),
(5, 'eve', 'Eve', 'Hacker', '1995-05-05', 'https://example.com/avatars/eve.jpg', 'Security researcher', FALSE, 'active'),
(6, 'frank', 'Frank', 'Ocean', '1994-06-06', 'https://example.com/avatars/frank.jpg', 'Music producer and artist', TRUE, 'active'),
(7, 'grace', 'Grace', 'Hopper', '1985-07-07', 'https://example.com/avatars/grace.jpg', 'Software engineer and mentor', TRUE, 'active'),
(8, 'henry', 'Henry', 'Ford', '1986-08-08', 'https://example.com/avatars/henry.jpg', 'Automotive engineer', FALSE, 'active'),
(9, 'ivy', 'Ivy', 'Green', '1993-09-09', 'https://example.com/avatars/ivy.jpg', 'Environmental activist', TRUE, 'active'),
(10, 'jack', 'Jack', 'Black', '1990-10-10', 'https://example.com/avatars/jack.jpg', 'Actor and musician', TRUE, 'active')
ON CONFLICT (id) DO NOTHING;

-------------------------------------------------------
-- Auth Records (fake passwords)
-------------------------------------------------------
INSERT INTO auth_user (user_id, email, password_hash)
SELECT id, LOWER(username) || '@example.com', 'hash'
FROM users
WHERE id IN (1,2,3,4,5,6,7,8,9,10)
ON CONFLICT (user_id) DO NOTHING;

-------------------------------------------------------
-- Follow relationships (direct)
-------------------------------------------------------
INSERT INTO follows (follower_id, following_id)
VALUES
(1, 3), -- Alice → Charlie
(1, 4), -- Alice → Diana
(3, 1), -- Charlie → Alice (mutual)
(4, 3), -- Diana → Charlie
(6, 1), -- Frank → Alice
(9, 1), -- Ivy → Alice
(10, 3) -- Jack → Charlie
ON CONFLICT (follower_id, following_id) DO NOTHING;

-------------------------------------------------------
-- Follow requests (for private profiles)
-------------------------------------------------------
INSERT INTO follow_requests (requester_id, target_id, status)
VALUES
(1, 2, 'pending'), -- Alice → Bob (private)
(3, 2, 'accepted'), -- Charlie → Bob
(4, 5, 'pending'), -- Diana → Eve (private)
(7, 8, 'rejected') -- Grace → Henry (private)
ON CONFLICT (requester_id, target_id) DO NOTHING;

-------------------------------------------------------
-- Groups
-------------------------------------------------------
INSERT INTO groups (id, group_owner, group_title, group_description, group_image)
OVERRIDING SYSTEM VALUE
VALUES
(1, 1, 'Nature Lovers', 'A group for nature enthusiasts', 'https://example.com/groups/nature.jpg'),
(2, 3, 'Gamers Unite', 'All about gaming', 'https://example.com/groups/gaming.jpg'),
(3, 6, 'Music Fans', 'People who love music', 'https://example.com/groups/music.jpg')
ON CONFLICT (id) DO NOTHING;

-------------------------------------------------------
-- Group Members
-------------------------------------------------------
-- Group 1 (Nature Lovers)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(1, 1, 'owner'),
(1, 3, 'member'),
(1, 4, 'member')
ON CONFLICT (group_id, user_id) DO NOTHING;

-- Group 2 (Gamers Unite)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(2, 3, 'owner'),
(2, 6, 'member'),
(2, 7, 'member')
ON CONFLICT (group_id, user_id) DO NOTHING;

-- Group 3 (Music Fans)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(3, 6, 'owner'),
(3, 1, 'member'),
(3, 9, 'member')
ON CONFLICT (group_id, user_id) DO NOTHING;

-------------------------------------------------------
-- Group Join Requests
-------------------------------------------------------
INSERT INTO group_join_requests (group_id, user_id, status)
VALUES
(1, 5, 'pending'),
(2, 1, 'rejected'),
(3, 10, 'accepted')
ON CONFLICT (group_id, user_id) DO NOTHING;

-------------------------------------------------------
-- Group Invites
-------------------------------------------------------
INSERT INTO group_invites (group_id, sender_id, receiver_id, status)
VALUES
(1, 1, 2, 'pending'),
(2, 3, 5, 'declined'),
(3, 6, 8, 'accepted')
ON CONFLICT (group_id, receiver_id) DO NOTHING;
