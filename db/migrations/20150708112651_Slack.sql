-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE accounts ADD COLUMN bridge_slack_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE accountS ADD COLUMN bridge_slack_token VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE accounts DROP COLUMN bridge_slack_enabled;
ALTER TABLE accounts DROP COLUMN bridge_slack_token;
