CREATE TYPE record_status_type as ENUM ('now', 'later', 'process');

CREATE TABLE IF NOT EXISTS records (
    record_id SERIAL PRIMARY KEY ,
    timeout integer not null default 60 check (timeout > 0),
    created_at  timestamp not null DEFAULT CURRENT_TIMESTAMP,
    status record_status_type not null default 'now'
);
CREATE TABLE IF NOT EXISTS entries (
    id SERIAL PRIMARY KEY ,
    record_id INTEGER REFERENCES records(record_id) ON DELETE CASCADE,
    data jsonb not null ,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
    );
CREATE TABLE IF NOT EXISTS user (
    id SERIAL PRIMARY KEY ,
    name VARCHAR(50) ,
    surname varchar(50),
    email varchar(50)
)