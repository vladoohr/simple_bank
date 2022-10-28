package gapi

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/pb"
	"github.com/vladoohr/simple_bank/util"
	"github.com/vladoohr/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser validates the request and updates new user
func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	updateUserParams := db.UpdateUserParams{
		Username: req.GetUsername(),
	}

	if req.GetFullname() != nil {
		updateUserParams.FullName = sql.NullString{
			String: req.Fullname.Value,
			Valid:  true,
		}
	}

	if req.GetEmail() != nil {
		updateUserParams.Email = sql.NullString{
			String: req.GetEmail().Value,
			Valid:  req.GetEmail() != nil,
		}
	}

	if req.GetPassword() != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword().Value)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		updateUserParams.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		updateUserParams.PasswordChangeAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

	}

	user, err := server.store.UpdateUser(ctx, updateUserParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
	}

	UpdateUserResponse := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return UpdateUserResponse, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	username := req.GetUsername()
	password := req.GetPassword()
	email := req.GetEmail()
	fullname := req.GetFullname()

	if err := val.ValidateUsername(username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if password != nil {
		if err := val.ValidatePassword(password.Value); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if email != nil {
		if err := val.ValidateEmail(email.Value); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if fullname != nil {
		if err := val.ValidateFullName(fullname.Value); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}

	return violations
}
