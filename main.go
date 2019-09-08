package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/golang/glog"
)

func main() {
	var parameters WhSvrParameters

	flag.IntVar(&parameters.port, "port", 443, "Webhook server port.")
	flag.StringVar(&parameters.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&parameters.dnsCfgFile, "dnsCfgFile", "/etc/webhook/config/dnsconfig.yaml", "File containing the mutation configuration.")
	flag.Parse()
	
	dnsConfig, err := loadConfig(parameters.dnsCfgFile)
	if err != nil {
		glog.Errorf("Failed to load configuration: %v", err)
	}
	
	pair, err := tls.LoadX509KeyPair(parameters.certFile, parameters.keyFile)
	if err != nil {
		glog.Errorf("Failed to load key pair: %v", err)
	}
	
	whsvr := &WebhookServer {
		dnsConfig:    dnsConfig,
		server:           &http.Server {
			Addr:        fmt.Sprintf(":%v", parameters.port),
			TLSConfig:   &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}
	
	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.serve)
	whsvr.server.Handler = mux
	glog.Infof("Starting wenhook server ...")
	
	// start webhook server in new rountine
	go func() {
		if err := whsvr.server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()
	
	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	
	glog.Infof("Got OS shutdown signal, shutting down wenhook server gracefully...")
	whsvr.server.Shutdown(context.Background())
}
