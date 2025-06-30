CREATE TABLE IF NOT EXISTS students (
    students_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    age INT,
    curs INT
    );
CREATE TABLE IF NOT EXISTS records (
    record_id SERIAL PRIMARY KEY ,
    timeout integer,
    handler_type varchar(100) ,
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
    status varchar(50)
);
CREATE TABLE IF NOT EXISTS entries (
    id SERIAL PRIMARY KEY ,
    record_id INTEGER REFERENCES records(record_id) ON DELETE CASCADE,
    data jsonb ,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
    );