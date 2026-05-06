package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/relief-atlas/backend/models"
)

func (h *Handler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx,
		`SELECT admin_id, admin_name, email FROM Admins ORDER BY admin_id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.Admin{}
	for rows.Next() {
		var a models.Admin
		if err := rows.Scan(&a.ID, &a.Name, &a.Email); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		out = append(out, a)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var a models.Admin
	err = h.db.QueryRowContext(ctx,
		`SELECT admin_id, admin_name, email FROM Admins WHERE admin_id=:1`, id).
		Scan(&a.ID, &a.Name, &a.Email)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "admin not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var a models.Admin
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if a.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newID int64
	_, err = h.db.ExecContext(ctx,
		`INSERT INTO Admins (admin_name, email, password)
		 VALUES (:1, :2, :3)
		 RETURNING admin_id INTO :4`,
		a.Name, a.Email, string(hash),
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.ID = newID
	a.Password = ""
	writeJSON(w, http.StatusCreated, a)
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var a models.Admin
	var hashed string
	err := h.db.QueryRowContext(ctx,
		`SELECT admin_id, admin_name, email, password FROM Admins WHERE email=:1`,
		req.Email).Scan(&a.ID, &a.Name, &a.Email, &hashed)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var a models.Admin
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var res sql.Result
	if a.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		res, err = h.db.ExecContext(ctx,
			`UPDATE Admins SET admin_name=:1, email=:2, password=:3 WHERE admin_id=:4`,
			a.Name, a.Email, string(hash), id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		res, err = h.db.ExecContext(ctx,
			`UPDATE Admins SET admin_name=:1, email=:2 WHERE admin_id=:3`,
			a.Name, a.Email, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "admin not found")
		return
	}
	a.ID = id
	a.Password = ""
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Admins WHERE admin_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "admin not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
