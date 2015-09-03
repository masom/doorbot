
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE people ADD COLUMN account_type INTEGER  DEFAULT 0;
ALTER TABLE accounts ADD COLUMN contact_email_confirmed BOOLEAN DEFAULT false;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE people DROP COLUMN account_type;
ALTER TABLE accounts DROP COLUMN contact_email_confirmed;
