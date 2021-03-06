package database

import "golang.org/x/crypto/bcrypt"

// To allow for testing the error case (bad random is hard to trigger)
var passwordFunction = bcrypt.GenerateFromPassword

func hashPassword(password string) (string, error) {
	bytes, err := passwordFunction([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", newPasswordError(err)
	}
	return string(bytes), nil
}

func checkPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return newPasswordError(err)
	}
	return nil
}
