package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sae-core/internal/api"
	"sae-core/internal/engine"
	"sae-core/internal/storage"
	"sae-core/internal/ueba"

	"sae-core/internal/cortex"
	"sae-core/internal/falco"
	"sae-core/internal/garak"
	"sae-core/internal/kubearmor"
	"sae-core/internal/scoutsuite"
	"sae-core/internal/shuffle"
	"sae-core/internal/suricata"
	"sae-core/internal/thehive"
	"sae-core/internal/trivy"
	"sae-core/internal/wazuh"
	"sae-core/internal/zeek"
	
	"github.com/google/uuid"
)

func main() {
	fmt.Println("Starting SAE Unified Core Platform...")

	// 1. Initialize Storage (Postgres & Redis)
	store, err := storage.InitStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.PG.Close()
	defer store.Redis.Close()
	fmt.Println("Connected to PostgreSQL and Redis successfully.")

	// 2. Initialize Correlation & AI Engine
	shuffleWebhook := "http://localhost:5001/api/v1/hooks/webhook_sae_action"
	if w := os.Getenv("SHUFFLE_WEBHOOK_URL"); w != "" {
		shuffleWebhook = w
	}
	thehiveURL := "http://localhost:9000"
	cortexURL := "http://localhost:9001"
	eng := engine.NewEngine(store, shuffleWebhook, thehiveURL, cortexURL)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go eng.Start(ctx)
	fmt.Println("Correlation Engine & AI Pipeline started.")

	// 2.5 Start UEBA Engine
	uebaEng := ueba.NewUEBAEngine(store)
	go uebaEng.Start(ctx, 1*time.Minute)
	fmt.Println("UEBA Behavioral Analytics Engine started.")

	// 3. Initialize REST API
	server := api.NewServer(store)
	go func() {
		fmt.Println("REST API listening on :8080")
		if err := server.Start(":8080"); err != nil {
			log.Printf("REST API stopped: %v", err)
		}
	}()

	// 4. Start Ingestion Adapters
	startIngestion(ctx, store)

	// Wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("\nShutting down SAE gracefully...")
	cancel()
	time.Sleep(1 * time.Second) // wait for cleanup
}

func startIngestion(ctx context.Context, store *storage.Storage) {
	// Wazuh
	wazuhChan := make(chan []byte, 100)
	go func() {
		wazuhConsumer := wazuh.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\engine\\logs\\alerts\\alerts.json")
		_ = wazuhConsumer.Start(ctx, wazuhChan)
	}()
	go func() {
		for raw := range wazuhChan {
			ocsf, err := wazuh.MapAlertToOCSF(raw)
			if err == nil {
				ocsf.EventID = uuid.New().String()
				store.PublishEvent(ctx, *ocsf)
			}
		}
	}()

	// Zeek
	zeekChan := make(chan []byte, 100)
	go func() {
		zeekConsumer := zeek.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\engine\\logs\\zeek\\conn.log")
		_ = zeekConsumer.Start(ctx, zeekChan)
	}()
	go func() {
		for raw := range zeekChan {
			ocsf, err := zeek.MapZeekConnToOCSF(raw)
			if err == nil {
				ocsf.EventID = uuid.New().String()
				store.PublishEvent(ctx, *ocsf)
			}
		}
	}()

	// Suricata
	suriChan := make(chan []byte, 100)
	go func() {
		suriConsumer := suricata.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\engine\\logs\\suricata\\eve.json")
		_ = suriConsumer.Start(ctx, suriChan)
	}()
	go func() {
		for raw := range suriChan {
			ocsf, err := suricata.MapEVEToOCSF(raw)
			if err == nil {
				ocsf.EventID = uuid.New().String()
				store.PublishEvent(ctx, *ocsf)
			}
		}
	}()

	// Trivy
	trivyChan := make(chan []byte, 100)
	go func() {
		trivyConsumer := trivy.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\trivy\\logs\\trivy_report.json")
		_ = trivyConsumer.Start(ctx, trivyChan)
	}()
	go func() {
		for raw := range trivyChan {
			ocsfs, err := trivy.MapReportToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()

	// Other tools (Falco, KubeArmor, ScoutSuite, Shuffle, TheHive, Cortex) would follow the identical pattern here.
	// (Shortened for test brevity since we're using Wazuh, Zeek, Suricata, and Trivy for the E2E correlation test)
	
	// Falco
	falcoChan := make(chan []byte, 100)
	go func() {
		falcoConsumer := falco.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\containers_security\\logs\\falco_events.json")
		_ = falcoConsumer.Start(ctx, falcoChan)
	}()
	go func() {
		for raw := range falcoChan {
			ocsf, err := falco.MapFalcoToOCSF(raw)
			if err == nil {
				ocsf.EventID = uuid.New().String()
				store.PublishEvent(ctx, *ocsf)
			}
		}
	}()
	
	// KubeArmor
	kaChan := make(chan []byte, 100)
	go func() {
		kaConsumer := kubearmor.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\cloudguard\\logs\\kubearmor_alerts.json")
		_ = kaConsumer.Start(ctx, kaChan)
	}()
	go func() {
		for raw := range kaChan {
			ocsf, err := kubearmor.MapAlertToOCSF(raw)
			if err == nil {
				ocsf.EventID = uuid.New().String()
				store.PublishEvent(ctx, *ocsf)
			}
		}
	}()
	
	// ScoutSuite
	scoutChan := make(chan []byte, 100)
	go func() {
		scoutConsumer := scoutsuite.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\scoutsuite\\logs\\scoutsuite_results.json")
		_ = scoutConsumer.Start(ctx, scoutChan)
	}()
	go func() {
		for raw := range scoutChan {
			ocsfs, err := scoutsuite.MapReportToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()
	
	// Shuffle
	shuffleChan := make(chan []byte, 100)
	go func() {
		shuffleConsumer := shuffle.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\shuffle\\logs\\shuffle_results.json")
		_ = shuffleConsumer.Start(ctx, shuffleChan)
	}()
	go func() {
		for raw := range shuffleChan {
			ocsfs, err := shuffle.MapResultToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()
	
	// TheHive
	thehiveChan := make(chan []byte, 100)
	go func() {
		thehiveConsumer := thehive.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\thehive\\logs\\thehive_results.json")
		_ = thehiveConsumer.Start(ctx, thehiveChan)
	}()
	go func() {
		for raw := range thehiveChan {
			ocsfs, err := thehive.MapCaseToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()
	
	// Cortex
	cortexChan := make(chan []byte, 100)
	go func() {
		cortexConsumer := cortex.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\cortex\\logs\\cortex_results.json")
		_ = cortexConsumer.Start(ctx, cortexChan)
	}()
	go func() {
		for raw := range cortexChan {
			ocsfs, err := cortex.MapJobToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()

	// Garak
	garakChan := make(chan []byte, 100)
	go func() {
		garakConsumer := garak.NewConsumer("E:\\New folder\\SAE_Tools\\SAE\\garak\\logs\\garak_report.jsonl")
		_ = garakConsumer.Start(ctx, garakChan)
	}()
	go func() {
		for raw := range garakChan {
			ocsfs, err := garak.ParseGarakEvalToOCSF(raw)
			if err == nil {
				for _, ocsf := range ocsfs {
					ocsf.EventID = uuid.New().String()
					store.PublishEvent(ctx, ocsf)
				}
			}
		}
	}()
}
