package handlers

import (
	"encoding/json"
	"github.com/Uranury/exploreMicro/service1/internal/models"
	"github.com/Uranury/exploreMicro/service1/internal/storage"
	"net/http"
	"strconv"
)

type User struct {
	store storage.Store
}

func NewUser(storage storage.Store) *User {
	return &User{
		store: storage,
	}
}

func (h *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.store.Save(&user)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *User) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	idStr := query.Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, ok := h.store.Get(uint(id))
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
	}
	_ = json.NewEncoder(w).Encode(user)
}

func (h *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := h.store.List()
	_ = json.NewEncoder(w).Encode(result)
}

type UpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Age  *uint   `json:"age,omitempty"`
}

func (h *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := h.store.Get(uint(id))
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	user.Balance = req.Balance
	h.store.Save(user)

	json.NewEncoder(w).Encode(user)
}
