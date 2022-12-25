package controllers

import (
	"encoding/json"
	"fmt"
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

func(m *Magazine) MagazineById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "magazineId")
	magazine, err := m.ms.FindById(id)
	fmt.Printf("magazine: %v", magazine)
	fmt.Printf("err: %v", err)


	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Document not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazine)
}
