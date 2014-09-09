package modinputs

import "encoding/xml"

type InputConfig struct {
	XMLName       xml.Name      `xml:"input"`
	ServerHost    string        `xml:"server_host"`
	ServerUri     string        `xml:"server_uri"`
	SessionKey    string        `xml:"session_key"`
	CheckpointDir string        `xml:"checkpoint_dir"`
	Stanzas       []InputStanza `xml:"configuration>stanza"`
}

type InputStanza struct {
	Name     string  `xml:"name,attr"`
	Params   []Param `xml:"param"`
	paramMap map[string]Param
}

func (stanza InputStanza) GetParam(name string) *Param {
	p, found := stanza.ParamMap()[name]
	if found {
		return &p
	} else {
		return nil
	}
}

func (stanza InputStanza) GetParamValue(name string, defaultValue string) string {
	if param := stanza.GetParam(name); param != nil {
		return param.Value
	} else {
		return defaultValue
	}
}

func (stanza InputStanza) HasParam(name string) bool {
	return stanza.GetParam(name) != nil
}

func (stanza InputStanza) ParamMap() map[string]Param {
	if stanza.paramMap == nil {
		stanza.paramMap = make(map[string]Param)
		for _, param := range stanza.Params {
			stanza.paramMap[param.Name] = param
		}
	}
	return stanza.paramMap
}

func (stanza InputStanza) MakeEvent() Event {
	event := new(Event)
	event.Stanza = stanza.Name
	if p := stanza.GetParam("source"); p != nil {
		event.Source = p.Value
	}
	if p := stanza.GetParam("sourcetype"); p != nil {
		event.Sourcetype = p.Value
	}
	if p := stanza.GetParam("host"); p != nil {
		event.Host = p.Value
	}
	if p := stanza.GetParam("index"); p != nil {
		event.Index = p.Value
	}
	return *event
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
