package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"ha-video-parser/pkg/image"
	"ha-video-parser/pkg/translate"
	"ha-video-parser/pkg/video"

	"github.com/google/uuid"
)

type (
	ParseRequest struct {
		Path string `json:"path"`
	}

	ParseResponse struct {
		Path string `json:"path"`
	}

	Service struct {
	}
)

func New() *Service {
	return &Service{}
}

func (s *Service) RegisterHttpHandlers() {
	http.HandleFunc("/parse", s.handle)
}

func (s *Service) handle(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Failed to read request")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var request ParseRequest
	if err = json.Unmarshal(body, &request); err != nil {
		log.Printf("Failed to parse request body [%v]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id := uuid.New()
	frames := fmt.Sprintf("output/%s/frames", id.String())

	if err = os.MkdirAll(frames, os.ModePerm); err != nil {
		log.Printf("Failed to create directory for frame extraction [%s]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if ok := video.ExtractFrames(request.Path, frames); !ok {
		log.Printf("Failed to extract frames [%s]", request.Path)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	crop := fmt.Sprintf("output/%s/crop", id.String())

	if err = os.MkdirAll(crop, os.ModePerm); err != nil {
		log.Printf("Failed to create director for crop [%s]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if ok := image.BatchCrop(frames, crop); !ok {
		log.Printf("Failed to crop frames [%s]", frames)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	uniques := image.NonEmptyUnique(crop, fmt.Sprintf("output/%s", id.String()))

	translate.BatchTranslate(uniques)

	var response ParseResponse
	response.Path = fmt.Sprintf("%v", uniques)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to marshal response [%v]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

}
