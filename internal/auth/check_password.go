package auth

import "github.com/alexedwards/argon2id"

func CheckPasswordHash(password string, hash string) (bool, error) {
	isSame, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false, err
	}
	return isSame, nil
}
