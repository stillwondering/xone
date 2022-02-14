CREATE TABLE `person` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `public_id` TEXT NOT NULL UNIQUE,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL
);

CREATE TABLE `person_history` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `created_at` TEXT NOT NULL,
    `person_id` INTEGER REFERENCES `person`(`id`) ON DELETE CASCADE,
    `public_id` TEXT NOT NULL,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL
);

CREATE TRIGGER update_history_after_insert_person
    AFTER INSERT ON person
BEGIN
    INSERT INTO person_history (
        created_at,
        person_id,
        public_id,
        first_name,
        last_name,
        date_of_birth
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.public_id,
        NEW.first_name,
        NEW.last_name,
        NEW.date_of_birth
    );
END;

CREATE TRIGGER update_history_after_update_person
    AFTER UPDATE ON person
BEGIN
    INSERT INTO person_history (
        created_at,
        person_id,
        public_id,
        first_name,
        last_name,
        date_of_birth
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.public_id,
        NEW.first_name,
        NEW.last_name,
        NEW.date_of_birth
    );
END;

CREATE TABLE `users` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL
);
