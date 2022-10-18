package gapi

import (
	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		CreatedAt:        timestamppb.New(user.CreatedAt),
		PasswordChangeAt: timestamppb.New(user.PasswordChangeAt),
	}
}
