package service

import (
	"context"
	"fmt"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockUser := fixture.GetMockUser()
		mockUser.ID = uid

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", uid).Return(mockUser, nil)

		u, err := us.Get(uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUser)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", uid).Return(nil, fmt.Errorf("some error down the call chain"))

		u, err := us.Get(uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByEmail", mockUser.Email).Return(mockUser, nil)

		u, err := us.GetByEmail(mockUser.Email)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUser)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		email := fixture.Email()

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByEmail", email).Return(nil, fmt.Errorf("some error down the call chain"))

		u, err := us.GetByEmail(email)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestRegister(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockUser := fixture.GetMockUser()

		initial := &model.User{
			Username: mockUser.Username,
			Email:    mockUser.Email,
			Password: mockUser.Password,
		}

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", initial).
			Run(func(args mock.Arguments) {
				mockUser.ID = uid
			}).Return(mockUser, nil)

		user, err := us.Register(initial)

		assert.NoError(t, err)

		// assert user now has a userID
		assert.Equal(t, uid, mockUser.ID)
		assert.Equal(t, user, mockUser)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "bob@bob.com",
			Username: "bobby",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", "bob@bob.com")

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mockUser).
			Return(nil, mockErr)

		user, err := us.Register(mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())
		assert.Nil(t, user)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	validPW := "howdyhoneighbor!"
	hashedValidPW, _ := hashPassword(validPW)
	invalidPW := "howdyhodufus!"

	mockUserRepository := new(mocks.UserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		mockUser.Password = hashedValidPW

		mockUserRepository.
			On("FindByEmail", mockUser.Email).Return(mockUser, nil)

		user, err := us.Login(mockUser.Email, validPW)

		assert.NoError(t, err)
		assert.Equal(t, user, mockUser)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockUser.Email)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUserResp := fixture.GetMockUser()
		mockUserResp.ID = uid
		mockUserResp.Password = hashedValidPW

		mockArgs := mock.Arguments{
			mockUserResp.Email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		user, err := us.Login(mockUserResp.Email, invalidPW)

		assert.Error(t, err)
		assert.EqualError(t, err, apperrors.InvalidCredentials)
		assert.Nil(t, user)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}

func TestUpdateDetails(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := fixture.GetMockUser()
		mockUser.ID = uid

		mockArgs := mock.Arguments{
			mockUser,
		}

		mockUserRepository.
			On("Update", mockArgs...).Return(nil)

		err := us.UpdateAccount(mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})

	t.Run("Failure", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := fixture.GetMockUser()
		mockUser.ID = uid

		mockArgs := mock.Arguments{
			mockUser,
		}

		mockError := apperrors.NewInternal()

		mockUserRepository.
			On("Update", mockArgs...).Return(mockError)

		err := us.UpdateAccount(mockUser)
		assert.Error(t, err)

		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.Internal, apperror.Type)

		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})
}

func TestUserService_ChangeAvatar(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockFileRepository := new(mocks.FileRepository)

	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
		FileRepository: mockFileRepository,
	})

	t.Run("Successful new image", func(t *testing.T) {
		uid, _ := GenerateId()

		// does not have have imageURL
		mockUser := fixture.GetMockUser()
		mockUser.ID = uid
		mockUser.Image = ""

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "test_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		imageURL := "https://imageurl.com/jdfkj34kljl"

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockUpdatedUser := &model.User{
			BaseModel: model.BaseModel{
				ID:        mockUser.ID,
				CreatedAt: mockUser.CreatedAt,
				UpdatedAt: mockUser.UpdatedAt,
			},
			Email:    mockUser.Email,
			Username: mockUser.Username,
			Image:    imageURL,
			Password: mockUser.Password,
		}

		mockUserRepository.
			On("Update", updateArgs...).
			Return(nil)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		mockUser.Image = url

		err = us.UpdateAccount(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, mockUser)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})

	t.Run("Successful update image", func(t *testing.T) {
		imageURL := "https://imageurl.com/jdfkj34kljl"
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
			Image:    imageURL,
		}
		mockUser.ID = uid

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "test_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		deleteImageArgs := mock.Arguments{
			imageURL,
		}

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		mockFileRepository.
			On("DeleteImage", deleteImageArgs...).
			Return(nil)

		mockUpdatedUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
			Image:    imageURL,
		}
		mockUpdatedUser.ID = uid

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockUserRepository.
			On("Update", updateArgs...).
			Return(nil)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		err = us.DeleteImage(mockUser.Image)
		assert.NoError(t, err)

		mockUser.Image = url
		err = us.UpdateAccount(mockUser)
		assert.NoError(t, err)

		assert.Equal(t, mockUpdatedUser, mockUser)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockFileRepository.AssertCalled(t, "DeleteImage", imageURL)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})

	t.Run("FileRepository Error", func(t *testing.T) {
		// need to create a new UserService and repository
		// because testify has no way to overwrite a mock's
		// "On" call.
		mockUserRepository := new(mocks.UserRepository)
		mockFileRepository := new(mocks.FileRepository)

		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
			FileRepository: mockFileRepository,
		})

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "file_directory"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		mockError := apperrors.NewInternal()
		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return("", mockError)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.Equal(t, "", url)
		assert.Error(t, err)

		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("UserRepository UpdateImage Error", func(t *testing.T) {
		uid, _ := GenerateId()
		imageURL := "https://imageurl.com/jdfkj34kljl"

		// has imageURL
		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "A New Bob!",
			Image:    imageURL,
		}
		mockUser.ID = uid

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "file_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockError := apperrors.NewInternal()
		mockUserRepository.
			On("Update", updateArgs...).
			Return(mockError)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		assert.Equal(t, imageURL, url)

		err = us.UpdateAccount(mockUser)

		assert.Error(t, err)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})
}

func TestUserService_ChangePassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		currentPassword := mockUser.Password

		hashedPassword, err := hashPassword(currentPassword)
		assert.NoError(t, err)
		mockUser.Password = hashedPassword
		newPassword := fixture.RandStringRunes(10)

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{UserRepository: mockUserRepository})

		mockUserRepository.On("Update", mockUser).Return(nil)

		err = us.ChangePassword(currentPassword, newPassword, mockUser)
		assert.NoError(t, err)

		assert.NotEqual(t, mockUser.Password, newPassword)
		assert.NotEqual(t, mockUser.Password, hashedPassword)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error verifying password", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		currentPassword := mockUser.Password
		newPassword := fixture.RandStringRunes(10)

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{UserRepository: mockUserRepository})

		err := us.ChangePassword(currentPassword, newPassword, mockUser)
		assert.Error(t, err)

		assert.Equal(t, err, apperrors.NewInternal())

		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("Current Password is incorrect", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		currentPassword := fixture.RandStringRunes(10)

		hashedPassword, err := hashPassword(mockUser.Password)
		assert.NoError(t, err)
		mockUser.Password = hashedPassword
		newPassword := fixture.RandStringRunes(10)

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{UserRepository: mockUserRepository})

		err = us.ChangePassword(currentPassword, newPassword, mockUser)
		assert.Error(t, err)

		assert.Equal(t, err, apperrors.NewAuthorization(apperrors.InvalidOldPassword))

		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("Error returned from the repository", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		currentPassword := mockUser.Password

		hashedPassword, err := hashPassword(currentPassword)
		assert.NoError(t, err)
		mockUser.Password = hashedPassword
		newPassword := fixture.RandStringRunes(10)

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{UserRepository: mockUserRepository})

		mockError := apperrors.NewInternal()
		mockUserRepository.On("Update", mockUser).Return(mockError)

		err = us.ChangePassword(currentPassword, newPassword, mockUser)
		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestUserService_ForgotPassword(t *testing.T) {
	mockUser := fixture.GetMockUser()
	token := fixture.RandStringRunes(10)

	t.Run("Success", func(t *testing.T) {
		mockRedisRepository := new(mocks.RedisRepository)
		mockMailRepository := new(mocks.MailRepository)

		us := NewUserService(&USConfig{
			RedisRepository: mockRedisRepository,
			MailRepository:  mockMailRepository,
		})

		mockRedisRepository.On("SetResetToken", mock.Anything, mockUser.ID).Return(token, nil)
		mockMailRepository.On("SendResetMail", mockUser.Email, token).Return(nil)

		err := us.ForgotPassword(context.TODO(), mockUser)
		assert.NoError(t, err)

		mockRedisRepository.AssertExpectations(t)
		mockMailRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRedisRepository := new(mocks.RedisRepository)
		mockMailRepository := new(mocks.MailRepository)

		us := NewUserService(&USConfig{
			RedisRepository: mockRedisRepository,
			MailRepository:  mockMailRepository,
		})

		mockError := apperrors.NewInternal()
		mockRedisRepository.On("SetResetToken", mock.Anything, mockUser.ID).Return("", mockError)

		err := us.ForgotPassword(context.TODO(), mockUser)
		assert.Error(t, err)

		mockRedisRepository.AssertExpectations(t)
		mockMailRepository.AssertNotCalled(t, "SendResetMail")
	})
}

func TestUserService_ResetPassword(t *testing.T) {
	mockUser := fixture.GetMockUser()
	password := fixture.RandStr(10)
	token := fixture.RandStr(10)

	t.Run("Success", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		mockRedisRepository := new(mocks.RedisRepository)

		us := NewUserService(&USConfig{
			UserRepository:  mockUserRepository,
			RedisRepository: mockRedisRepository,
		})

		mockRedisRepository.On("GetIdFromToken", mock.Anything, token).Return(mockUser.ID, nil)
		mockUserRepository.On("FindByID", mockUser.ID).Return(mockUser, nil)
		mockUserRepository.On("Update", mockUser).Return(nil)

		user, err := us.ResetPassword(context.TODO(), password, token)
		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockUserRepository.AssertExpectations(t)
		mockRedisRepository.AssertExpectations(t)
	})

	t.Run("No id found", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		mockRedisRepository := new(mocks.RedisRepository)

		us := NewUserService(&USConfig{
			UserRepository:  mockUserRepository,
			RedisRepository: mockRedisRepository,
		})

		mockError := apperrors.NewInternal()
		mockRedisRepository.On("GetIdFromToken", mock.Anything, token).Return("", mockError)

		user, err := us.ResetPassword(context.TODO(), password, token)
		assert.Error(t, err)
		assert.Nil(t, user)

		mockRedisRepository.AssertCalled(t, "GetIdFromToken", mock.Anything, token)
		mockUserRepository.AssertNotCalled(t, "FindByID")
		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("No user found", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		mockRedisRepository := new(mocks.RedisRepository)
		id := fixture.RandID()

		us := NewUserService(&USConfig{
			UserRepository:  mockUserRepository,
			RedisRepository: mockRedisRepository,
		})

		mockError := apperrors.NewInternal()
		mockRedisRepository.On("GetIdFromToken", mock.Anything, token).Return(id, nil)
		mockUserRepository.On("FindByID", id).Return(nil, mockError)

		user, err := us.ResetPassword(context.TODO(), password, token)
		assert.Error(t, err)
		assert.Nil(t, user)

		mockRedisRepository.AssertCalled(t, "GetIdFromToken", mock.Anything, token)
		mockUserRepository.AssertCalled(t, "FindByID", id)
		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("Error returned from the repository", func(t *testing.T) {
		mockUserRepository := new(mocks.UserRepository)
		mockRedisRepository := new(mocks.RedisRepository)
		mockError := apperrors.NewInternal()

		us := NewUserService(&USConfig{
			UserRepository:  mockUserRepository,
			RedisRepository: mockRedisRepository,
		})

		mockRedisRepository.On("GetIdFromToken", mock.Anything, token).Return(mockUser.ID, nil)
		mockUserRepository.On("FindByID", mockUser.ID).Return(mockUser, nil)
		mockUserRepository.On("Update", mockUser).Return(mockError)

		user, err := us.ResetPassword(context.TODO(), password, token)
		assert.Error(t, err)
		assert.Nil(t, user)

		mockUserRepository.AssertExpectations(t)
		mockRedisRepository.AssertExpectations(t)
	})
}
