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

func (h *Handler) ListDistributions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := `SELECT distribution_id, material_name, quantity, distribution_date, area_id, admin_id
	          FROM Distribution ORDER BY distribution_id`
	args := []any{}

	if aid := r.URL.Query().Get("area_id"); aid != "" {
		query = `SELECT distribution_id, material_name, quantity, distribution_date, area_id, admin_id
		         FROM Distribution WHERE area_id=:1 ORDER BY distribution_id`
		args = append(args, aid)
	}

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.Distribution{}
	for rows.Next() {
		var d models.Distribution
		var adminID sql.NullInt64
		if err := rows.Scan(&d.ID, &d.MaterialName, &d.Quantity,
			&d.DistributionDate, &d.AreaID, &adminID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if adminID.Valid {
			d.AdminID = &adminID.Int64
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetDistribution(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var d models.Distribution
	var adminID sql.NullInt64
	err = h.db.QueryRowContext(ctx,
		`SELECT distribution_id, material_name, quantity, distribution_date, area_id, admin_id
		 FROM Distribution WHERE distribution_id=:1`, id).
		Scan(&d.ID, &d.MaterialName, &d.Quantity, &d.DistributionDate, &d.AreaID, &adminID)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "distribution not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if adminID.Valid {
		d.AdminID = &adminID.Int64
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) CreateDistribution(w http.ResponseWriter, r *http.Request) {
	var d models.Distribution
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if d.DistributionDate.IsZero() {
		d.DistributionDate = time.Now()
	}

	var adminVal any
	if d.AdminID != nil {
		adminVal = *d.AdminID
	}

	var newID int64
	_, err := h.db.ExecContext(ctx,
		`INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
		 VALUES (:1, :2, :3, :4, :5)
		 RETURNING distribution_id INTO :6`,
		d.MaterialName, d.Quantity, d.DistributionDate, d.AreaID, adminVal,
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	d.ID = newID
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handler) UpdateDistribution(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var d models.Distribution
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var adminVal any
	if d.AdminID != nil {
		adminVal = *d.AdminID
	}

	res, err := h.db.ExecContext(ctx,
		`UPDATE Distribution SET material_name=:1, quantity=:2,
		 distribution_date=:3, area_id=:4, admin_id=:5
		 WHERE distribution_id=:6`,
		d.MaterialName, d.Quantity, d.DistributionDate, d.AreaID, adminVal, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "distribution not found")
		return
	}
	d.ID = id
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) DeleteDistribution(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Distribution WHERE distribution_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "distribution not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
