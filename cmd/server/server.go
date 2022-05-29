package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	stop     chan struct{}
	waitStop *sync.WaitGroup
	once     sync.Once

	http *http.Server
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() {
	s.init()

	s.stop = make(chan struct{})
	s.waitStop = new(sync.WaitGroup)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		sig := <-interrupt
		fmt.Printf("server receive signal %q, closing\n", sig)
		close(s.stop)
	}()

	s.startHTTP()

	s.waitStop.Wait()
	fmt.Println("server existed")
}

func (s *Server) startHTTP() {
	fmt.Println("http server: start at address", s.http.Addr)

	s.waitStop.Add(1)

	go func() {
		<-s.stop
		fmt.Println("http server: closing")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.http.Shutdown(ctx); err != nil {
			fmt.Println("http server: Shutdown failed:", err)
		}
	}()

	go func() {
		defer s.waitStop.Done()
		if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("http server: ListenAndServe failed:", err)
			return
		}
		fmt.Println("http server: closed successfully")
	}()
}
