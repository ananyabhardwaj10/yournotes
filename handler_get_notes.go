package main 
import(
	"net/http"
	"time"
	"encoding/json"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
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

func (cfg *apiConfig) handlerUpdateNote(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Title *string `json:"title"`
		Contents *string `json:"contents"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading the request body", err)
		return 
	}

	noteIDStr := req.PathValue("noteID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get note ID from path", err)
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

	title := sql.NullString{}
	contents := sql.NullString{}

	if params.Title != nil {
		title.String = *params.Title 
		title.Valid = true
	}

	if params.Contents != nil {
		contents.String = *params.Contents
		contents.Valid = true
	}


	note, err := cfg.db.UpdateNote(req.Context(), database.UpdateNoteParams{
		ID: noteID,
		UserID: userID,
		Title: title,
		Body: contents,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update note. Please try again", err)
		return 
	}

	respondWithJSON(w, http.StatusOK, NoteResp{
		ID: note.ID,
		UserID: userID,
		Title: note.Title,
		Contents: note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	})

}