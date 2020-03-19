package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/auth"
	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/models"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

/**************
Customer routes
***************/

func listAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := models.ListAllCustomers(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"]) // This parameter is always an int (Regex in mux route)

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := models.Customer{
		CustomerOut: models.CustomerOut{
			Id: id,
		},
	}
	err = c.GetCustomer(db.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Customer not found"})
		default:
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	cOut := c.CustomerOut
	utils.ResponseJSON(w, http.StatusOK, cOut)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var c models.Customer
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	id, err := auth.GetUserIdFromJWT(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.CreatedByUserId = id

	err = c.CreateCustomer(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	cOut := c.CustomerOut
	utils.ResponseJSON(w, http.StatusCreated, cOut)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["customerId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	var c models.Customer
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&c)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	id, err := auth.GetUserIdFromJWT(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.LastModifiedByUserId = id

	c.Id = userId
	err = c.UpdateCustomer(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	//c.GetCustomer(db.DB)

	cOut := c.CustomerOut
	utils.ResponseJSON(w, http.StatusOK, cOut)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := models.Customer{
		CustomerOut: models.CustomerOut{
			Id: id,
		},
	}
	err = c.DeleteCustomer(db.DB)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
