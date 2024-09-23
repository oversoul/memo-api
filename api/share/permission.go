package share

import (
	"fmt"
	"memo/api/notes/models"
)

func CanWrite(note *models.EmbeddedNote, userId string) bool {
	if note == nil {
		fmt.Println("[CanWrite] Note is nil")
		return false
	}
	if note.UserId.Hex() != userId {
		hasPermission := false
		for _, sharedUser := range note.SharedWith {
			if sharedUser.UserID.Hex() == userId && sharedUser.Permission == models.PermissionWrite {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return false
		}
	}

	return true
}

func CanRead(note *models.EmbeddedNote, userId string) bool {
	if note == nil {
		fmt.Println("[CanRead] Note is nil")
		return false
	}

	if note.UserId.Hex() != userId {
		hasPermission := false
		for _, sharedUser := range note.SharedWith {
			if sharedUser.UserID.Hex() == userId {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return false
		}
	}

	return true
}
