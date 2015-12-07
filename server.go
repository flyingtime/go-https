package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpsHandler struct {
}

func (*httpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "golang https server!!!")
}

func main() {
	pool := x509.NewCertPool()
	caCertPath := "root.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatal("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    ":8080",
		Handler: &httpsHandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	if err = s.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatal("ListenAndServeTLS err:", err)
	}
}
