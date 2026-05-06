package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	go_ora "github.com/sijms/go-ora/v2"

	"github.com/relief-atlas/backend/models"
)

func (h *Handler) ListShelters(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := `SELECT shelter_id, shelter_name, capacity, location, occupied_number, contact_number, area_id
	          FROM Shelter ORDER BY shelter_id`
	args := []any{}

	if aid := r.URL.Query().Get("area_id"); aid != "" {
		query = `SELECT shelter_id, shelter_name, capacity, location, occupied_number, contact_number, area_id
		         FROM Shelter WHERE area_id=:1 ORDER BY shelter_id`
		args = append(args, aid)
	}

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.Shelter{}
	for rows.Next() {
		var s models.Shelter
		if err := rows.Scan(&s.ID, &s.Name, &s.Capacity, &s.Location,
			&s.OccupiedNumber, &s.ContactNumber, &s.AreaID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out = append(out, s)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetShelter(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var s models.Shelter
	err = h.db.QueryRowContext(ctx,
		`SELECT shelter_id, shelter_name, capacity, location, occupied_number, contact_number, area_id
		 FROM Shelter WHERE shelter_id=:1`, id).
		Scan(&s.ID, &s.Name, &s.Capacity, &s.Location,
			&s.OccupiedNumber, &s.ContactNumber, &s.AreaID)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "shelter not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) CreateShelter(w http.ResponseWriter, r *http.Request) {
	var s models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newID int64
	_, err := h.db.ExecContext(ctx,
		`INSERT INTO Shelter (shelter_name, capacity, location, occupied_number, contact_number, area_id)
		 VALUES (:1, :2, :3, :4, :5, :6)
		 RETURNING shelter_id INTO :7`,
		s.Name, s.Capacity, s.Location, s.OccupiedNumber, s.ContactNumber, s.AreaID,
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.ID = newID
	writeJSON(w, http.StatusCreated, s)
}

func (h *Handler) UpdateShelter(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var s models.Shelter
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`UPDATE Shelter SET shelter_name=:1, capacity=:2, location=:3,
		 occupied_number=:4, contact_number=:5, area_id=:6
		 WHERE shelter_id=:7`,
		s.Name, s.Capacity, s.Location, s.OccupiedNumber,
		s.ContactNumber, s.AreaID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "shelter not found")
		return
	}
	s.ID = id
	writeJSON(w, http.StatusOK, s)
}

func (h *Handler) DeleteShelter(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Shelter WHERE shelter_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "shelter not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
