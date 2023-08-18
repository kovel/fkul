package main

import (
	"database/sql"
	"flag"
	"github.com/kovel/fkul/internal/controller"
	"github.com/kovel/fkul/internal/fkul/client"
	_ "github.com/marcboeker/go-duckdb"
	"log"
	"strings"
)

func main() {
	log.SetFlags(log.Flags() | log.Llongfile)

	helpFlag := flag.Bool("help", false, "print help")
	collectFlag := flag.Bool("collect", false, "collecting data from football.kulichki.net to DuckDB")
	printFlag := flag.Bool("print", false, "print data about bombardiers from DuckDB")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	// creating DuckDB dabase
	db, err := sql.Open("duckdb", "./fkul.db")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()

	// testing db connection
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
		return
	}

	// create FK client and testing connection
	c := client.NewClient()
	responseText, err := c.Ping()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(strings.TrimSpace(responseText))

	// controllers mapping
	var ctl controller.IController
	switch true {
	case *collectFlag:
		ctl = controller.NewCollectorController(db, c)
	case *printFlag:
		ctl = controller.NewPrintController(db, c)
	default:
		log.Fatalln("Cannot find controller")
	}
	ctl.Run()
}
