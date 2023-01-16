package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tmlamb/document-api/cmd/document-api/model"
	"github.com/tmlamb/document-api/cmd/document-api/repository"
)

// ReturnAllAccounts ...
func ReturnAllAccounts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: returnAllAccounts")
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	accounts, err := repository.FindAllAccounts()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	err = encoder.Encode(accounts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// CreateNewAccount ...
func CreateNewAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: createNewAccount")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var account model.Account
	err := json.Unmarshal(reqBody, &account)
	if err != nil {
		log.Printf("Error unmarshalling request body (%q): %v\n", reqBody, err)
		w.WriteHeader(400)
		return
	}

	newAccount, err := repository.SaveAccount(account)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err = encoder.Encode(newAccount)
	if err != nil {
		log.Printf("Error encoding response: %v\n", err)
		w.WriteHeader(500)
		return
	}
}

// ReturnAccount ...
func ReturnAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: returnAccount")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	account, err := repository.FindOneAccount(accountId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err = encoder.Encode(account)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// UpdateAccount ...
func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: updateAccount")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updatedAccount model.Account
	err = json.Unmarshal(reqBody, &updatedAccount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	updatedAccount.AccountID = accountId

	err = repository.UpdateAccount(updatedAccount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// DeleteAccount ...
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: deleteAccount")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	err = repository.DeleteAccount(accountId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}
