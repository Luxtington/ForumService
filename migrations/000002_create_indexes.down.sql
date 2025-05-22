DROP INDEX IF EXISTS idx_comments_post_id;
DROP INDEX IF EXISTS idx_comments_author_id;
DROP INDEX IF EXISTS idx_comments_created_at;

DROP INDEX IF EXISTS idx_posts_thread_id;
DROP INDEX IF EXISTS idx_posts_author_id;
DROP INDEX IF EXISTS idx_posts_created_at;

DROP INDEX IF EXISTS idx_threads_author_id;
DROP INDEX IF EXISTS idx_threads_created_at;

DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_created_at;