package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/sakti/o11ygo/api"
	"github.com/sakti/o11ygo/pkg/config"
	"github.com/sakti/o11ygo/pkg/service"
	"github.com/sakti/o11ygo/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run_auth(conf *config.Configuration) {
	authSvc, err := service.NewAuthService(conf.DBPath)
	checkErr(err)
	lisAuth, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.AuthsvcPort))
	checkErr(err)

	authServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	api.RegisterAuthsvcServer(authServer, authSvc)
	log.Printf("auth server listening at %v", lisAuth.Addr())
	err = authServer.Serve(lisAuth)
	checkErr(err)

}

func run_adder(conf *config.Configuration) {
	adderSvc, err := service.NewAdderService()
	checkErr(err)
	lisAdder, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.AddersvcPort))
	checkErr(err)

	adderServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	api.RegisterAddersvcServer(adderServer, adderSvc)
	log.Printf("adder server listening at %v", lisAdder.Addr())
	err = adderServer.Serve(lisAdder)
	checkErr(err)
}

func run_multiply(conf *config.Configuration) {
	multiplySvc, err := service.NewMultiplyService(conf.AddersvcPort)
	checkErr(err)
	lisMultiply, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.MultiplysvcPort))
	checkErr(err)

	multiplyServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	api.RegisterMultiplysvcServer(multiplyServer, multiplySvc)
	log.Printf("multiply server listening at %v", lisMultiply.Addr())
	err = multiplyServer.Serve(lisMultiply)
	checkErr(err)
}

func main() {
	conf, err := config.LoadConfig()
	checkErr(err)
	log.Printf("%+v", conf)

	var serviceName string
	flag.StringVar(&serviceName, "service", "auth", "a service name")
	flag.Parse()

	// tracing
	shutdown := tracing.SetupTracing(serviceName, conf.OTLPEndpoint)
	defer shutdown()

	// run service
	switch serviceName {
	case "auth":
		run_auth(conf)
	case "adder":
		run_adder(conf)
	case "multiply":
		run_multiply(conf)
	default:
		log.Panic("invalid service name")
	}
}
