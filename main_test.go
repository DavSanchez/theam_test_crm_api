package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"theam.io/jdavidsanchez/test_crm_api/api"
	"theam.io/jdavidsanchez/test_crm_api/db"
)

func TestMain(m *testing.M) {
	code := m.Run()

	clearCustomersTable()

	os.Exit(code)
}

func TestAPI_listAllCustomers(t *testing.T) {

	t.Run("Testing getting an empty customer list", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != "[]" {
			t.Errorf("Expected an empty array. Got %s", body)
		}
	})
}

func TestAPI_getCustomer(t *testing.T) {

	t.Run("Testing getting a non existing customer", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/11", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		if m["error"] != "Customer not found" {
			t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
		}
	})
}

func TestAPI_createCustomer(t *testing.T) {
	// TODO: not implemented
}

func TestAPI_updateCustomer(t *testing.T) {
	// TODO: not implemented
}

func TestAPI_deleteCustomer(t *testing.T) {
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
	api.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
