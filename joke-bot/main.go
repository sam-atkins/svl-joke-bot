package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router := router()
	adapter := chiadapter.New(router)

	return adapter.ProxyWithContext(context.Background(), request)
}

func router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", healthCheckView)

	return router
}
