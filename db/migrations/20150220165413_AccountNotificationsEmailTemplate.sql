
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE accounts ADD COLUMN notifications_email_message_template TEXT DEFAULT null;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE accounts DROP COLUMN notifications_email_message_template;
