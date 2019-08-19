package service

import (
	"context"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/router"
	pb "github.com/micro/go-micro/router/proto"
)

type table struct {
	table    pb.TableService
	callOpts []client.CallOption
}

// Create new route in the routing table
func (t *table) Create(r router.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Create(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// Delete deletes existing route from the routing table
func (t *table) Delete(r router.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Delete(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// Update updates route in the routing table
func (t *table) Update(r router.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Update(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// List returns the list of all routes in the table
func (t *table) List() ([]router.Route, error) {
	resp, err := t.table.List(context.Background(), &pb.Request{}, t.callOpts...)
	if err != nil {
		return nil, err
	}

	routes := make([]router.Route, len(resp.Routes))
	for i, route := range resp.Routes {
		routes[i] = router.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int(route.Metric),
		}
	}

	return routes, nil
}

// Lookup looks up routes in the routing table and returns them
func (t *table) Query(q router.Query) ([]router.Route, error) {
	// call the router
	resp, err := t.table.Query(context.Background(), &pb.QueryRequest{
		Query: &pb.Query{
			Service: q.Options().Service,
			Gateway: q.Options().Gateway,
			Network: q.Options().Network,
		},
	}, t.callOpts...)

	// errored out
	if err != nil {
		return nil, err
	}

	routes := make([]router.Route, len(resp.Routes))
	for i, route := range resp.Routes {
		routes[i] = router.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int(route.Metric),
		}
	}

	return routes, nil
}
