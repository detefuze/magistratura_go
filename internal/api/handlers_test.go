package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"klimash.detefuze/pz3-http/internal/storage"
)

func newTestHandlers() *Handlers {
	return NewHandlers(storage.NewMemoryStore())
}

func TestCreateTaskOK(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", rec.Code)
	}
	var got storage.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == 0 || got.Title != "hello" || got.Done {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	h := newTestHandlers()
	long := strings.Repeat("a", 141)
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"`+long+`"}`))
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422, got %d", rec.Code)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	h := newTestHandlers()
	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	rec := httptest.NewRecorder()

	h.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}

func TestPatchTaskSetsDone(t *testing.T) {
	h := newTestHandlers()
	created := h.Store.Create("do me")

	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	h.PatchTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var got storage.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != created.ID || !got.Done {
		t.Fatalf("task not marked done: %+v", got)
	}
}

func TestDeleteTask(t *testing.T) {
	h := newTestHandlers()
	created := h.Store.Create("remove me")

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	h.DeleteTask(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rec.Code)
	}
	if _, err := h.Store.Get(created.ID); err != storage.ErrNotFound {
		t.Fatalf("task still present, err=%v", err)
	}
}
