package notes

import (
	"net/http"

	"memo/api/notes/repository"
	"memo/pkg/logger"
	"memo/pkg/response"
	"memo/pkg/validation"
)

func HandleUpdateMovie(logger logger.Logger, repo repository.MovieNotesRepository) http.HandlerFunc {
	type movieRequest struct {
		Year     int    `json:"year" validate:"required|numeric"`
		Watched  bool   `json:"watched" validate:"required|boolean"`
		Director string `json:"director"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data, problems := validation.DecodeValid[*movieRequest](r)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		updates := map[string]any{
			"year":     data.Year,
			"watched":  data.Watched,
			"director": data.Director,
		}

		err := repo.Update(id, updates, r.Context())
		if err != nil {
			logger.Error("update movie issue " + err.Error())
			response.RespondErr(w, response.BadRequest())
			return
		}

		response.RespondSuccess(w)
	})
}
