package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ccr1m3/osmc/main/static"
)

type Player struct {
	DiscordID      string `db:"discordID"`
	PlayerUsername string `db:"playerusername"`
	Elo            int    `db:"elo"`
	LastRankUpdate int    `db:"lastrankupdate"`
}

type Players []*Player

func CreatePlayerWithID(ctx context.Context, discordID string) (*Player, error) {
	_, err := GetPlayerByID(ctx, discordID)
	if err != nil && !errors.Is(err, static.ErrNotFound) {
		return nil, err
	} else if err == nil {
		return nil, static.ErrAlreadyExists
	}
	_, err = db.Exec("INSERT INTO players (discordID) VALUES (?)", discordID)
	if err != nil {
		return nil, static.ErrDB(err)
	}
	return &Player{DiscordID: discordID, Elo: 0}, nil
}

func GetPlayerByID(ctx context.Context, discordID string) (*Player, error) {
	var player Player
	err := db.Get(&player, "SELECT * FROM players WHERE discordID=?", discordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	return &player, nil
}

func GetPlayerByUsername(ctx context.Context, username string) (*Player, error) {
	var player Player
	err := db.Get(&player, "SELECT * FROM players WHERE playerusername=?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, static.ErrNotFound
		}
		return nil, static.ErrDB(err)
	}
	return &player, nil
}

func GetOrCreatePlayerByID(ctx context.Context, discordID string) (*Player, error) {
	p, err := GetPlayerByID(ctx, discordID)
	if err != nil && errors.Is(err, static.ErrNotFound) {
		return CreatePlayerWithID(ctx, discordID)
	} else if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Player) SetPlayerUsername(ctx context.Context, Username string) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET playerusername=? WHERE discordID=?", Username, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.PlayerUsername = Username
	return nil
}

func (p *Player) SetElo(ctx context.Context, elo int) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE players SET elo=? WHERE discordID=?", elo, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.Elo = elo
	return nil
}

func (p *Player) SetLastUpdate(ctx context.Context) error {
	_, err := GetPlayerByID(ctx, p.DiscordID)
	if err != nil {
		return err
	}
	lastRankUpdate := int(time.Now().Unix())
	_, err = db.Exec("UPDATE players SET lastrankupdate=? WHERE discordID=?", lastRankUpdate, p.DiscordID)
	if err != nil {
		return static.ErrDB(err)
	}
	p.LastRankUpdate = lastRankUpdate
	return nil
}
