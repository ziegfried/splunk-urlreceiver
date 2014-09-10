package modinputs

import "encoding/xml"

type EndpointArg struct {
	Name             string `xml:"name,attr"`
	Title            string `xml:"title,omitempty"`
	Description      string `xml:"description,omitempty"`
	DataType         string `xml:"data_type,omitempty"`
	Validation       string `xml:"validation,omitempty"`
	RequiredOnEdit   bool   `xml:"required_on_edit"`
	RequiredOnCreate bool   `xml:"required_on_create"`
}

type Scheme struct {
	XMLName xml.Name `xml:"scheme"`

	Title                 string `xml:"title"`
	Description           string `xml:"description"`
	UseExternalValidation bool   `xml:"use_external_validation"`
	StreamingMode         string `xml:"streaming_mode"`
	UseSingleInstance     bool   `xml:"use_single_instance"`

	EndpointArguments []EndpointArg `xml:"endpoint>args>arg"`
}
