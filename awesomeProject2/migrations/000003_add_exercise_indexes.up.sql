CREATE INDEX IF NOT EXISTS exercise_title_idx ON exercise USING GIN (to_tsvector('simple', title));

