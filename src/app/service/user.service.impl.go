package service

import (
	"coba_01/database/umongo"
	"coba_01/pkg/response"
	"coba_01/src/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserServiceImpl struct {
	collectionName string
	ctx            context.Context
	dbutil         *umongo.MongoDBUtil
}

func NewServiceImpl() UserServiceInterface {
	ctx := context.Background()
	dbutil := umongo.NewMongoDBUtil("User")

	return &UserServiceImpl{
		collectionName: "User",
		ctx:            ctx,
		dbutil:         dbutil,
	}
}

// ================>[ CURD ]<================\\
func (us UserServiceImpl) Signup(user *models.User) (string, error) {
	isEmailExist, err := us.dbutil.IsEmailExist(user.Email)
	if err != nil {
		return "", err
	}
	if isEmailExist {
		return "", errors.New("Email Sudah Ada")
	}
	hashPw, _ := response.HasPasswordBcrypt(user.Password)
	user.Password = hashPw
	user.ID = primitive.NewObjectID()

	err = us.dbutil.Insert(user)
	if err != nil {
		return "", err
	}
	token, err := response.CreateToken(user, 24*time.Hour)

	return token, err
}

func (us UserServiceImpl) Login(user *models.User) (string, error) {
	foundUser, err := us.dbutil.GetEmailLogin(user.Email)
	if err != nil {
		return "", errors.New("Login Tidak Diketahui")
	}

	if !response.CheckPasswordFromHash(user.Password, foundUser.Password) {
		return "", errors.New("Login Tidak Diketahui")
	}

	token, err := response.GenerateToken(foundUser.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us UserServiceImpl) Create(user *models.User) error {
	if user.Role == "" {
		user.Role = "user"
	}
	exists, err := us.dbutil.IsEmailExist(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email sudah ada")
	}

	// hash pasword pake bcrpty
	hashpw, err := response.HasPasswordBcrypt(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashpw
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// simpan ke database
	return us.dbutil.Insert(user)
}

func (us UserServiceImpl) GetAll(page int, perPage int) ([]models.User, error) {
	skip := (page - 1) * perPage
	limit := perPage

	users, err := us.dbutil.GetAllWithPaging(skip, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us UserServiceImpl) Update(user *models.User) error {
	emailUsed, err := us.dbutil.IsEmailUsedByAnotherUser(user.Email, user.ID)

	if user.Role == "" {
		user.Role = "user"
	}

	if err != nil {
		return err
	}
	if emailUsed {
		return errors.New("Email Sudah Digunakan")
	}

	if user.Password != "" {
		hashpw, err := response.HasPasswordBcrypt(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashpw
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"address":    user.Address,
			"role":       user.Role, // Tambahkan ini
			"updated_at": time.Now(),
		},
	}

	err = us.dbutil.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (us UserServiceImpl) Delete(id string) error {
	deleteCount, err := us.dbutil.Delete(id)
	if err != nil {
		return err
	}
	if deleteCount == 0 {
		return errors.New("User Tidak Ditemukan")
	}
	return nil
}

// ================>[ END CURD ]<================\\

// ================>[ Pencarian ID - Email ]<================\\
func (us UserServiceImpl) FindByID(id string) (*models.User, error) {
	return us.dbutil.GetUserByID(id)
}
func (us UserServiceImpl) FindByEmail(email string) ([]models.User, error) {
	return us.dbutil.GetUserByEmail(email)

}
