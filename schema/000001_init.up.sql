
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

-- EXTENSIONS --

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --
CREATE TABLE users
(
    id            serial       not null unique,
    name          TEXT not null,
    username      TEXT not null unique,
    password_hash TEXT not null
);

CREATE TABLE todo_lists
(
    id          serial       not null unique,
    title       TEXT not null,
    description TEXT
);

CREATE TABLE users_lists
(
    id      serial                                           not null unique,
    user_id int references users (id) on delete cascade      not null,
    list_id int references todo_lists (id) on delete cascade not null
);

CREATE TABLE todo_items
(
    id          serial       not null unique,
    title       TEXT not null,
    description TEXT,
    done        boolean      not null default false
);


CREATE TABLE lists_items
(
    id      serial                                           not null unique,
    item_id int references todo_items (id) on delete cascade not null,
    list_id int references todo_lists (id) on delete cascade not null
);