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

func (h *Handler) ListDonations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := `SELECT donation_id, amount, donation_date, user_id
	          FROM Donations ORDER BY donation_id`
	args := []any{}

	if uid := r.URL.Query().Get("user_id"); uid != "" {
		query = `SELECT donation_id, amount, donation_date, user_id
		         FROM Donations WHERE user_id=:1 ORDER BY donation_id`
		args = append(args, uid)
	}

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.Donation{}
	for rows.Next() {
		var d models.Donation
		var userID sql.NullInt64
		if err := rows.Scan(&d.ID, &d.Amount, &d.DonationDate, &userID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if userID.Valid {
			d.UserID = &userID.Int64
		}
		out = append(out, d)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetDonation(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var d models.Donation
	var userID sql.NullInt64
	err = h.db.QueryRowContext(ctx,
		`SELECT donation_id, amount, donation_date, user_id
		 FROM Donations WHERE donation_id=:1`, id).
		Scan(&d.ID, &d.Amount, &d.DonationDate, &userID)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "donation not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userID.Valid {
		d.UserID = &userID.Int64
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) CreateDonation(w http.ResponseWriter, r *http.Request) {
	var d models.Donation
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if d.DonationDate.IsZero() {
		d.DonationDate = time.Now()
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var userVal any
	if d.UserID != nil {
		userVal = *d.UserID
	}

	var newID int64
	_, err := h.db.ExecContext(ctx,
		`INSERT INTO Donations (amount, donation_date, user_id)
		 VALUES (:1, :2, :3)
		 RETURNING donation_id INTO :4`,
		d.Amount, d.DonationDate, userVal,
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	d.ID = newID
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handler) DeleteDonation(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Donations WHERE donation_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "donation not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
