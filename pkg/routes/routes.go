package routes

import (
	"database/sql"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/chat"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/post"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func RegisterRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	PostRoutes(r, db, logger)
}

func PostRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	postRepo := post.NewRepository(db, logger)
	postCtrl := post.NewControllerPost(postRepo, logger)

	r.HandleFunc("/api/create", postCtrl.Create).Methods("POST")
	r.HandleFunc("/api/get", postCtrl.Get).Methods("GET")
	r.HandleFunc("/ws", chat.HandleConnections)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
}
