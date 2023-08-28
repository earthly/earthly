package base

import (
	"math/rand"
	"time"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/conslogging"
	"github.com/pkg/errors"
)

func ShowSatelliteLoading(console conslogging.ConsoleLogger, satName string, out chan cloud.SatelliteStatusUpdate) error {
	loadingMsgs := getSatelliteLoadingMessages()
	var (
		loggedSleep      bool
		loggedStop       bool
		loggedStart      bool
		loggedUpdating   bool
		loggedOffline    bool
		loggedDestroying bool
		loggedCreating   bool
		shouldLogLoading bool
	)
	for o := range out {
		if o.Err != nil {
			return errors.Wrap(o.Err, "failed processing satellite status")
		}
		shouldLogLoading = true
		switch o.State {
		case cloud.SatelliteStatusSleep:
			if !loggedSleep {
				console.Printf("%s is waking up. Please wait...", satName)
				loggedSleep = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStopping:
			if !loggedStop {
				console.Printf("%s is currently falling asleep. Waiting to send wake up signal...", satName)
				loggedStop = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStarting:
			if !loggedStart && !loggedSleep {
				console.Printf("%s is starting. Please wait...", satName)
				loggedStart = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusUpdating:
			if !loggedUpdating {
				console.Printf("%s is updating. It may take a few minutes to be ready...", satName)
				loggedUpdating = true
			}
		case cloud.SatelliteStatusDestroying:
			if !loggedDestroying {
				console.Printf("%s is going offline. It may take a few minutes to be ready...", satName)
				loggedDestroying = true
			}
		case cloud.SatelliteStatusOffline:
			if !loggedOffline {
				console.Printf("%s is coming online. Please wait...", satName)
				loggedOffline = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusCreating:
			if !loggedCreating {
				console.Printf("%s is creating. Please wait...", satName)
				loggedCreating = true
			}
		case cloud.SatelliteStatusOperational:
			if loggedSleep || loggedStop || loggedStart || loggedUpdating {
				// Satellite was in a different state previously but is now online
				console.Printf("...System online.")
			}
			shouldLogLoading = false
		default:
			// In case there's a new state later which we didn't expect here,
			// we'll still try to inform the user as best we can.
			// Note the state might just be "Unknown" if it maps to an gRPC enum we don't know about.
			console.Printf("%s state is: %s", satName, o)
			shouldLogLoading = false
		}
		if shouldLogLoading {
			var msg string
			msg, loadingMsgs = nextSatelliteLoadingMessage(loadingMsgs)
			console.Printf("...%s...", msg)
		}
	}
	return nil
}

func nextSatelliteLoadingMessage(msgs []string) (nextMsg string, remainingMsgs []string) {
	if len(msgs) == 0 {
		msgs = getSatelliteLoadingMessages()
	}
	return msgs[0], msgs[1:]
}

func getSatelliteLoadingMessages() []string {
	baseMessages := []string{
		"tracking orbit",
		"adjusting course",
		"deploying solar array",
		"aligning solar panels",
		"calibrating guidance system",
		"establishing transponder uplink",
		"testing signal quality",
		"fueling thrusters",
		"amplifying transmission signal",
		"checking thermal controls",
		"stabilizing trajectory",
		"contacting mission control",
		"testing antennas",
		"reporting fuel levels",
		"scanning surroundings",
		"avoiding debris",
		"taking solar reading",
		"reporting thermal conditions",
		"testing system integrity",
		"checking battery levels",
		"calibrating transponders",
		"modifying downlink frequency",
		"reticulating splines",
		"perturbing matrices",
		"synthesizing gravity",
		"iterating cellular automata",
	}
	msgs := baseMessages
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(msgs), func(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] })
	return msgs
}
