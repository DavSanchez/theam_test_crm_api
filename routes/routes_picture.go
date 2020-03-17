package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/models"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

/*************
Picture routes
**************/

func getPicturePath(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["pictureId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid picture ID"})
		return
	}

	p := models.PicturePath{
		Id: id,
	}
	err = p.GetPicturePath(db.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Picture not found"})
		default:
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	utils.ResponseJSON(w, http.StatusOK, p)
}

func addPicture(w http.ResponseWriter, r *http.Request) {

	var p models.PicturePath
	imageName, err := utils.FileUpload(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	p.Path = imageName
	err = p.AddPicture(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, p)
}
