package response

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

func Write(w http.ResponseWriter, code int, resp interface{}) {
	w.WriteHeader(code)
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf(`failed to marshal error: [response %+v, err: %s]`, resp, err))
	}
	if _, err = w.Write(data); err != nil {
		log.Println(fmt.Sprintf(`failed to write response: [data %s, err: %s]`, data, err))
	}
}

func StreamResponse(w http.ResponseWriter, cl grpc.ClientStream, newMessageObjectFunc func() interface{}) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		Error(w, http.StatusUpgradeRequired, "client should support streaming   response")
		return
	}
	w.WriteHeader(http.StatusOK)
	en := json.NewEncoder(w)
	proxyGRPCMessageToHTTPClient(flusher, en, cl, newMessageObjectFunc)
}

func proxyGRPCMessageToHTTPClient(f http.Flusher, en *json.Encoder, cl grpc.ClientStream, newMessageObjectFunc func() interface{}) {
	for {
		msg := newMessageObjectFunc()
		err := cl.RecvMsg(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(fmt.Sprintf(`failed to receive message: [err: %s]`, err))
			break
		}
		if err = en.Encode(msg); err != nil {
			log.Println(fmt.Sprintf(`failed to encode message: [err: %s]`, err))
			break
		}
		f.Flush()
	}
}
