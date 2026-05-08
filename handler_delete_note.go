package main 
import(
	"net/http"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
)

func (cfg *apiConfig) handlerDeleteNote(w http.ResponseWriter, req *http.Request) {
	noteIDStr := req.PathValue("noteID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get note id from path", err)
		return 
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get access token from request header", err)
		return 
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate jwt", err)
		return 
	}

	err = cfg.db.DeleteNote(req.Context(), database.DeleteNoteParams{
		ID: noteID,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete note. Please try again", err)
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}