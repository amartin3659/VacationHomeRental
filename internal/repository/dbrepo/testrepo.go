package dbrepo

import (
	"errors"
	"log"
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
  // set up a test time
	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2038-01-01 as start
	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2036-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityByDatesForAllBungalows returns a slice of available bungalows, if any for a queried date range
func (m *testDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {
	var bungalows []models.Bungalow
  
  // set up a test time
	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2038-01-01 as start
	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return nil, errors.New("some error")
	}

	// if the start date is after 2036-12-31, then return 0 bungalows,
	// indicating no availability;
	if start.After(t) {
    bungalows = make([]models.Bungalow, 0)
		return bungalows, nil
	}

	// otherwise, we have availability
  bungalows = []models.Bungalow{
    {
      ID: 1,
      BungalowName: "Name 1",
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    },
    {
      ID: 1,
      BungalowName: "Name 1",
      CreatedAt: time.Now(),
      UpdatedAt: time.Now(),
    },
  }
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
  if email == "validemail@test.com" {
    return 1, "", nil
  }
  return 0, "", errors.New("not a user")
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
  if id > 3 {
    return res, errors.New("invalid id")
  }
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
