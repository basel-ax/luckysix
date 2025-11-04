package amqprpc

import (
	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter creates a new AMQP RPC router.
func NewRouter(services *service.Services) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, services.Translation)
		newTasksRoutes(routes, services.Tasks)
	}

	return routes
}
