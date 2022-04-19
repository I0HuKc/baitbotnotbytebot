CREATE TABLE desctb(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    text VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now(),   

    UNIQUE (text)
);