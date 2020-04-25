package main

import (
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro"
	pb "github.com/tonymj76/shippy-service-vessel/vessel/proto/vessel"
)

const (
	defaultHost = "localhost:27017"
)

func createDummyData(d Repository) {
	defer d.Close()
	vessels := []*pb.Vessel{
		{Id: 1, Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		d.Create(v)
	}
}

func main() {
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}
	session, err := CreateSession(host)
	if err != nil {
		log.Fatalf("Error connecting to datastore: %v", err)
	}

	dockyard := &VesselRepository{session.Copy()}
	createDummyData(dockyard)

	srv := micro.NewService(
		micro.Name("shippy.service.vessel"),
		micro.Version("0.1"),
	)
	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{session})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
