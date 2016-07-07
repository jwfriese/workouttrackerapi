\set ON_ERROR_STOP true

CREATE TABLE workouts(
        id serial PRIMARY KEY,
        name varchar(256),
        timestamp timestamptz,
        lifts integer[]
); 

INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle one','2016-03-07 06:26:34.0-0800','{1,2,3}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle two','2016-03-09 06:04:44.0-0800','{4,5}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle three','2016-03-11 05:57:56-0800','{6,7,8}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle four','2016-03-15 05:57:00-0700','{9,10,11}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle five','2016-03-16 05:46:51-0700','{12}');

