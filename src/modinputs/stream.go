package modinputs

import "fmt"
import "encoding/xml"
import "os"

type Event struct {
	XMLName xml.Name `xml:"event"`

	Data       string `xml:"data"`
	Stanza     string `xml:"stanza,attr,omitempty"`
	Time       string `xml:"time,omitempty"`
	Host       string `xml:"host,omitempty"`
	Source     string `xml:"source,omitempty"`
	Sourcetype string `xml:"sourcetype,omitempty"`
	Index      string `xml:"index,omitempty"`
}

type Stream struct {
	buffer  chan Event
	done    chan bool
	stopped chan bool
}

func NewStream(bufferSize int) *Stream {
	stream := new(Stream)
	stream.buffer = make(chan Event, bufferSize)
	stream.done = make(chan bool)
	stream.stopped = make(chan bool)
	return stream
}

func (stream Stream) Start() {
	fmt.Println("<stream>")
	go stream.streamEvents()
}

func (stream Stream) Send(event Event) {
	stream.buffer <- event
}

func (stream Stream) Stop() chan bool {
	stream.done <- true
	return stream.stopped
}

func (stream Stream) streamEvents() {
	enc := xml.NewEncoder(os.Stdout)
	active := true
	for active {
		select {
		case <-stream.done:
			active = false
			break
		case event := <-stream.buffer:
			if err := enc.Encode(event); err != nil {
				Log.Error("Error sending event: %s", err)
			}
			os.Stdout.Write([]byte("\n"))
		}
	}
	fmt.Println("</stream>")
	stream.stopped <- true
}
