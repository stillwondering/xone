INSERT INTO person (
    id,
    public_id,
    first_name,
    last_name,
    date_of_birth,
    email,
    phone,
    mobile
) VALUES
(1, "1", "Harry", "Potter", "1980-07-31", "harry.potter@hogwarts.co.uk", "", ""),
(2, "2", "Ron", "Weasley", "", "ron.weasley@hogwarts.co.uk", "", ""),
(3, "3", "Hermione", "Granger", "1979-09-19", "hermione.granger@hogwarts.co.uk", "", "");

INSERT INTO membership_type (
    id,
    name
) VALUES
(1, "active"),
(2, "passive");

INSERT INTO membership (
    id,
    type_id,
    person_id,
    effective_from
) VALUES
(1, 1, 1, "1998-07-31"),
(2, 2, 1, "2045-07-31")
;