package volumes

import (
	"encoding/json"
	"fmt"
)

type VolumeSource string
type VolumeType string

const (
	NewVolume VolumeSource = "new-volume"
	Image     VolumeSource = "image"
	Snapshot  VolumeSource = "snapshot"
	Standard  VolumeType   = "standard"
	SsdHiIops VolumeType   = "ssd_hiiops"
	Cold      VolumeType   = "cold"
)

func (vs VolumeSource) IsValid() error {
	switch vs {
	case NewVolume, Image, Snapshot:
		return nil
	}
	return fmt.Errorf("invalid VolumeSource type: %v", vs)
}

func (vs VolumeSource) ValidOrNil() (*VolumeSource, error) {
	if vs.String() == "" {
		return nil, nil
	}
	err := vs.IsValid()
	if err != nil {
		return &vs, err
	}
	return &vs, nil
}

func (vs VolumeSource) String() string {
	return string(vs)
}

func (vs VolumeSource) List() []VolumeSource {
	return []VolumeSource{NewVolume, Image, Snapshot}
}

func (vs VolumeSource) StringList() []string {
	var s []string
	for _, v := range vs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for VolumeSource
func (vs *VolumeSource) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := VolumeSource(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*vs = v
	return nil
}

// MarshalJSON - implements Marshaler interface for VolumeSource
func (vs *VolumeSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(vs.String())
}

func (vt VolumeType) String() string {
	return string(vt)
}

func (vt VolumeType) List() []VolumeType {
	return []VolumeType{Standard, SsdHiIops, Cold}
}

func (vt VolumeType) StringList() []string {
	var s []string
	for _, v := range vt.List() {
		s = append(s, v.String())
	}
	return s
}

func (vt VolumeType) IsValid() error {
	switch vt {
	case Standard, Cold, SsdHiIops:
		return nil
	}
	return fmt.Errorf("invalid VolumeType type: %v", vt)
}

func (vt VolumeType) ValidOrNil() (*VolumeType, error) {
	if vt.String() == "" {
		return nil, nil
	}
	err := vt.IsValid()
	if err != nil {
		return &vt, err
	}
	return &vt, nil
}

// UnmarshalJSON - implements Unmarshaler interface for VolumeType
func (vt *VolumeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := VolumeType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*vt = v
	return nil
}

// MarshalJSON - implements Marshaler interface for VolumeType
func (vt *VolumeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(vt.String())
}
