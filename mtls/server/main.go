package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"mtls/utils"
	"net/http"
)

func main() {
	server := getServer()
	http.HandleFunc("/", myHandler)
	server.ListenAndServeTLS("", "")
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Request")
	w.Write([]byte("Hey GopherCon!"))
}

func getServer() *http.Server {
	data, _ := ioutil.ReadFile("../ca/minica.pem")
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)

	tls := &tls.Config{
		ClientCAs:             cp,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		GetCertificate:        utils.CertReqFunc("cert.pem", "key.pem"),
		VerifyPeerCertificate: utils.CertificateChains,
	}

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tls,
	}

	return server
}
