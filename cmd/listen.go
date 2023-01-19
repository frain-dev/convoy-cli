package main

import (
	"fmt"
	"net/url"
	"time"

	convoyCli "github.com/frain-dev/convoy-cli"
	"github.com/frain-dev/convoy-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addListenCommand() *cobra.Command {
	var since string
	var source string
	// var events string
	var forwardTo string

	cmd := &cobra.Command{
		Use:   "listen",
		Short: "Starts a websocket client that listens to events streamed by the server",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := convoyCli.LoadConfig()
			if err != nil {
				log.Fatal("Error loading config file:", err)
			}

			if util.IsStringEmpty(forwardTo) {
				log.Fatal("flag forward-to cannot be empty")
			}

			// TODO(all): enforce the source filter?
			// if util.IsStringEmpty(source) {
			// 	log.Fatal("flag source cannot be empty")
			// }

			hostInfo, err := url.Parse(c.Host)
			if err != nil {
				log.Fatal("Error parsing host URL: ", err)
			}

			if !util.IsStringEmpty(since) {
				var sinceTime time.Time
				sinceTime, err = time.Parse(time.RFC3339, since)
				if err != nil {
					log.WithError(err).Error("since is not a valid timestamp, will try time duration")

					dur, err := time.ParseDuration(since)
					if err != nil {
						log.WithError(err).Fatal("since is neither a valid time duration or timestamp, see the listen command help menu for a valid since value")
					} else {
						since = fmt.Sprintf("since|duration|%v", since)
						sinceTime = time.Now().Add(-dur)
					}
				} else {
					since = fmt.Sprintf("since|timestamp|%v", since)
				}

				log.Printf("will resend all discarded events after: %v", sinceTime)
			}

			listenRequest := convoyCli.ListenRequest{
				HostName:  c.Host,
				DeviceID:  c.ActiveDeviceID,
				SourceID:  source,
				Since:     since,
				ForwardTo: forwardTo,
			}

			l := convoyCli.NewListener(c)
			l.Listen(&listenRequest, hostInfo)
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "The source id of the source you want to receive events from (only applies to incoming projects)")
	cmd.Flags().StringVar(&since, "since", "", "Send discarded events since a timestamp (e.g. 2013-01-02T13:23:37Z) or relative time (e.g. 42m for 42 minutes)")
	cmd.Flags().StringVar(&forwardTo, "forward-to", "", "The host/web server you want to forward events to")
	// cmd.Flags().StringVar(&events, "events", "*", "Events types")

	return cmd
}
