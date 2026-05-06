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

func (h *Handler) ListAreas(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := `SELECT area_id, area_name, severity, population, disaster_id
	          FROM Affected_Areas ORDER BY area_id`
	args := []any{}

	// Optional filter: ?disaster_id=N
	if did := r.URL.Query().Get("disaster_id"); did != "" {
		query = `SELECT area_id, area_name, severity, population, disaster_id
		         FROM Affected_Areas WHERE disaster_id=:1 ORDER BY area_id`
		args = append(args, did)
	}

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.AffectedArea{}
	for rows.Next() {
		var a models.AffectedArea
		if err := rows.Scan(&a.ID, &a.Name, &a.Severity, &a.Population, &a.DisasterID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out = append(out, a)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetArea(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var a models.AffectedArea
	err = h.db.QueryRowContext(ctx,
		`SELECT area_id, area_name, severity, population, disaster_id
		 FROM Affected_Areas WHERE area_id=:1`, id).
		Scan(&a.ID, &a.Name, &a.Severity, &a.Population, &a.DisasterID)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "area not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) CreateArea(w http.ResponseWriter, r *http.Request) {
	var a models.AffectedArea
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newID int64
	_, err := h.db.ExecContext(ctx,
		`INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
		 VALUES (:1, :2, :3, :4)
		 RETURNING area_id INTO :5`,
		a.Name, a.Severity, a.Population, a.DisasterID,
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.ID = newID
	writeJSON(w, http.StatusCreated, a)
}

func (h *Handler) UpdateArea(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var a models.AffectedArea
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`UPDATE Affected_Areas SET area_name=:1, severity=:2, population=:3, disaster_id=:4
		 WHERE area_id=:5`,
		a.Name, a.Severity, a.Population, a.DisasterID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "area not found")
		return
	}
	a.ID = id
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) DeleteArea(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Affected_Areas WHERE area_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "area not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
