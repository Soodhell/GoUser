package jwtWork

import (
	"User/interal/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtConfig(t *testing.T) {
	var user entity.User
	user.Name = "test"
	user.Email = "test@test.com"
	user.ID = 1
	user.NameImage = "test"
	user.PathImage = "test"
	user.Roles = 1

	token, err := CreateToken(user)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	verifyToken, err := VerifyToken(token)
	assert.Nil(t, err)
	assert.NotEmpty(t, verifyToken)

}
