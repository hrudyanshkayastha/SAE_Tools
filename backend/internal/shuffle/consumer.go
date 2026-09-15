package shuffle

import (
	"context"
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
	var lastMod time.Time

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			info, err := os.Stat(c.LogPath)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			if info.ModTime().After(lastMod) && info.Size() > 0 {
				lastMod = info.ModTime()
				
				data, err := os.ReadFile(c.LogPath)
				if err == nil && len(data) > 0 {
					output <- data
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
