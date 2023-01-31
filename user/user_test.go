package user

import (
	"github.com/MicBun/go-100-coverage-docker-crud/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAuthUser(t *testing.T) {
	database.RunTest(func(db *gorm.DB) {
		a := AdminAuth(db)
		u, err := a.RegisterUser("foo@bar.com", "securePassword", "Foo Bar")
		assert.NoError(t, err)

		assert.Equal(t, "foo@bar.com", u.Username)
		assert.NotEqual(t, "securePassword", u.Password)
		assert.Equal(t, "Foo Bar", u.Name)

		_, err = a.AuthenticateUser("foo@bar.com", "notPassword")
		assert.Error(t, err)

		_, err = a.AuthenticateUser("bar@bar.com", "notPassword")
		assert.Error(t, err)

		u2, err := a.AuthenticateUser("foo@bar.com", "securePassword")
		assert.NoError(t, err)
		assert.Equal(t, "foo@bar.com", u2.Username)
		assert.Equal(t, u.ID, u2.ID)
		assert.Equal(t, u.Name, u2.Name)

		_, err = a.GetUser(u.ID)
		assert.NoError(t, err)

		_, err = a.GetUser(2)
		assert.Error(t, err)

		_, err = a.UpdateUser(u.ID, "bar@foo.com", "passwordSecure", "Bar Foo")
		assert.NoError(t, err)

		u3, err := a.GetUser(u.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Bar Foo", u3.Name)

		users, err := a.ListUsers()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))

		err = a.DeleteUser(u.ID)
		assert.NoError(t, err)

		err = a.DeleteUser(u.ID)
		assert.Error(t, err)

		_, err = a.GetUser(u.ID)
		assert.Error(t, err)

		err = a.SaveToken(u.ID, "token")
		assert.Error(t, err)

		_, err = a.UpdateUser(u.ID, "", "", "Bar Foo")
		assert.Error(t, err)

		users, err = a.ListUsers()
		assert.Error(t, err)
		assert.Equal(t, 0, len(users))

		u, err = a.RegisterUser("new@user.com", "securePassword", "New User")
		assert.NoError(t, err)

		err = a.SaveToken(u.ID, "token")
		assert.NoError(t, err)

		//
		//userAuth := UserAuth(db)
		//newUser, err := userAuth.GetUserByID(u.ID)
		//assert.NoError(t, err)
		//assert.Equal(t, u.ID, newUser.ID)
		//assert.Equal(t, u.Username, newUser.Username)
	})
}
