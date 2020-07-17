package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router := Router()
	adapter := chiadapter.New(router)

	return adapter.ProxyWithContext(context.Background(), request)
}

// Router provides routing to the API endpoints
func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", healthCheckView)

	return router
}

// HealthCheck struct
type HealthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthCheckView provides an endpoint to check the service is running
func healthCheckView(w http.ResponseWriter, r *http.Request) {
	data := HealthCheck{
		Status:  "OK",
		Message: "Healthy",
	}
	applyResponseJSON(200, data, w)
}

func applyResponseJSON(code int, data interface{}, w http.ResponseWriter) {
	if code != 200 {
		w.WriteHeader(code)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		se := ServerError{}
		se.ApplyServerError(err, "Server error", w)
		return
	}
	w.Write(jsonData)
	w.Header().Set("Content-Type", "application/json")
}

// ServerError struct
type ServerError struct {
	// do something like errors.New("some error message") and remove ErrorMessage string?
	Err          error  `json:"error"`
	ErrorMessage string `json:"message"`
}

// ApplyServerError handles 5xx errors for dev only.
// TODO: #6 Refactor for Prod so it logs errors instead of passing on full stack trace to the
// client
func (serverError *ServerError) ApplyServerError(err error, errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	serverErr := &ServerError{
		Err:          err,
		ErrorMessage: errorMessage,
	}
	jsonData, err := json.Marshal(serverErr)
	w.Write(jsonData)
}
