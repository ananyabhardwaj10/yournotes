package main
import(
	"net/http"
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
)

type response struct {
	ID uuid.UUID  `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string  `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong. Please try again", err)
		return 
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get user details using email", err)
		return 
	}

	match, err := auth.CheckHashedPassword(params.Password, user.HashedPassword)
	if !match || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return 
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create jwt", err)
		return 
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 30),
		UserID: user.ID,
	})


	respondWithJSON(w, http.StatusOK, response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name: user.Name,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	})

}