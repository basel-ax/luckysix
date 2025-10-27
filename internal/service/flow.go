package service

import (
	"github.com/basel-ax/luckysix/config"
	log "github.com/basel-ax/luckysix/pkg/logger"
)

const (
// BaseURL              = "https://api-key.fusionbrain.ai/key/api/v1/"
)

type FlowService struct {
	cfg *config.Config
	log *log.Logger
}

// example
type Flow struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

func NewFlowService(cfg *config.Config, l *log.Logger) *FlowService {
	return &FlowService{
		cfg: cfg,
		log: l,
	}
}

func (service *FlowService) Generating() (*string, error) {
	/*
		db, err := gorm.Open(postgres.Open(cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}
	*/

	return nil, nil
}
