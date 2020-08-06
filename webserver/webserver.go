package webserver

import (
	"context"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/webserver/api"
	"log"

	"net/http"
	"strconv"
	"time"
)

var webserver *http.Server

func StartWebDaemon() {
	mux := http.NewServeMux()

	cfg := config.GetCBtlm()
	mux.Handle(cfg.GetpurchaseWebPath(), &api.PurchaseLicense{})
	mux.Handle(cfg.GetListMinersWebPath(), &api.ListMiners{})
	mux.Handle(cfg.GetRegisterMinerWebPath(), &api.MinerRegister{})

	addr := ":" + strconv.Itoa(config.GetCBtlm().HttpServerPort)

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)

	log.Println("Web Server Stopped")
}
