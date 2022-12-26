package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jgsheppa/mongo-go/models"
)

type Magazine struct {
	ms models.MagazineService
}

func NewMagazine(ms models.MagazineService) *Magazine {
	return &Magazine{
		ms,
	}
}

func (m *Magazine) MagazineById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "magazineId")
	magazine, err := m.ms.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Document not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazine)
}

func (m *Magazine) GetAllMagazines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	magazines, err := m.ms.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Document not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazines)
}

func (m *Magazine) DeleteMagazine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "magazineId")
	_, err := m.ms.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Document not found"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Document removed"))
}