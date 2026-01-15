package user

import "time"

type User struct {
	ID        int
	Name      string
	Phone     string
	Rating    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
