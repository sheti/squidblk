package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	. "./configuration"
	"bufio"
)

func in_array(val string, array []string) bool {

	for _, v := range array {
		if val == v {
			return true

		}
	}

	return false
}

func main() {
	var configfile= flag.String("conf", "conf.toml", "Config file name")
	flag.Parse()
	_, err := os.Stat(*configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", *configfile)
	}
	var config= ReadConfig(*configfile)

	files, err := ioutil.ReadDir(config.Blacklistdir)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO blocklist(val, category, `type`) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, dir := range files {
		if dir.IsDir() && in_array(dir.Name(), config.Category) {
			fmt.Println(dir.Name())
			files, err := ioutil.ReadDir(config.Blacklistdir + "/" + dir.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				if (file.Name() == "domains") || (file.Name() == "urls") {
					lines, err := os.Open(config.Blacklistdir + "/" + dir.Name() + "/" + file.Name())
					if err != nil {
						lines.Close()
						log.Fatal(err)
					}

					scanner := bufio.NewScanner(lines)
					for scanner.Scan() {
						_, err := stmt.Exec(scanner.Text(), dir.Name(), file.Name())
						if err != nil {
							log.Fatal(err)
						}
					}

					if err := scanner.Err(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}
