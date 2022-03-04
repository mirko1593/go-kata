package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Server implements an http.Hadnler that acts as either a reverse proxy or
// a simple file server, as determined by a rule set.
type Server struct {
	mu    sync.Mutex
	last  time.Time
	rules []*Rule
}

// Rule represents a rule in a configuration file.
type Rule struct {
	Host    string
	Forward string
	Serve   string

	handler http.Handler
}

// NewServer ...
func NewServer(file string, poll time.Duration) (*Server, error) {
	s := new(Server)
	if err := s.loadRules(file); err != nil {
		return nil, err
	}

	go s.refreshRules(file, poll)

	return s, nil
}

func (s *Server) loadRules(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	mtime := fi.ModTime()
	if !mtime.After(s.last) && s.rules != nil {
		return nil
	}

	rules, err := parseRules(file)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.last = mtime
	s.rules = rules
	s.mu.Unlock()

	return nil
}

func (s *Server) refreshRules(file string, poll time.Duration) error {
	for {
		if err := s.loadRules(file); err != nil {
			log.Println(err)
		}
		time.Sleep(poll)
	}
}

func parseRules(file string) ([]*Rule, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rules []*Rule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, err
	}

	for _, r := range rules {
		r.handler = makeHandler(r)
		if r.handler == nil {
			log.Printf("bad rule: %#v", r)
		}
	}

	return rules, nil
}

func makeHandler(r *Rule) http.Handler {
	if h := r.Forward; h != "" {
		return &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = h
			},
		}
	}

	if d := r.Serve; d != "" {
		return http.FileServer(http.Dir(d))
	}

	return nil
}

// hostPolicy implements autocert.HostPolicy by consulting
// the rules list for a matching host name.
func (s *Server) hostPolicy(ctx context.Context, host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, rule := range s.rules {
		if host == rule.Host || host == "www."+rule.Host {
			return nil
		}
	}

	return fmt.Errorf("unrecognized host %q", host)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Not found.", http.StatusNotFound)
}

func (s *Server) handler(r *http.Request) http.Handler {
	s.mu.Lock()
	defer s.mu.Unlock()
	h := r.Host

	if i := strings.Index(h, ":"); i > 0 {
		h = h[:i]
	}

	for _, r := range s.rules {
		if h == r.Host || strings.HasSuffix(h, "."+r.Host) {
			hitCounter.With(prometheus.Labels{"host": r.Host}).Inc()
			return r.handler
		}
	}

	return nil
}
