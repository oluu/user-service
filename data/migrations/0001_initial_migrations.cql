CREATE KEYSPACE IF NOT EXISTS user WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};

CREATE TABLE IF NOT EXISTS user (
    id text,
    email text,
    first_name text,
    last_name text,
    PRIMARY KEY(id)
);