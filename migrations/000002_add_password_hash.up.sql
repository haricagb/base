-- migrations/000002_add_password_hash.up.sql

ALTER TABLE users
    ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '';
