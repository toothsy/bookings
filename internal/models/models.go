package models

import "time"

// holds the room model
type RoomType int

const (
	Godly = iota + 1
	Empororly
	Saintly
	Kingly
)

func (rt RoomType) String() string {
	switch rt {
	case Godly:
		return "Godly"
	case Empororly:
		return "Empororly"
	case Saintly:
		return "Saintly"
	case Kingly:
		return "Kingly"
	}
	return "unknown-room-type"
}

// Users holds User model

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms holds Room model

type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reservation holds Reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomId    int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// Restriction holds Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// RoomRestriction holds RoomRestriction model
type RoomRestriction struct {
	ID            int
	RoomId        int
	ReservationId int
	RestrictionId int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// Putting in variables of other types so I can hold them if needed
