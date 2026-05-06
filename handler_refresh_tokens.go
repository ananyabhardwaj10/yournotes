package main
import(
	"net/http"
	"time"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
)

type refToken struct {
	AccessToken string `json:"access_token"`
}

func (cfg *apiConfig) handlerRefreshTokens(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to exctract refresh token from request header", err)
		return 
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting user from refresh token", err)
		return 
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create a new access token", err)
		return 
	}

	respondWithJSON(w, http.StatusOK, refToken{
		AccessToken: accessToken,
	})

}

func (cfg *apiConfig) handlerRevokeRefreshTokens(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header) 
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to extract refresh token from request header", err)
		return 
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke the refresh token", err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}
