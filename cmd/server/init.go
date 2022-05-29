package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sampleBackend/internal/api"
	"sampleBackend/internal/product"
	"sampleBackend/internal/storage/memory"
	"sampleBackend/internal/user"
)

func (s *Server) init() {
	s.once.Do(func() {
		// Init API server
		userStorage := memory.NewUserStorage()
		userSvc := user.NewService(userStorage)

		prdStorage := memory.NewProductStorage()
		prdSvc := product.NewService(prdStorage)
		a := api.NewAPI(userSvc, prdSvc)

		gin.SetMode(gin.ReleaseMode)

		e := gin.New()
		e.Use(func(c *gin.Context) {
			c.Next()
			if len(c.Errors) > 0 {
				fmt.Println(c.Errors.String())
			}
		})
		e.Use(gin.Recovery())

		a.Route(e)

		addr := ":8080"
		s.http = &http.Server{
			Addr:    addr,
			Handler: e,
		}
	})
}
