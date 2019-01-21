package configserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HTTPConfigServer struct {
	storage Storage
}

func (h *HTTPConfigServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")
	log.Printf("Request: Path %v", tmp)
	if len(tmp) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
	identifier := tmp[1]
	if r.Method == http.MethodGet {
		h.Get(identifier, w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.Set(identifier, w, r)
		return
	}
	if r.Method == http.MethodOptions {
		h.LastUpdate(identifier, w, r)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func (h *HTTPConfigServer) Get(identifier string, w http.ResponseWriter, r *http.Request) {
	c, err := h.storage.Get(identifier)
	if err == NoSuchFile {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%v", NoSuchFile)))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	b, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (h *HTTPConfigServer) Set(identifier string, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	conf := Config{
		Checksum: "notSortedOut",
		Config:   string(data),
	}
	err = h.storage.Set(identifier, conf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}

func (h *HTTPConfigServer) LastUpdate(identifier string, w http.ResponseWriter, r *http.Request) {
	c, err := h.storage.LastUpdate(identifier)
	if err == NoSuchFile {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%v", NoSuchFile)))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	log.Printf("Returning %s", c.String())
	ret, err := c.MarshalText()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
	return
}

func NewHTTPConfigServer(storage Storage) *HTTPConfigServer {
	return &HTTPConfigServer{storage: storage}
}
