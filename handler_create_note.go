package main 
import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
)

type ResponseNote struct {
	UserID uuid.UUID `json:"userID"`
	Title string `json:"title"`
	Contents string `json:"contents"`
}

func (cfg *apiConfig) handlerCreateNote(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Title string `json:"title"`
		Contents string `json:"contents"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode json", err)
		return 
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving the tokenString", err)
		return 
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to Validate JWT", err)
		return 
	}

	note, err := cfg.db.CreateNote(req.Context(), database.CreateNoteParams{
		Title: params.Title, 
		Body: params.Contents,
		UserID: userID, 
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create a note", err)
		return 
	}

	respondWithJSON(w, http.StatusCreated, ResponseNote{
		UserID: userID,
		Title: note.Title,
		Contents: note.Body,
	})

}