package adding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/leveldorado/experiment/api/lib/response"

	"github.com/julienschmidt/httprouter"
	"github.com/leveldorado/experiment/grpc/portspb"
)

func RegisterEndpoints(r *httprouter.Router, serviceClient portspb.AddingServiceClient) {
	r.Handle(http.MethodPost, "/api/v1/ports", MakePOSTPortEndpoint(serviceClient))
}

/*
MakeGetPortEndpoint creates handler for POST /api/v1/ports endpoint
*/
func MakePOSTPortEndpoint(serviceClient portspb.AddingServiceClient) func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handleInputJSONObject(w, r.Body, func(key string, d *json.Decoder) error {
			port := portspb.Port{}
			if err := d.Decode(&port); err != nil {
				return err
			}
			port.Id = key
			_, err := serviceClient.Save(r.Context(), &port)
			return err
		})
	}
}

func prepareJSONDecoder(r io.Reader) (*json.Decoder, error) {
	d := json.NewDecoder(r)
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, errors.New("invalid json object")
	}
	return d, nil
}

func handleInputJSONObject(w http.ResponseWriter, r io.Reader, handleObjectFunc func(key string, d *json.Decoder) error) {
	d, err := prepareJSONDecoder(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	for d.More() {
		key, err := readJSONObjectKey(d)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		err = handleObjectFunc(key, d)
		if response.RespondErrorIfNeeded(w, err) {
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}

func readJSONObjectKey(d *json.Decoder) (string, error) {
	token, err := d.Token()
	if err != nil {
		return "", err
	}
	key, ok := token.(string)
	if !ok {
		return "", fmt.Errorf(`token key is not string: [token: %s]`, token)
	}
	return key, nil
}
