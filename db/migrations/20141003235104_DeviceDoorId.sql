
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE devices ADD COLUMN door_id INTEGER DEFAULT NULL;
ALTER TABLE devices ADD FOREIGN KEY (door_id) REFERENCES doors(id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE devices DROP COLUMN door_id;
