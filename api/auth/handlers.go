package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"memo/pkg/response"
	"memo/pkg/upload"
	"memo/pkg/validation"
)

func HandleLogin(store AuthStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, problems := validation.DecodeValid[*loginRequest](r)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		if token, err := store.Authenticate(data, r.Context()); err != nil {
			response.ErrMessage(w, "Can't login", http.StatusUnauthorized)
		} else {
			response.Respond(w, token, http.StatusOK)
		}
	})
}

func HandleProfile(store AuthStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user").(string)

		if u, err := store.GetUserById(userId, r.Context()); err != nil {
			response.ErrMessage(w, "User not found", http.StatusNotFound)
		} else {
			response.Respond(w, u, http.StatusOK)
		}
	})
}

func HandleProfileUpdate(store AuthStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if strings.TrimSpace(name) == "" {
			response.ValidationErr(w, map[string]string{"name": "required"})
			return
		}

		userId := r.Context().Value("user").(string)
		u, err := store.GetUserById(userId, r.Context())
		if err != nil {
			response.ErrMessage(w, "User not found", http.StatusNotFound)
			return
		}

		u.Name = name

		fu := upload.NewUpload(r, "image")

		fu.SetMaxSize(1 * 1000 * 1000) // 1Mb
		fu.SetAllowedTypes("png", "jpg")

		file, err := fu.ValidateAndUpload("public/images")
		if err != nil && err != http.ErrMissingFile {
			response.ErrMessage(w, "Can't upload image", http.StatusNotFound)
			return
		}

		if file != nil {
			// delete old image
			if u.Image != "" {
				if err := os.Remove(u.Image); err != nil {
					fmt.Println("Image was not deleted.")
				}
			}

			u.Image = file.Name
		}

		err = store.UpdateUserInfo(u, r.Context())
		if err != nil {
			response.ErrMessage(w, "User not updated", http.StatusBadRequest)
			return
		}

		response.Respond(w, map[string]string{"name": u.Name, "image": u.Image}, http.StatusOK)
	})
}

func HandleRegister(store AuthStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, problems := validation.DecodeValid[*registerRequest](r)
		if len(problems) > 0 {
			response.ValidationErr(w, problems)
			return
		}

		if err := store.CreateUser(data, r.Context()); err != nil {
			response.ErrMessage(w, "Can't register", http.StatusUnauthorized)
			return
		}

		if token, err := store.Authenticate(data.AsLogin(), r.Context()); err != nil {
			response.ErrMessage(w, "Couldn't create token", http.StatusUnauthorized)
		} else {
			response.Respond(w, token, http.StatusOK)
		}
	})
}

func HandleLogout(store AuthStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user").(string)

		if err := store.DeleteToken(userId, r.Context()); err != nil {
			response.ErrMessage(w, "Can't logout", http.StatusNotModified)
		} else {
			response.Respond(w, map[string]any{"message": "Logged out"}, http.StatusOK)
		}
	})
}
