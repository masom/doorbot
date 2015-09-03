
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE accounts RENAME COLUMN notifications_sms_enabled TO notifications_twilio_enabled;
ALTER TABLE accounts RENAME COLUMN notifications_sms_source_phone_number TO notifications_twilio_source_phone_number;

ALTER TABLE accounts RENAME COLUMN notifications_email_enabled TO notifications_mailgun_enabled;

ALTER TABLE accounts ADD COLUMN notifications_postmark_enabled BOOLEAN DEFAULT FALSE;

ALTER TABLE accounts ADD COLUMN notifications_slack_token VARCHAR(128) DEFAULT FALSE;
ALTER TABLE accounts ADD COLUMN notifications_slack_enabled BOOLEAN DEFAULT FALSE;

ALTER TABLE accounts ADD COLUMN notifications_nexmo_token VARCHAR(128) DEFAULT FALSE;
ALTER TABLE accounts ADD COLUMN notifications_nexmo_enabled BOOLEAN DEFAULT FALSE;

ALTER TABLE people ADD COLUMN notifications_chat_enabled BOOLEAN DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE accounts RENAME COLUMN notifications_twilio_enabled TO notifications_sms_enabled;
ALTER TABLE accounts RENAME COLUMN notifications_twilio_source_phone_number TO notifications_sms_source_phone_number;

ALTER TABLE accounts RENAME COLUMN notifications_mailgun_enabled TO notifications_email_enabled;
ALTER TABLE accounts DROP COLUMN notifications_postmark_enabled;

ALTER TABLE accounts DROP COLUMN notifications_slack_enabled;
ALTER TABLE accounts DROP COLUMN notifications_slack_token;

ALTER TABLE accounts DROP COLUMN notifications_nexmo_enabled;
ALTER TABLE accounts DROP COLUMN notifications_nexmo_token;


ALTER TABLE people DROP COLUMN notifications_chat_enabled;
