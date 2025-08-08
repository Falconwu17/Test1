CREATE TYPE record_status_type AS ENUM ('now', 'later', 'process');
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name_ VARCHAR(50),
    surname VARCHAR(50),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS records (
    record_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    timeout INTEGER NOT NULL DEFAULT 60 CHECK (timeout > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status record_status_type NOT NULL DEFAULT 'now'
);

CREATE TABLE IF NOT EXISTS entries (
    id SERIAL PRIMARY KEY,
    record_id INTEGER REFERENCES records(record_id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE auto_clean_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE  ,
    enabled BOOLEAN NOT NULL DEFAULT true,
    interval_seconds INTEGER NOT NULL DEFAULT 604800,
    last_cleaned_at TIMESTAMP DEFAULT NULL
);

INSERT INTO users (name_, surname, email, password)
VALUES ('Sagir', 'Jusupov', 'sagir.jusupov2004@gmail.com', '123456')
ON CONFLICT DO NOTHING;
