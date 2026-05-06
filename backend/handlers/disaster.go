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

func (h *Handler) ListDisasters(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx,
		`SELECT disaster_id, disaster_name, disaster_type, start_date, status
		 FROM Disaster ORDER BY disaster_id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.Disaster{}
	for rows.Next() {
		var d models.Disaster
		var dtype, status sql.NullString
		if err := rows.Scan(&d.ID, &d.Name, &dtype, &d.StartDate, &status); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		d.Type = dtype.String
		d.Status = status.String
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetDisaster(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var d models.Disaster
	var dtype, status sql.NullString
	err = h.db.QueryRowContext(ctx,
		`SELECT disaster_id, disaster_name, disaster_type, start_date, status
		 FROM Disaster WHERE disaster_id = :1`, id).
		Scan(&d.ID, &d.Name, &dtype, &d.StartDate, &status)
	d.Type = dtype.String
	d.Status = status.String
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "disaster not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) CreateDisaster(w http.ResponseWriter, r *http.Request) {
	var d models.Disaster
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newID int64
	_, err := h.db.ExecContext(ctx,
		`INSERT INTO Disaster (disaster_name, disaster_type, start_date, status)
		 VALUES (:1, :2, :3, :4)
		 RETURNING disaster_id INTO :5`,
		d.Name, d.Type, d.StartDate, d.Status,
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	d.ID = newID
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handler) UpdateDisaster(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var d models.Disaster
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`UPDATE Disaster SET disaster_name=:1, disaster_type=:2, start_date=:3, status=:4
		 WHERE disaster_id=:5`,
		d.Name, d.Type, d.StartDate, d.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "disaster not found")
		return
	}
	d.ID = id
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) DeleteDisaster(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Disaster WHERE disaster_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "disaster not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
