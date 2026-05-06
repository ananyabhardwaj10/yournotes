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

	mux.HandleFunc("POST /api/register", apiCfg.handlerRegisterUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshTokens)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshTokens)
	
	mux.HandleFunc("GET /api/notes", apiCfg.handlerGetAllNotes)
	//get notes :id (single note)
	mux.HandleFunc("POST /api/createnote", apiCfg.handlerCreateNote)
	//put notes :id
	//delete notes :id



	server.ListenAndServe()
}