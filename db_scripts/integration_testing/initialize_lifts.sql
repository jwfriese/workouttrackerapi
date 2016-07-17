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

INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',1,'weight/reps','{1,2,3}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle press',1,'weight/reps','{4,5,6}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle push',1,'weight/reps','{7,8,9}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle press',2,'weight/reps','{10,11,12}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle cleans',2,'weight/reps','{13,14,15}');

INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,100.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,110.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,120.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,155.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,165.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,175.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,210.5,NULL,NULL,3);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,215.5,NULL,NULL,3);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,220.5,NULL,NULL,2);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,165.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,175.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,185.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,135.0,NULL,NULL,5);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,155.0,NULL,NULL,5);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,175.0,NULL,NULL,4);
