package types

import (
	"encoding/json"
	"fmt"
)

type StackResourceStatus string
type StackResourceAction string

const (
	StackResourceStatusCompete    StackResourceStatus = "COMPLETE"
	StackResourceStatusFailed     StackResourceStatus = "FAILED"
	StackResourceStatusInProgress StackResourceStatus = "IN_PROGRESS"
	StackResourceActionAdopt      StackResourceAction = "ADOPT"
	StackResourceActionCheck      StackResourceAction = "CHECK"
	StackResourceActionCreate     StackResourceAction = "CREATE"
	StackResourceActionDelete     StackResourceAction = "DELETE"
	StackResourceActionInit       StackResourceAction = "INIT"
	StackResourceActionRestore    StackResourceAction = "RESTORE"
	StackResourceActionResume     StackResourceAction = "RESUME"
	StackResourceActionRollback   StackResourceAction = "ROLLBACK"
	StackResourceActionSnapshot   StackResourceAction = "SNAPSHOT"
	StackResourceActionSuspend    StackResourceAction = "SUSPEND"
	StackResourceActionUpdated    StackResourceAction = "UPDATE"
)

func (srs StackResourceStatus) IsValid() error {
	switch srs {
	case StackResourceStatusCompete,
		StackResourceStatusFailed,
		StackResourceStatusInProgress:
		return nil
	}
	return fmt.Errorf("invalid StackResourceStatus type: %v", srs)
}

func (srs StackResourceStatus) ValidOrNil() (*StackResourceStatus, error) {
	if srs.String() == "" {
		return nil, nil
	}
	err := srs.IsValid()
	if err != nil {
		return &srs, err
	}
	return &srs, nil
}

func (srs StackResourceStatus) String() string {
	return string(srs)
}

func (srs StackResourceStatus) List() []StackResourceStatus {
	return []StackResourceStatus{
		StackResourceStatusCompete,
		StackResourceStatusFailed,
		StackResourceStatusInProgress,
	}
}

func (srs StackResourceStatus) StringList() []string {
	var s []string
	for _, v := range srs.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (srs *StackResourceStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := StackResourceStatus(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*srs = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (srs *StackResourceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(srs.String())
}

func (sra StackResourceAction) IsValid() error {
	switch sra {
	case StackResourceActionAdopt,
		StackResourceActionCheck,
		StackResourceActionCreate,
		StackResourceActionDelete,
		StackResourceActionInit,
		StackResourceActionRestore,
		StackResourceActionResume,
		StackResourceActionRollback,
		StackResourceActionSnapshot,
		StackResourceActionSuspend,
		StackResourceActionUpdated:
		return nil
	}
	return fmt.Errorf("invalid StackResourceAction type: %v", sra)
}

func (sra StackResourceAction) ValidOrNil() (*StackResourceAction, error) {
	if sra.String() == "" {
		return nil, nil
	}
	err := sra.IsValid()
	if err != nil {
		return &sra, err
	}
	return &sra, nil
}

func (sra StackResourceAction) String() string {
	return string(sra)
}

func (sra StackResourceAction) List() []StackResourceAction {
	return []StackResourceAction{
		StackResourceActionAdopt,
		StackResourceActionCheck,
		StackResourceActionCreate,
		StackResourceActionDelete,
		StackResourceActionInit,
		StackResourceActionRestore,
		StackResourceActionResume,
		StackResourceActionRollback,
		StackResourceActionSnapshot,
		StackResourceActionSuspend,
		StackResourceActionUpdated,
	}
}

func (sra StackResourceAction) StringList() []string {
	var s []string
	for _, v := range sra.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (sra *StackResourceAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := StackResourceAction(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*sra = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (sra *StackResourceAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(sra.String())
}
