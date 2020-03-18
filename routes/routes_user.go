package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"theam.io/jdavidsanchez/test_crm_api/auth"
	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/models"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

/***************
User auth routes
****************/

func registerUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	if len := utf8.RuneCountInString(string(u.Password)); len < 12 {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Password less than 12 characters"})
		return
	}

	err = u.CreateUser(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, map[string]string{"result": "Success"})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	err = u.LoginUser(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	auth.SetJWT(u.Username, w, r)
}

// func logoutUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Logout user")
// }
