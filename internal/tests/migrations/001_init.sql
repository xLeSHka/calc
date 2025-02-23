-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS expressions (
   id BIGSERIAL PRIMARY KEY,
   expression VARCHAR NOT NULL,
   status VARCHAR NOT NULL,
   result DECIMAL default null
);

-- Add DROP statements for the UP migration here