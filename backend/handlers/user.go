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

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctx,
		`SELECT user_id, user_name, user_email, user_phoneno FROM Users ORDER BY user_id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	out := []models.User{}
	for rows.Next() {
		var u models.User
		var phone sql.NullString
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &phone); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		u.PhoneNo = phone.String
		out = append(out, u)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var u models.User
	var phone sql.NullString
	err = h.db.QueryRowContext(ctx,
		`SELECT user_id, user_name, user_email, user_phoneno FROM Users WHERE user_id=:1`, id).
		Scan(&u.ID, &u.Name, &u.Email, &phone)
	u.PhoneNo = phone.String
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if u.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newID int64
	_, err = h.db.ExecContext(ctx,
		`INSERT INTO Users (user_name, user_email, user_phoneno, password)
		 VALUES (:1, :2, :3, :4)
		 RETURNING user_id INTO :5`,
		u.Name, u.Email, u.PhoneNo, string(hash),
		go_ora.Out{Dest: &newID},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.ID = newID
	u.Password = ""
	writeJSON(w, http.StatusCreated, u)
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var u models.User
	var hashed string
	var phone sql.NullString
	err := h.db.QueryRowContext(ctx,
		`SELECT user_id, user_name, user_email, user_phoneno, password
		 FROM Users WHERE user_email=:1`, req.Email).
		Scan(&u.ID, &u.Name, &u.Email, &phone, &hashed)
	u.PhoneNo = phone.String
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
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var res sql.Result
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to hash password")
			return
		}
		res, err = h.db.ExecContext(ctx,
			`UPDATE Users SET user_name=:1, user_email=:2, user_phoneno=:3, password=:4
			 WHERE user_id=:5`,
			u.Name, u.Email, u.PhoneNo, string(hash), id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		res, err = h.db.ExecContext(ctx,
			`UPDATE Users SET user_name=:1, user_email=:2, user_phoneno=:3 WHERE user_id=:4`,
			u.Name, u.Email, u.PhoneNo, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	u.ID = id
	u.Password = ""
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx,
		`DELETE FROM Users WHERE user_id=:1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
