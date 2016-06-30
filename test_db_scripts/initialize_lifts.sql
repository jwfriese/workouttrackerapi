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
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',3,'weight/reps','{16,17,18}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle press',3,'weight/reps','{19,20,21}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle push',3,'weight/reps','{22,23,24}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle cleans',4,'weight/reps','{25,26,27}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',4,'weight/reps','{28,29,30}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle box jumps',4,'height/reps','{31,32,33}');
INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle hangs',5,'time_in_seconds','{34}');

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
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',6,110.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',6,120.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',6,130.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',7,160.0,NULL,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',7,170.0,NULL,NULL,7);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',7,180.0,NULL,NULL,6);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',8,220.5,NULL,NULL,3);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',8,225.5,NULL,NULL,3);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',8,240.5,NULL,NULL,2);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',9,140.0,NULL,NULL,5);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',9,160.0,NULL,NULL,5);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',9,175.0,NULL,NULL,4);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',10,120.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',10,130.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',10,140.0,NULL,NULL,10);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('height/reps',11,NULL,36.0,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('height/reps',11,NULL,36.0,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('height/reps',11,NULL,42.0,NULL,8);
INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('time_in_seconds',12,NULL,NULL,35.0,NULL);

