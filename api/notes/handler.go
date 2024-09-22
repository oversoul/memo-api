package notes

import (
	"encoding/json"
	"net/http"
	"time"

	"memo/api/notes/models"
	"memo/api/notes/repository"
	"memo/pkg/logger"
	"memo/pkg/response"
	"memo/pkg/validation"
)

func HandleAll(logger logger.Logger, repo repository.NotesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sort := r.URL.Query().Get("sort")
		userId := r.Context().Value("user").(string)
		nType := r.URL.Query().Get("type")

		filter := repository.FetchFilter{
			Count:  10,
			Sort:   sort,
			Type:   nType,
			UserId: userId,
		}

		notes, err := repo.List(filter, r.Context())
		if err != nil {
			logger.Error(err.Error())
			response.RespondErr(w, response.NotFound())
			return
		}

		response.Respond(w, notes, http.StatusOK)
	})
}

func HandleGet(logger logger.Logger, repo repository.NotesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user").(string)
		if note, err := repo.GetById(r.PathValue("id"), userId, r.Context()); err != nil {
			response.RespondErr(w, response.NotFound())
		} else {
			response.Respond(w, note, http.StatusOK)
		}
	})
}

func HandleAdd(logger logger.Logger, repo repository.NotesRepository) http.HandlerFunc {
	type noteRequest struct {
		Type  string   `json:"type" validate:"required|in:movie,todo,text"`
		Title string   `json:"title" validate:"required"`
		Tags  []string `json:"tags" validate:"array"`
	}

	type textRequest struct {
		Content string `json:"content" validate:"required"`
	}

	type movieRequest struct {
		Year     int    `json:"year" validate:"required|numeric"`
		Watched  bool   `json:"watched" validate:"required|boolean"`
		Director string `json:"director"`
	}

	type todoRequest struct {
		Tasks []string `json:"tasks" validate:"required|array"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			// NOTE: does this error needs to halt.
			// the error here can happen if the body is not valid json, or empty.
		}

		note, problems := validation.Valid[*noteRequest](data)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		var movieInfo *movieRequest
		if note.Type == "movie" {
			movieInfo, problems = validation.Valid[*movieRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		var todoInfo *todoRequest
		if note.Type == "todo" {
			todoInfo, problems = validation.Valid[*todoRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		var textInfo *textRequest
		if note.Type == "text" {
			textInfo, problems = validation.Valid[*textRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		embeddedNote := models.EmbeddedNote{
			BaseNote: models.BaseNote{
				Type:      note.Type,
				Title:     note.Title,
				Tags:      note.Tags,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		if note.Type == "text" && textInfo != nil {
			embeddedNote.TextNote = &models.TextNoteData{Content: textInfo.Content}
		}

		if note.Type == "todo" && todoInfo != nil {
			todo := models.TodoNoteData{}
			for _, task := range todoInfo.Tasks {
				todo.Tasks = append(todo.Tasks, models.Task{
					Content:     task,
					IsCompleted: false,
					CompletedAt: nil,
				})
			}

			embeddedNote.TodoNote = &todo
		}

		if note.Type == "movie" && movieInfo != nil {
			embeddedNote.MovieNote = &models.MovieNoteData{
				Year:     movieInfo.Year,
				Watched:  movieInfo.Watched,
				Director: movieInfo.Director,
			}
		}

		userId := r.Context().Value("user").(string)
		_, err := repo.Add(embeddedNote, userId, r.Context())
		if err != nil {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		response.RespondSuccess(w)
	})
}

func HandleUpdate(logger logger.Logger, repo repository.NotesRepository) http.HandlerFunc {
	type noteRequest struct {
		Title string   `json:"title" validate:"required"`
		Tags  []string `json:"tags" validate:"array"`
	}

	type textRequest struct {
		Content string `json:"content" validate:"required"`
	}

	type movieRequest struct {
		Year     int    `json:"year" validate:"required|numeric"`
		Director string `json:"director"`
	}

	type todoRequest struct {
		// Tasks []string `json:"tasks" validate:"required|array"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			// NOTE: does this error needs to halt.
			// the error here can happen if the body is not valid json, or empty.
		}

		userId := r.Context().Value("user").(string)

		oldNote, err := repo.GetById(id, userId, r.Context())
		if err != nil {
			response.RespondErr(w, response.NotFound())
			return
		}

		note, problems := validation.Valid[*noteRequest](data)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		var movieInfo *movieRequest
		if oldNote.Type == "movie" {
			movieInfo, problems = validation.Valid[*movieRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		var todoInfo *todoRequest
		if oldNote.Type == "todo" {
			todoInfo, problems = validation.Valid[*todoRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		var textInfo *textRequest
		if oldNote.Type == "text" {
			textInfo, problems = validation.Valid[*textRequest](data)
			if len(problems) > 0 {
				response.ValidationErr(w, problems)
				return
			}
		}

		oldNote.Title = note.Title

		if oldNote.Type == "text" && textInfo != nil {
			oldNote.TextNote.Content = textInfo.Content
		}

		if oldNote.Type == "todo" && todoInfo != nil {
			// TODO: this is probably shouldn't be allowed
		}

		if oldNote.Type == "movie" && movieInfo != nil {
			oldNote.MovieNote.Year = movieInfo.Year
			oldNote.MovieNote.Director = movieInfo.Director
		}

		err = repo.Update(oldNote, r.Context())
		if err != nil {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		response.RespondSuccess(w)
	})
}

func HandleDelete(logger logger.Logger, repo repository.NotesRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		userId := r.Context().Value("user").(string)

		err := repo.Delete(id, userId, r.Context())
		if err != nil {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		response.RespondSuccess(w)
	})
}
