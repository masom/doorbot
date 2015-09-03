
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE devices ADD COLUMN token CHAR(40) DEFAULT '' NOT NULL;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE devices DROP COLUMN token;
