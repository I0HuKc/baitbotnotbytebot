CREATE TABLE IF NOT EXISTS changedesc(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    groupid BIGINT NOT NULL,
    groupname VARCHAR(255) NOT NULL,
    nextdescchange TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    UNIQUE (groupid)
);