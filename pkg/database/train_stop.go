package database

import (
	"git.adyxax.org/adyxax/trains/pkg/model"
)

func (env *DBEnv) CountTrainStops() (i int, err error) {
	query := `SELECT count(*) from train_stops;`
	err = env.db.QueryRow(query).Scan(&i)
	if err != nil {
		return 0, newQueryError("Could not run database query: most likely the schema is corrupted", err)
	}
	return
}

func (env *DBEnv) ReplaceAndImportTrainStops(trainStops []model.TrainStop) error {
	pre_query := `DELETE FROM train_stops;`
	query := `
		INSERT INTO train_stops
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
	for i := 0; i < len(trainStops); i++ {
		_, err = tx.Exec(
			query,
			trainStops[i].Id,
			trainStops[i].Name,
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
