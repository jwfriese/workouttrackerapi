\set ON_ERROR_STOP true

CREATE TABLE workouts(
        id serial PRIMARY KEY,
        name varchar(256),
        timestamp timestamptz,
        lifts integer[]
); 

INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle one','2016-03-07 06:26:34.0-0800','{1,2,3}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle two','2016-03-09 06:04:44.0-0800','{4,5}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle three','2016-03-11 06:12:56.0-0800','{}');


