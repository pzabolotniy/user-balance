-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS source_types (
                                            id uuid PRIMARY KEY,
                                            name varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
                                     id uuid PRIMARY KEY,
                                     balance double precision NOT NULL
);

CREATE TABLE IF NOT EXISTS tx_states (
                                         id uuid PRIMARY KEY,
                                         name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions (
                                            id uuid PRIMARY KEY,
                                            external_tx_id varchar(255) NOT NULL UNIQUE,
                                            user_id uuid REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
                                            tx_state_id uuid NOT NULL REFERENCES tx_states(id) ON DELETE RESTRICT ON UPDATE CASCADE,
                                            amount double precision NOT NULL,
                                            received_at timestamp without time zone NOT NULL
);

CREATE INDEX IF NOT EXISTS transactions_external_tx_id ON transactions USING hash (external_tx_id);

CREATE TABLE IF NOT EXISTS canceled_txs (
                                            id uuid PRIMARY KEY,
                                            tx_id uuid REFERENCES transactions(id) ON DELETE RESTRICT ON UPDATE CASCADE,
                                            canceled_at timestamp without time zone NOT NULL
);

CREATE INDEX IF NOT EXISTS canceled_txs_tx_id ON canceled_txs USING hash (tx_id);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS canceled_txs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tx_states;
DROP TABLE IF EXISTS source_types;
-- +migrate StatementEnd
