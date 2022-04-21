package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModel_NewAccessToken(t *testing.T) {
	u, err := NewUser(
		"test@test.test",
		"test_password",
		false,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
	)
	if err != nil {
		t.Fatal(err)
	}

	token, err := NewAccessToken(u)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestModel_NewRefreshToken(t *testing.T) {
	u, err := NewUser(
		"test@test.test",
		"test_password",
		false,
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
	)
	if err != nil {
		t.Fatal(err)
	}

	token, err := NewRefreshToken(u)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
