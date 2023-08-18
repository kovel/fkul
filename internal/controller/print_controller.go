package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/kovel/fkul/internal/dto"
	"github.com/kovel/fkul/internal/fkul/client"
	"log"
)

type printController struct {
	IController
	db     *sql.DB
	client *client.Client
}

func NewPrintController(db *sql.DB, c *client.Client) *printController {
	return &printController{nil, db, c}
}

func (c *printController) Run() {
	rows, err := c.db.Query(
		"SELECT footballer, goals, year FROM fkul_world_championship_bombardier ORDER BY footballer",
	)
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	bombardiers := make([]dto.Bombardier, 0)
	for rows.Next() {
		var bombardier dto.Bombardier
		err := rows.Scan(&bombardier.Name, &bombardier.Goals, &bombardier.Year)
		if err != nil {
			log.Println(err)
			return
		}
		bombardiers = append(bombardiers, bombardier)
	}

	data, err := json.MarshalIndent(bombardiers, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
}
