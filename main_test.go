package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/models"
	"theam.io/jdavidsanchez/test_crm_api/routes"
)

// Integration test (test the API and its connection with the database)

func TestMain(m *testing.M) {
	code := m.Run()

	clearCustomersTable()

	os.Exit(code)
}

func Test_Route_Customer_listAllCustomers(t *testing.T) {
	t.Run("Get an empty customer list", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != "[]" {
			t.Errorf("Expected an empty array. Got %s", body)
		}
	})
	t.Run("Get one customer", func(t *testing.T) {
		clearCustomersTable()
		// Add one customer
		newCustomer := models.Customer{
			Name:                 "Test_Name",
			Surname:              "Test_Surname",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		executeRequest(t, req)

		req, _ = http.NewRequest("GET", "/customers/all", nil)
		response := executeRequest(t, req)

		want := "[{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1}]"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("Get two customers", func(t *testing.T) {
		clearCustomersTable()
		// Add two customers
		newCustomer := models.Customer{
			Name:                 "Test_Name",
			Surname:              "Test_Surname",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		anotherCustomer := models.Customer{
			Name:                 "Test_Name_2",
			Surname:              "Test_Surname_2",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}

		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		executeRequest(t, req)
		data, _ = json.Marshal(anotherCustomer)
		req, _ = http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		executeRequest(t, req)

		req, _ = http.NewRequest("GET", "/customers/all", nil)
		response := executeRequest(t, req)

		want := "[{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1},{\"id\":2,\"name\":\"Test_Name_2\",\"surname\":\"Test_Surname_2\",\"pictureId\":1,\"lastModifiedByUserId\":1}]"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
}

func Test_Route_Customer_getCustomer(t *testing.T) {
	t.Run("Get a non existing customer", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/22", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		if m["error"] != "Customer not found" {
			t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
		}
	})
	t.Run("Get one customer", func(t *testing.T) {
		clearCustomersTable()
		// Add one customer
		newCustomer := models.Customer{
			Name:                 "Test_Name",
			Surname:              "Test_Surname",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		executeRequest(t, req)

		req, _ = http.NewRequest("GET", "/customers/1", nil)
		response := executeRequest(t, req)

		want := "{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1}"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("Non valid ID parameter", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/hola", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		got := response.Body.Bytes()
		want := []byte("404 page not found")
		if reflect.DeepEqual(got, want) {
			t.Errorf("Expected %s. Got '%s'", want, got)
		}
	})
}

func Test_Route_Customer_createCustomer(t *testing.T) {
	// TODO: not implemented
}

func Test_Route_Customer_updateCustomer(t *testing.T) {
	// TODO: not implemented
}

func Test_Route_Customer_deleteCustomer(t *testing.T) {
	// TODO: not implemented
}

func clearCustomersTable() {
	_, err := db.DB.Exec("DELETE FROM customers")
	if err != nil {
		fmt.Printf(err.Error())
	}
	_, err = db.DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func executeRequest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	routes.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
