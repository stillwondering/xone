CREATE TABLE `person` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL,
    `gender` TEXT NOT NULL
);

CREATE TABLE `users` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL
);
