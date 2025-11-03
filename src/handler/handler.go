package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Gergenus/GoMockServer/src/config"
)

type Server struct {
	Config *config.Config
	log    *slog.Logger
}

func NewServer(cfg *config.Config, log *slog.Logger) *Server {
	s := &Server{
		Config: cfg,
		log:    log,
	}

	return s
}

func (s *Server) HandleRequests(w http.ResponseWriter, r *http.Request) {
	for _, ep := range s.Config.Endpoints {
		if ep.Path == r.URL.String() && r.Method == ep.Method {
			s.log.Debug("matched endpoint",
				slog.String("path", ep.Path),
				slog.String("method", ep.Method),
			)
			if filepath.Ext(ep.ResponsePath) == ".xml" {
				w.Header().Set("Content-Type", "application/xml")
			} else if filepath.Ext(ep.ResponsePath) == ".html" {
				w.Header().Set("Content-Type", "text/html")
			}

			resBody, err := s.getResponse(ep.ResponsePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`<ValCurs> Internal server error </ValCurs>`))
				return
			}
			w.WriteHeader(ep.Status)
			w.Write(resBody)
			return
		}

	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<ValCurs> Not found </ValCurs>"))
}

func (s *Server) getResponse(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		s.log.Error("error getting response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("file not found: %w", err)
	}
	defer file.Close()
	response, err := io.ReadAll(file)
	if err != nil {
		s.log.Error("error reading response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return response, nil
}
