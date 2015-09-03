
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE accounts ADD COLUMN notifications_sms_source_phone_number VARCHAR(20) DEFAULT null;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE accounts DROP COLUMN notifications_sms_source_phone_number;
