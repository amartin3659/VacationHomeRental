package models

import "time"

// Reservation contains reservaton data
type Reservation struct {
	Name  string
	Email string
	Phone string
}

// User is the model of user data
type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bungalow is the model of a bugalow
type Bungalow struct {
	ID           int
	BungalowName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Restrictions is the model of a restriction
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations is the model of a reservation
type Reservations struct {
	ID         int
	FullName   string
	Email      string
	Phone      string
	StartDate  time.Time
	EndDate    time.Time
	BungalowID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Bungalow   Bungalow
	Processed  int
}

type BungalowRestrictions struct {
  ID int
  StartDate time.Time
  EndDate time.Time
  BungalowID int
  ReservationID int
  RestrictionID int
  CreatedAt time.Time
  UpdatedAt time.Time
  Bungalow Bungalow
  Reservation Reservations
  Restriction Restrictions
}
