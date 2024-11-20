package routes

import (
	"database/sql"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/auth"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/chat"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/dialog"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/post"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
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
	ChatRoutes(r, controller, wsManager)
	UserRoutes(r, db, logger)
	HomePage(r, db, logger)
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

func HomePage(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	r.Handle("/api/dialogs/create", middleware.JWTMiddleware(logger)(http.HandlerFunc((serveIndex))))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./internal/ui/index.html") // Путь к вашему HTML-файлу
}

func DialogRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	dialogRepo := dialog.NewRepository(db, logger)
	dialogCtrl := dialog.NewControllerDialog(dialogRepo, logger)

	r.Handle("/api/dialogs/create", middleware.JWTMiddleware(logger)(http.HandlerFunc((dialogCtrl.Create))))
	r.Handle("/api/dialogs", middleware.JWTMiddleware(logger)(http.HandlerFunc((dialogCtrl.Get)))).Methods("GET")
}

func ChatRoutes(r *mux.Router, controller *chat.ControllerChat, wsManager *chat.WebSocketManager) {
	r.HandleFunc("/api/dialogs", controller.SendMessage).Methods("POST")         // POST: Отправка сообщения
	r.HandleFunc("/api/messages/history", controller.GetMessages).Methods("GET") // GET: История сообщений

	// Обработчик WebSocket
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string) // Получение user_id из контекста (если он был установлен через middleware JWT)
		conn, err := chat.Upgrader.Upgrade(w, r, nil)   // апгрейд соединения с HTTP на WebSocket
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		// Регистрация WebSocket-соединения в WebSocketHub
		wsManager.Register(userID, conn)
	}).Methods("GET")
}

func UserRoutes(r *mux.Router, db *sql.DB, logger *slog.Logger) {
	userRepo := user.NewUserRepository(db, logger)
	userCtrl := user.NewUserController(userRepo, logger)

	r.Handle("/api/dialogs", middleware.JWTMiddleware(logger)(http.HandlerFunc((userCtrl.GetSearch))))

}
