package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jgsheppa/mongo-go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Magazine struct {
	ms models.MagazineService
}

type Response struct {
	Message      string
	Error        bool
	ErrorMessage error
	StatusCode   int
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
	res, err := m.ms.Delete(id)

	jsonMessage := Response{}

	if err != nil {
		jsonMessage = Response{
			Message:      "Document not Found",
			Error:        true,
			ErrorMessage: err,
			StatusCode:   http.StatusNotFound,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonMessage)
		return
	}

	deletedCount := strconv.Itoa(int(res.DeletedCount))

	jsonMessage = Response{
		Message:      "Documents deleted: " + deletedCount,
		Error:        false,
		ErrorMessage: nil,
		StatusCode:   http.StatusAccepted,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(jsonMessage)

}

func (m *Magazine) CreateMagazine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	title := chi.URLParam(r, "title")
	price := chi.URLParam(r, "price")
	magazine := models.Magazine{
		ID:    primitive.NewObjectID(),
		Title: title,
		Price: price,
	}
	res, err := m.ms.Create(magazine)

	jsonMessage := Response{}

	if err != nil {
		jsonMessage = Response{
			Message:      "Document not found",
			Error:        true,
			ErrorMessage: err,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonMessage)
		return
	}

	jsonMessage = Response{
		Message:      "Document created. Unique id is : " + res.InsertedID.(primitive.ObjectID).Hex(),
		Error:        false,
		ErrorMessage: nil,
		StatusCode:   http.StatusAccepted,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(jsonMessage)
}

func (m *Magazine) UpdateMagazine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonMessage := Response{}

	id := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		jsonMessage = Response{
			Message:      "Error converting primitve to string",
			Error:        true,
			ErrorMessage: err,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonMessage)
		return
	}

	title := chi.URLParam(r, "title")
	price := chi.URLParam(r, "price")
	magazine := models.Magazine{
		ID:    objectId,
		Title: title,
		Price: price,
	}

	res, err := m.ms.UpdateById(magazine)
	if err != nil {
		jsonMessage = Response{
			Message:      "Document not found",
			Error:        true,
			ErrorMessage: err,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonMessage)
		return
	}

	count := strconv.Itoa(int(res.UpsertedCount))

	jsonMessage = Response{
		Message:      "Documents created: " + count,
		Error:        false,
		ErrorMessage: nil,
		StatusCode:   http.StatusAccepted,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(jsonMessage)
}
