package main

import (
	"errors"
	"log"
	"strings"

	pb "github.com/Sh4d1/wat-user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type service struct {
	repo         Repository
	tokenService Authable
}

func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	log.Println(req.Id)
	user, err := s.repo.Get(req.Id)
	if err != nil {
		var err pb.Error
		log.Println("No user with id: ", req.Id)
		err.Code = 1
		err.Description = "User does not exist"
		res.Errors = append(res.Errors, &err)
		return nil
	}
	res.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := s.repo.GetByEmail(req.Email)
	log.Println(user, err)
	if err != nil {
		var err pb.Error
		err.Code = 1
		err.Description = "User does not exist, or the password is wrong"
		res.Errors = append(res.Errors, &err)

		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		var err pb.Error
		err.Code = 2
		err.Description = "User does not exist, or the password is wrong"
		res.Errors = append(res.Errors, &err)

		return nil
	}

	token, err := s.tokenService.Encode(user)
	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	req.Email = strings.Trim(req.Email, " ")
	req.Name = strings.Trim(req.Name, " ")

	if req.Email == "" {
		var err pb.Error
		err.Code = 5
		err.Description = "Email cannot be empty"
		res.Errors = append(res.Errors, &err)
		return nil
	}
	if req.Name == "" {
		var err pb.Error
		err.Code = 5
		err.Description = "Name cannot be empty"
		res.Errors = append(res.Errors, &err)
		return nil
	}

	if req.Password == "" {
		var err pb.Error
		err.Code = 5
		err.Description = "Password cannot be empty"
		res.Errors = append(res.Errors, &err)
		return nil
	}
	_, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		var err pb.Error
		err.Code = 4
		err.Description = "This email is already taken"
		res.Errors = append(res.Errors, &err)
		return nil
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPass)

	log.Println("User name:", req.Name, " Email:", req.Email, " Password:", req.Password)

	if err := s.repo.Create(req); err != nil {
		return err
	}
	res.User = req
	return nil
}

func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := s.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
