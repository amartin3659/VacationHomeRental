package repository

import (
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertBungalowRestriction(r models.BungalowRestriction) error
	SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error)
	SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error)
}
