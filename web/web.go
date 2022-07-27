package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"monster/algoutil"
	"monster/web/api"
)

var (
	logger = log.WithField("source", "http")
)

func Startup() {
	var (
		addr      = viper.GetString("webserver.addr")
		cert      = viper.GetString("webserver.certificates.cert")
		key       = viper.GetString("webserver.certificates.key")
		enableSSL = viper.GetBool("webserver.enable_ssl")
	)
	logger.Infof("Web service addr: %s(enable ssl:%v)", addr, enableSSL)
	go func() {
		mux := startupService()
		if enableSSL {
			log.Fatal(http.ListenAndServeTLS(addr, cert, key, mux))
		} else {
			log.Fatal(http.ListenAndServe(addr, mux))
		}
	}()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case s := <-sg:
		log.Infof("got signal: %s", s.String())
	}
}

func logRequest(ctx context.Context, r *http.Request) (context.Context, error) {
	if uri := r.RequestURI; uri != "/ping" {
		logger.Debugf("Method=%s, RemoteAddr=%s URL=%s", r.Method, r.RemoteAddr, uri)
	}
	return ctx, nil
}

type ErrorMessage struct {
	Code string `json:"code"`
}

func startupService() http.Handler {
	var (
		mux    = http.NewServeMux()
		webDir = viper.GetString("webserver.static_dir")
	)
	nex.SetErrorEncoder(func(err error) interface{} {
		return &ErrorMessage{Code: err.Error()}
	})
	nex.Before(logRequest)

	mux.Handle("/user/", api.MakeUserService())

	mux.Handle("/monster/", api.MakeMonsterService())

	mux.Handle("/embattle/", api.MakeEmbattleService())

	mux.Handle("/prop/", api.MakePropService())

	mux.Handle("/matching/", api.MakeMatchingService())

	mux.Handle("/friend/", api.MakeFriendService())

	mux.Handle("/ranking/", api.MakeRankingService())

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(webDir))))

	mux.Handle("/ping", nex.Handler(pongHandler))

	return algoutil.AccessControl(algoutil.OptionControl(mux))
}

func pongHandler() (string, error) {
	return "pong", nil
}
