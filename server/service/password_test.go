package service

import (
	"github.com/sentrionic/valkyrie/model/fixture"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := fixture.RandStringRunes(10)

	hashedPassword1, err := hashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	valid, err := comparePasswords(hashedPassword1, password)
	require.NoError(t, err)
	require.True(t, valid)

	wrongPassword := fixture.RandStringRunes(10)
	valid, err = comparePasswords(hashedPassword1, wrongPassword)
	require.NoError(t, err)
	require.False(t, valid)

	hashedPassword2, err := hashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)

	valid, err = comparePasswords(password, hashedPassword1)
	require.Error(t, err)
	require.EqualError(t, err, "did not provide a valid hash")
	require.False(t, valid)
}
