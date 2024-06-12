package networkapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/envmission/template-api/common/http/types"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

const maximumRequestLength = 1 << 20
const maximumRequestLengthAttachment = 100 << 20

type UserProfileAPI struct {
	// dependensi disini
}

// Simple CRUD untuk user profile
func New() *UserProfileAPI {
	return &UserProfileAPI{
		// taro dependensi disini / usecase
	}
}

// Tambahkan user baru
func (n *UserProfileAPI) Put(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Baca dan validasi request body
	r.Body = http.MaxBytesReader(w, r.Body, maximumRequestLength)
	input, err := io.ReadAll(r.Body)
	if err != nil {
		errMessage := serializeError(&types.CommonError{
			Errors: []types.Error{
				{Message: "Failed to read all body", Code: "SERVER_ERROR"},
			},
		},
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMessage)
		return
	}

	// Validasi

	// Lakukan sesuatu dari input. Jalankan business logic, masuk ke "Usecase" layer.
	_ = input

	response := []byte(`{"response": "good good"}`)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (n *UserProfileAPI) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userID string
	var ID string

	// Dapat param user_id
	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		d := serializeError(&types.CommonError{
			Errors: []types.Error{
				{HTTPCode: http.StatusBadRequest, Code: "EMPTY_USER_ID", Message: "Please specify 'user_id'"},
			},
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	ID = r.URL.Query().Get("id")

	// validasi

	// lakukan sesuatu ke ID, masukan ke usecase
	_ = userID
	_ = ID

	result := "mantap"

	// Balikan response
	payload, err := json.Marshal(&types.CommonResponse{
		Success: result,
	})

	if err != nil {
		log.Err(err).Msgf("Failed to parse payload")
		errMessage := serializeError(&types.CommonError{
			Errors: []types.Error{
				{Message: "Failed to parse response", Code: "SERVER_ERROR"},
			},
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMessage)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// Delete network data
func (n *UserProfileAPI) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userID string

	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		d := serializeError(&types.CommonError{
			Errors: []types.Error{
				{HTTPCode: http.StatusBadRequest, Code: "EMPTY_USER_ID", Message: "Please specify 'user_id'"},
			},
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	ID := r.URL.Query().Get("id")

	// validasi

	// lakukan sesuatu ke ID, masukan ke usecase
	_ = userID
	_ = ID

	result := "mantap"

	// Balikan response

	payload, err := json.Marshal(&types.CommonResponse{
		Success: &result,
	})

	if err != nil {
		log.Err(err).Msgf("Failed to parse payload")
		errMessage := serializeError(&types.CommonError{
			Errors: []types.Error{
				{Message: "Failed to parse response", Code: "SERVER_ERROR"},
			},
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMessage)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)

}

func serializeError(err *types.CommonError) []byte {
	d, errMarshal := json.Marshal(&types.CommonResponse{
		Error: err,
	})
	if errMarshal != nil {
		log.Err(errMarshal).Msgf("Failed to parse err")
	}
	return d
}
