
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE accounts ADD COLUMN contact_name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE accounts ADD COLUMN contact_email VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE accounts ADD COLUMN contact_phone_number VARCHAR(20) NOT NULL DEFAULT '';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE accounts DROP COLUMN contact_name;
ALTER TABLE accounts DROP COLUMN contact_email;
ALTER TABLE accounts DROP COLUMN contact_phone_number;
