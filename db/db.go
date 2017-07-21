package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ZeaLoVe/alarm/g"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/common/model"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxIdleConns(g.Config().MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
	go UpdateHostMap()
}

//func getFormatTime() string {
//	cur := time.Now()
//	year, month, day := cur.UTC().Date()
//	hour, minute, second := cur.UTC().Clock()

//	return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", year, month, day, hour, minute, second)
//}

func RecordEvent(event *model.Event) error {
	sql := fmt.Sprintf("insert into auto_server_status Set level='%v', end_point='%v', description='%v', status='%v', times='%v', ip='%v', error_time='%v'",
		event.Priority(),
		event.Endpoint,
		event.Note(),
		event.Status,
		event.CurrentStep,
		HostMap.GetIPByEndpoint(event.Endpoint),
		event.FormattedTime(),
	)
	//	log.Println(sql)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
	return err
}
