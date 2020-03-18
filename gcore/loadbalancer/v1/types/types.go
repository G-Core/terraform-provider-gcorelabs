package types

import (
	"encoding/json"
	"fmt"
)

type ProvisioningStatus string
type OperatingStatus string
type LoadBalancerAlgorithm string
type PersistenceType string
type ProtocolType string
type HealthMonitorType string
type HTTPMethod string

type ItemID struct {
	ID string `json:"id"`
}

const (
	ProvisioningStatusActive              ProvisioningStatus    = "ACTIVE"
	ProvisioningStatusDeleted             ProvisioningStatus    = "DELETED"
	ProvisioningStatusError               ProvisioningStatus    = "ERROR"
	ProvisioningStatusPendingCreate       ProvisioningStatus    = "PENDING_CREATE"
	ProvisioningStatusPendingUpdate       ProvisioningStatus    = "PENDING_UPDATE"
	ProvisioningStatusPendingDelete       ProvisioningStatus    = "PENDING_DELETE"
	OperatingStatusOnline                 OperatingStatus       = "ONLINE"
	OperatingStatusDraining               OperatingStatus       = "DRAINING"
	OperatingStatusOffline                OperatingStatus       = "OFFLINE"
	OperatingStatusDegraded               OperatingStatus       = "DEGRADED"
	OperatingStatusOperatingError         OperatingStatus       = "ERROR"
	OperatingStatusNoMonitor              OperatingStatus       = "NO_MONITOR"
	LoadBalancerAlgorithmRoundRobin       LoadBalancerAlgorithm = "ROUND_ROBIN"
	LoadBalancerAlgorithmLeastConnections LoadBalancerAlgorithm = "LEAST_CONNECTIONS"
	LoadBalancerAlgorithmSourceIP         LoadBalancerAlgorithm = "SOURCE_IP"
	LoadBalancerAlgorithmSourceIPPort     LoadBalancerAlgorithm = "SOURCE_IP_PORT"
	PersistenceTypeAppCookie              PersistenceType       = "APP_COOKIE"
	PersistenceTypeHTTPCookie             PersistenceType       = "HTTP_COOKIE"
	PersistenceTypeSourceIP               PersistenceType       = "SOURCE_IP"
	ProtocolTypeHTTP                      ProtocolType          = "HTTP"
	ProtocolTypeHTTPS                     ProtocolType          = "HTTPS"
	ProtocolTypeTCP                       ProtocolType          = "TCP"
	ProtocolTypeTerminatedHTTPS           ProtocolType          = "TERMINATED_HTTPS"
	ProtocolTypeUDP                       ProtocolType          = "UDP"
	HealthMonitorTypeHTTP                 HealthMonitorType     = "HTTP"
	HealthMonitorTypeHTTPS                HealthMonitorType     = "HTTPS"
	HealthMonitorTypePING                 HealthMonitorType     = "PING"
	HealthMonitorTypeTCP                  HealthMonitorType     = "TCP"
	HealthMonitorTypeTLSHello             HealthMonitorType     = "TLS-HELLO"
	HealthMonitorTypeUDPConnect           HealthMonitorType     = "UDP-CONNECT"
	HTTPMethodCONNECT                     HTTPMethod            = "CONNECT"
	HTTPMethodDELETE                      HTTPMethod            = "DELETE"
	HTTPMethodGET                         HTTPMethod            = "GET"
	HTTPMethodHEAD                        HTTPMethod            = "HEAD"
	HTTPMethodOPTIONS                     HTTPMethod            = "OPTIONS"
	HTTPMethodPATCH                       HTTPMethod            = "PATCH"
	HTTPMethodPOST                        HTTPMethod            = "POST"
	HTTPMethodPUT                         HTTPMethod            = "PUT"
	HTTPMethodTRACE                       HTTPMethod            = "TRACE"
)

func (ps ProvisioningStatus) IsValid() error {
	switch ps {
	case ProvisioningStatusActive,
		ProvisioningStatusDeleted,
		ProvisioningStatusError,
		ProvisioningStatusPendingCreate,
		ProvisioningStatusPendingDelete,
		ProvisioningStatusPendingUpdate:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", ps)
}

func (ps ProvisioningStatus) String() string {
	return string(ps)
}

func (ps ProvisioningStatus) List() []ProvisioningStatus {
	return []ProvisioningStatus{
		ProvisioningStatusActive,
		ProvisioningStatusDeleted,
		ProvisioningStatusError,
		ProvisioningStatusPendingCreate,
		ProvisioningStatusPendingDelete,
		ProvisioningStatusPendingUpdate,
	}
}

func (ps ProvisioningStatus) StringList() []string {
	var s []string
	for _, v := range ps.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (ps *ProvisioningStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ProvisioningStatus(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*ps = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (ps *ProvisioningStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

func (os OperatingStatus) IsValid() error {
	switch os {
	case OperatingStatusOnline,
		OperatingStatusDraining,
		OperatingStatusOperatingError,
		OperatingStatusOffline,
		OperatingStatusDegraded,
		OperatingStatusNoMonitor:
		return nil
	}
	return fmt.Errorf("invalid OperatingStatus type: %v", os)
}

func (os OperatingStatus) String() string {
	return string(os)
}

func (os OperatingStatus) List() []OperatingStatus {
	return []OperatingStatus{
		OperatingStatusOnline,
		OperatingStatusDraining,
		OperatingStatusOperatingError,
		OperatingStatusOffline,
		OperatingStatusDegraded,
		OperatingStatusNoMonitor,
	}
}

func (os OperatingStatus) StringList() []string {
	var s []string
	for _, v := range os.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (os *OperatingStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := OperatingStatus(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*os = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (os *OperatingStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

func (lba LoadBalancerAlgorithm) IsValid() error {
	switch lba {
	case LoadBalancerAlgorithmRoundRobin,
		LoadBalancerAlgorithmLeastConnections,
		LoadBalancerAlgorithmSourceIP,
		LoadBalancerAlgorithmSourceIPPort:
		return nil
	}
	return fmt.Errorf("invalid LoadBalancerAlgorithm type: %v", lba)
}

func (lba LoadBalancerAlgorithm) String() string {
	return string(lba)
}

func (lba LoadBalancerAlgorithm) List() []LoadBalancerAlgorithm {
	return []LoadBalancerAlgorithm{
		LoadBalancerAlgorithmRoundRobin,
		LoadBalancerAlgorithmLeastConnections,
		LoadBalancerAlgorithmSourceIP,
		LoadBalancerAlgorithmSourceIPPort,
	}
}

func (lba LoadBalancerAlgorithm) StringList() []string {
	var s []string
	for _, v := range lba.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (lba *LoadBalancerAlgorithm) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := LoadBalancerAlgorithm(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*lba = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (lba *LoadBalancerAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(lba.String())
}

func (lbspt PersistenceType) IsValid() error {
	switch lbspt {
	case PersistenceTypeAppCookie,
		PersistenceTypeHTTPCookie,
		PersistenceTypeSourceIP:
		return nil
	}
	return fmt.Errorf("invalid PersistenceType type: %v", lbspt)
}

func (lbspt PersistenceType) String() string {
	return string(lbspt)
}

func (lbspt PersistenceType) List() []PersistenceType {
	return []PersistenceType{
		PersistenceTypeAppCookie,
		PersistenceTypeHTTPCookie,
		PersistenceTypeSourceIP,
	}
}

func (lbspt PersistenceType) StringList() []string {
	var s []string
	for _, v := range lbspt.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (lbspt *PersistenceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := PersistenceType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*lbspt = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (lbspt *PersistenceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(lbspt.String())
}

func (pt ProtocolType) IsValid() error {
	switch pt {
	case ProtocolTypeHTTP, ProtocolTypeHTTPS, ProtocolTypeTCP, ProtocolTypeTerminatedHTTPS, ProtocolTypeUDP:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", pt)
}

func (pt ProtocolType) String() string {
	return string(pt)
}

func (pt ProtocolType) List() []ProtocolType {
	return []ProtocolType{ProtocolTypeHTTP, ProtocolTypeHTTPS, ProtocolTypeTCP, ProtocolTypeTerminatedHTTPS, ProtocolTypeUDP}
}

func (pt ProtocolType) StringList() []string {
	var s []string
	for _, v := range pt.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (pt *ProtocolType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ProtocolType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*pt = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (pt *ProtocolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

func (hm HealthMonitorType) IsValid() error {
	switch hm {
	case HealthMonitorTypeHTTP,
		HealthMonitorTypeHTTPS,
		HealthMonitorTypeTCP,
		HealthMonitorTypePING,
		HealthMonitorTypeTLSHello,
		HealthMonitorTypeUDPConnect:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", hm)
}

func (hm HealthMonitorType) String() string {
	return string(hm)
}

func (hm HealthMonitorType) List() []HealthMonitorType {
	return []HealthMonitorType{
		HealthMonitorTypeHTTP,
		HealthMonitorTypeHTTPS,
		HealthMonitorTypeTCP,
		HealthMonitorTypePING,
		HealthMonitorTypeTLSHello,
		HealthMonitorTypeUDPConnect,
	}
}

func (hm HealthMonitorType) StringList() []string {
	var s []string
	for _, v := range hm.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (hm *HealthMonitorType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := HealthMonitorType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*hm = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (hm *HealthMonitorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(hm.String())
}

func (m HTTPMethod) IsValid() error {
	switch m {
	case HTTPMethodCONNECT,
		HTTPMethodHEAD,
		HTTPMethodGET,
		HTTPMethodOPTIONS,
		HTTPMethodPOST,
		HTTPMethodPATCH,
		HTTPMethodPUT,
		HTTPMethodDELETE,
		HTTPMethodTRACE:
		return nil
	}
	return fmt.Errorf("invalid ProvisioningStatus type: %v", m)
}

func (m HTTPMethod) String() string {
	return string(m)
}

func (m HTTPMethod) List() []HTTPMethod {
	return []HTTPMethod{
		HTTPMethodCONNECT,
		HTTPMethodHEAD,
		HTTPMethodGET,
		HTTPMethodOPTIONS,
		HTTPMethodPOST,
		HTTPMethodPATCH,
		HTTPMethodPUT,
		HTTPMethodDELETE,
		HTTPMethodTRACE,
	}
}

func (m HTTPMethod) StringList() []string {
	var s []string
	for _, v := range m.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface
func (m *HTTPMethod) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := HTTPMethod(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*m = v
	return nil
}

// MarshalJSON - implements Marshaler interface
func (m *HTTPMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
