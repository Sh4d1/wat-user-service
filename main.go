package main

import (
	"log"
	"os"

	pb "github.com/Sh4d1/wat-user-service/proto/user"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	var srv micro.Service

	if os.Getenv("DEV") == "true" {
		srv = micro.NewService(
			micro.Name("wat.user"),
		)
	} else {
		srv = k8s.NewService(
			micro.Name("go.micro.api.user"),
		)
	}
	srv.Init()

	pb.RegisterUserServiceHandler(srv.Server(), &service{
		repo:         repo,
		tokenService: tokenService,
	})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
