package main

import (
	pb "github.com/Sh4d1/wat-user-service/proto/user"
	"golang.org/x/net/context"
)

type service struct {
	db           Database
	tokenService Authable
}

func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.db.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.db.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.db.GetByEmailAndPassword(req)
	if err != nil {
		return err
	}
	res.Token = "test"
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	if err := s.db.Create(req); err != nil {
		return err
	}
	res.User = req
	return nil
}

func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	return nil
}
