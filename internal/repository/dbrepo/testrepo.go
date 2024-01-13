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
