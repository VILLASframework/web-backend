package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Verify that you can connect to the database
func TestDBConnection(t *testing.T) {
	db := InitDB()
	defer db.Close()

	assert.NoError(t, VerifyConnection(db), "DB must ping")
}

// Verify that the associations between every model are done properly
func TestDummyDB(t *testing.T) {
	test_db := DummyInitDB()
	defer test_db.Close()

	DummyPopulateDB(test_db)

}
