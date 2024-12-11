package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB()
	assert.NotNil(t, db)
	assert.Nil(t, err)
}
