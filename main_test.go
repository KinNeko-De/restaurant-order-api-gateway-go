package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/order", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response = w.Result()

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Errorf("Result can not be closed: %e.", err)
		}
	}(response.Body)

	type CreateOrderResponse struct {
		Id uuid.UUID
	}

	actualResponse := readResponse[CreateOrderResponse](t, response)

	assert.NotNil(t, actualResponse.Id)
	assert.IsType(t, uuid.UUID{}, actualResponse.Id)
}

func TestGetOrder(t *testing.T) {
	type GetOrderResponse struct {
		Id      uuid.UUID
		Version int
	}
	var expectedVersion = 1
	var expectedId = uuid.New()

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/order/"+expectedId.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response = w.Result()

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Errorf("Result can not be closed: %v.", err)
		}
	}(response.Body)

	actualResponse := readResponse[GetOrderResponse](t, response)

	assert.NotNil(t, actualResponse.Id)
	assert.IsType(t, uuid.UUID{}, actualResponse.Id)

	assert.Equal(t, expectedVersion, actualResponse.Version)
}

func TestGetOrderBadFormatedOrderId(t *testing.T) {
	var expectedId = "869dacc9-1962"

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/order/"+expectedId, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Errorf("Result can not be closed: %v.", err)
		}
	}(response.Body)
}

func readResponse[K any](t *testing.T, response *http.Response) K {
	data := ReadAllBytes(t, response)
	actualResponse := decodeToJson[K](t, data)
	return actualResponse
}

func ReadAllBytes(t *testing.T, response *http.Response) []byte {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Unable to read response body %v", err)
	}
	return data
}

func decodeToJson[K any](t *testing.T, data []byte) K {
	var actualResponse K
	err := json.Unmarshal(
		data,
		&actualResponse,
	)
	if err != nil {
		str1 := string(data[:])
		t.Errorf("Response can not be read to expected response %v. Raw: %s", err, str1)
	}
	return actualResponse
}
