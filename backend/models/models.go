package models

import "time"

type Disaster struct {
	ID        int64     `json:"disaster_id"`
	Name      string    `json:"disaster_name"`
	Type      string    `json:"disaster_type"`
	StartDate time.Time `json:"start_date"`
	Status    string    `json:"status"`
}

type AffectedArea struct {
	ID         int64  `json:"area_id"`
	Name       string `json:"area_name"`
	Severity   string `json:"severity"`
	Population int64  `json:"population"`
	DisasterID int64  `json:"disaster_id"`
}

type Shelter struct {
	ID             int64  `json:"shelter_id"`
	Name           string `json:"shelter_name"`
	Capacity       int64  `json:"capacity"`
	Location       string `json:"location"`
	OccupiedNumber int64  `json:"occupied_number"`
	ContactNumber  string `json:"contact_number"`
	AreaID         int64  `json:"area_id"`
}

type Admin struct {
	ID       int64  `json:"admin_id"`
	Name     string `json:"admin_name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Distribution struct {
	ID               int64     `json:"distribution_id"`
	MaterialName     string    `json:"material_name"`
	Quantity         int64     `json:"quantity"`
	DistributionDate time.Time `json:"distribution_date"`
	AreaID           int64     `json:"area_id"`
	AdminID          *int64    `json:"admin_id"`
}

type User struct {
	ID       int64  `json:"user_id"`
	Name     string `json:"user_name"`
	Email    string `json:"user_email"`
	PhoneNo  string `json:"user_phoneno"`
	Password string `json:"password,omitempty"`
}

type Donation struct {
	ID           int64     `json:"donation_id"`
	Amount       float64   `json:"amount"`
	DonationDate time.Time `json:"donation_date"`
	UserID       *int64    `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
