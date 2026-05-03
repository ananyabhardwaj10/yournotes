package main 
import(
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
)

type User struct {
	ID uuid.UUID           `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name string          `json:"name"`
	Email string         `json:"email"`
	Password string      `json:"-"`
}

func (cfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Password string   `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding the request body", err)
		return 
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing the password", err)
		return 
	}

	user := database.User{}

	user, err = cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Name: params.Name,
		Email: params.Email,
		HashedPassword: hashed_password,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating the user", err)
		return 
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name: user.Name,
		Email: user.Email,
	})

}