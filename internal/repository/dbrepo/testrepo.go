package dbrepo

import (
	"errors"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
  if res.BungalowID == 99 {
    return 0, errors.New("some error")
  }
	return 1, nil
}

// InsertBungalowRestriction places a restriction in the database
func (m *testDBRepo) InsertBungalowRestriction(r models.BungalowRestriction) error {
  if r.BungalowID == 999 {
    return errors.New("some error")
  }
	return nil
}

// SearchAvailabilityByDatesByBungalowID returns true if there is availability for a bungalow between date range, false if not
func (m *testDBRepo) SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried date range
func (m *testDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {
	var bungalows []models.Bungalow
	return bungalows, nil
}

// GetBungalowByID gets a bungalow by id
func (m *testDBRepo) GetBungalowByID(id int) (models.Bungalow, error) {
  var bungalow models.Bungalow

  if id > 3 {
    return bungalow, errors.New("an error occured")
  }

  return bungalow, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
  var u models.User

  return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
  return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
  return 1, "", nil
}

// AllReservations builds and returns a slice of all reservations from the database
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
  var reservations []models.Reservation

  return reservations, nil
}

// AllNewReservations builds and returns a slice of all new reservations from the database
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
  var reservations []models.Reservation

  return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
  var res models.Reservation

  return res, nil
}

func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
  return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
  return nil
}

func (m *testDBRepo) UpdateStatusOfReservation(id int, status int) error {
  return nil
}

func (m *testDBRepo) AllBungalows() ([]models.Bungalow, error) {
  var bungalows []models.Bungalow

  return bungalows, nil
}

func (m *testDBRepo) GetRestrictionsForBungalowByDate(bungalowID int, start, end time.Time) ([]models.BungalowRestriction, error) {
  var restrictions []models.BungalowRestriction

  return restrictions, nil
}

func (m *testDBRepo) InsertBlockForBungalow(id int, startDate time.Time) error {
  return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {
  return nil
}
