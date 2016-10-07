package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

/*Kapacitors kapacitors

swagger:model Kapacitors
*/
type Kapacitors struct {

	/* kapacitors

	Required: true
	*/
	Kapacitors []*Kapacitor `json:"kapacitors"`
}

// Validate validates this kapacitors
func (m *Kapacitors) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateKapacitors(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Kapacitors) validateKapacitors(formats strfmt.Registry) error {

	if err := validate.Required("kapacitors", "body", m.Kapacitors); err != nil {
		return err
	}

	for i := 0; i < len(m.Kapacitors); i++ {

		if swag.IsZero(m.Kapacitors[i]) { // not required
			continue
		}

		if m.Kapacitors[i] != nil {

			if err := m.Kapacitors[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}