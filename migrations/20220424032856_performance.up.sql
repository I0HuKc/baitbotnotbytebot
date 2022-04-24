CREATE TABLE IF NOT EXISTS performance(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    groupid BIGINT NOT NULL,
    groupname VARCHAR(255) NOT NULL,
    nextjoke TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    UNIQUE (groupid)
);