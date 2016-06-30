\set ON_ERROR_STOP true

CREATE TABLE workouts(
        id serial PRIMARY KEY,
        name varchar(256),
        timestamp timestamptz,
        lifts varchar(128)[]
); 

INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle one','2016-03-07 06:26:34.0-0800','{turtle lift,turtle press,turtle push}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle two','2016-03-09 06:04:44.0-0800','{turtle press,turtle cleans}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle three','2016-03-11 05:57:56-0800','{turtle lift,turtle press,turtle push}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle four','2016-03-15 05:57:00-0700','{turtle cleans,turtle lift,turtle box jumps}');
INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle five','2016-03-16 05:46:51-0700','{turtle hangs}');

