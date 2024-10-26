package routes

import (
	"database/sql"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/auth"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/chat"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/post"
	"github.com/Fyefhqdishka/deadlock_v.1/pkg/middleware"
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
	authRepo := auth.NewRepository(db, logger)
	authCtrl := auth.NewControllerAuth(authRepo, logger)

	r.Handle("/api/create", middleware.JWTMiddleware(logger)(http.HandlerFunc(postCtrl.Create))).Methods("POST")
	r.Handle("/api/get", middleware.JWTMiddleware(logger)(http.HandlerFunc(postCtrl.Get))).Methods("POST")
	r.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	r.HandleFunc("/ws", chat.HandleConnections)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
}
