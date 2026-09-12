package garak

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

func NewConsumer(path string) *Consumer {
	return &Consumer{
		LogPath: path,
	}
}

func (c *Consumer) Start(ctx context.Context, outChan chan<- []byte) error {
	file, err := os.Open(c.LogPath)
	if err != nil {
		return fmt.Errorf("failed to open garak log: %w", err)
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
