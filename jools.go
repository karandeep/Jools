package main

import (
	"./controller"
	"./model"
	"./webadmin"
	"encoding/gob"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8081", "Port number on which to run")
	flag.Parse()

	gob.Register(&model.User{})
	r := mux.NewRouter()
	r.HandleFunc("/", controller.HomeHandler)
	r.HandleFunc("/comment/{action:[A-Za-z0-9-]+/?}", controller.CommentHandler)
	r.HandleFunc("/contest/{action:[A-Za-z0-9-]+/?}", controller.ContestHandler)
	r.HandleFunc("/fetch/{item:[A-Za-z0-9-]+/?}", controller.FetchHandler)
	r.HandleFunc("/inspiration/{page:[A-Za-z0-9-]+/?}", controller.InspirationHandler)
	r.HandleFunc("/product/{action:[A-Za-z0-9-]+/?}", controller.ProductHandler)
	r.HandleFunc("/rate", controller.RatingHandler)
	r.HandleFunc("/rate/skip", controller.RatingSkipHandler)
	r.HandleFunc("/sync/{item:[A-Za-z0-9-]+/?}", controller.SyncHandler)
	r.HandleFunc("/track", controller.TrackingHandler)
	r.HandleFunc("/user/{page:[A-Za-z0-9-]+/?}", controller.UserHandler)
	r.HandleFunc("/service/{action:[A-Za-z0-9-]+/?}", controller.ServiceHandler)
	r.HandleFunc("/webadmin/{action:[A-Za-z0-9-]+/?}", webadmin.WebAdminHandler)

	r.HandleFunc("/news/{action:[A-Za-z0-9-]+/?}", controller.NewsHandler)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println(err)
	}
}
