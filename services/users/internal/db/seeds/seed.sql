-------------------------------------------------------
-- Development Seed Data
-- Safe for non-production environments only
-------------------------------------------------------

-- Clear tables (order matters due to FKs)
TRUNCATE group_invites, group_join_requests, group_members, groups,
         follow_requests, follows, auth_user, users
RESTART IDENTITY CASCADE;

-------------------------------------------------------
-- Users
-------------------------------------------------------
INSERT INTO users (username, first_name, last_name, date_of_birth, profile_public)
VALUES
('alice', 'Alice', 'Wonder', '1990-01-01', TRUE),
('bob', 'Bob', 'Builder', '1992-02-02', FALSE),
('charlie', 'Charlie', 'Day', '1991-03-03', TRUE),
('diana', 'Diana', 'Prince', '1988-04-04', TRUE),
('eve', 'Eve', 'Hacker', '1995-05-05', FALSE),
('frank', 'Frank', 'Ocean', '1994-06-06', TRUE),
('grace', 'Grace', 'Hopper', '1985-07-07', TRUE),
('henry', 'Henry', 'Ford', '1986-08-08', FALSE),
('ivy', 'Ivy', 'Green', '1993-09-09', TRUE),
('jack', 'Jack', 'Black', '1990-10-10', TRUE);

-------------------------------------------------------
-- Auth Records (fake passwords)
-------------------------------------------------------
INSERT INTO auth_user (user_id, email, password_hash, salt)
SELECT id, LOWER(username) || '@example.com', 'hash', 'salt'
FROM users;

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
(10, 3); -- Jack → Charlie

-------------------------------------------------------
-- Follow requests (for private profiles)
-------------------------------------------------------
INSERT INTO follow_requests (requester_id, target_id, status)
VALUES
(1, 2, 'pending'), -- Alice → Bob (private)
(3, 2, 'accepted'), -- Charlie → Bob
(4, 5, 'pending'), -- Diana → Eve (private)
(7, 8, 'rejected'); -- Grace → Henry (private)

-------------------------------------------------------
-- Groups
-------------------------------------------------------
INSERT INTO groups (group_owner, group_title, group_description)
VALUES
(1, 'Nature Lovers', 'A group for nature enthusiasts'),
(3, 'Gamers Unite', 'All about gaming'),
(6, 'Music Fans', 'People who love music');

-------------------------------------------------------
-- Group Members
-------------------------------------------------------
-- Group 1 (Nature Lovers)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(1, 1, 'owner'),
(1, 3, 'member'),
(1, 4, 'member');

-- Group 2 (Gamers Unite)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(2, 3, 'owner'),
(2, 6, 'member'),
(2, 7, 'member');

-- Group 3 (Music Fans)
INSERT INTO group_members (group_id, user_id, role)
VALUES
(3, 6, 'owner'),
(3, 1, 'member'),
(3, 9, 'member');

-------------------------------------------------------
-- Group Join Requests
-------------------------------------------------------
INSERT INTO group_join_requests (group_id, user_id, status)
VALUES
(1, 5, 'pending'),
(2, 1, 'rejected'),
(3, 10, 'accepted');

-------------------------------------------------------
-- Group Invites
-------------------------------------------------------
INSERT INTO group_invites (group_id, sender_id, receiver_id, status)
VALUES
(1, 1, 2, 'pending'),
(2, 3, 5, 'declined'),
(3, 6, 8, 'accepted');
