package token

import (
	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) *mux.Router {
	router.HandleFunc("/login", HandlerGenerateTokenPair).Methods("POST")
	router.HandleFunc("/refresh", HandlerRefreshTokenPair).Methods("POST")

	return router
}
