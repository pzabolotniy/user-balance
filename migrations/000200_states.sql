-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO tx_states (id, name) VALUES (uuid_generate_v4(), 'win');
INSERT INTO tx_states (id, name) VALUES (uuid_generate_v4(), 'lost');
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
-- +migrate StatementEnd
