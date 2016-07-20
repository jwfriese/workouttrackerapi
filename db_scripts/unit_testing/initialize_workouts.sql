\set ON_ERROR_STOP true

CREATE TABLE workouts(
        id serial PRIMARY KEY,
        name varchar(256),
        timestamp timestamptz,
        lifts integer[]
); 

