package share

import "memo/api/notes/models"

func CanWrite(note *models.BaseNote, userId string) bool {
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

func CanRead(note *models.BaseNote, userId string) bool {
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
