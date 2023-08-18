package controller

import (
	"database/sql"
	"github.com/kovel/fkul/internal/fkul/client"
	"log"
)

type collectorController struct {
	IController
	db     *sql.DB
	client *client.Client
}

func NewCollectorController(db *sql.DB, c *client.Client) IController {
	return &collectorController{nil, db, c}
}

func (c *collectorController) Run() {
	log.Println("Collecting data about footballers from football.kulichki.net")

	bombardiers, err := c.client.Bombardiers2022WC()
	if err != nil {
		log.Fatalln(err)
	}

	for _, bombardier := range bombardiers {
		_, err := c.db.Exec(
			"INSERT INTO fkul_world_championship_bombardier (footballer, goals, year) VALUES  ($1, $2, $3)",
			bombardier.Name, bombardier.Goals, bombardier.Year,
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
