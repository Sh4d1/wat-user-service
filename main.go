package main

import (
	"log"

	pb "github.com/Sh4d1/wat-user-service/proto/user"
	micro "github.com/micro/go-micro"
)

func main() {
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	database := &UserDatabase{db}

	tokenService := &TokenService{database}

	service := micro.NewService(
		micro.Name("wat.user"),
	)

	service.Init()

	pb.RegisterUserHandler(service.Server(), &service{database, tokenService})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
