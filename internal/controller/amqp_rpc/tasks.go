package amqprpc

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/server"
)

type tasksRoutes struct {
	tasks service.Tasks
}

func newTasksRoutes(routes map[string]server.CallHandler, t service.Tasks) {
	r := &tasksRoutes{t}
	routes["checkRabbit"] = r.checkRabbit()
	routes["generateLuckyTwo"] = r.generateLuckyTwo()
	routes["generateLuckyFive"] = r.generateLuckyFive()
	routes["generateLuckySix"] = r.generateLuckySix()
	routes["processWallets"] = r.processWallets()
}

type rabbitResponse struct {
	Message string `json:"message"`
}

func (r *tasksRoutes) checkRabbit() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		answer := r.tasks.CheckRabbitTask()
		return rabbitResponse{Message: answer}, nil
	}
}

func (r *tasksRoutes) generateLuckyTwo() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		r.tasks.GenerateLuckyTwo()
		return rabbitResponse{Message: "GenerateLuckyTwo task started in the background."}, nil
	}
}

func (r *tasksRoutes) generateLuckyFive() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		r.tasks.GenerateLuckyFive()
		return rabbitResponse{Message: "GenerateLuckyFive task started in the background."}, nil
	}
}

func (r *tasksRoutes) generateLuckySix() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		r.tasks.GenerateLuckySix()
		return rabbitResponse{Message: "GenerateLuckySix task started in the background."}, nil
	}
}

func (r *tasksRoutes) processWallets() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		r.tasks.ProcessWallets()
		return rabbitResponse{Message: "ProcessWallets task started in the background."}, nil
	}
}
