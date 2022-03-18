package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"feature/api"
	"feature/data"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var port = flag.Int("port", 50051, "the port to serve in")

type ecServer struct {
	api.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()
	log.Printf("server starting on port...\n", *port)

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/key_cert.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	ca := x509.NewCertPool()
	caFilePath := data.Path("x509/client_ca_cert.pem")
	caBytes, err := ioutil.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	api.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Serve(lis)

}
