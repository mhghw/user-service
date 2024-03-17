package ports

import (
	genUser "github.com/mhghw/user-service/pb/gen"
	userPkg "github.com/mhghw/user-service/pkg/domain/user"
)

func EncodeUser(user userPkg.User) *genUser.User {
	return &genUser.User{
		Id:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
}
