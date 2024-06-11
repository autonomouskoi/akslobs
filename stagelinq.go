package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/icedream/go-stagelinq"
)

const (
	appName    = "AKSLOBS"
	appVersion = "0.0.1"
	timeout    = 5 * time.Second
)

var stateValues = []string{
	stagelinq.EngineDeck1.Play(),
	stagelinq.EngineDeck1.PlayState(),
	stagelinq.EngineDeck1.PlayStatePath(),
	stagelinq.EngineDeck1.TrackArtistName(),
	stagelinq.EngineDeck1.TrackTrackNetworkPath(),
	stagelinq.EngineDeck1.TrackSongLoaded(),
	stagelinq.EngineDeck1.TrackSongName(),
	stagelinq.EngineDeck1.TrackTrackData(),
	stagelinq.EngineDeck1.TrackTrackName(),

	stagelinq.EngineDeck2.Play(),
	stagelinq.EngineDeck2.PlayState(),
	stagelinq.EngineDeck2.PlayStatePath(),
	stagelinq.EngineDeck2.TrackArtistName(),
	stagelinq.EngineDeck2.TrackTrackNetworkPath(),
	stagelinq.EngineDeck2.TrackSongLoaded(),
	stagelinq.EngineDeck2.TrackSongName(),
	stagelinq.EngineDeck2.TrackTrackData(),
	stagelinq.EngineDeck2.TrackTrackName(),

	stagelinq.EngineDeck3.Play(),
	stagelinq.EngineDeck3.PlayState(),
	stagelinq.EngineDeck3.PlayStatePath(),
	stagelinq.EngineDeck3.TrackArtistName(),
	stagelinq.EngineDeck3.TrackTrackNetworkPath(),
	stagelinq.EngineDeck3.TrackSongLoaded(),
	stagelinq.EngineDeck3.TrackSongName(),
	stagelinq.EngineDeck3.TrackTrackData(),
	stagelinq.EngineDeck3.TrackTrackName(),

	stagelinq.EngineDeck4.Play(),
	stagelinq.EngineDeck4.PlayState(),
	stagelinq.EngineDeck4.PlayStatePath(),
	stagelinq.EngineDeck4.TrackArtistName(),
	stagelinq.EngineDeck4.TrackTrackNetworkPath(),
	stagelinq.EngineDeck4.TrackSongLoaded(),
	stagelinq.EngineDeck4.TrackSongName(),
	stagelinq.EngineDeck4.TrackTrackData(),
	stagelinq.EngineDeck4.TrackTrackName(),
}

func makeStateMap() map[string]bool {
	retval := map[string]bool{}
	for _, value := range stateValues {
		retval[value] = false
	}
	return retval
}

func allStateValuesReceived(v map[string]bool) bool {
	for _, value := range v {
		if !value {
			return false
		}
	}
	return true
}

func maim(ctx context.Context, handle func(*stagelinq.State)) error {
	listener, err := stagelinq.ListenWithConfiguration(&stagelinq.ListenerConfiguration{
		DiscoveryTimeout: timeout,
		SoftwareName:     appName,
		SoftwareVersion:  appVersion,
		Name:             "testing",
	})
	if err != nil {
		return fmt.Errorf("listening: %w", err)
	}
	defer listener.Close()

	listener.AnnounceEvery(time.Second)

	deadline := time.After(timeout)
	foundDevices := []*stagelinq.Device{}

	log.Printf("Listening for StagelinQ devices for %s", timeout)

discoveryLoop:
	for {
		select {
		case <-ctx.Done():
			break discoveryLoop
		case <-deadline:
			break discoveryLoop
		default:
			device, deviceState, err := listener.Discover(timeout)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
				continue discoveryLoop
			}
			if device == nil {
				continue
			}
			// ignore device leaving messages since we do a one-off list
			if deviceState != stagelinq.DevicePresent {
				continue discoveryLoop
			}
			// check if we already found this device before
			for _, foundDevice := range foundDevices {
				if foundDevice.IsEqual(device) {
					continue discoveryLoop
				}
			}
			foundDevices = append(foundDevices, device)
			if err := handleDevice(ctx, device, listener.Token(), handle); err != nil {
				log.Printf("error handling device: %v", err)
				continue
			}
		}
	}

	log.Printf("Found devices: %d", len(foundDevices))
	return nil
}

func handleDevice(ctx context.Context, device *stagelinq.Device, token stagelinq.Token, handle func(*stagelinq.State)) error {
	log.Printf("%s %q %q %q", device.IP.String(), device.Name, device.SoftwareName, device.SoftwareVersion)

	// discover provided services
	log.Println("\tattempting to connect to this device…")
	deviceConn, err := device.Connect(token, []*stagelinq.Service{})
	if err != nil {
		return fmt.Errorf("connecting to device: %w", err)
	}
	defer deviceConn.Close()
	log.Println("\trequesting device data services…")
	services, err := deviceConn.RequestServices()
	if err != nil {
		return fmt.Errorf("requesting services: %w", err)
	}

	for _, service := range services {
		log.Printf("\toffers %s at port %d", service.Name, service.Port)
		if service.Name == "StateMap" {
			smh := stateMapHandler{
				device:  device,
				service: service,
				token:   token,
			}
			if err := smh.handleStateMap(ctx, handle); err != nil {
				log.Print("WARNING: handling StateMap: ", err)
			}
		}
	}

	log.Println("\tend of list of device data services")
	return nil
}

type stateMapHandler struct {
	device  *stagelinq.Device
	service *stagelinq.Service
	token   stagelinq.Token
}

func (smh stateMapHandler) handleStateMap(ctx context.Context, handle func(*stagelinq.State)) error {
	stateMapTCPConn, err := smh.device.Dial(smh.service.Port)
	if err != nil {
		return fmt.Errorf("creating stateMapTCPConn: %w", err)
	}
	defer stateMapTCPConn.Close()
	stateMapConn, err := stagelinq.NewStateMapConnection(stateMapTCPConn, smh.token)
	if err != nil {
		return fmt.Errorf("creating stateMapConn: %w", err)
	}

	m := makeStateMap()
	for _, stateValue := range stateValues {
		stateMapConn.Subscribe(stateValue)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-stateMapConn.ErrorC():
			return fmt.Errorf("in state map connection: %w", err)
		case state := <-stateMapConn.StateC():
			handle(state)
			m[state.Name] = true
			if allStateValuesReceived(m) {
				return nil
			}
		}
	}
}
