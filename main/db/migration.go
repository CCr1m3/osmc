package db

import (
	"fmt"
	"strings"

	"github.com/ccr1m3/osmc/main/static"
	log "github.com/sirupsen/logrus"
)

var migrations = []string{
	migration0,
	migration1,
}

func migrate() error {
	var start int
	_, err := db.Exec(migrations[0])
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") && !strings.Contains(err.Error(), "1062") {
		return static.ErrDB(err)
	}
	start, err = getLatestMigration()
	if err != nil {
		return static.ErrDB(err)
	}
	for i := start + 1; i < len(migrations); i++ {
		log.Info(fmt.Sprintf("applying migration %d", i))
		_, err = db.Exec(migrations[i])
		if err != nil {
			return static.ErrDB(err)
		}
		_, err = db.Exec(`INSERT INTO migrations (version) VALUES (?)`, i)
		if err != nil {
			return static.ErrDB(err)
		}
	}
	return nil
}

func getLatestMigration() (int, error) {
	ver := 0
	row := db.QueryRow(`SELECT version
	FROM migrations
	ORDER BY version DESC
	LIMIT 1`)
	err := row.Scan(&ver)
	if err != nil {
		return 0, static.ErrDB(err)
	}
	return ver, err
}

var migration0 = `CREATE TABLE IF NOT EXISTS migrations (
	version INTEGER,
	PRIMARY KEY (version)
	);
INSERT INTO migrations (version) VALUES (0);`

var migration1 = `CREATE TABLE players (
    discordID VARCHAR(100) UNIQUE NOT NULL,
		elo INTEGER DEFAULT 1500 NOT NULL,
		playerusername VARCHAR(50) DEFAULT "" NOT NULL,
		lastrankupdate INT NOT NULL DEFAULT 0,
		PRIMARY KEY (discordID)
	);`
