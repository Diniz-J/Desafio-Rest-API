package handler

import (
	"encoding/json"
	"net/http"

	"github.com/DinizJ/desafio/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

// Adjust Layers
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

// Cria nova tafera

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	//Parse request
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//Validar entrada
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	//Chama service
	task, err := h.service.CreateTask(r.Context(), req.Title, req.Description)
	//r.Context() é cancelado se o cliente fechar a conexão ou der timeout!
	if err != nil {
		//erro em service
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//Erro ja tratado!
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get task", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	task, err := h.service.Update(r.Context(), id, req.Title, req.Description, req.Status, req.Priority)
	if err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
