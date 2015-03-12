package ec2query

import (
	"encoding/xml"
	"io"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/internal/protocol/xml/xmlutil"
)

func Unmarshal(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() {
		err := xmlutil.UnmarshalXML(r.Data, xml.NewDecoder(r.HTTPResponse.Body))
		if err != nil && err != io.EOF {
			r.Error = err
			return
		}
	}
}

func UnmarshalMeta(r *aws.Request) {
	// TODO implement unmarshaling of request IDs
}

type xmlErrorResponse struct {
	XMLName   xml.Name `xml:"Response"`
	Code      string   `xml:"Errors>Error>Code"`
	Message   string   `xml:"Errors>Error>Message"`
	RequestID string   `xml:"RequestId"`
}

func UnmarshalError(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()

	resp := &xmlErrorResponse{}
	err := xml.NewDecoder(r.HTTPResponse.Body).Decode(resp)
	if err != nil && err != io.EOF {
		r.Error = err
	} else {
		apiErr := r.Error.(aws.APIError)
		apiErr.Code = resp.Code
		apiErr.Message = resp.Message
		r.Error = apiErr
	}
}
