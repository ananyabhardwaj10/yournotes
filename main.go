package main
import(
	"net/http"
	"os"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
)

type apiConfig struct {
	db *database.Queries
	jwtSecretKey string
}

func main() {
	godotenv.Load()
	dbUrl:= os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("error opening the database: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":8085",
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Working just fine"))
	})

	apiCfg := apiConfig{
		db: dbQueries,
		jwtSecretKey: jwtSecretKey,
	}

	mux.HandleFunc("POST /api/auth/register", apiCfg.handlerRegisterUser)
	mux.HandleFunc("POST /api/auth/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/auth/refresh", apiCfg.handlerRefreshTokens)
	mux.HandleFunc("POST /api/auth/revoke", apiCfg.handlerRevokeRefreshTokens)	
	mux.HandleFunc("GET /api/notes", apiCfg.handlerGetAllNotes)
	mux.HandleFunc("GET /api/notes/{noteID}", apiCfg.handlerGetSingleNote)
	mux.HandleFunc("POST /api/notes", apiCfg.handlerCreateNote)
	mux.HandleFunc("PATCH /api/notes/{noteID}", apiCfg.handlerUpdateNote)
	mux.HandleFunc("DELETE /api/notes/{noteID}", apiCfg.handlerDeleteNote)



	server.ListenAndServe()
}