-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO users (id, balance) VALUES ('de152cc3-9cbf-45c6-9081-7dff96708254', 0);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
-- +migrate StatementEnd
