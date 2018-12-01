package api_server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type MLH2018ApiServer struct {
	muxRouter   *mux.Router
	appId       string
	endpointKey string
	region      string
	port        string
}

func NewAPIServer(listening int, secret string, client *Client) *MLH2018ApiServer {

	router := mux.NewRouter()
	router.HandleFunc("/get_intent/{secret}/{phrase}", GetIntent(client, secret)).Methods("GET")

	return &MLH2018ApiServer{
		muxRouter: router,
		port:      fmt.Sprintf(":%d", listening),
	}
}

func (m *MLH2018ApiServer) Run() {
	log.Fatal(http.ListenAndServe(m.port, m.muxRouter))
}

func GetIntent(client *Client, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if secret != params["secret"] {
			w.WriteHeader(400)
			log.Println("invalid secret")
			return
		}
		msg := params["phrase"]
		predictionResult, err := client.Predict(msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		_ = json.NewEncoder(w).Encode(predictionResult)
	}

}
