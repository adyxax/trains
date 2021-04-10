package database

import (
	"git.adyxax.org/adyxax/trains/pkg/model"
)

// Creates a new user in the database
// a QueryError is return if the username already exists (database constraints not met)
func (env *DBEnv) CreateUser(reg *model.UserRegistration) (*model.User, error) {
	hash, err := hashPassword(reg.Password)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO users
			(username, password, email)
		VALUES
			($1, $2, $3);`
	tx, err := env.db.Begin()
	if err != nil {
		return nil, newTransactionError("Could not Begin()", err)
	}
	result, err := tx.Exec(
		query,
		reg.Username,
		hash,
		reg.Email,
	)
	if err != nil {
		tx.Rollback()
		return nil, newQueryError("Could not run database query, most likely the username already exists", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, newTransactionError("Could not get LastInsertId, the database driver does not support this feature", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, newTransactionError("Could not commit transaction", err)
	}
	user := model.User{
		Id:       int(id),
		Username: reg.Username,
		Email:    reg.Email,
	}
	return &user, nil
}
