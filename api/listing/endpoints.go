package listing

import (
	"net/http"

	"github.com/leveldorado/experiment/api/lib/response"

	"github.com/julienschmidt/httprouter"
	"github.com/leveldorado/experiment/api/lib/headers"
	"github.com/leveldorado/experiment/grpc/portspb"
)

func RegisterEndpoints(r *httprouter.Router, serviceClient portspb.ListingServiceClient) {
	r.Handle(http.MethodGet, "/api/v1/ports/:id", MakeGetPortEndpoint(serviceClient))
	r.Handle(http.MethodGet, "/api/v1/ports", MakeListPortEndpoint(serviceClient))
}

/*
MakeGetPortEndpoint creates handler for GET /api/v1/ports/:id endpoint
*/
func MakeGetPortEndpoint(serviceClient portspb.ListingServiceClient) func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		headers.SetContentType(w.Header(), headers.ApplicationJSON)
		id := params.ByName("id")
		if id == "" {
			response.Error(w, http.StatusBadRequest, "missing id")
			return
		}
		port, err := serviceClient.Get(r.Context(), &portspb.GetPortRequest{Id: id})
		if response.RespondErrorIfNeeded(w, err) {
			return
		}
		response.Write(w, http.StatusOK, port)
	}
}

/*
MakeListPortsEndpoint creates handler for GET /api/v1/ports endpoint
*/
func MakeListPortEndpoint(serviceClient portspb.ListingServiceClient) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		headers.SetContentType(w.Header(), headers.ApplicationJSON)
		cl, err := serviceClient.List(r.Context(), &portspb.ListPortsRequest{})
		if response.RespondErrorIfNeeded(w, err) {
			return
		}
		response.StreamResponse(w, cl, func() interface{} { return &portspb.Port{} })
	}
}
