package main

import (
	pb "github.com/Sh4d1/wat-user-service/proto/user"
	"github.com/jinzhu/gorm"
)

type Database interface {
	GetAll() ([]*pb.User, error)
	Get(id string) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmailAndPassword(user *pb.User) (*pb.User, error)
}

type UserDatabase struct {
	db *gorm.DB
}

func (db *UserDatabase) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	if err := db.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (db *UserDatabase) Get(id string) (*pb.User, error) {
	var user *pb.User
	user.Id = id
	if err := db.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *UserDatabase) GetByEmailAndPassword(user *pb.User) (*pb.User, error) {
	if err := db.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *UserDatabase) Create(user *pb.User) error {
	if err := db.db.Create(user).Error; err != nil {
		return err
	}
}
