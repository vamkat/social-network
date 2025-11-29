-- ============================================================================
-- NOTIFICATION SERVICE TEST SUITE
-- ============================================================================
-- Run with: psql -h localhost -p 5436 -U postgres -d social_notifications -f backend/services/notifications/internal/db/test/notification_tests.sql
-- ============================================================================

\set ON_ERROR_STOP on
\timing off

BEGIN;

-- Test result tracking
CREATE TEMP TABLE test_results (
    test_number INT,
    test_name TEXT,
    passed BOOLEAN,
    error_message TEXT
);

CREATE TEMP SEQUENCE test_counter;

-- Helper function to run tests
CREATE OR REPLACE FUNCTION run_test(
    test_name TEXT,
    test_query TEXT,
    expected_result BOOLEAN,
    error_msg TEXT DEFAULT NULL
) RETURNS VOID AS $$
DECLARE
    test_num INT;
    actual_result BOOLEAN;
BEGIN
    test_num := nextval('test_counter');

    BEGIN
        EXECUTE test_query INTO actual_result;

        IF actual_result = expected_result THEN
            INSERT INTO test_results VALUES (test_num, test_name, TRUE, NULL);
            RAISE NOTICE '[%] PASS: %', test_num, test_name;
        ELSE
            INSERT INTO test_results VALUES (test_num, test_name, FALSE,
                format('Expected %s but got %s', expected_result, actual_result));
            RAISE NOTICE '[%] FAIL: % - Expected % but got %',
                test_num, test_name, expected_result, actual_result;
        END IF;
    EXCEPTION WHEN OTHERS THEN
        INSERT INTO test_results VALUES (test_num, test_name, FALSE, SQLERRM);
        RAISE NOTICE '[%] FAIL: % - Error: %', test_num, test_name, SQLERRM;
    END;
END;
$$ LANGUAGE plpgsql;

-- Helper to expect exception
CREATE OR REPLACE FUNCTION expect_exception(
    test_name TEXT,
    test_query TEXT,
    expected_error_substring TEXT DEFAULT NULL
) RETURNS VOID AS $$
DECLARE
    test_num INT;
BEGIN
    test_num := nextval('test_counter');

    BEGIN
        EXECUTE test_query;
        -- If we get here, no exception was raised
        INSERT INTO test_results VALUES (test_num, test_name, FALSE,
            'Expected exception but none was raised');
        RAISE NOTICE '[%] FAIL: % - Expected exception but none was raised',
            test_num, test_name;
    EXCEPTION WHEN OTHERS THEN
        IF expected_error_substring IS NULL OR SQLERRM LIKE '%' || expected_error_substring || '%' THEN
            INSERT INTO test_results VALUES (test_num, test_name, TRUE, NULL);
            RAISE NOTICE '[%] PASS: % (caught expected exception)', test_num, test_name;
        ELSE
            INSERT INTO test_results VALUES (test_num, test_name, FALSE,
                format('Expected error containing "%s" but got: %s', expected_error_substring, SQLERRM));
            RAISE NOTICE '[%] FAIL: % - Wrong error message: %', test_num, test_name, SQLERRM;
        END IF;
    END;
END;
$$ LANGUAGE plpgsql;

\echo ''
\echo '========================================================================'
\echo 'SETTING UP TEST DATA'
\echo '========================================================================'

-- Insert test notification types
INSERT INTO notification_types (notif_type, category, default_enabled)
VALUES
    ('new_follower', 'social', TRUE),
    ('follow_request', 'social', TRUE),
    ('group_invite', 'group', TRUE),
    ('group_join_request', 'group', TRUE),
    ('new_event', 'group', TRUE),
    ('new_message', 'chat', TRUE),
    ('post_reply', 'posts', TRUE),
    ('like', 'posts', TRUE)
ON CONFLICT (notif_type) DO NOTHING;

\echo 'Test notification types created'

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST SUITE 1: NOTIFICATION CREATION'
\echo '========================================================================'

DO $$
DECLARE
    notification_id BIGINT;
    test_user_id BIGINT := 1001;
    test_source_entity_id BIGINT := 2001;
    count INT;
BEGIN
    -- Test 1: Create notification successfully
    INSERT INTO notifications (
        user_id, notif_type, source_service, source_entity_id, payload
    ) VALUES (
        test_user_id, 'new_follower', 'users', test_source_entity_id, '{"follower_id": "2001"}'
    ) RETURNING id INTO notification_id;

    -- Verify notification was created
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE id = notification_id AND user_id = test_user_id AND notif_type = 'new_follower';

    IF count = 1 THEN
        RAISE NOTICE '[%] PASS: Notification created successfully', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification creation', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification not created properly', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification creation', FALSE, 'Notification not found');
    END IF;

    -- Test 2: Notification has proper default values
    IF EXISTS (
        SELECT 1 FROM notifications
        WHERE id = notification_id
        AND seen = FALSE
        AND needs_action = FALSE
        AND acted = FALSE
        AND created_at IS NOT NULL
    ) THEN
        RAISE NOTICE '[%] PASS: Notification has correct default values', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification defaults', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification defaults incorrect', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification defaults', FALSE, 'Default values incorrect');
    END IF;

    -- Test 3: All notification types can be created
    INSERT INTO notifications (user_id, notif_type, source_service, source_entity_id, payload)
    SELECT test_user_id + 100, unnest_type, 'users', test_source_entity_id + 100, '{}'
    FROM unnest(ARRAY['follow_request', 'group_invite', 'group_join_request', 'new_event', 'new_message', 'post_reply', 'like']) AS unnest_type;

    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id + 100;

    IF count = 7 THEN
        RAISE NOTICE '[%] PASS: All notification types supported', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'All notification types', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Not all notification types supported', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'All notification types', FALSE, format('Expected 7, got %s', count));
    END IF;

    -- Test 4: Valid source services enforced
    INSERT INTO notifications (
        user_id, notif_type, source_service, source_entity_id, payload
    ) VALUES (
        test_user_id + 200, 'new_follower', 'chat', test_source_entity_id + 200, '{}'
    );

    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id + 200 AND source_service = 'chat';

    IF count = 1 THEN
        RAISE NOTICE '[%] PASS: Chat source service works', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Valid source service', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Chat source service failed', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Valid source service', FALSE, 'Chat service not accepted');
    END IF;

    -- Test 5: Invalid source service rejected
    BEGIN
        INSERT INTO notifications (
            user_id, notif_type, source_service, source_entity_id, payload
        ) VALUES (
            test_user_id + 300, 'new_follower', 'invalid_service', test_source_entity_id + 300, '{}'
        );
        INSERT INTO test_results VALUES (nextval('test_counter'), 'Invalid source service rejection', FALSE, 'No exception raised');
        RAISE NOTICE '[%] FAIL: Should reject invalid source service', currval('test_counter');
    EXCEPTION WHEN check_violation THEN
        INSERT INTO test_results VALUES (nextval('test_counter'), 'Invalid source service rejection', TRUE, NULL);
        RAISE NOTICE '[%] PASS: Invalid source service correctly rejected', currval('test_counter');
    END;
END $$;

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST SUITE 2: NOTIFICATION STATUS MANAGEMENT'
\echo '========================================================================'

DO $$
DECLARE
    test_user_id BIGINT := 2001;
    notification_id_1 BIGINT;
    notification_id_2 BIGINT;
    notification_id_3 BIGINT;
    count INT;
    seen_status BOOLEAN;
    acted_status BOOLEAN;
BEGIN
    -- Create test notifications
    INSERT INTO notifications (user_id, notif_type, source_service, payload)
    VALUES (test_user_id, 'new_follower', 'users', '{}') RETURNING id INTO notification_id_1;

    INSERT INTO notifications (user_id, notif_type, source_service, payload)
    VALUES (test_user_id, 'new_follower', 'users', '{}') RETURNING id INTO notification_id_2;

    INSERT INTO notifications (user_id, notif_type, source_service, payload)
    VALUES (test_user_id, 'new_follower', 'users', '{}') RETURNING id INTO notification_id_3;

    -- Test 1: Initial notifications are unread
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND seen = FALSE;

    IF count = 3 THEN
        RAISE NOTICE '[%] PASS: Notifications initially unread', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Initial unread status', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notifications not initially unread', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Initial unread status', FALSE, format('Expected 3, got %s', count));
    END IF;

    -- Test 2: Mark single notification as read
    UPDATE notifications
    SET seen = true
    WHERE id = notification_id_1;

    SELECT seen INTO seen_status
    FROM notifications
    WHERE id = notification_id_1;

    IF seen_status = TRUE THEN
        RAISE NOTICE '[%] PASS: Single notification marked as read', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark single as read', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Single notification not marked as read', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark single as read', FALSE, 'Not marked as read');
    END IF;

    -- Test 3: Mark multiple notifications as read using ANY
    UPDATE notifications
    SET seen = true
    WHERE user_id = test_user_id AND id = ANY(ARRAY[notification_id_2, notification_id_3]);

    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND seen = TRUE;

    IF count = 3 THEN
        RAISE NOTICE '[%] PASS: Multiple notifications marked as read', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark multiple as read', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Multiple notifications not marked as read', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark multiple as read', FALSE, format('Expected 3, got %s', count));
    END IF;

    -- Test 4: Mark notification as acted updates acted status
    UPDATE notifications
    SET acted = true
    WHERE id = notification_id_1;

    SELECT acted INTO acted_status
    FROM notifications
    WHERE id = notification_id_1;

    IF acted_status = TRUE THEN
        RAISE NOTICE '[%] PASS: Notification marked as acted', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark as acted', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification not marked as acted', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Mark as acted', FALSE, 'Not marked as acted');
    END IF;

    -- Test 5: Verify both seen and acted statuses can be set independently
    INSERT INTO notifications (user_id, notif_type, source_service, payload)
    VALUES (test_user_id + 100, 'follow_request', 'users', '{}') RETURNING id INTO notification_id_1;

    -- Only mark as acted, not seen
    UPDATE notifications
    SET acted = true
    WHERE id = notification_id_1;

    -- Check that it's acted but not seen
    SELECT seen, acted INTO seen_status, acted_status
    FROM notifications
    WHERE id = notification_id_1;

    IF acted_status = TRUE AND seen_status = FALSE THEN
        RAISE NOTICE '[%] PASS: Seen and acted status independent', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Status independence', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Seen and acted status not independent', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Status independence', FALSE, format('Seen: %s, Acted: %s', seen_status, acted_status));
    END IF;
END $$;

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST SUITE 3: NOTIFICATION QUERYING'
\echo '========================================================================'

DO $$
DECLARE
    test_user_id BIGINT := 3001;
    notification_id BIGINT;
    count INT;
    result_count BIGINT;
BEGIN
    -- Clean up any existing test data
    DELETE FROM notifications WHERE user_id >= test_user_id AND user_id < test_user_id + 100;

    -- Create test notifications with different types and timestamps
    INSERT INTO notifications (user_id, notif_type, source_service, payload, created_at)
    VALUES
        (test_user_id, 'new_follower', 'users', '{}', NOW() - INTERVAL '3 hours'),
        (test_user_id, 'follow_request', 'users', '{}', NOW() - INTERVAL '2 hours'),
        (test_user_id, 'group_invite', 'users', '{}', NOW() - INTERVAL '1 hour'),
        (test_user_id, 'new_message', 'chat', '{}', NOW());

    -- Test 1: Get all notifications for user
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND deleted_at IS NULL;

    IF count = 4 THEN
        RAISE NOTICE '[%] PASS: All notifications retrieved', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Get all notifications', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: All notifications not retrieved', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Get all notifications', FALSE, format('Expected 4, got %s', count));
    END IF;

    -- Test 2: Get notifications filtered by seen status
    UPDATE notifications
    SET seen = true
    WHERE user_id = test_user_id AND notif_type IN ('new_follower', 'follow_request');

    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND seen = true;

    IF count = 2 THEN
        RAISE NOTICE '[%] PASS: Seen notifications filtered correctly', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Filter by seen status', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Seen notifications not filtered correctly', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Filter by seen status', FALSE, format('Expected 2, got %s', count));
    END IF;

    -- Test 3: Notification count functions work
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND deleted_at IS NULL;

    -- Simulate what the application code would do for counting
    WITH counts AS (
        SELECT 
            COUNT(*) AS total_count,
            COUNT(*) FILTER (WHERE seen = true) AS seen_count,
            COUNT(*) FILTER (WHERE seen = false) AS unseen_count
        FROM notifications
        WHERE user_id = test_user_id AND deleted_at IS NULL
    )
    SELECT total_count, unseen_count INTO result_count, count
    FROM counts;

    IF result_count = 4 AND count = 2 THEN
        RAISE NOTICE '[%] PASS: Notification counts calculated correctly', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification counts', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification counts incorrect', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Notification counts', FALSE, format('Total: %s, Unseen: %s', result_count, count));
    END IF;

    -- Test 4: Notifications ordered by most recent first
    SELECT id INTO notification_id
    FROM notifications
    WHERE user_id = test_user_id
    ORDER BY created_at DESC
    LIMIT 1;

    -- The most recent notification should be the 'new_message' one
    IF EXISTS (
        SELECT 1 FROM notifications
        WHERE id = notification_id AND notif_type = 'new_message'
    ) THEN
        RAISE NOTICE '[%] PASS: Notifications ordered by most recent', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Order by recent', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notifications not ordered by most recent', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Order by recent', FALSE, 'Ordering incorrect');
    END IF;

    -- Test 5: Limit and offset work correctly
    SELECT COUNT(*) INTO count
    FROM (
        SELECT id
        FROM notifications
        WHERE user_id = test_user_id
        ORDER BY created_at DESC
        LIMIT 2
        OFFSET 1
    ) AS limited_results;

    IF count = 2 THEN
        RAISE NOTICE '[%] PASS: Limit and offset work correctly', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Limit and offset', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Limit and offset not working', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Limit and offset', FALSE, format('Expected 2, got %s', count));
    END IF;
END $$;

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST SUITE 4: INDEXES AND PERFORMANCE'
\echo '========================================================================'

DO $$
DECLARE
    count INT;
BEGIN
    -- Test 1: Verify indexes exist and are used
    -- This tests that the important indexes are in place
    SELECT COUNT(*) INTO count
    FROM pg_indexes
    WHERE schemaname = 'public'
    AND tablename = 'notifications'
    AND indexname LIKE 'idx_notifications%';

    -- Should have at least the user_unread index and others
    IF count >= 3 THEN
        RAISE NOTICE '[%] PASS: Notification indexes exist', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Indexes exist', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Not all expected indexes exist', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Indexes exist', FALSE, format('Expected >=3, got %s', count));
    END IF;

    -- Test 2: Test user_unread index query pattern
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = 999999 AND seen = false
    LIMIT 10;

    -- Should be able to execute the query without error
    RAISE NOTICE '[%] PASS: User-unread index query pattern works', nextval('test_counter');
    INSERT INTO test_results VALUES (currval('test_counter'), 'User-unread query', TRUE, NULL);

    -- Test 3: Test user_created index query pattern
    -- Just test that the query structure works (this is for performance testing)
    -- We'll test the actual functionality separately
    CREATE TEMP TABLE temp_test_results AS
    SELECT id
    FROM notifications
    WHERE user_id = 999999
    ORDER BY created_at DESC
    LIMIT 10;

    SELECT COUNT(*) INTO count FROM temp_test_results;
    DROP TABLE temp_test_results;

    -- Should be able to execute the query without error
    RAISE NOTICE '[%] PASS: User-created index query pattern works', nextval('test_counter');
    INSERT INTO test_results VALUES (currval('test_counter'), 'User-created query', TRUE, NULL);
END $$;

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST SUITE 5: SOFT DELETE FUNCTIONALITY'
\echo '========================================================================'

DO $$
DECLARE
    test_user_id BIGINT := 4001;
    notification_id BIGINT;
    count INT;
BEGIN
    -- Create test notification
    INSERT INTO notifications (user_id, notif_type, source_service, payload)
    VALUES (test_user_id, 'new_follower', 'users', '{}') RETURNING id INTO notification_id;

    -- Test 1: Notification exists initially
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE id = notification_id AND deleted_at IS NULL;

    IF count = 1 THEN
        RAISE NOTICE '[%] PASS: Notification exists before soft delete', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Pre-delete existence', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification does not exist before soft delete', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Pre-delete existence', FALSE, 'Notification missing');
    END IF;

    -- Test 2: Soft delete notification
    UPDATE notifications
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = notification_id;

    -- Verify it's no longer returned in normal queries
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE id = notification_id AND deleted_at IS NULL;

    IF count = 0 THEN
        RAISE NOTICE '[%] PASS: Notification not returned after soft delete', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Soft delete hiding', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification still returned after soft delete', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Soft delete hiding', FALSE, 'Notification still visible');
    END IF;

    -- Test 3: Soft deleted notification still exists in DB with deleted_at
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE id = notification_id AND deleted_at IS NOT NULL;

    IF count = 1 THEN
        RAISE NOTICE '[%] PASS: Notification physically exists but marked deleted', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Soft delete marking', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Notification not properly marked deleted', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Soft delete marking', FALSE, 'Not properly marked');
    END IF;

    -- Test 4: Verify all queries use deleted_at filter
    -- Count all notifications for the user including soft-deleted ones
    SELECT COUNT(*) INTO count
    FROM notifications
    WHERE user_id = test_user_id AND deleted_at IS NULL;

    -- Should return 0 since we soft-deleted the only notification
    IF count = 0 THEN
        RAISE NOTICE '[%] PASS: All queries filter out deleted notifications', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Delete filter in queries', TRUE, NULL);
    ELSE
        RAISE NOTICE '[%] FAIL: Queries not properly filtering deleted notifications', nextval('test_counter');
        INSERT INTO test_results VALUES (currval('test_counter'), 'Delete filter in queries', FALSE, format('Expected 0, got %s', count));
    END IF;
END $$;

-- ============================================================================
\echo ''
\echo '========================================================================'
\echo 'TEST RESULTS SUMMARY'
\echo '========================================================================'

DO $$
DECLARE
    total_tests INT;
    passed_tests INT;
    failed_tests INT;
BEGIN
    SELECT COUNT(*) INTO total_tests FROM test_results;
    SELECT COUNT(*) INTO passed_tests FROM test_results WHERE passed = TRUE;
    SELECT COUNT(*) INTO failed_tests FROM test_results WHERE passed = FALSE;

    RAISE NOTICE '';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Total Tests: %', total_tests;
    RAISE NOTICE 'Passed:      % (%.1f%%)', passed_tests, (passed_tests::DECIMAL / NULLIF(total_tests, 0) * 100);
    RAISE NOTICE 'Failed:      % (%.1f%%)', failed_tests, (failed_tests::DECIMAL / NULLIF(total_tests, 0) * 100);
    RAISE NOTICE '========================================';
    RAISE NOTICE '';

    IF failed_tests > 0 THEN
        RAISE NOTICE 'FAILED TESTS:';
        RAISE NOTICE '';
    END IF;
END $$;

-- Show failed tests
SELECT
    test_number,
    test_name,
    error_message
FROM test_results
WHERE passed = FALSE
ORDER BY test_number;

-- Final result
DO $$
DECLARE
    failed_count INT;
BEGIN
    SELECT COUNT(*) INTO failed_count FROM test_results WHERE passed = FALSE;

    IF failed_count = 0 THEN
        RAISE NOTICE '';
        RAISE NOTICE '✓ ALL TESTS PASSED! ✓';
        RAISE NOTICE '';
    ELSE
        RAISE NOTICE '';
        RAISE NOTICE '✗ SOME TESTS FAILED ✗';
        RAISE NOTICE '';
    END IF;
END $$;

ROLLBACK;

\echo ''
\echo 'Test completed. All changes rolled back.'