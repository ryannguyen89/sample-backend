package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sampleBackend/internal/user"
)

type API struct {
	userSvc *user.Service
}

func NewAPI(userSvc *user.Service) *API {
	return &API{
		userSvc: userSvc,
	}
}

func (api *API) Route(route gin.IRouter) {
	g := route.Group("/api")
	g.POST("/register", api.handleUserRegister())
	g.POST("/auth/login", api.handleUserLogin())
}

func (api *API) handleUserRegister() gin.HandlerFunc {
	type (
		request struct {
			Email    string `form:"email" binding:"required""`
			Password string `form:"password" binding:"required"`
		}
	)

	return func(c *gin.Context) {
		var (
			r   request
			ctx = c.Request.Context()
		)

		err := c.ShouldBind(&r)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, NewError(fmt.Sprintf("parse request: %v", err)))
			return
		}
		fmt.Printf("register user: %v\n", r.Email)

		// Do create user
		err = api.userSvc.CreateUser(ctx, user.User{
			Email:    r.Email,
			Password: r.Password,
		})
		if err != nil {
			_ = c.Error(err)
			if user.IsErrUserExist(err) {
				c.JSON(http.StatusBadRequest, NewError(fmt.Sprintf("user already exist")))
				return
			}

			c.JSON(http.StatusInternalServerError, NewError(fmt.Sprintf("%v", err)))
			return
		}

		c.Status(http.StatusCreated)
	}
}

func (api *API) handleUserLogin() gin.HandlerFunc {
	type (
		request struct {
			Email    string `form:"email" binding:"required""`
			Password string `form:"password" binding:"required"`
		}
		response struct {
			Token string `json:"token"`
		}
	)

	return func(c *gin.Context) {
		var (
			r   request
			ctx = c.Request.Context()
		)

		err := c.ShouldBind(&r)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, NewError(fmt.Sprintf("parse request: %v", err)))
			return
		}
		fmt.Printf("user login: %v\n", r.Email)

		// Do log in
		resp, err := api.userSvc.Login(ctx, user.User{
			Email:    r.Email,
			Password: r.Password,
		})
		if err != nil {
			_ = c.Error(err)
			if user.IsErrUserInvalid(err) {
				c.JSON(http.StatusBadRequest, NewError(fmt.Sprintf("user invalid")))
				return
			}

			c.JSON(http.StatusInternalServerError, NewError(fmt.Sprintf("%v", err)))
			return
		}

		c.JSON(http.StatusOK, response{Token: resp.Token})
	}
}
