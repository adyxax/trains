package database

import (
	"git.adyxax.org/adyxax/trains/pkg/model"
)

func (env *DBEnv) CountStops() (i int, err error) {
	query := `SELECT count(*) from stops;`
	err = env.db.QueryRow(query).Scan(&i)
	if err != nil {
		return 0, newQueryError("Could not run database query: most likely the schema is corrupted", err)
	}
	return
}

func (env *DBEnv) GetStop(id string) (*model.Stop, error) {
	query := `SELECT name FROM stops WHERE id = $1;`
	stop := model.Stop{Id: id}
	err := env.db.QueryRow(
		query,
		id,
	).Scan(
		&stop.Name,
	)
	if err != nil {
		return nil, newQueryError("Could not run database query", err)
	}
	return &stop, nil
}

func (env *DBEnv) GetStops() (stops []model.Stop, err error) {
	query := `SELECT id, name FROM stops;`
	rows, err := env.db.Query(query)
	if err != nil {
		return nil, newQueryError("Could not run database query", err)
	}
	defer rows.Close()
	for rows.Next() {
		var stop model.Stop
		if err := rows.Scan(&stop.Id, &stop.Name); err != nil {
			return nil, newQueryError("Could not run database query", err)
		}
		stops = append(stops, stop)
	}
	return
}

func (env *DBEnv) ReplaceAndImportStops(stops []model.Stop) error {
	pre_query := `DELETE FROM stops;`
	query := `
		INSERT INTO stops
			(id, name)
		VALUES
			($1, $2);`
	tx, err := env.db.Begin()
	if err != nil {
		return newTransactionError("Could not Begin()", err)
	}
	_, err = tx.Exec(pre_query)
	if err != nil {
		tx.Rollback()
		return newQueryError("Could not run database query: most likely the schema is corrupted", err)
	}
	for i := 0; i < len(stops); i++ {
		_, err = tx.Exec(
			query,
			stops[i].Id,
			stops[i].Name,
		)
		if err != nil {
			tx.Rollback()
			return newQueryError("Could not run database query", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return newTransactionError("Could not commit transaction", err)
	}
	return nil
}
