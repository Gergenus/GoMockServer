package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Gergenus/GoMockServer/src/config"
)

type Server struct {
	Config *config.Config
	Logger *slog.Logger
}

func NewServer(cfg *config.Config, log *slog.Logger) *Server {
	s := &Server{
		Config: cfg,
		Logger: log,
	}

	return s
}

func (s *Server) matchPath(pattern, path string) (bool, map[string]string) {
	if pattern != path {
		return false, nil
	}

	queryParamValues := map[string]string{}

	patternParts := strings.Split(pattern, "?")
	if len(patternParts) == 1 {
		return true, nil
	}
	params := strings.Split(patternParts[1], "&")
	for _, param := range params {
		data := strings.Split(param, "=")
		queryParamValues[data[0]] = data[1]
	}

	return true, queryParamValues
}

func (s *Server) HandleRequests(w http.ResponseWriter, r *http.Request) {
	for _, ep := range s.Config.Endpoints {
		match, _ := s.matchPath(ep.Path, r.URL.String())
		if match && r.Method == ep.Method {
			w.Header().Set("Content-Type", "application/xml")
			resBody, err := s.getXMLResonse(ep.XMLPath)
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

func (s *Server) getXMLResonse(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		s.Logger.Error("error getting xml response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("file not found: %w", err)
	}
	defer file.Close()
	response, err := io.ReadAll(file)
	if err != nil {
		s.Logger.Error("error reading xml response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return response, nil
}
