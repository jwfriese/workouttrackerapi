\set ON_ERROR_STOP true

CREATE TYPE lift_data_template AS ENUM(
    'weight/reps',
    'height/reps',
    'time_in_seconds',
    'weight/time_in_seconds'
);

CREATE TABLE lifts(
    id serial PRIMARY KEY,
    name varchar(256),
    workout integer NOT NULL,
    data_template lift_data_template NOT NULL,
    sets integer[]
);

CREATE TABLE sets(
    id serial PRIMARY KEY,
    data_template lift_data_template NOT NULL,
    lift integer NOT NULL,
    weight real,
    height real,
    time_in_seconds real,
    reps integer
);

