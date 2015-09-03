
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE people ADD COLUMN is_visible BOOLEAN DEFAULT TRUE;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE people DROP COLUMN is_visible;
