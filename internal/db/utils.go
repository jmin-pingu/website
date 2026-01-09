package db

import "github.com/google/uuid"

func Hash(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(value))
}
