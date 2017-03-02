package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"strings"
	. "./configuration"
)

func main()  {
	var configfile= flag.String("conf", "conf.toml", "Config file name")
	var blocktype= flag.String("t", "domains", "Type of blocks")
	flag.Parse()
	_, err := os.Stat(*configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", *configfile)
	}
	var config= ReadConfig(*configfile)

	db, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT COUNT(id) FROM blocklist WHERE val = ? AND `type` = ? AND category IN(\"" + strings.Join(config.Blockcategory,"\",\"") + "\")")
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
		err = stmt.QueryRow(n, *blocktype).Scan(&Count)
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
