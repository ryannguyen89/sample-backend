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

func TestAPIProductUpdate(t *testing.T) {
	pathAdd := "/api/item/add"
	pathUpdate := "/api/item/update"

	validReq := func() url.Values {
		data := url.Values{}
		data.Add("sku", "ABT-001")
		data.Add("name", "ABT-Sehat01")
		data.Add("qty", fmt.Sprintf("%v", 100))
		data.Add("price", fmt.Sprintf("%v", 100000))
		data.Add("unit", "Carton")
		data.Add("status", fmt.Sprintf("%v", 1))

		return data
	}

	t.Run("should return error when item not exist", func(t *testing.T) {
		t.Parallel()

		data := validReq()
		data.Set("sku", "XBT-001")
		api := makeAPI(t)
		w := postForm(t, api, pathUpdate, data, bearer)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("add and update", func(t *testing.T) {
		t.Parallel()

		data := validReq()
		api := makeAPI(t)
		w := postForm(t, api, pathAdd, data, bearer)
		assert.Equal(t, http.StatusCreated, w.Code)

		data.Set("qty", fmt.Sprintf("%v", 95))
		w = postForm(t, api, pathUpdate, data, bearer)
		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestAPIProductDelete(t *testing.T) {
	path := "/api/item/delete"
	pathAdd := "/api/item/add"

	validReq := func() url.Values {
		data := url.Values{}
		data.Add("sku", "DBT-001")
		data.Add("name", "DBT-Sehat01")
		data.Add("qty", fmt.Sprintf("%v", 100))
		data.Add("price", fmt.Sprintf("%v", 100000))
		data.Add("unit", "Carton")
		data.Add("status", fmt.Sprintf("%v", 1))

		return data
	}

	t.Run("should return error when item not exist", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Set("sku", "XBT-001")
		api := makeAPI(t)
		w := postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("add and update", func(t *testing.T) {
		t.Parallel()

		data := validReq()
		api := makeAPI(t)
		w := postForm(t, api, pathAdd, data, bearer)
		assert.Equal(t, http.StatusCreated, w.Code)

		data = url.Values{}
		data.Set("sku", "DBT-001")
		w = postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAPIProductList(t *testing.T) {
	path := "/api/items"
	pathAdd := "/api/item/add"

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

	validReq := func() url.Values {
		data := url.Values{}
		data.Add("sku", "DBT-001")
		data.Add("name", "DBT-Sehat01")
		data.Add("qty", fmt.Sprintf("%v", 100))
		data.Add("price", fmt.Sprintf("%v", 100000))
		data.Add("unit", "Carton")
		data.Add("status", fmt.Sprintf("%v", 1))

		return data
	}

	data := validReq()
	api := makeAPI(t)
	w := postForm(t, api, pathAdd, data, bearer)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = get(t, api, path, bearer)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := response{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Data), 1)
}

func TestAPIProductSearch(t *testing.T) {
	path := "/api/item/search"
	pathAdd := "/api/item/add"

	validReq := func() url.Values {
		data := url.Values{}
		data.Add("sku", "SBT-001")
		data.Add("name", "SBT-Sehat01")
		data.Add("qty", fmt.Sprintf("%v", 100))
		data.Add("price", fmt.Sprintf("%v", 100000))
		data.Add("unit", "Carton")
		data.Add("status", fmt.Sprintf("%v", 1))

		return data
	}

	t.Run("search not exist item should return error", func(t *testing.T) {
		t.Parallel()

		data := url.Values{}
		data.Set("sku", "1BT-001")

		api := makeAPI(t)
		w := postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("add and search", func(t *testing.T) {
		t.Parallel()

		data := validReq()
		api := makeAPI(t)
		w := postForm(t, api, pathAdd, data, bearer)
		assert.Equal(t, http.StatusCreated, w.Code)

		data = url.Values{}
		data.Set("sku", "SBT-001")
		w = postForm(t, api, path, data, bearer)
		assert.Equal(t, http.StatusOK, w.Code)
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
