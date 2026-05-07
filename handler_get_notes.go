package main 
import(
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
)

type NoteResp struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Title string `json:"title"`
	Contents string `json:"contents"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerGetAllNotes(w http.ResponseWriter, req *http.Request) {
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to extract access token from request header", err)
		return 
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User ID mismatch", err)
		return 
	}

	notes, err := cfg.db.GetNotesByUserID(req.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get all notes", err)
		return 
	}

	var allNotes []NoteResp

	for _, note := range notes {
		allNotes = append(allNotes, NoteResp{
			ID: note.ID,
			UserID: note.UserID,
			Title: note.Title,
			Contents: note.Body,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, allNotes)
}

func (cfg *apiConfig) handlerGetSingleNote(w http.ResponseWriter, req *http.Request) {
	noteIDStr := req.PathValue("noteID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get note id from path", err)
		return 
	}

	note, err := cfg.db.GetNoteByID(req.Context(), noteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get note using id", err)
		return 
	}

	respondWithJSON(w, http.StatusOK, NoteResp{
		ID: note.ID,
		UserID: note.UserID,
		Title: note.Title,
		Contents: note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	})
}