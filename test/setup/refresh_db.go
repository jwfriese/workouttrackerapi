package setup

import (
	"database/sql"
	"log"
)

func RefreshDatabase(openConnection *sql.DB) {
	_, _ = openConnection.Exec("TRUNCATE sets")
	_, _ = openConnection.Exec("TRUNCATE lifts")
	_, err := openConnection.Exec("TRUNCATE workouts")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = openConnection.Exec("ALTER SEQUENCE sets_id_seq RESTART WITH 1")
	_, _ = openConnection.Exec("ALTER SEQUENCE lifts_id_seq RESTART WITH 1")
	_, err = openConnection.Exec("ALTER SEQUENCE workouts_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = openConnection.Exec("INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle one','2016-03-07T06:26:34-08:00','{1,2,3}')")
	_, _ = openConnection.Exec("INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle two','2016-03-09T06:04:44-08:00','{4,5}')")
	_, err = openConnection.Exec("INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle three','2016-03-11T06:12:56-08:00','{}')")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',1,'weight/reps','{1,2,3}')")
	_, _ = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle press',1,'weight/reps','{4,5,6}')")
	_, _ = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle push',1,'weight/reps','{7,8,9}')")
	_, _ = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle press',2,'weight/reps','{10,11,12}')")
	_, err = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle cleans',2,'weight/reps','{13,14,15}')")
	if err != nil {
		log.Fatal(err)
	}

	_, _ = openConnection.Exec("")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,100.0,NULL,NULL,10)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,110.0,NULL,NULL,10)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',1,120.0,NULL,NULL,10)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,155.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,165.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',2,175.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,210.5,NULL,NULL,3)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,215.5,NULL,NULL,3)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',3,220.5,NULL,NULL,2)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,165.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,175.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',4,185.0,NULL,NULL,8)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,135.0,NULL,NULL,5)")
	_, _ = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,155.0,NULL,NULL,5)")
	_, err = openConnection.Exec("INSERT INTO sets(data_template,lift,weight,height,time_in_seconds,reps) VALUES ('weight/reps',5,175.0,NULL,NULL,4)")
	if err != nil {
		log.Fatal(err)
	}
}
