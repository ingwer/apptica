package httpserver

import (
	"apptica/commander"
	"apptica/model"

	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	addr         = "0.0.0.0:8080"
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second

	appID       = 1421444
	countryCode = 1
)

func NewServer(commander commander.Commander) *http.Server {
	mux := http.NewServeMux()
	handler := topCategoryPositionsHandlerCreator(commander)
	mux.Handle("/appTopCategory", handler)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server listen error: %s\n", err)
		}
	}()

	return srv
}

type endpointResponse struct {
	StatusCode int                               `json:"status_code"`
	Message    string                            `json:"message"`
	Data       map[model.Category]model.Position `json:"data"`
}

func topCategoryPositionsHandlerCreator(cmd commander.Commander) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		date := request.URL.Query().Get("date")
		if date == "" {
			log.Println("missing date parameter")

			err := writeResponse(writer, endpointResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "fail",
				Data:       nil,
			})

			if err != nil {
				log.Println("failed to write response: ", err)
			}

			return
		}

		res, err := cmd.GetTopPositions(
			request.Context(),
			commander.Request{
				CountryCode: countryCode,
				AppID:       appID,
				Date:        model.Date(date),
			},
		)

		if err != nil {
			log.Println("failed to get response: ", err)

			err := writeResponse(writer, endpointResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "fail",
				Data:       nil,
			})

			if err != nil {
				log.Println("failed to write response: ", err)
			}

			return
		}

		err = writeResponse(writer, endpointResponse{
			StatusCode: http.StatusOK,
			Message:    "ok",
			Data:       res,
		})
		if err != nil {
			log.Println("failed to write response: ", err)
		}
	}
}

func writeResponse(writer http.ResponseWriter, resp endpointResponse) error {
	writer.WriteHeader(resp.StatusCode)

	body, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = writer.Write(body)

	return err
}
