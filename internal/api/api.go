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

	g.GET("/items", api.authorizationMiddleware(), api.handleProductList())
	prdGroup := g.Group("/item", api.authorizationMiddleware())
	prdGroup.POST("/add", api.handleProductAdd())
	prdGroup.POST("/update", api.handleProductUpdate())
	prdGroup.POST("/delete", api.handleProductDelete())
	prdGroup.POST("/search", api.handleProductSearch())
}

func (api *API) handleUserRegister() gin.HandlerFunc {
	type (
		request struct {
			Email    string `form:"email" binding:"required"`
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
			Email    string `form:"email" binding:"required"`
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
			Quantity uint32 `form:"qty"`
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
			Quantity uint32 `form:"qty"`
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
				c.Status(http.StatusNotFound)
				return
			}
			c.JSON(http.StatusInternalServerError, NewError(err.Error()))
			return
		}

		c.Status(http.StatusOK)
	}
}

func (api *API) handleProductDelete() gin.HandlerFunc {
	type (
		request struct {
			SKU string `form:"sku" binding:"required"`
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
		fmt.Printf("product delete: %#v\n", r)

		err = api.prdSvc.DeleteProduct(ctx, r.SKU)
		if err != nil {
			_ = c.Error(err)
			if product.IsErrNotFound(err) {
				c.Status(http.StatusNotFound)
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
		item struct {
			SKU      string `json:"sku"`
			Name     string `json:"name"`
			Quantity uint32 `json:"qty"`
			Price    uint64 `json:"price"`
			Unit     string `json:"unit"`
			Status   uint8  `json:"status"`
		}
		response struct {
			Data []*item `json:"data"`
		}
	)
	return func(c *gin.Context) {
		var (
			ctx  = c.Request.Context()
			data []*item
		)

		resp, err := api.prdSvc.ListProduct(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewError(err.Error()))
			return
		}

		for _, i := range resp {
			data = append(data, &item{
				SKU:      i.SKU,
				Name:     i.Name,
				Quantity: i.Quantity,
				Price:    i.Price,
				Unit:     i.Unit,
				Status:   i.Status,
			})
		}

		c.JSON(http.StatusOK, response{Data: data})
	}
}

func (api *API) handleProductSearch() gin.HandlerFunc {
	type (
		request struct {
			SKU string `form:"sku" binding:"required"`
		}
		item struct {
			SKU      string `json:"sku"`
			Name     string `json:"name"`
			Quantity uint32 `json:"qty"`
			Price    uint64 `json:"price"`
			Unit     string `json:"unit"`
			Status   uint8  `json:"status"`
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
		fmt.Printf("product search: %#v\n", r)

		prd, err := api.prdSvc.SearchProduct(ctx, r.SKU)
		if err != nil {
			_ = c.Error(err)
			if product.IsErrNotFound(err) {
				c.Status(http.StatusNotFound)
				return
			}
			c.JSON(http.StatusInternalServerError, NewError(err.Error()))
			return
		}

		c.JSON(http.StatusOK, item{
			SKU:      prd.SKU,
			Name:     prd.Name,
			Quantity: prd.Quantity,
			Price:    prd.Price,
			Unit:     prd.Unit,
			Status:   prd.Status,
		})
	}
}
