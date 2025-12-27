DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    id         UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    username   TEXT        NOT NULL UNIQUE,
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS packages CASCADE;
CREATE TABLE packages
(
    id          UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name        TEXT        NOT NULL UNIQUE,
    description TEXT        NOT NULL,
    visibility  INTEGER     NOT NULL,
    owner_id    UUID        NOT NULL REFERENCES users (id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS package_versions CASCADE;
CREATE TABLE package_versions
(
    id         UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    package_id UUID        NOT NULL REFERENCES packages (id),
    version    TEXT        NOT NULL,
    checksum   TEXT        NOT NULL,
    size_bytes BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (package_id, version)
);

DROP TABLE IF EXISTS dependencies CASCADE;
CREATE TABLE dependencies
(
    version_id      UUID NOT NULL REFERENCES package_versions (id),
    dependency_name TEXT NOT NULL,
    constraint_expr TEXT NOT NULL
);
