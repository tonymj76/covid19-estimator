package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/micro/go-micro"
	"github.com/tonymj76/shippy-service-vessel/vessel/proto/vessel"
)

type dockyard interface {
	Find(*vessel.Specification) (*vessel.Vessel, error)
}

// Dockyard where vessels are stored
type Dockyard struct {
	vessels []*vessel.Vessel
}

// Find a particule vessels in a Dockyard
func (d *Dockyard) Find(spec *vessel.Specification) (*vessel.Vessel, error) {
	for _, vessel := range d.vessels {
		if vessel.Capacity >= spec.Capacity && vessel.MaxWeight >= spec.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("vessels not found with that particuler specs")
}

//Service of Vessels
type Service struct {
	docky dockyard
}

// FindAvailable vessels in the dockyard and return it as resp
func (s *Service) FindAvailable(ctx context.Context, req *vessel.Specification, res *vessel.Response) error {
	vessel, err := s.docky.Find(req)
	res.Vessel = vessel
	return err
}

func main() {
	vessels := []*vessel.Vessel{
		&vessel.Vessel{
			Id:        "vessel001",
			Name:      "Boaty McBoatface",
			MaxWeight: 200000,
			Capacity:  500,
		},
	}

	dockyard := &Dockyard{vessels}
	srv := micro.NewService(
		micro.Name("shippy.service.vessel"),
	)
	srv.Init()

	vessel.RegisterVesselServiceHandler(srv.Server(), &Service{dockyard})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
