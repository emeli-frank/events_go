DROP TABLE IF EXISTS users;
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

DROP TABLE IF EXISTS events;
CREATE TABLE events
(
    id SERIAL,
    title VARCHAR(64),
    description VARCHAR (512),
    is_virtual BOOLEAN,
    address VARCHAR(128),
    link VARCHAR(128),
    number_of_seats INT,
    start_time timestamptz,
    end_time timestamptz,
    welcome_message VARCHAR (256),
    is_published BOOLEAN,
    host INT NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (host)
        REFERENCES users (id)
        ON DELETE CASCADE
);

DROP TABLE IF EXISTS user_events;
CREATE TABLE user_events
(
    user_id INT,
    event_id INT,

    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,
    FOREIGN KEY (event_id)
        REFERENCES events (id)
        ON DELETE CASCADE
);
