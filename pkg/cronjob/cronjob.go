package cronjob

import (
	"github.com/recative/recative-backend-sdk/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

type Cron interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	AddJob(spec string, cmd cron.Job) (cron.EntryID, error)
	// Run in current thread
	Run()
	// Start in goroutine
	Start()
	Raw() *cron.Cron
}

type cron_ struct {
	*cron.Cron
}

func (c *cron_) Raw() *cron.Cron {
	return c.Cron
}

var _ Cron = &cron_{}

type Config struct {
	LocationString string `env:"CRON_LOCATION" envDefault:"Asia/Shanghai"`
}

func New(config Config) Cron {
	location, err := time.LoadLocation(config.LocationString)
	if err != nil {
		logger.Panic("cron failed to parse location string", zap.Error(err))
	}

	c := cron.New(
		cron.WithLocation(location),
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
	)

	return &cron_{c}
}
