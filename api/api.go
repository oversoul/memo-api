package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"

	"memo/api/auth"
	"memo/api/notes"
	"memo/api/notes/repository"
	"memo/api/share"
	"memo/pkg/logger"
)

type DI struct {
	Logger    logger.Logger
	NoteRepo  repository.NotesRepository
	TodoRepo  repository.TodoNotesRepository
	MovieRepo repository.MovieNotesRepository
	ShareRepo share.ShareRepository
	AuthStore auth.AuthStore
}

func FileServer(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fullPath := filepath.Join("./public", path)

		// Check if path is a directory
		fi, err := os.Stat(fullPath)
		if err != nil || fi.IsDir() || os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	})
}

func addRoutes(mux *http.ServeMux, di DI) {
	middleware := auth.AuthMiddleware{
		Mux:   mux,
		Store: di.AuthStore,
	}

	// serve files.
	fs := FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	// auth
	mux.Handle("POST /api/v1/login", auth.HandleLogin(di.AuthStore))
	mux.Handle("POST /api/v1/register", auth.HandleRegister(di.AuthStore))

	middleware.Handle("POST /api/v1/logout", auth.HandleLogout(di.AuthStore))

	middleware.Handle("GET /api/v1/profile", auth.HandleProfile(di.AuthStore))
	middleware.Handle("PUT /api/v1/profile", auth.HandleProfileUpdate(di.AuthStore))

	middleware.Handle("GET /api/v1/notes", notes.HandleAll(di.Logger, di.NoteRepo))
	middleware.Handle("GET /api/v1/notes/{id}", notes.HandleGet(di.Logger, di.NoteRepo))

	middleware.Handle("POST /api/v1/notes", notes.HandleAdd(di.Logger, di.NoteRepo))
	middleware.Handle("PUT /api/v1/notes/{id}", notes.HandleUpdate(di.Logger, di.NoteRepo))
	middleware.Handle("DELETE /api/v1/notes/{id}", notes.HandleDelete(di.Logger, di.NoteRepo))

	middleware.Handle("PUT /api/v1/notes/todo/{id}", notes.HandleUpdateTodo(di.Logger, di.TodoRepo))
	middleware.Handle("POST /api/v1/notes/todo/{id}", notes.HandleCreateTodo(di.Logger, di.TodoRepo))

	middleware.Handle("PUT /api/v1/notes/movie/{id}", notes.HandleUpdateMovie(di.Logger, di.MovieRepo))

	middleware.Handle("GET /api/v1/shared-notes", share.HandleGetShared(di.Logger, di.ShareRepo))
	middleware.Handle("POST /api/v1/notes/share", share.HandleShareNote(di.Logger, di.ShareRepo))
}

func New(di DI) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, di)

	var handler http.Handler = mux

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	return corsMiddleware(handler)
}
