package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "sampleBackend/internal/api"
	"sampleBackend/internal/product"
	"sampleBackend/internal/storage/memory"
	"sampleBackend/internal/user"
)

const (
	registeredUser = "user@gmail.com"
	password       = "password"

	bearer = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsidXNlckBnbWFpbC5jb20iXSwiaWF0IjoxNjUzODAwMjE3fQ.ugHUydqAuhxAUCqkbLEXhKn531rqSGkT0MGd333aFRg"
)

func TestAPIUserRegister(t *testing.T) {
	path := "/api/register"

	t.Run("nil body should return bad request", func(t *testing.T) {
		t.Parallel()

		api := makeAPI(t)
		w := postForm(t, api, path, nil, "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("request lack email should return bad request", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("password", "123456")

		api := makeAPI(t)
		w := postForm(t, api, path, data, "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("request lack password should return bad request", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", "abc@gmail.com")

		api := makeAPI(t)
		w := postForm(t, api, path, data, "")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return error when user exist", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", "abc@gmail.com")
		data.Add("password", "123456")

		api := makeAPI(t)
		w := postForm(t, api, path, data, "")
		assert.Equal(t, http.StatusCreated, w.Code)

		w = postForm(t, api, path, data, "")
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user already exist")
	})
}

func TestAPIUserLogin(t *testing.T) {
	path := "/api/auth/login"

	type (
		response struct {
			Token string `json:"token"`
		}
	)

	t.Run("should return error when user invalid", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", registeredUser)
		data.Add("password", "123456")

		api := makeAPI(t)
		w := postForm(t, api, path, data, "")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("happy case", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Add("email", registeredUser)
		data.Add("password", password)

		api := makeAPI(t)
		w := postForm(t, api, path, data, "")
		assert.Equal(t, http.StatusOK, w.Code)

		resp := response{}
		err := json.Unmarshal([]byte(w.Body.String()), &resp)
		require.NoError(t, err)
		require.NotEqual(t, "", resp.Token)
	})
}

func TestAPIProductAdd(t *testing.T) {
	path := "/api/item/add"

	validReq := func() url.Values {
		data := url.Values{}
		data.Add("sku", "OBT-001")
		data.Add("name", "OBT-Sehat01")
		data.Add("qty", fmt.Sprintf("%v", 100))
		data.Add("price", fmt.Sprintf("%v", 100000))
		data.Add("unit", "Carton")
		data.Add("status", fmt.Sprintf("%v", 1))

		return data
	}

	t.Run("nil body should return bad request", func(t *testing.T) {
		t.Parallel()

		api := makeAPI(t)
		w := postForm(t, api, path, nil, bearer)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validate request", func(t *testing.T) {
		tests := map[string]struct {
			alterData func(data url.Values)
		}{
			"lack sku": {
				alterData: func(data url.Values) {
					data.Del("sku")
				},
			},
			"lack name": {
				alterData: func(data url.Values) {
					data.Del("name")
				},
			},
		}
		for name, test := range tests {
			test := test

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				data := validReq()
				test.alterData(data)

				api := makeAPI(t)
				w := postForm(t, api, path, nil, bearer)

				require.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})

	t.Run("create and duplicated", func(t *testing.T) {
		data := validReq()
		api := makeAPI(t)
		w := postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusCreated, w.Code)

		w = postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func makeAPI(t *testing.T) http.Handler {
	userStorage := memory.NewUserStorage()
	err := userStorage.Create(context.Background(), user.User{
		Email:    registeredUser,
		Password: password,
	})
	require.NoError(t, err)

	prdStorage := memory.NewProductStorage()

	userSvc := user.NewService(userStorage)
	prdSvc := product.NewService(prdStorage)
	api := NewAPI(userSvc, prdSvc)
	e := gin.New()
	e.Use(func(c *gin.Context) {
		log.Println(c.Errors.String())
	})

	api.Route(e)
	return e
}
