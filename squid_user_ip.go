package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "./configuration"
)

func main() {
	var configfile = flag.String("conf", "conf.toml", "Config file name")
	var groupId = flag.Int("g", 0, "Group ID")
	flag.Parse()
	_, err := os.Stat(*configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", *configfile)
	}
	var config = ReadConfig(*configfile)

	db, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT COUNT(id) FROM usersip WHERE ip = INET_ATON(?) AND status = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var Count int = 0
	var n string
	for true {
		_, err := fmt.Scanln(&n)
		if err != nil {
			log.Fatal(err)
		}
		err = stmt.QueryRow(n, *groupId).Scan(&Count)
		if err != nil {
			log.Fatal(err)
		}
		if Count > 0 {
			fmt.Println("OK")
		} else {
			fmt.Println("ERR")
		}
		Count = 0
	}
	fmt.Println(err)
}