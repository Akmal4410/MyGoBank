package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	storage       Storage
}

func NewApiServer(listenAddress string, storage Storage) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
		storage:       storage,
	}
}

func (server *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFucn(server.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFucn(server.handleAccountById))
	router.HandleFunc("/transfer", makeHTTPHandleFucn(server.handleTrasferAccount))

	fmt.Println("Go Bank Running on port : ", server.listenAddress)
	log.Fatal(http.ListenAndServe(server.listenAddress, router))

}

func (server *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return server.handleGetAccout(w, r)
	case "POST":
		return server.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (server *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return server.handleAccountById(w, r)
	case "DELETE":
		return server.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (server *ApiServer) handleGetAccout(w http.ResponseWriter, r *http.Request) error {
	accounts, err := server.storage.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (server *ApiServer) handleGetAccoutById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("Enter a valid ID")
	}
	account, err := server.storage.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (server *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := server.storage.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (server *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		return err
	}
	if id == 0 {
		return fmt.Errorf("Enter a valid ID")
	}
	if err := server.storage.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (server *ApiServer) handleTrasferAccount(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferReq)
}

func GetId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFucn(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
