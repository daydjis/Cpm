
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       telegram_id BIGINT NOT NULL UNIQUE,
                       username VARCHAR(50),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE(user_id, category_id)
);

CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,
                               category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
                               item_name VARCHAR(255) NOT NULL,
                               item_url TEXT NOT NULL,
                               sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE(category_id, item_url)
);
