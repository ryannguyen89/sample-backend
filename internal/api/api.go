package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sampleBackend/internal/product"
	"sampleBackend/internal/user"
)

type API struct {
	userSvc *user.Service
	prdSvc  *product.Service
}

func NewAPI(userSvc *user.Service, prdSvc *product.Service) *API {
	return &API{
		userSvc: userSvc,
		prdSvc:  prdSvc,
	}
}

func (api *API) Route(route gin.IRouter) {
	g := route.Group("/api")
	g.POST("/register", api.handleUserRegister())
	g.POST("/auth/login", api.handleUserLogin())

	prdGroup := g.Group("/item", api.authorizationMiddleware())
	prdGroup.POST("/add", api.handleProductAdd())
	prdGroup.POST("/update", api.handleProductUpdate())
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

func (api *API) handleProductAdd() gin.HandlerFunc {
	type (
		request struct {
			SKU      string `form:"sku" binding:"required"`
			Name     string `form:"name" binding:"required"`
			Quantity uint32 `form:"quantity"`
			Price    uint64 `form:"price" binding:"required"`
			Unit     string `form:"unit" binding:"required"`
			Status   uint8  `form:"status"`
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
		fmt.Printf("product add: %#v\n", r)

		err = api.prdSvc.AddProduct(ctx, product.Product{
			SKU:      r.SKU,
			Name:     r.Name,
			Quantity: r.Quantity,
			Price:    r.Price,
			Unit:     r.Unit,
			Status:   r.Status,
		})
		if err != nil {
			_ = c.Error(err)
			if product.IsErrExist(err) {
				c.JSON(http.StatusBadRequest, NewError(err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, NewError(err.Error()))
			return
		}

		c.Status(http.StatusCreated)
	}
}

func (api *API) handleProductUpdate() gin.HandlerFunc {
	type (
		request struct {
			SKU      string `form:"sku" binding:"required"`
			Name     string `form:"name" binding:"required"`
			Quantity uint32 `form:"quantity"`
			Price    uint64 `form:"price" binding:"required"`
			Unit     string `form:"unit" binding:"required"`
			Status   uint8  `form:"status"`
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
		fmt.Printf("product update: %#v\n", r)

		err = api.prdSvc.UpdateProduct(ctx, product.Product{
			SKU:      r.SKU,
			Name:     r.Name,
			Quantity: r.Quantity,
			Price:    r.Price,
			Unit:     r.Unit,
			Status:   r.Status,
		})
		if err != nil {
			_ = c.Error(err)
			if product.IsErrNotFound(err) {
				c.JSON(http.StatusBadRequest, NewError(err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, NewError(err.Error()))
			return
		}

		c.Status(http.StatusOK)
	}
}

func (api *API) handleProductList() gin.HandlerFunc {
	type (
		request struct {
		}
	)
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
