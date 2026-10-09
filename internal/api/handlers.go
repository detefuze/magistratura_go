package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"klimash.detefuze/pz3-http/internal/storage"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

// GET /tasks
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()

	// Поддержка простых фильтров через query: ?q=text
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

// POST /tasks
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		BadRequest(w, "title is required")
		return
	}
	// зв.2: длина считается в символах (rune), а не в байтах — иначе кириллица даёт двойной счёт.
	if len([]rune(req.Title)) > 140 {
		UnprocessableEntity(w, "title length must be 1..140")
		return
	}

	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)
}

// GET /tasks/{id} (простой path-парсер без стороннего роутера)
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		BadRequest(w, "invalid id")
		return
	}

	t, err := h.Store.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			NotFound(w, "task not found")
			return
		}
		Internal(w, "unexpected error")
		return
	}
	JSON(w, http.StatusOK, t)
}

// parseID вытаскивает id из пути /tasks/{id} (без стороннего роутера).
func parseID(r *http.Request) (int64, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		return 0, errors.New("invalid path")
	}
	return strconv.ParseInt(parts[1], 10, 64)
}

func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		BadRequest(w, "invalid id")
		return
	}

	t, err := h.Store.SetDone(id, true)
	if err != nil {
		NotFound(w, "task not found")
		return
	}
	JSON(w, http.StatusOK, t)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		BadRequest(w, "invalid id")
		return
	}

	if err := h.Store.Delete(id); err != nil {
		NotFound(w, "task not found")
		return
	}
	NoContent(w)
}
