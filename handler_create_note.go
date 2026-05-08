package main 
import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
)

type ResponseNote struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userID"`
	Title string `json:"title"`
	Contents string `json:"contents"`
	Pinned bool `json:"isPinned"`
}

func (cfg *apiConfig) handlerCreateNote(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Title string `json:"title"`
		Contents string `json:"contents"`
		Pinned bool `json:"isPinned"`
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

	if params.Pinned {
		pinCount, err := cfg.db.CountPinnedNotes(req.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error counting total pinned notes", err)
			return 
		}

		if pinCount >= 3 {
		respondWithError(w, http.StatusConflict, "Cannot Pin More than 3 notes", nil)
		return 
		}
	}

	note, err := cfg.db.CreateNote(req.Context(), database.CreateNoteParams{
		Title: params.Title, 
		Body: params.Contents,
		UserID: userID, 
		IsPinned: params.Pinned,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create a note", err)
		return 
	}

	respondWithJSON(w, http.StatusCreated, ResponseNote{
		ID: note.ID,
		UserID: userID,
		Title: note.Title,
		Contents: note.Body,
		Pinned: note.IsPinned,
	})

}