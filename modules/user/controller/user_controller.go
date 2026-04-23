package controller

import (
	"net/http"

	"github.com/Mobilizes/materi-be-alpro/modules/user/service"
	"github.com/Mobilizes/materi-be-alpro/modules/user/validation"
	"github.com/Mobilizes/materi-be-alpro/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
    service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
    return &UserController{service: service}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    req, err := validation.ValidateCreateUser(c)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    user, err := ctrl.service.CreateUser(req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "User berhasil dibuat", user)
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
    id := c.Param("id")

    user, err := ctrl.service.GetUserByID(id)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "User tidak ditemukan")
        // utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Berhasil ambil user", user)
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
    users, err := ctrl.service.GetAllUsers()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data user")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Berhasil ambil semua user", users)
}
