package api_test

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	. "sampleBackend/internal/api"
	"sampleBackend/internal/storage/memory"
	"sampleBackend/internal/user"
)

func TestAPIRegister(t *testing.T) {
	path := "/api/register"

	t.Run("nil body should return bad request", func(t *testing.T) {
		t.Parallel()

		api := makeAPI()
		w := postForm(t, api, path, nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("request lack email should return bad request", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("password", "123456")

		api := makeAPI()
		w := postForm(t, api, path, data)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("request lack password should return bad request", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", "abc@gmail.com")

		api := makeAPI()
		w := postForm(t, api, path, data)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when user exist", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", "abc@gmail.com")
		data.Add("password", "123456")

		api := makeAPI()
		w := postForm(t, api, path, data)
		assert.Equal(t, http.StatusCreated, w.Code)

		w = postForm(t, api, path, data)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user already exist")
	})
}

func makeAPI() http.Handler {
	userStorage := memory.NewUserStorage()
	userSvc := user.NewService(userStorage)
	api := NewAPI(userSvc)
	e := gin.New()
	e.Use(func(c *gin.Context) {
		log.Println(c.Errors.String())
	})

	api.Route(e)
	return e
}
