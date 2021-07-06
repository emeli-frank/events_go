CREATE TABLE users
(
    id SERIAL,
    names VARCHAR(64) NOT NULL,
    email VARCHAR (128) NOT NULL,
    password CHAR(60) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,

    PRIMARY KEY (id),
    UNIQUE (email)
);

CREATE TABLE invitations
(
    id SERIAL,
    title VARCHAR(64) NOT NULL,
    description VARCHAR (128) NOT NULL,
    is_virtual BOOLEAN,
    address VARCHAR(128),
    link VARCHAR(128),
    seat_number INT,
    start_time timestamptz,
    end_time timestamptz,
    welcome_message text,

    PRIMARY KEY (id)
);

    CREATE TABLE user_invitations
    (
        user_id INT,
        invitation_id INT,

        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,
        FOREIGN KEY (invitation_id)
            REFERENCES invitations (id)
            ON DELETE CASCADE
    );
