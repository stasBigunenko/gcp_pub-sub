DROP TABLE IF EXISTS user_activities;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS actions;

CREATE TABLE products
(
    id          CHAR(36) PRIMARY KEY NOT NULL,
    name        VARCHAR(50)          NOT NULL,
    description VARCHAR(500)         NOT NULL,
    price       FLOAT(8)             NOT NULL,
    category_id CHAR(36)             NOT NULL
);

CREATE TABLE categories
(
    id   CHAR(36) PRIMARY KEY NOT NULL,
    name VARCHAR(50)          NOT NULL
);

CREATE TABLE actions
(
    id   CHAR(36) PRIMARY KEY NOT NULL,
    name VARCHAR(50)          NOT NULL
);

CREATE TABLE user_activities
(
    id         CHAR(36) PRIMARY KEY NOT NULL,
    action_id  CHAR(36)             NOT NULL,
    product_id CHAR(36),
    created_at TIMESTAMP,
    FOREIGN KEY (action_id) REFERENCES actions (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE
);