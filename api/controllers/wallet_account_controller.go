package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/LucatonyRaudales/hindra-pay-accounts-svc/api/auth"
	"github.com/LucatonyRaudales/hindra-pay-accounts-svc/api/models"
	"github.com/LucatonyRaudales/hindra-pay-accounts-svc/api/responses"
	"github.com/LucatonyRaudales/hindra-pay-accounts-svc/api/utils/formaterror"
)

func (server *Server) CreateWalletAccount(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	walletAccount := models.WalletAccount{}
	err = json.Unmarshal(body, &walletAccount)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	walletAccount.Prepare()
	err = walletAccount.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != walletAccount.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	walletAccountCreated, err := walletAccount.SaveWalletAccount(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, walletAccountCreated.ID))
	responses.JSON(w, http.StatusCreated, walletAccountCreated)
}

func (server *Server) GetWalletAccounts(w http.ResponseWriter, r *http.Request) {

	walletAccount := models.WalletAccount{}

	walletAccounts, err := walletAccount.FindAllWalletAccounts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, walletAccount)
}

func (server *Server) GetWalletAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	walletAccount := models.WalletAccount{}

	walletAccountReceived, err := walletAccount.FindWalletAccountByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, walletAccountReceived)
}

func (server *Server) UpdateWalletAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the WalletAccount id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the WalletAccount exist
	walletAccount := models.WalletAccount{}
	err = server.DB.Debug().Model(models.WalletAccount{}).Where("id = ?", pid).Take(&walletAccount).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("WalletAccount not found"))
		return
	}

	// If a user attempt to update a WalletAccount not belonging to him
	if uid != walletAccount.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data WalletAccounted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	walletAccountUpdate := models.WalletAccount{}
	err = json.Unmarshal(body, &walletAccountUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != walletAccountUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	walletAccountUpdate.Prepare()
	err = walletAccountUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	walletAccountUpdate.ID = walletAccount.ID //this is important to tell the model the walletAccount id to update, the other update field are set above

	walletAccountUpdated, err := walletAccountUpdate.UpdateAWalletAccount(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, walletAccountUpdated)
}

func (server *Server) DeleteWalletAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid WalletAccount id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the WalletAccount exist
	walletAccount := models.WalletAccount{}
	err = server.DB.Debug().Model(models.WalletAccount{}).Where("id = ?", pid).Take(&walletAccount).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this WalletAccount?
	if uid != walletAccount.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = walletAccount.DeleteAWalletAccount(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}