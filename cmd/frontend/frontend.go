package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/sakti/o11ygo/api"
	"github.com/sakti/o11ygo/pkg/client"
	"github.com/sakti/o11ygo/pkg/config"
	"github.com/sakti/o11ygo/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

const (
	gracefulTimeout = time.Second * 10
	serviceName     = "frontend"
)

var (
	homeTmpl = template.Must(template.ParseFiles("templates/home.html"))
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewFrontend(conf *config.Configuration) (*Frontend, error) {
	multiplyClient, err := client.NewMultiplyClient(conf.MultiplysvcPort)
	if err != nil {
		return nil, err
	}
	authClient, err := client.NewAuthClient(conf.AuthsvcPort)
	if err != nil {
		return nil, err
	}
	return &Frontend{multiplyClient: multiplyClient, authClient: authClient}, nil
}

type PageData struct {
	Token  string
	A      int64
	B      int64
	Result int64
}

type Frontend struct {
	multiplyClient api.MultiplysvcClient
	authClient     api.AuthsvcClient
}

func (f *Frontend) home(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	homeTmpl.Execute(w, data)
}

func (f *Frontend) process_form(w http.ResponseWriter, r *http.Request) {
	log.Println("processing form")
	data := PageData{
		Token:  r.FormValue("token"),
		Result: 42,
	}
	a, err := strconv.ParseInt(r.FormValue("a"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := strconv.ParseInt(r.FormValue("b"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate token
	authResponse, err := f.authClient.Enforce(r.Context(), &api.EnforceRequest{Token: data.Token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authResponse.Allowed {
		http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
		return
	}
	data.A = a
	data.B = b
	response, err := f.multiplyClient.Multiply(r.Context(), &api.MultiplyRequest{A: data.A, B: data.B})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(response)
	data.Result = response.Result
	homeTmpl.Execute(w, data)
}

func main() {
	log.Println("frontend...")
	conf, err := config.LoadConfig()
	checkErr(err)
	log.Println(conf)
	// tracing
	shutdown := tracing.SetupTracing(serviceName, conf.OTLPEndpoint)

	fe, err := NewFrontend(conf)
	checkErr(err)
	// mux
	r := mux.NewRouter()
	r.Use(otelmux.Middleware(serviceName))
	r.HandleFunc("/", fe.home).Methods("GET")
	r.HandleFunc("/", fe.process_form).Methods("POST")

	// http server
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", conf.FEPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		log.Printf("fe server listening at %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	// cleanup
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()
	log.Println("shutting down")
	srv.Shutdown(ctx)
	shutdown()
	os.Exit(0)
}
