package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Db struct {
	Db *sql.DB
}

func NewDb() *Db {
	return &Db{}
}

func (d *Db) Conectar() *sql.DB {
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes_dolar (
        id_cotacao INTEGER PRIMARY KEY AUTOINCREMENT,
        valor_cotacao VARCHAR(50) NOT NULL,
     	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
