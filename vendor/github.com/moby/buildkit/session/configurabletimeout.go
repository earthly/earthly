package session

import (
	"context"
	"time"

	"github.com/moby/buildkit/util/bklog"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func configurableMonitorHealth(ctx context.Context, cc *grpc.ClientConn, cancelConn func(), healthCfg ManagerHealthCfg) {
	defer cancelConn()
	defer cc.Close()

	ticker := time.NewTicker(healthCfg.frequency)
	defer ticker.Stop()
	healthClient := grpc_health_v1.NewHealthClient(cc)

	consecutiveFailures := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			timeoutStart := time.Now().UTC()

			ctx, cancel := context.WithTimeout(ctx, healthCfg.timeout)
			_, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
			cancel()

			logFields := logrus.Fields{
				"timeout":        healthCfg.timeout,
				"actualDuration": time.Since(timeoutStart),
			}

			if err != nil {
				consecutiveFailures++

				logFields["allowedFailures"] = healthCfg.allowedFailures
				logFields["consecutiveFailures"] = consecutiveFailures
				bklog.G(ctx).WithFields(logFields).Warn("healthcheck failed")

				if consecutiveFailures >= healthCfg.allowedFailures {
					bklog.G(ctx).Error("healthcheck failed too many times")
					return
				}
			} else {
				bklog.G(ctx).WithFields(logFields).Debug("healthcheck completed")
				consecutiveFailures = 0
			}
		}
	}
}
