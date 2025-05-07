CREATE TABLE threads (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(255) NOT NULL,
                         created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       thread_id INTEGER REFERENCES threads(id) ON DELETE CASCADE,
                       user_id INTEGER NOT NULL,
                       content TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
                          user_id INTEGER NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW()
);