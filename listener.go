package convoy_cli

import (
	"encoding/json"
	"github.com/frain-dev/convoy-cli/net"
	"github.com/frain-dev/convoy-cli/util"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const (
	// Time allowed to write a message to the server.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the server.
	pongWait = 10 * time.Second

	// Send pings to server with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Listener struct {
	done      chan interface{} // Channel to indicate that the receiverHandler is done
	interrupt chan os.Signal   // Channel to listen for interrupt signal to terminate gracefully
	c         *Config
}

func NewListener(c *Config) *Listener {
	return &Listener{
		c:         c,
		done:      make(chan interface{}),
		interrupt: make(chan os.Signal),
	}
}

func (l *Listener) Listen(listenRequest *ListenRequest, hostInfo *url.URL) {
	signal.Notify(l.interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	body, err := json.Marshal(listenRequest)
	if err != nil {
		log.Fatal("Error marshalling json:", err)
	}

	url := url.URL{
		Scheme: "ws",
		Host:   hostInfo.Host,
		Path:   "/stream/listen",
	}

	conn, response, err := websocket.DefaultDialer.Dial(url.String(), http.Header{
		"Authorization": []string{"Bearer " + l.c.ActiveApiKey},
		"Body":          []string{string(body)},
	})

	if err != nil {
		if response != nil {
			buf, e := io.ReadAll(response.Body)
			if e != nil {
				log.Fatal("Error parsing request body", e)
			}
			defer response.Body.Close()
			log.Fatalln("websocket dialer failed with response: ", string(buf))
		}

		log.Fatal(err)
	}

	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return
	}

	defer conn.Close()

	if !util.IsStringEmpty(listenRequest.Since) {
		// Send a message to the server to resend unsuccessful events to the device
		err := conn.WriteMessage(websocket.TextMessage, []byte(listenRequest.Since))
		if err != nil {
			log.WithError(err).Errorln("an error occurred sending 'since' message")
		}
	}

	go l.HandleMessage(conn, listenRequest.ForwardTo)
	l.PingUntilInterrupt(conn)
}

func (l *Listener) PingUntilInterrupt(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-ticker.C:
			err := conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.WithError(err).Errorln("failed to set write deadline")
			}

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithError(err).Errorln("failed to set write ping message")
				return
			}

		case <-l.interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Send a message to set the device to offline
			err := conn.WriteMessage(websocket.TextMessage, []byte("disconnect"))
			if err != nil {
				log.WithError(err).Errorln("error during closing websocket")
				return
			}

			// Close our websocket connection
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.WithError(err).Errorln("error during closing websocket")
				return
			}

			select {
			case <-l.done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

func (l *Listener) HandleMessage(connection *websocket.Conn, url string) {
	defer close(l.done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				return
			}

			log.Error("an error occurred in the receive handler:", err)
			return
		}

		var event CLIEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			log.Error("an error occurred in unmarshalling json:", err)
			continue
		}

		// send request to the recipient
		d, err := net.NewDispatcher(time.Second*10, "")
		if err != nil {
			log.Error("an error occurred while forwarding the event", err)
			continue
		}

		res, err := d.ForwardCliEvent(url, http.MethodPost, event.Data, event.Headers)
		if err != nil {
			log.Error("an error occurred while forwarding the event", err)
			continue
		}

		// set the event delivery status to Success when we successfully forward the event
		ack := &AckEventDelivery{UID: event.UID}
		mb, err := json.Marshal(ack)
		if err != nil {
			log.Error("an error occurred in marshalling json:", err)
			continue
		}

		// write an ack message back to the connection here
		err = connection.WriteMessage(websocket.TextMessage, mb)
		if err != nil {
			log.Error("an error occurred while acknowledging the event", err)
		}

		log.Println(string(res.Body))
	}
}
