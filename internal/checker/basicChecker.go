package checker

import (
	"context"
	"time"
	"vip_patroni/internal/config"
	log "vip_patroni/internal/logging"
)

func InitChecker(ctx context.Context, conf *config.Config, out chan<- bool) error {
checkLoop:
	for {
		patroni_status, err := GetPatroniStatus(conf.PatroniURL, conf.PatroniTimeoutMillis)
		if err != nil {
			log.Error("no response from request (pkg->checker): %s", err)
			if ctx.Err() != nil {
				break checkLoop
			}
			out <- false
			time.Sleep(time.Duration(conf.Interval) * time.Millisecond)
			continue
		}

		state := GetRole(patroni_status.Role)
		select {
		case <-ctx.Done():
			break checkLoop
		case out <- state:
			log.Debug("Check Patroni API")
			time.Sleep(time.Duration(conf.Interval) * time.Millisecond)
			continue
		}
	}
	return ctx.Err()
}

func GetRole(role string) bool {
	return role == "master"
}
