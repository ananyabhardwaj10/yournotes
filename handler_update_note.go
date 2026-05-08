package main
import(
	"net/http"
	"encoding/json"
	"time"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/yournotes/internal/auth"
	"github.com/ananyabhardwaj10/yournotes/internal/database"
)

type Resp struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Title string `json:"title"`
	Contents string `json:"contents"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Pinned bool `json:"isPinned"`
}

func (cfg *apiConfig) handlerUpdateNote(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Title *string `json:"title"`
		Contents *string `json:"contents"`
		Pinned *bool `json:"isPinned"`
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
	isPinned := sql.NullBool{}

	if params.Title != nil {
		title.String = *params.Title 
		title.Valid = true
	}

	if params.Contents != nil {
		contents.String = *params.Contents
		contents.Valid = true
	}

	if params.Pinned != nil {
		isPinned.Bool = true
		isPinned.Valid = true
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

	respondWithJSON(w, http.StatusOK, Resp{
		ID: note.ID,
		UserID: userID,
		Title: note.Title,
		Contents: note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	})

}