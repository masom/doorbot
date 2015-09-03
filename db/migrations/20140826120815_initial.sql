
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE NOT NULL,
    host VARCHAR(255) UNIQUE NOT NULL,

    notifications_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notifications_email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notifications_sms_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    notifications_sms_message_template varchar(140) DEFAULT null,

    bridge_hub_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    bridge_hub_url VARCHAR(255),
    bridge_hub_token VARCHAR(255),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE people (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL  REFERENCES accounts(id),
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,

    notifications_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notifications_sms_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notifications_email_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    notifications_app_enabled BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX people_account_id ON people (account_id);
CREATE UNIQUE INDEX people_account_id_email ON people (account_id, email);

CREATE TABLE bridge_users (
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    bridge_id INTEGER NOT NULL,
    person_id INTEGER NOT NULL REFERENCES people(id),
    user_id VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (account_id, bridge_id, person_id, user_id)
);

CREATE TABLE administrators (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE administrator_authentications (
    administrator_id INTEGER NOT NULL  REFERENCES administrators (id) ON DELETE CASCADE,
    provider_id      INTEGER NOT NULL,
    token            TEXT,
    last_used_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX administrator_authentications_administrator_id_provider_id ON administrator_authentications (administrator_id, provider_id);

CREATE TABLE authentications (
    id SERIAL NOT NULL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    person_id INTEGER NOT NULL REFERENCES people(id),
    provider_id INTEGER NOT NULL,
    token TEXT
);

CREATE UNIQUE INDEX authentications_account_id_person_id_provider_id ON authentications (account_id, person_id, provider_id);

CREATE TABLE doors (
    id SERIAL PRIMARY KEY,
    account_id INTEGER  NOT NULL  REFERENCES accounts (id) ,
    name VARCHAR(255) NOT NULL
);

CREATE INDEX doors_account_id ON doors (account_id);

CREATE TABLE devices (
    id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER NOT NULL  REFERENCES accounts (id),
    name varchar(255),
    device_id varchar(255),
    make varchar(255),
    description varchar(255),
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX devices_account_id_device_id ON devices (account_id, device_id);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    account_id INTEGER  REFERENCES accounts(id) NOT NULL,
    door_id INTEGER  REFERENCES doors(id),
    device_id INTEGER  REFERENCES devices(id),
    event_id INTEGER NOT NULL,
    person_id INTEGER  REFERENCES people(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX events_account_id ON events (account_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX events_account_id;
DROP TABLE events;

DROP INDEX devices_account_id_device_id;
DROP TABLE devices;

DROP INDEX doors_account_id;
DROP TABLE doors;

DROP INDEX authentications_account_id_person_id_provider_id;

DROP TABLE authentications;

DROP TABLE administrator_authentications;

DROP TABLE administrators;

DROP TABLE bridge_users;

DROP INDEX people_account_id_email;
DROP INDEX people_account_id;
DROP TABLE people;

DROP TABLE accounts;
