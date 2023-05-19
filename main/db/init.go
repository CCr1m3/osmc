package db

import (
	"time"

	"github.com/ccr1m3/osmc/main/env"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func GetInstance() *sqlx.DB {
	if db == nil {
		var err error
		if env.DB.Type == env.MYSQL {
			db, err = sqlx.Open("mysql", env.DB.Path)
			if err != nil {
				log.Fatal(err)
			}
			db.SetConnMaxLifetime(time.Minute * 5)
			db.SetMaxOpenConns(1)
			db.SetMaxIdleConns(1)
		}
		return db
	} else {
		return db
	}
}

func Init() {
	log.Info("starting db service")
	GetInstance()
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}

func Clear() {
	GetInstance()
	if db != nil {
		if env.DB.Type == env.MYSQL {
			tx, err := db.Beginx()
			if err != nil {
				log.Error("error starting transaction:" + err.Error())
			}
			rows := []struct {
				Tables_in_euos string `db:"Tables_in_euos"`
			}{}
			tx.Exec("SET foreign_key_checks = 0")
			err = tx.Select(&rows, "SHOW TABLES in euos")
			if err != nil {
				log.Error("error getting database table:" + err.Error())
			}
			for _, row := range rows {
				_, err := tx.Exec("DROP TABLE " + row.Tables_in_euos)
				if err != nil {
					log.Error("error dropping table:" + err.Error())
				}
			}
			tx.Exec("SET foreign_key_checks = 1")
			tx.Commit()
			err = db.Close()
			if err != nil {
				log.Error("failed to close db: " + err.Error())
			}
		}
	}
	db = nil
}
