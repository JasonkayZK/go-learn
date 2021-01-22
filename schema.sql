CREATE DATABASE IF NOT EXISTS `go_graphql_db`;

USE `go_graphql_db`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  id serial PRIMARY KEY,
  name VARCHAR (50) NOT NULL,
  age INT NOT NULL,
  profession VARCHAR (50) NOT NULL,
  friendly BOOLEAN NOT NULL
) ENGINE=InnoDB CHARSET=utf8mb4;

INSERT INTO users VALUES
  (1, 'kevin', 35, 'waiter', true),
  (2, 'angela', 21, 'concierge', true),
  (3, 'alex', 26, 'zoo keeper', false),
  (4, 'becky', 67, 'retired', false),
  (5, 'kevin', 15, 'in school', true),
  (6, 'frankie', 45, 'teller', true);
