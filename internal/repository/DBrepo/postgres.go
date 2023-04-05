package dbrepo

import (
	"context"
	"time"

	"github.com/toothsy/bookings-app/internal/models"
)

var TIMEOUT = 3 * time.Second

func (m *postrgesDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation and returns an id in-case of successful exec
func (m *postrgesDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// incase the client looses connection during the transaction, we dont want it to run it during that time
	// to avoid that situation, we use context package that helps us to keep track of such events.

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	sqlStmt := `insert into reservations(first_name,last_name,email,phone,start_date,end_date,room_id,
		created_at,updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	var resId int
	err := m.DB.QueryRowContext(ctx, sqlStmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now()).Scan(&resId)

	if err != nil {
		return 0, err
	}

	return resId, nil
}

// InsertRoomReservation inserts a room restriction into the database
func (m *postrgesDBRepo) InsertRoomReservation(restr models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	sqlStmt := `insert into room_restrictions(start_date,end_date,room_id,reservation_id,
		created_at,updated_at,restriction_id) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, sqlStmt,
		restr.StartDate,
		restr.EndDate,
		restr.RoomId,
		restr.ReservationId,
		time.Now(),
		time.Now(),
		restr.RestrictionId,
	)

	if err != nil {
		return err
	}
	return nil
}

// SearchRoomReservationByRoomID searches for room availibilty with roomID start, end date true if no restrictions exists, false otherwise
func (m *postrgesDBRepo) SearchRoomReservationByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	sqlStmt := `select count(id ) 	
			from room_restrictions rr 
			where room_id=$1
			and $2 <= rr.end_date  and $3>= rr.start_date`
	rows, err := m.DB.QueryContext(ctx, sqlStmt, roomID, end, start)
	if err != nil {
		return false, nil
	}
	var count int
	rows.Scan(&count)
	if count != 0 {
		return false, nil
	}
	return true, nil

}

// SearchAvailability searches availability based on start and end date
func (m *postrgesDBRepo) SearchAvailability(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	var rooms []models.Room
	sqlQuery := `
			select r.id ,r.room_name 
			from rooms r 
			where r.id 
			not in 
			(    select rr.id
				from room_restrictions rr 
				where $1<= rr.end_date  and $2>= rr.start_date
			)`

	rows, err := m.DB.QueryContext(ctx, sqlQuery, end, start)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}
