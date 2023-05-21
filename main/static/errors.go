package static

import (
	"database/sql"
	"errors"
	"fmt"
)

// DB Errors
func ErrDB(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRowsNotFound
	}
	return fmt.Errorf("database error: %w", err)
}

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrRowsNotFound = errors.New("rows not found")

// Prometheus Errors
var ErrRankUpdateTooFast = errors.New("update too fast")
var ErrUserAlreadyLinked = errors.New("user already linked")
var ErrUsernameAlreadyLinked = errors.New("username already linked")
var ErrUsernameNotFound = errors.New("username not found")
var ErrUsernameNotOnGlobal = errors.New("user cannot be found on global leaderboard")
var ErrUserNotLinked = errors.New("user not linked")
