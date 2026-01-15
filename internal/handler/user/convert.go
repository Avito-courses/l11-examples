package user

import "github.com/Avito-courses/l11-examples/internal/model/user"

func ModelToResponse(user user.User) User {
	return User{
		ID:     user.ID,
		Name:   user.Name,
		Phone:  user.Phone,
		Rating: user.Rating,
	}
}
