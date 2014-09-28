package modinputs

import "encoding/xml"
import "fmt"
import "os"
import "strings"

type EndpointArg struct {
	Name             string   `xml:"name,attr"`
	Title            string   `xml:"title,omitempty"`
	Description      string   `xml:"description,omitempty"`
	DataType         string   `xml:"data_type,omitempty"`
	EnumValues       []string `xml:"-"`
	Default          string   `xml:"-"`
	Validation       string   `xml:"validation,omitempty"`
	RequiredOnEdit   bool     `xml:"required_on_edit"`
	RequiredOnCreate bool     `xml:"required_on_create"`
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

func PrintScheme(scheme Scheme) {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(scheme); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	fmt.Println()
}

func PrintSpec(scheme Scheme, name string) {
	fmt.Printf("[%s://<name>]\n\n", name)

	for _, param := range scheme.EndpointArguments {
		placeholder := "<string>"

		switch param.DataType {
		case "boolean":
			placeholder = "[true|false]"
		case "integer":
			placeholder = "<integer>"
		}

		if len(param.EnumValues) > 0 {
			placeholder = fmt.Sprintf("[%s]", strings.Join(param.EnumValues, "|"))
		}

		fmt.Printf("%s = %s\n", param.Name, placeholder)

		fmt.Print("* ")
		if param.Title != "" {
			fmt.Printf("%s ", param.Title)
		}
		if param.RequiredOnCreate {
			fmt.Print("(required")
		} else {
			fmt.Print("(optional")
		}

		if param.Default != "" {
			fmt.Printf(", default: %s", param.Default)
		}
		fmt.Println(")")

		if param.Description != "" {
			for _, line := range strings.Split(param.Description, "\n") {
				fmt.Printf("* %s\n", line)
			}
		}

		fmt.Println()
	}
}
