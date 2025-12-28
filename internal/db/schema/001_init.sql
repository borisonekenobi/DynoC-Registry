DROP TYPE IF EXISTS VISIBILITY CASCADE;
CREATE TYPE VISIBILITY AS ENUM ('public', 'private');

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    id         UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    username   TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    password   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_user_id PRIMARY KEY (id),
    CONSTRAINT ak_username UNIQUE (username),
    CONSTRAINT ak_email UNIQUE (email)
);

DROP TABLE IF EXISTS packages CASCADE;
CREATE TABLE packages
(
    id          UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL,
    visibility  VISIBILITY  NOT NULL,
    owner_id    UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_package_id PRIMARY KEY (id),
    CONSTRAINT fk_package_owner FOREIGN KEY (owner_id) REFERENCES users (id),
    CONSTRAINT ak_package_name UNIQUE (name)
);

DROP TABLE IF EXISTS package_versions CASCADE;
CREATE TABLE package_versions
(
    id                  UUID        NOT NULL DEFAULT GEN_RANDOM_UUID(),
    package_id          UUID        NOT NULL,
    version             TEXT        NOT NULL,
    checksum            TEXT        NOT NULL,
    size_bytes          BIGINT      NOT NULL,
    location            TEXT        NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_package_version_id PRIMARY KEY (id),
    CONSTRAINT fk_package_id FOREIGN KEY (package_id) REFERENCES packages (id),
    CONSTRAINT ak_package_version UNIQUE (package_id, version)
);

DROP TABLE IF EXISTS dependencies CASCADE;
CREATE TABLE dependencies
(
    version_id      UUID NOT NULL,
    dependency_name TEXT NOT NULL,
    constraint_expr TEXT NOT NULL,

    CONSTRAINT pk_dependency PRIMARY KEY (version_id, dependency_name),
    CONSTRAINT fk_dependency_version FOREIGN KEY (version_id) REFERENCES package_versions (id)
);
