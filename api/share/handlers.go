package share

import (
	"net/http"

	"memo/pkg/logger"
	"memo/pkg/response"
	"memo/pkg/validation"
)

func HandleGetShared(logger logger.Logger, repo ShareRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user").(string)
		notes, err := repo.List(userId, r.Context())
		if err != nil {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusBadGateway,
				Message: err.Error(),
			})
			return
		}

		response.Respond(w, notes, http.StatusOK)
	})
}

func HandleShareNote(logger logger.Logger, repo ShareRepository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUserId := r.Context().Value("user").(string)

		data, problems := validation.DecodeValid[*shareRequest](r)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		if currentUserId == data.UserID {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusBadGateway,
				Message: "Can't share the note with yourself.",
			})
			return
		}

		if err := repo.ShareNote(data, r.Context()); err != nil {
			response.RespondErr(w, response.ErrorResponse{
				Status:  http.StatusBadGateway,
				Message: err.Error(),
			})
			return
		}

		response.RespondSuccess(w)
	})
}
