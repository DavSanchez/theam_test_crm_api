package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strconv"
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

func Test_Non_Auth_Customer_Routes(t *testing.T) {
	// Unauthenticated requests
	clearCustomersTable()
	t.Run("NO_AUTH Create customer", func(t *testing.T) {
		// Add one customer
		newCustomer := models.Customer{
			Name:                 "Test_Name",
			Surname:              "Test_Surname",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

		got := response.Body.String()
		want := "Unauthorized\n"
		if got != want {
			t.Errorf("Expected %q response. Got %q", want, got)
		}
	})
	t.Run("NO_AUTH Get all customers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

		got := response.Body.String()
		want := "Unauthorized\n"
		if got != want {
			t.Errorf("Expected %q response. Got %q", want, got)
		}
	})
	t.Run("NO_AUTH Get customer", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/1", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

		got := response.Body.String()
		want := "Unauthorized\n"
		if got != want {
			t.Errorf("Expected %q response. Got %q", want, got)
		}
	})
	t.Run("NO_AUTH Update customer", func(t *testing.T) {
		newCustomer := models.Customer{
			Name:                 "Test_Name_MODIFIED",
			Surname:              "Test_Surname_MODIFIED",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBufferString(string(data)))
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

		got := response.Body.String()
		want := "Unauthorized\n"
		if got != want {
			t.Errorf("Expected %q response. Got %q", want, got)
		}
	})
	t.Run("NO_AUTH Delete customer", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/customers/1", nil)
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

		got := response.Body.String()
		want := "Unauthorized\n"
		if got != want {
			t.Errorf("Expected %q response. Got %q", want, got)
		}
	})
}

func Test_Auth_Customer_Routes(t *testing.T) {
	clearCustomersTable()
	var token string
	// Authenticating and getting token
	t.Run("Authenticate existing user", func(t *testing.T) {
		user := models.User{
			Username: "Admin",
			Password: "hunter2",
		}
		response := authenticateUser(t, user)

		want := fmt.Sprintf(`{"result":"success", "token":"%s"}`, regexp.QuoteMeta(`^[a-zA-Z0-9-_=]+?.[a-zA-Z0-9-_=]+?.[a-zA-Z0-9-_.+/=]*?$`))
		got := response.Body.String()

		if matched, _ := regexp.MatchString(want, got); !matched {
			t.Logf("Response %v does not match expected format", got)
			t.Fail()
		}

		m := make(map[string]string)
		err := json.NewDecoder(response.Body).Decode(&m)
		if err != nil {
			t.Fatalf("Error decoding response body: %q", err.Error())
		}
		token = m["token"]
	})

	// Authenticated requests
	t.Run("AUTH Get list with no customers", func(t *testing.T) {
		clearCustomersTable()
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != "[]" {
			t.Errorf("Expected an empty array. Got %s", body)
		}
	})
	t.Run("AUTH Create one customer", func(t *testing.T) {
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusCreated, response.Code)
		want := "{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1}"

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("AUTH Get list with one customer", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(t, req)

		want := "[{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1}]"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("AUTH Create another customer", func(t *testing.T) {
		// Add one customer
		newCustomer := models.Customer{
			Name:                 "Test_Name_2",
			Surname:              "Test_Surname_2",
			PictureId:            1,
			LastModifiedByUserId: 1,
		}
		data, _ := json.Marshal(newCustomer)
		req, _ := http.NewRequest("POST", "/customers/create", bytes.NewBufferString(string(data)))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusCreated, response.Code)
		want := "{\"id\":2,\"name\":\"Test_Name_2\",\"surname\":\"Test_Surname_2\",\"pictureId\":1,\"lastModifiedByUserId\":1}"

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("AUTH Get list with two customers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/all", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(t, req)

		want := "[{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1},{\"id\":2,\"name\":\"Test_Name_2\",\"surname\":\"Test_Surname_2\",\"pictureId\":1,\"lastModifiedByUserId\":1}]"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})

	t.Run("AUTH Get a non existing customer", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/22", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		if m["error"] != "Customer not found" {
			t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
		}
	})
	t.Run("AUTH Get one customer", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/1", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(t, req)

		want := "{\"id\":1,\"name\":\"Test_Name\",\"surname\":\"Test_Surname\",\"pictureId\":1,\"lastModifiedByUserId\":1}"

		checkResponseCode(t, http.StatusOK, response.Code)

		if body := response.Body.String(); body != want {
			t.Errorf("Expected %s. Got %s", want, body)
		}
	})
	t.Run("AUTH Get a non valid ID parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/customers/hola", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		response := executeRequest(t, req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		got := response.Body.Bytes()
		want := []byte("404 page not found")
		if reflect.DeepEqual(got, want) {
			t.Errorf("Expected %s. Got '%s'", want, got)
		}
	})
}

func Test_Auth_User_Routes(t *testing.T) {
	t.Run("Authenticate existing user", func(t *testing.T) {
		user := models.User{
			Username: "Admin",
			Password: "hunter2",
		}
		response := authenticateUser(t, user)

		want := "{\"result\":\"success\", \"token\":\"^[A-Za-z0-9-_=]+.[A-Za-z0-9-_=]+.?[A-Za-z0-9-_.+/=]*$\"}"
		got := response.Body.String()

		if matched, _ := regexp.MatchString(want, got); !matched {
			t.Fatalf("Response %v does not match expected format", got)
		}
	})

	t.Run("Authenticate invalid user", func(t *testing.T) {
		user := models.User{
			Username: "Admin_NOT_EXISTS",
			Password: "hunter2",
		}
		response := authenticateUser(t, user)

		want := "{\"error\":\"Invalid credentials\"}"
		got := response.Body.String()

		if got != want {
			t.Fatalf("Expected response was %q, got %q", want, got)
		}
	})

	t.Run("Register new user", func(t *testing.T) {
		anotherUser := models.User{
			Username: "Admin_ANOTHER",
			Password: "hunter2_ANOTHER",
		}
		data, _ := json.Marshal(anotherUser)
		req, _ := http.NewRequest("POST", "/users/register", bytes.NewBufferString(string(data)))
		response := executeRequest(t, req)

		want := "{\"result\":\"success\"}"
		got := response.Body.String()

		if got != want {
			t.Fatalf("Expected response was %q, got %q", want, got)
		}
	})
	clearAdditionalUsers()
}

func clearCustomersTable() {
	_, err := db.DB.Exec("DELETE FROM customers")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = db.DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
	if err != nil {
		fmt.Print(err.Error())
	}
}

func clearAdditionalUsers() {
	_, err := db.DB.Exec("DELETE FROM users WHERE id > 1")
	if err != nil {
		fmt.Print(err.Error())
	}
	// _, err = db.DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }
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

func authenticateUser(t *testing.T, u models.User) *httptest.ResponseRecorder {
	data, _ := json.Marshal(u)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBufferString(string(data)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
	return executeRequest(t, req)
}
