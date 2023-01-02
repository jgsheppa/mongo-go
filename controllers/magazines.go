package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jgsheppa/mongo-go/errors"
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
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazine)
}

func (m *Magazine) MagazineBySlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := chi.URLParam(r, "magazineSlug")
	magazine, err := m.ms.FindBySlug(slug)

	if err != nil {
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazine)
}

func (m *Magazine) GetAllMagazines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	magazines, err := m.ms.FindAll()

	if err != nil {
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(magazines)
}

func (m *Magazine) DeleteMagazine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "magazineId")
	res, err := m.ms.Delete(id)

	if err != nil {
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	deletedCount := strconv.Itoa(int(res.DeletedCount))

	jsonMessage := Response{
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
	jsonMessage := Response{}

	title := chi.URLParam(r, "title")
	price := chi.URLParam(r, "price")
	priceInt, err := primitive.ParseDecimal128(price)
	if err != nil {
		responseError := errors.InternalError("Conversion from string to float failed", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}
	magazine := models.Magazine{
		ID:    primitive.NewObjectID(),
		Title: title,
		Price: priceInt,
	}
	res, err := m.ms.Create(magazine)

	if err != nil {
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
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
		responseError := errors.InternalError("Error converting primitve to string", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	title := chi.URLParam(r, "title")
	price := chi.URLParam(r, "price")
	priceInt, err := primitive.ParseDecimal128(price)
	if err != nil {
		responseError := errors.InternalError("Conversion from string to float failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	magazine := models.Magazine{
		ID:    objectId,
		Title: title,
		Price: priceInt,
	}

	res, err := m.ms.UpdateById(magazine)
	if err != nil {
		errorResponse := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
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

func (m *Magazine) AggregateMagazinePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	price := chi.URLParam(r, "price")

	res, err := m.ms.AggregateByPrice(price)
	if err != nil {
		errorResponse := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

func (m *Magazine) SearchMagazines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	term := chi.URLParam(r, "term")
	field := chi.URLParam(r, "field")

	res, err := m.ms.Search(field, term)
	if err != nil {
		errorResponse := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}
