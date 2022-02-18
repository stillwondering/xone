CREATE TABLE `person` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `public_id` TEXT NOT NULL UNIQUE,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL,
    `email` TEXT DEFAULT '',
    `phone` TEXT DEFAULT '',
    `mobile` TEXT DEFAULT ''
);

CREATE TABLE `person_history` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `created_at` TEXT NOT NULL,
    `person_id` INTEGER REFERENCES `person`(`id`) ON DELETE CASCADE,
    `public_id` TEXT NOT NULL,
    `first_name` TEXT NOT NULL,
    `last_name` TEXT NOT NULL,
    `date_of_birth` TEXT NOT NULL,
    `email` TEXT DEFAULT '',
    `phone` TEXT DEFAULT '',
    `mobile` TEXT DEFAULT ''
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
        date_of_birth,
        email,
        phone,
        mobile
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.public_id,
        NEW.first_name,
        NEW.last_name,
        NEW.date_of_birth,
        NEW.email,
        NEW.phone,
        NEW.mobile
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
        date_of_birth,
        email,
        phone,
        mobile
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.public_id,
        NEW.first_name,
        NEW.last_name,
        NEW.date_of_birth,
        NEW.email,
        NEW.phone,
        NEW.mobile
    );
END;

--
-- Membership type
--
CREATE TABLE `membership_type` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `name` TEXT NOT NULL UNIQUE
);

CREATE TABLE `membership_type_history` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `created_at` TEXT NOT NULL,
    `membership_type_id` INTEGER NOT NULL,
    `name` TEXT NOT NULL
);

CREATE TRIGGER membership_type_insert
    AFTER INSERT ON membership_type
BEGIN
    INSERT INTO membership_type_history (
        created_at,
        membership_type_id,
        name
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.name
    );
END;

CREATE TRIGGER membership_type_update
    AFTER UPDATE ON membership_type
BEGIN
    INSERT INTO membership_type_history (
        created_at,
        membership_type_id,
        name
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.name
    );
END;

--
-- Membership
--
CREATE TABLE `membership` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `type_id` INTEGER NOT NULL REFERENCES `membership_type`(`id`) ON DELETE RESTRICT,
    `person_id` INTEGER NOT NULL REFERENCES `person`(`id`) ON DELETE CASCADE,
    `effective_from` TEXT NOT NULL DEFAULT (date())
);

CREATE TABLE `membership_history` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `created_at` TEXT NOT NULL,
    `membership_id` INTEGER NOT NULL REFERENCES `membership`(`id`) ON DELETE CASCADE,
    `type_id` INTEGER NOT NULL,
    `person_id` INTEGER NOT NULL,
    `effective_from` TEXT NOT NULL
);

CREATE TRIGGER membership_insert
    AFTER INSERT ON membership
BEGIN
    INSERT INTO membership_history (
        created_at,
        membership_id,
        type_id,
        person_id,
        effective_from
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.type_id,
        NEW.person_id,
        NEW.effective_from
    );
END;

CREATE TRIGGER membership_update
    AFTER UPDATE ON membership
BEGIN
    INSERT INTO membership_history (
        created_at,
        membership_id,
        type_id,
        person_id,
        effective_from
    ) VALUES (
        datetime(),
        NEW.id,
        NEW.type_id,
        NEW.person_id,
        NEW.effective_from
    );
END;

CREATE TABLE `users` (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `email` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL
);
