CREATE TABLE `person` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `pid` TEXT NOT NULL UNIQUE,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL,
    `gender` TEXT NOT NULL
);