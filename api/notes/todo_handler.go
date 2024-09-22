package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"memo/api/notes/repository"
	"memo/pkg/logger"
	"memo/pkg/response"
	"memo/pkg/validation"
)

func HandleCreateTodo(logger logger.Logger, repo repository.TodoNotesRepository) http.HandlerFunc {
	type taskRequest struct {
		Content string `json:"content" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data, problems := validation.DecodeValid[*taskRequest](r)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		id, err := repo.Create(id, data.Content, r.Context())
		if err != nil {
			logger.Error("creation issue " + err.Error())
			response.RespondErr(w, response.BadRequest())
			return
		}

		response.Respond(w, map[string]string{"task_id": id}, http.StatusOK)
	})
}

func HandleUpdateTodo(logger logger.Logger, repo repository.TodoNotesRepository) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			// NOTE: does this error needs to halt.
			// the error here can happen if the body is not valid json, or empty.
		}

		taskId, ok := data["task_id"]
		if !ok {
			logger.Error("Task id not provided")
			response.RespondErr(w, response.NotFound())
			return
		}

		updates := make(map[string]any)
		if content, ok := data["content"]; ok {
			updates["content"] = content
		}

		if is_completed, ok := data["is_completed"]; ok {
			switch is_completed.(type) {
			case bool:
				value := is_completed.(bool)
				updates["is_completed"] = value
				if value {
					updates["completed_at"] = time.Now()
				} else {
					updates["completed_at"] = nil
				}
				break
			}
		}

		if len(updates) == 0 {
			logger.Error("Nothing to update")
			response.RespondErr(w, response.BadRequest())
			return
		}

		err := repo.Update(id, taskId.(string), updates, r.Context())
		if err != nil {
			logger.Error("update issue " + err.Error())
			logger.Info(fmt.Sprintf("id: %s, taskId: %s", id, taskId.(string)))
			response.RespondErr(w, response.BadRequest())
			return
		}

		response.RespondSuccess(w)
	})
}
