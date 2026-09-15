package wazuh

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"
)

// Consumer tails the Wazuh alerts.json file and passes raw JSON to the mapper.
// In a full production SAE environment, this would read from a Unix socket
// or Kafka topic for high throughput and isolation.
type Consumer struct {
	AlertLogPath string
}

func NewConsumer(path string) *Consumer {
	return &Consumer{
		AlertLogPath: path,
	}
}

// Start opens the file and tails it for new lines, passing mapped OCSF findings to a channel.
func (c *Consumer) Start(ctx context.Context, outChan chan<- []byte) error {
	file, err := os.Open(c.AlertLogPath)
	if err != nil {
		return fmt.Errorf("failed to open wazuh alerts log: %w", err)
	}
	defer file.Close()

	// Seek to end of file to act like 'tail -f'
	file.Seek(0, 2)

	reader := bufio.NewReader(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				// EOF, wait and retry
				time.Sleep(1 * time.Second)
				continue
			}
			
			if len(line) > 0 {
				outChan <- line
			}
		}
	}
}
