package beater

//var regex = (?P<time>\d*\.\d*) (?P<hostname>\S*) (?P<process>\S*) (?P<message>.*)

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/deedubs/pipebeat/config"
)

// Pipebeat - struct for Pipebeat functionality
type Pipebeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// New - Creates new Pipebeat
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Pipebeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

// Run - Start Pipebeat
func (bt *Pipebeat) Run(b *beat.Beat) error {
	logp.Info("pipebeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"message":    scanner.Text(),
		}

		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}

	return scanner.Err()
}

// Stop - Stop Pipebeat
func (bt *Pipebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
