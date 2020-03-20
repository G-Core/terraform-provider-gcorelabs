package types

import (
	"encoding/json"
	"fmt"
)

type AddressType string

type ItemName struct {
	Name string `json:"name"`
}

const (
	AddressTypeFixed    AddressType = "fixed"
	AddressTypeFloating AddressType = "floating"
)

func (at AddressType) IsValid() error {
	switch at {
	case AddressTypeFixed, AddressTypeFloating:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", at)
}

func (at AddressType) ValidOrNil() (*AddressType, error) {
	if at.String() == "" {
		return nil, nil
	}
	err := at.IsValid()
	if err != nil {
		return &at, err
	}
	return &at, nil
}

func (at AddressType) String() string {
	return string(at)
}

func (at AddressType) List() []AddressType {
	return []AddressType{
		AddressTypeFixed,
		AddressTypeFloating,
	}
}

func (at AddressType) StringList() []string {
	var s []string
	for _, v := range at.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (at *AddressType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := AddressType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*at = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (at *AddressType) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.String())
}
