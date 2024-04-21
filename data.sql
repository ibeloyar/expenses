DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INT AUTO_INCREMENT NOT NULL,
    username VARCHAR(128) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email_confirmed BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (`id`)
);

INSERT INTO users (id, username, email, password, email_confirmed)
VALUES (0, 'admin', 'admin@email.com', 'hashpassword23124', false);
