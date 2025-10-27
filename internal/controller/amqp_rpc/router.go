package amqprpc

import (
	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/server"
)

// NewTranslationRouter -.
func NewTranslationRouter(t service.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}

// NewTasksRouter -.
func NewTasksRouter(t service.Tasks) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTasksRoutes(routes, t)
	}

	return routes
}
