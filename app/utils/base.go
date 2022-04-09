package utils

import (
	"testing"

	"com.pikpo.events/app/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Inital test
func BaseTest(t *testing.T) {
	if database.Database.DB == nil || database.Database.DB.Stats().OpenConnections == 1 {
		godotenv.Load("../../.env")
		database.Database.Open()
		_, err := database.Database.DB.Exec("DELETE FROM events")
		assert.Equal(t, nil, err)
	}
}
