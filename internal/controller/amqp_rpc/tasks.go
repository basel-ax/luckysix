package amqprpc

import (
	"github.com/basel-ax/luckysix/internal/service"
	"github.com/basel-ax/luckysix/pkg/rabbitmq/rmq_rpc/server"
)

type tasksRoutes struct {
	tasks service.Tasks
}

func newTasksRoutes(routes map[string]server.CallHandler, t service.Tasks) {
	r := &tasksRoutes{t}
	{
		//routes["checkProfit"] = r.checkProfit()
		routes["CheckRabbitTask"] = r.CheckRabbitTask()
	}
}

type rabbitResponse struct {
	answer string `json:"answer"`
}

// not ended, obsolete
/*
func (r *tasksRoutes) checkProfit() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		pools, err := r.tasksService.CheckProfit(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - tasksRoutes - checkProfit - r.tasksService.CheckProfit: %w", err)
		}
		response := poolResponse{pools}

		return response, nil
	}
}
*/

// simple RabbitMQ task, with logging
func (r *tasksRoutes) CheckRabbitTask() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		answer := r.tasks.CheckRabbitTask()

		response := rabbitResponse{answer}

		return response, nil
	}
}
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
