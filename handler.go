package main

import (
	"context"

	"github.com/tonymj76/shippy-service-vessel/vessel/proto/vessel"
	"gopkg.in/mgo.v2"
)

// Our gRPC service handler
type service struct {
	session *mgo.Session
}

func (s *service) GetRepo() Repository {
	return &VesselRepository{s.session.Clone()}
}

// FindAvailable vessels in the dockyard and return it as resp
func (s *service) FindAvailable(ctx context.Context, req *vessel.Specification, res *vessel.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	vessel, err := repo.Find(req)
	res.Vessel = vessel
	return err
}

func (s *service) CreateVessel(ctx context.Context, req *vessel.Vessel, res *vessel.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	if err := repo.Create(req); err != nil {
		return err
	}
	res.Vessel = req
	res.Created = "true"
	return nil
}
