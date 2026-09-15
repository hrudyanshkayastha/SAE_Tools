package kubearmor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"
)

type Consumer struct {
	LogPath string
}

func NewConsumer(logPath string) *Consumer {
	return &Consumer{LogPath: logPath}
}

func (c *Consumer) Start(ctx context.Context, output chan<- []byte) error {
	file, err := os.Open(c.LogPath)
	if err != nil {
		return fmt.Errorf("failed to open kubearmor log: %w", err)
	}
	defer file.Close()

	file.Seek(0, 2)
	reader := bufio.NewReader(file)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			if len(line) > 0 {
				output <- line
			}
		}
	}
}
