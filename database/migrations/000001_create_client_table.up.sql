-- Enable the uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    username TEXT NOT NULL
);
