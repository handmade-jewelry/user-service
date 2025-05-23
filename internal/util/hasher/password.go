package hasher

import "golang.org/x/crypto/bcrypt"

const cost = bcrypt.DefaultCost

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(password, usersPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(usersPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
