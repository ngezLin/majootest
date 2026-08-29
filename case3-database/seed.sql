-- ============================================================================
-- CASE 3: SAMPLE SEED DATA
-- ============================================================================

USE `social_media_db`;

-- 1. Seed Users
INSERT INTO `users` (`id`, `username`, `email`, `password_hash`, `full_name`, `avatar_url`, `bio`) VALUES
(1, 'alice',   'alice@mail.com',   '$2a$10$abcdef...', 'Alice Johnson',   'https://cdn.example.com/alice.png',   'Tech enthusiast and coffee lover ☕'),
(2, 'bob',     'bob@mail.com',     '$2a$10$abcdef...', 'Bob Smith',       'https://cdn.example.com/bob.png',     'Software Engineer @ Majoo | Golang & PHP'),
(3, 'charlie', 'charlie@mail.com', '$2a$10$abcdef...', 'Charlie Brown',   'https://cdn.example.com/charlie.png', 'Photographer & traveler 📷✈️'),
(4, 'diana',   'diana@mail.com',   '$2a$10$abcdef...', 'Diana Prince',    'https://cdn.example.com/diana.png',   'Product Manager | Building great tools'),
(5, 'evan',    'evan@mail.com',    '$2a$10$abcdef...', 'Evan Wright',     'https://cdn.example.com/evan.png',    'DevOps & Cloud Architect');

-- 2. Seed Follows
-- Alice follows Bob (2) and Charlie (3)
-- Bob follows Alice (1) and Diana (4)
-- Charlie follows Alice (1), Bob (2), Diana (4)
-- Diana follows Bob (2)
INSERT INTO `follows` (`follower_id`, `following_id`, `created_at`) VALUES
(1, 2, '2026-08-01 08:00:00'),
(1, 3, '2026-08-01 08:30:00'),
(2, 1, '2026-08-01 09:00:00'),
(2, 4, '2026-08-01 09:15:00'),
(3, 1, '2026-08-02 10:00:00'),
(3, 2, '2026-08-02 10:30:00'),
(3, 4, '2026-08-02 11:00:00'),
(4, 2, '2026-08-02 11:30:00'),
(5, 1, '2026-08-02 12:00:00');

-- 3. Seed Posts
INSERT INTO `posts` (`id`, `user_id`, `content`, `like_count`, `comment_count`, `created_at`) VALUES
(1, 2, 'Excited to announce our new Golang microservice architecture! #golang #backend', 3, 2, '2026-08-03 09:00:00'),
(2, 3, 'Sunset over Mount Bromo yesterday. What a breathtaking view! 🌄', 5, 1, '2026-08-03 14:30:00'),
(3, 1, 'Just finished setting up our multi-tenant reporting engine with Laravel 11. 🚀', 4, 1, '2026-08-04 10:00:00'),
(4, 4, 'Looking for feedback on our new dashboard UI. Check it out!', 2, 0, '2026-08-04 15:00:00'),
(5, 2, 'Tips for optimizing MySQL composite indexes: always check the leftmost prefix rule.', 6, 2, '2026-08-05 11:00:00');

-- 4. Seed Post Media
INSERT INTO `post_media` (`id`, `post_id`, `media_url`, `media_type`, `display_order`) VALUES
(1, 2, 'https://cdn.example.com/media/bromo_sunset.jpg', 'IMAGE', 1),
(2, 2, 'https://cdn.example.com/media/bromo_crater.jpg', 'IMAGE', 2),
(3, 3, 'https://cdn.example.com/media/dashboard_preview.png', 'IMAGE', 1),
(4, 5, 'https://cdn.example.com/media/mysql_indexing_explained.mp4', 'VIDEO', 1);

-- 5. Seed Comments
INSERT INTO `comments` (`id`, `post_id`, `user_id`, `parent_id`, `content`, `created_at`) VALUES
(1, 1, 1, NULL, 'Awesome work Bob! How are you handling distributed tracing?', '2026-08-03 09:15:00'),
(2, 1, 2, 1,    'Thanks Alice! We are using OpenTelemetry with Jaeger.',       '2026-08-03 09:30:00'),
(3, 2, 1, NULL, 'Stunning shot Charlie! Which lens did you use?',              '2026-08-03 15:00:00'),
(4, 3, 2, NULL, 'Congrats Alice! The zero-fill revenue report is super clean.', '2026-08-04 10:30:00'),
(5, 5, 4, NULL, 'Very helpful post, thanks for sharing!',                      '2026-08-05 11:30:00');

-- 6. Seed Reactions
INSERT INTO `reactions` (`id`, `user_id`, `target_type`, `target_id`, `reaction_type`, `created_at`) VALUES
(1, 1, 'POST', 1, 'LOVE', '2026-08-03 09:10:00'),
(2, 3, 'POST', 1, 'LIKE', '2026-08-03 09:20:00'),
(3, 4, 'POST', 1, 'LIKE', '2026-08-03 09:45:00'),
(4, 1, 'POST', 2, 'LOVE', '2026-08-03 14:40:00'),
(5, 2, 'POST', 2, 'LIKE', '2026-08-03 14:50:00'),
(6, 2, 'POST', 3, 'LIKE', '2026-08-04 10:15:00'),
(7, 3, 'POST', 3, 'LOVE', '2026-08-04 10:20:00');

-- 7. Seed Conversations & Messages
INSERT INTO `conversations` (`id`, `is_group`, `title`, `created_at`) VALUES
(1, FALSE, NULL, '2026-08-01 10:00:00'), -- DM: Alice & Bob
(2, FALSE, NULL, '2026-08-02 11:00:00'); -- DM: Alice & Charlie

INSERT INTO `conversation_members` (`conversation_id`, `user_id`, `last_read_at`) VALUES
(1, 1, '2026-08-03 10:00:00'),
(1, 2, '2026-08-03 09:50:00'),
(2, 1, '2026-08-02 11:05:00'),
(2, 3, '2026-08-02 11:00:00');

INSERT INTO `messages` (`id`, `conversation_id`, `sender_id`, `content`, `created_at`) VALUES
(1, 1, 1, 'Hi Bob, do you have time to review the Go PR today?', '2026-08-03 09:45:00'),
(2, 1, 2, 'Hey Alice! Sure, will take a look after lunch.',      '2026-08-03 09:48:00'),
(3, 1, 1, 'Sounds good, thank you!',                             '2026-08-03 09:50:00'),
(4, 2, 3, 'Hey Alice, sent you the photo files via drive link.', '2026-08-02 11:02:00'),
(5, 2, 1, 'Got them! They look fantastic.',                       '2026-08-02 11:05:00');

-- 8. Seed Notifications
INSERT INTO `notifications` (`id`, `recipient_id`, `actor_id`, `type`, `target_id`, `is_read`, `created_at`) VALUES
(1, 2, 1, 'FOLLOW', 1, TRUE,  '2026-08-01 08:00:00'),
(2, 2, 1, 'LIKE_POST', 1, TRUE, '2026-08-03 09:10:00'),
(3, 2, 1, 'COMMENT_POST', 1, FALSE, '2026-08-03 09:15:00'),
(4, 3, 1, 'LIKE_POST', 2, TRUE, '2026-08-03 14:40:00'),
(5, 1, 2, 'LIKE_POST', 3, FALSE, '2026-08-04 10:15:00');
