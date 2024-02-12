package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"ytsruh.com/endtoend/models"
)

type UpdateProfileInput struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	ProfileImage  string `json:"profileImage"`
	ShowBooks     bool   `json:"showBooks"`
	ShowDocuments bool   `json:"showDocuments"`
}

func (api *API) GetProfile(c echo.Context) error {
	claims, err := api.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get user",
		})
	}
	user := models.User{}
	if err := api.DB.First(&user, "id = ?", claims.Id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to find user",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"profileImage":  user.ProfileImage,
		"showBooks":     user.ShowBooks,
		"showDocuments": user.ShowDocuments,
	})
}

func (api *API) PatchProfile(c echo.Context) error {
	input := new(UpdateProfileInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}
	claims, err := api.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get user",
		})
	}
	user := models.User{}
	tx := api.DB.Model(&user).Where("id = ?", claims.Id).Updates(models.User{
		Name:          input.Name,
		Email:         input.Email,
		ProfileImage:  &input.ProfileImage,
		ShowBooks:     input.ShowBooks,
		ShowDocuments: input.ShowDocuments,
	})
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to update user",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
