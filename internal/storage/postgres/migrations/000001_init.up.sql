-- Таблица ролей пользователей
-- На первом этапе предзаполнены, 2 роли user - id: 10, root - id: 1
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
INSERT INTO roles (id, name) VALUES (1, 'admin');
INSERT INTO roles (id, name) VALUES (10, 'user');

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL ,
    email_confirmed BOOLEAN DEFAULT FALSE,
    password VARCHAR(255) NOT NULL,
    role_id INTEGER DEFAULT 10,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT users_unique_login UNIQUE(login),
    CONSTRAINT users_unique_email UNIQUE(email),
    CONSTRAINT users_empty_login CHECK (LENGTH(login) > 0),
    CONSTRAINT users_empty_email CHECK (LENGTH(email) > 0),
    CONSTRAINT users_empty_password CHECK (LENGTH(password) > 3),
    CONSTRAINT users_role_id_key FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE SET DEFAULT
);

-- Таблица токенов пользователей
-- Нужно обдумать механизм отчистки, может по created_at раз в день?
-- created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id INTEGER NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(user_id, token)
);

-- Таблица категорий транзакций
-- Есть общие категории, есть созданные пользователем.
-- Общие предзаполнить
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER, -- Если не пользовательская NULL
    name VARCHAR(255) NOT NULL,
    description VARCHAR(2048) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT categories_user_category_name UNIQUE(user_id, name),
    CONSTRAINT categories_user_id_key FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица контрагентов
CREATE TABLE IF NOT EXISTS counterparties (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(2048) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT counterparties_user_category_name UNIQUE(user_id, name),
    CONSTRAINT counterparties_user_id_key FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    category_id INTEGER,
    counterparty_id INTEGER,
    type VARCHAR(20) NOT NULL,
    date TIMESTAMP NOT NULL,
    amount DECIMAL(11,2) NOT NULL,
    currency VARCHAR(20) NOT NULL,
    comment VARCHAR(2048),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    FOREIGN KEY (counterparty_id) REFERENCES counterparties(id) ON DELETE SET NULL,
    CONSTRAINT transactions_type CHECK (type IN ('in', 'out')),
    CONSTRAINT transactions_amount CHECK (amount >= 0)
);