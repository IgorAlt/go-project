package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"unrealProject/db"
)

func getUserIDFromURLParams(w http.ResponseWriter, r *http.Request) int {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return 0
	}
	return id
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := getUserIDFromURLParams(w, r)

	result, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
