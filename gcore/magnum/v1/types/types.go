package types

import (
	"encoding/json"
	"fmt"
)

type NodegroupRole string
type ClusterUpdateOperation string

const (
	NodegroupMasterRole           NodegroupRole          = "master"
	NodegroupWorkerRole           NodegroupRole          = "worker"
	ClusterUpdateOperationRemove  ClusterUpdateOperation = "remove"
	ClusterUpdateOperationReplace ClusterUpdateOperation = "replace"
	ClusterUpdateOperationAdd     ClusterUpdateOperation = "add"
)

func (ng NodegroupRole) IsValid() error {
	switch ng {
	case NodegroupMasterRole,
		NodegroupWorkerRole:
		return nil
	}
	return fmt.Errorf("invalid NodegroupRole type: %v", ng)
}

func (ng NodegroupRole) ValidOrNil() (*NodegroupRole, error) {
	if ng.String() == "" {
		return nil, nil
	}
	err := ng.IsValid()
	if err != nil {
		return &ng, err
	}
	return &ng, nil
}

func (ng NodegroupRole) String() string {
	return string(ng)
}

func (ng NodegroupRole) List() []NodegroupRole {
	return []NodegroupRole{
		NodegroupMasterRole,
		NodegroupWorkerRole,
	}
}

func (ng NodegroupRole) StringList() []string {
	var s []string
	for _, v := range ng.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ng *NodegroupRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := NodegroupRole(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ng = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ng *NodegroupRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ng.String())
}

func (co ClusterUpdateOperation) IsValid() error {
	switch co {
	case ClusterUpdateOperationAdd,
		ClusterUpdateOperationRemove,
		ClusterUpdateOperationReplace:
		return nil
	}
	return fmt.Errorf("invalid ClusterUpdateOperation type: %v", co)
}

func (co ClusterUpdateOperation) ValidOrNil() (*ClusterUpdateOperation, error) {
	if co.String() == "" {
		return nil, nil
	}
	err := co.IsValid()
	if err != nil {
		return &co, err
	}
	return &co, nil
}

func (co ClusterUpdateOperation) String() string {
	return string(co)
}

func (co ClusterUpdateOperation) List() []ClusterUpdateOperation {
	return []ClusterUpdateOperation{
		ClusterUpdateOperationAdd,
		ClusterUpdateOperationRemove,
		ClusterUpdateOperationReplace,
	}
}

func (co ClusterUpdateOperation) StringList() []string {
	var s []string
	for _, v := range co.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (co *ClusterUpdateOperation) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ClusterUpdateOperation(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*co = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (co *ClusterUpdateOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(co.String())
}
