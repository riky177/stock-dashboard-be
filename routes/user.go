package routes

import (
	"net/http"
	"stock-dashboard/models"
	"stock-dashboard/utils"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)

	if err != nil {
		response := models.NewErrorResponse("Invalid request format")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		response := models.NewErrorResponse("Invalid email or password")
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	token, err := utils.GenerateToken(user.Email, user.ID, user.Role)

	if err != nil {
		response := models.NewErrorResponse("Could not generate token")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"accessToken": token,
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
	}
	response := models.NewSuccessResponse(data, "Login successful")
	c.JSON(http.StatusOK, response)

}

func RegisterHandler(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		response := models.NewErrorResponse("Could not parse request data")
		context.JSON(http.StatusBadRequest, response)
		return
	}

	err = user.Save()
	if err != nil {
		response := models.NewErrorResponse("Could not save user")
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	}
	response := models.NewSuccessResponse(data, "User created successfully")
	context.JSON(http.StatusCreated, response)
}

func GetAllStaff(c *gin.Context) {
	data, err := models.GetAllStaff()
	if err != nil {
		response := models.NewErrorResponse("Failed to fetch staff")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := models.NewSuccessResponse(data, "Staff fetched successfully")
	c.JSON(http.StatusOK, response)
}

func DeleteStaff(c *gin.Context) {
	staffID := c.Param("id")

	if staffID == "" {
		response := models.NewErrorResponse("Staff ID is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user := models.User{ID: staffID}
	err := user.Delete()

	if err != nil {
		if err.Error() == "user not found" {
			response := models.NewErrorResponse("Staff member not found")
			c.JSON(http.StatusNotFound, response)
			return
		}
		response := models.NewErrorResponse("Failed to delete staff member")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.NewSuccessResponse(nil, "Staff member deleted successfully")
	c.JSON(http.StatusOK, response)
}
