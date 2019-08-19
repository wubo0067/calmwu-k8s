package handler

import (
	"net/http"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/api/handler"
	"github.com/micro/go-micro/api/handler/event"
	"github.com/micro/go-micro/api/router"
	"github.com/micro/go-micro/errors"

	// TODO: only import handler package
	aapi "github.com/micro/go-micro/api/handler/api"
	ahttp "github.com/micro/go-micro/api/handler/http"
	arpc "github.com/micro/go-micro/api/handler/rpc"
	aweb "github.com/micro/go-micro/api/handler/web"
)

type metaHandler struct {
	s micro.Service
	r router.Router
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError(m.r.Options().Namespace, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// TODO: don't do this ffs
	switch service.Endpoint.Handler {
	// web socket handler
	case aweb.Handler:
		aweb.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	// proxy handler
	case "proxy", ahttp.Handler:
		ahttp.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	// rpcx handler
	case arpc.Handler:
		arpc.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	// event handler
	case event.Handler:
		ev := event.NewHandler(
			handler.WithNamespace(m.r.Options().Namespace),
			handler.WithService(m.s),
		)
		ev.ServeHTTP(w, r)
	// api handler
	case aapi.Handler:
		aapi.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	// default handler: rpc
	default:
		arpc.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	}
}

// Meta is a http.Handler that routes based on endpoint metadata
func Meta(s micro.Service, r router.Router) http.Handler {
	return &metaHandler{
		s: s,
		r: r,
	}
}
