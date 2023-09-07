package controller

import (
	"coba_01/pkg/response"
	"coba_01/src/app/service"
	"coba_01/src/middleware"
	"coba_01/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type UserController struct {
	router  *gin.RouterGroup
	service service.UserServiceInterface
}

func NewUserController(router *gin.RouterGroup) *UserController {
	uc := &UserController{
		router:  router,
		service: service.NewServiceImpl(),
	}

	userRouter := uc.router.Group("/users")
	userRouter.POST("/signup", uc.SignUp)
	userRouter.POST("/login", uc.Login)

	authLoginGrup := userRouter.Group("/")
	authLoginGrup.Use(middleware.Authenticate)
	authLoginGrup.POST("/create", uc.CreateUser)
	authLoginGrup.GET("/all", uc.GetAllUser)
	authLoginGrup.DELETE("/delete/:id", uc.DeleteUserById)
	authLoginGrup.PATCH("/update/:id", uc.UpdateUser)
	authLoginGrup.GET("/search/:id", uc.GetUserByID)
	authLoginGrup.GET("/search", uc.GetUserByEmail)
	return uc
}

// ================>[ Auth ]<================\\

// @Summary Signup a user
// @Tags User Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Signup user"
// @Success 200 {object} map[string]string "OK"
// @Router /users/signup [post]
func (uc UserController) SignUp(c *gin.Context) {
	var user models.User
	var responSignUp response.SignupResponse

	if err := c.ShouldBind(&responSignUp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Username = responSignUp.Username
	user.Email = responSignUp.Email
	user.Password = responSignUp.Password

	token, err := uc.service.Signup(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Login a user
// @Tags User Auth
// @Accept json
// @Produce json
// @Param details body response.LoginResponse true "Login details"
// @Success 200 {object} map[string]string "OK"
// @Router /users/login [post]
func (uc UserController) Login(c *gin.Context) {
	var user models.User
	var userResponse response.LoginResponse

	if err := c.ShouldBind(&userResponse); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	user.Email = userResponse.Email
	user.Password = userResponse.Password

	token, err := uc.service.Login(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ================>[END Auth ]<================\\

// ================>[ CRUD ]<================\\

// @Summary Create a user
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 200 {object} map[string]string "OK"
// @Router /users/create [post]
func (uc UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err := uc.service.Create(&user)
	if err != nil {
		if err.Error() == "email sudah ada" {
			c.JSON(http.StatusConflict, gin.H{"message": "Email Sudah Ada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data Berhasil Dibuat",
		"data":    user,
	})
}

// @Summary Get all users
// @Security ApiKeyAuth
// @Tags User
// @Produce json
// @Success 200 {object} []models.User "OK"
// @Router /users/all [get]
func (uc UserController) GetAllUser(c *gin.Context) {
	username := c.MustGet("username")
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("perpage", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format Pag tidak valdi"})
		return
	}
	perpage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format Pag tidak valdi"})
		return
	}

	users, err := uc.service.GetAll(page, perpage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if users == nil || len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tidak Ada Data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "name": username})
}

// @Summary Delete a user
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "OK"
// @Router /users/delete/{id} [delete]
func (uc UserController) DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	err := uc.service.Delete(id)
	if err != nil {
		if err.Error() == "id tidak ada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil di hapus"})
}

// @Summary Update a user
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {object} map[string]string "OK"
// @Router /users/update/{id} [patch]
func (uc UserController) UpdateUser(c *gin.Context) {
	var user models.User

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data Tidak Ada"})
		return
	}
	user.ID = objID

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.service.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil Update", "data": user})
}

// ================>[ END CRUD ]<================\\

// ================>[ Pencarian ID - Email ]<================\\

// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Param id path string true "ID"
// @Produce json
// @Success 200 {object} []models.User "OK"
// @Router /users/search/{id} [get]
func (uc UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.service.FindByID(id)
	if err != nil {
		if err.Error() == "user tidak ada" {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("ID = %s ", id, " Tidak Ditemukan")})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary Find users by partial email
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param email query string true "Partial Email"
// @Success 200 {object} []models.User "OK"
// @Router /users/search [get]
func (uc UserController) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Email Harus Di Isi"})
		return
	}
	users, err := uc.service.FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Data Dengan Email : %s %s ", email, " Tidak Ditemukan")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}
