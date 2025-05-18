ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;
 
-- Создаем администратора (пароль: admin)
INSERT INTO users (username, password_hash, role) 
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'); 