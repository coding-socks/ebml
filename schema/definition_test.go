package schema

import (
	"encoding/xml"
	"testing"
)

func TestElementID_String(t *testing.T) {
	h, want := ElementID(0x81), "0x81"
	if got := h.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

func TestElementID_UnmarshalXMLAttr(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    ElementID
		wantErr bool
	}{
		{name: "invalid", data: []byte(`<foo id="a"></foo>`), wantErr: true},
		{name: "hex", data: []byte(`<foo id="0x81"></foo>`), want: ElementID(0x81)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foo := struct {
				ID ElementID `xml:"id,attr"`
			}{}
			if err := xml.Unmarshal(tt.data, &foo); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXMLAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := foo.ID; got != tt.want {
				t.Errorf("UnmarshalXMLAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
