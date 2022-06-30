package controllers

import "github.com/LucatonyRaudales/hindra-pay-accounts-svc/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	//WalletAccounts routes
	s.Router.HandleFunc("/WalletAccounts", middlewares.SetMiddlewareJSON(s.CreateWalletAccount)).Methods("POST")
	s.Router.HandleFunc("/WalletAccounts", middlewares.SetMiddlewareJSON(s.GetWalletAccounts)).Methods("GET")
	s.Router.HandleFunc("/WalletAccounts/{id}", middlewares.SetMiddlewareJSON(s.GetWalletAccount)).Methods("GET")
	s.Router.HandleFunc("/WalletAccounts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateWalletAccount))).Methods("PUT")
	s.Router.HandleFunc("/WalletAccounts/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteWalletAccount)).Methods("DELETE")
}