-- Удаляем все таблицы в правильном порядке
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;

-- Создаем тип user_role
CREATE TYPE user_role AS ENUM ('admin', 'user');

-- Создаем таблицу users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индексы для users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Создаем тестового пользователя
INSERT INTO users (id, username, email, password, role) 
VALUES (1, 'test_user', 'test@example.com', 'password123', 'user');

-- Создаем таблицу threads
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT threads_author_id_fkey FOREIGN KEY (author_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

-- Создаем индексы для threads
CREATE INDEX idx_threads_author_id ON threads(author_id);
CREATE INDEX idx_threads_created_at ON threads(created_at); 