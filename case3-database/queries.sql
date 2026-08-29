-- ============================================================================
-- CASE 3: COMPLEX SQL QUERIES & OPTIMIZATIONS
-- ============================================================================

USE `social_media_db`;

-- ----------------------------------------------------------------------------
-- QUERY 1: CHRONOLOGICAL NEWSFEED GENERATION (WITH CURSOR PAGINATION)
-- Scenario: Generate the home feed for Alice (user_id = 1), showing posts from
-- people she follows (Bob & Charlie), ordered by newest first.
--
-- Optimization:
-- - Uses Keyset / Cursor Pagination (WHERE p.created_at < :last_seen_timestamp)
--   instead of OFFSET to maintain O(1) performance as feed depth increases.
-- - Leverages composite index `idx_posts_user_created` (user_id, created_at DESC).
-- ----------------------------------------------------------------------------
SELECT 
    p.id AS post_id,
    p.user_id AS author_id,
    u.username AS author_username,
    u.avatar_url AS author_avatar,
    p.content,
    p.like_count,
    p.comment_count,
    p.created_at AS published_at,
    -- Aggregate media attachments into JSON array
    COALESCE(
        JSON_ARRAYAGG(
            JSON_OBJECT(
                'id', pm.id,
                'media_url', pm.media_url,
                'media_type', pm.media_type,
                'order', pm.display_order
            )
        ), 
        JSON_ARRAY()
    ) AS media_attachments,
    -- Check if Alice (user_id = 1) has liked this post
    EXISTS (
        SELECT 1 FROM `reactions` r 
        WHERE r.target_type = 'POST' 
          AND r.target_id = p.id 
          AND r.user_id = 1
    ) AS is_liked_by_viewer
FROM `posts` p
INNER JOIN `follows` f ON f.following_id = p.user_id AND f.follower_id = 1
INNER JOIN `users` u ON u.id = p.user_id
LEFT JOIN `post_media` pm ON pm.post_id = p.id
-- Cursor filter for pagination (e.g. posts created before '2026-08-05 12:00:00')
WHERE p.created_at <= '2026-08-05 12:00:00'
GROUP BY p.id, u.id
ORDER BY p.created_at DESC
LIMIT 10;

-- ----------------------------------------------------------------------------
-- QUERY 2: UNREAD DIRECT MESSAGES COUNT PER CONVERSATION
-- Scenario: Calculate unread message badge count for Alice (user_id = 1) across
-- all her active conversations.
--
-- Optimization:
-- - Uses composite index `idx_messages_conv_created` (conversation_id, created_at).
-- ----------------------------------------------------------------------------
SELECT 
    c.id AS conversation_id,
    c.is_group,
    c.title,
    other_user.username AS chat_partner,
    cm.last_read_at,
    COUNT(m.id) AS unread_message_count,
    MAX(m.created_at) AS last_message_time
FROM `conversations` c
INNER JOIN `conversation_members` cm 
    ON cm.conversation_id = c.id AND cm.user_id = 1
LEFT JOIN `conversation_members` other_cm 
    ON other_cm.conversation_id = c.id AND other_cm.user_id != 1 AND c.is_group = FALSE
LEFT JOIN `users` other_user 
    ON other_user.id = other_cm.user_id
LEFT JOIN `messages` m 
    ON m.conversation_id = c.id 
   AND (cm.last_read_at IS NULL OR m.created_at > cm.last_read_at)
   AND m.sender_id != 1
GROUP BY c.id, cm.user_id, other_user.id
ORDER BY last_message_time DESC;

-- ----------------------------------------------------------------------------
-- QUERY 3: MUTUAL FOLLOWERS (FRIEND GRAPH)
-- Scenario: Find mutual followers between Alice (user_id = 1) and Bob (user_id = 2)
-- ----------------------------------------------------------------------------
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar_url
FROM `users` u
WHERE u.id IN (
    -- Users that Alice follows
    SELECT following_id FROM `follows` WHERE follower_id = 1
)
AND u.id IN (
    -- Users that Bob follows
    SELECT following_id FROM `follows` WHERE follower_id = 2
);

-- ----------------------------------------------------------------------------
-- QUERY 4: USER PROFILE ENGAGEMENT STATS
-- Scenario: Aggregate profile statistics for Bob (user_id = 2):
-- Total posts, Followers, Following count, and Total Likes Received across all posts.
-- ----------------------------------------------------------------------------
SELECT 
    u.id AS user_id,
    u.username,
    u.full_name,
    (SELECT COUNT(*) FROM `posts` WHERE user_id = u.id) AS total_posts,
    (SELECT COUNT(*) FROM `follows` WHERE following_id = u.id) AS followers_count,
    (SELECT COUNT(*) FROM `follows` WHERE follower_id = u.id) AS following_count,
    COALESCE((SELECT SUM(like_count) FROM `posts` WHERE user_id = u.id), 0) AS total_likes_received
FROM `users` u
WHERE u.id = 2;

-- ----------------------------------------------------------------------------
-- QUERY 5: TRENDING POSTS DISCOVERY (ENGAGEMENT SCORING)
-- Scenario: Discover top trending posts in the last 7 days ranked by engagement.
-- Formula: (likes * 2 + comments * 3) / (age_in_hours + 2)^1.5 (Hacker News decay algorithm)
-- ----------------------------------------------------------------------------
SELECT 
    p.id,
    u.username,
    p.content,
    p.like_count,
    p.comment_count,
    p.created_at,
    (
        (p.like_count * 2 + p.comment_count * 3) / 
        POW(TIMESTAMPDIFF(HOUR, p.created_at, NOW()) + 2, 1.5)
    ) AS ranking_score
FROM `posts` p
INNER JOIN `users` u ON u.id = p.user_id
WHERE p.created_at >= NOW() - INTERVAL 7 DAY
ORDER BY ranking_score DESC
LIMIT 10;
