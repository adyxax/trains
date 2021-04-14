package database

import (
	"git.adyxax.org/adyxax/trains/pkg/model"
	"github.com/google/uuid"
)

func (env *DBEnv) CreateSession(user *model.User) (*string, error) {
	token := uuid.NewString()

	query := `
		INSERT INTO sessions
			(token, user_id)
		VALUES
			($1, $2);`
	tx, err := env.db.Begin()
	if err != nil {
		return nil, newTransactionError("Could not Begin()", err)
	}
	_, err = tx.Exec(
		query,
		token,
		user.Id,
	)
	if err != nil {
		tx.Rollback()
		return nil, newQueryError("Could not run database query: most likely the token already exists in database, or the user id does not exist", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, newTransactionError("Could not commit transaction", err)
	}
	return &token, nil
}

func (env *DBEnv) ResumeSession(token string) (*model.User, error) {
	user := model.User{}
	query := `
		SELECT
			id, username, email
		FROM
			users
		INNER JOIN
			sessions ON users.id = sessions.user_id
		WHERE
			sessions.token = $1;`
	err := env.db.QueryRow(
		query,
		token,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return nil, newQueryError("Could not run database query, most likely the token is invalid", err)
	}
	return &user, nil
}
