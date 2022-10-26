package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/pb"
	"github.com/vladoohr/simple_bank/util"
	"github.com/vladoohr/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser validates the request and creates new user
func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	createUserParams := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullname(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, createUserParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	createUserResponse := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return createUserResponse, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateFullName(req.Fullname); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}

	return violations
}
