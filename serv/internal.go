package serv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dosco/graphjin/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	configRoute = "/v1/config"
	reloadRoute = "/v1/reload"
)

func configHandler(s1 *Service) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		s := s1.Load().(*service)
		var b []byte
		b, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
		var c core.ExternalConfigRequest
		if err != nil {
			intErr(w, fmt.Sprintf("error loading external config body: %s", err.Error()))
			return
		}

		defer r.Body.Close()
		err = json.Unmarshal(b, &c)
		if err != nil {
			intErr(w, fmt.Sprintf("error parsing external config: %s", err.Error()))
			return
		}

		ec := s.gj.GetExternalConfig()
		err = ec.Store(c)
		if err != nil {
			intErr(w, fmt.Sprintf("error storing external config: %s", err.Error()))
			return
		}

		err = ec.Load()
		if err != nil {
			intErr(w, fmt.Sprintf("error loading external config: %s", err.Error()))
			return
		}

		io.WriteString(w, "config stored successfully")
	}

	return http.HandlerFunc(h)
}

func reloadHandler(s1 *Service) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		s := s1.Load().(*service)
		ec := s.gj.GetExternalConfig()
		err := ec.Load()
		if err != nil {
			intErr(w, fmt.Sprintf("error loading external config: %s", err.Error()))
			return
		}

		io.WriteString(w, "config loaded successfully")
	}

	return http.HandlerFunc(h)
}

func internalRouteHandler(s1 *Service, mux *http.ServeMux) (http.Handler, error) {
	s := s1.Load().(*service)

	s.log.Info("setting up config routes")
	mux.Handle(configRoute, configHandler(s1))
	mux.Handle(reloadRoute, reloadHandler(s1))

	return setServerHeader(mux), nil
}

func startInternalHTTP(s1 *Service) {
	s := s1.Load().(*service)

	routes, err := internalRouteHandler(s1, http.NewServeMux())
	if err != nil {
		s.log.Fatalf("error setting up routes: %s", err)
	}

	s.intSrv = &http.Server{
		Addr:           s.conf.internalHostPort,
		Handler:        routes,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.intSrv.Shutdown(context.Background()); err != nil {
			s.log.Warn("shutdown signal received")
		}
		close(idleConnsClosed)
	}()

	fields := []zapcore.Field{
		zap.String("host-port", s.conf.InternalHostPort),
		zap.String("app-name", s.conf.AppName),
		zap.String("env", os.Getenv("GO_ENV")),
		zap.Bool("production", s.conf.Core.Production),
		zap.Bool("secrets-used", (s.conf.Serv.SecretsFile != "")),
	}

	s.zlog.Info("Config endpoints initialized", fields...)

	l, err := net.Listen("tcp", s.conf.internalHostPort)
	if err != nil {
		s.log.Fatalf("failed to init internal port: %s", err)
	}

	// signal we are open for business.
	s.state = servListening

	if err := s.intSrv.Serve(l); err != http.ErrServerClosed {
		s.log.Fatalf("failed to start internal server: %s", err)
	}
	<-idleConnsClosed
}
