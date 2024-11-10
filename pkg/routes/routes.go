package routes

import (
	"database/sql"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/auth"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/chat"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/dialog"
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
	AuthRoutes(r, db, logger)
	DialogRoutes(r, db, logger)
	ChatRoutes(r, db, logger)
}

func PostRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	postRepo := post.NewRepository(db, logger)
	postCtrl := post.NewControllerPost(postRepo, logger)

	r.Handle("/api/create", middleware.JWTMiddleware(logger)(http.HandlerFunc(postCtrl.Create))).Methods("POST")
	r.Handle("/api/get", middleware.JWTMiddleware(logger)(http.HandlerFunc(postCtrl.Get))).Methods("POST")
}

func AuthRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	authRepo := auth.NewRepository(db, logger)
	authCtrl := auth.NewControllerAuth(authRepo, logger)

	r.HandleFunc("/api/register", authCtrl.Register).Methods("POST")
	r.HandleFunc("/api/login", authCtrl.Login).Methods("POST")
}

func DialogRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	dialogRepo := dialog.NewRepository(db, logger)
	dialogCtrl := dialog.NewControllerDialog(dialogRepo, logger)

	r.Handle("/api/dialog/create", middleware.JWTMiddleware(logger)(http.HandlerFunc((dialogCtrl.Create))))
	r.Handle("/api/dialogs", middleware.JWTMiddleware(logger)(http.HandlerFunc((dialogCtrl.Get)))).Methods("GET")
	//r.HandleFunc("/api/dialogs", dialogCtrl.Get).Methods("GET")
}

func ChatRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	chatCtrl := chat.ControllerChat{DB: db, Logger: logger}

	r.HandleFunc("/ws", chatCtrl.HandleConnections)
}
