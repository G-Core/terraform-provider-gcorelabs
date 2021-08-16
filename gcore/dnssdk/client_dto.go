package dnssdk

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ListZones dto to read list of zones from API
type ListZones struct {
	Zones []Zone `json:"zones"`
}

// Zone dto to read info from API
type Zone struct {
	Name    string       `json:"name"`
	Records []ZoneRecord `json:"records"`
}

// AddZone dto to create new zone
type AddZone struct {
	Name string `json:"name"`
}

// CreateResponse dto to create new zone
type CreateResponse struct {
	ID    uint64 `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// RRSet dto as part of zone info from API
type RRSet struct {
	TTL     int               `json:"ttl"`
	Records []ResourceRecords `json:"resource_records"`
}

// ResourceRecords dto describe records in RRSet
type ResourceRecords struct {
	Content []string               `json:"content"`
	Meta    map[string]interface{} `json:"meta"`
}

type RecordType interface {
	ToContent() []string
}

// RecordTypeMX as type of record
type RecordTypeMX string

// ToContent convertor
func (mx RecordTypeMX) ToContent() []string {
	parts := strings.Split(string(mx), " ")
	content := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		content = append(content, p)
	}
	if len(content) != 2 {
		return nil
	}
	return content
}

// RecordTypeCAA as type of record
type RecordTypeCAA string

// ToContent convertor
func (caa RecordTypeCAA) ToContent() []string {
	parts := strings.Split(string(caa), " ")
	content := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		content = append(content, p)
	}
	if len(content) != 3 {
		return nil
	}
	return content
}

// RecordTypeAny as type of record
type RecordTypeAny string

// ToContent convertor
func (any RecordTypeAny) ToContent() []string {
	return []string{string(any)}
}

// ToRecordType builder
func ToRecordType(rType, content string) RecordType {
	switch strings.ToLower(rType) {
	case "mx":
		return RecordTypeMX(content)
	case "caa":
		return RecordTypeCAA(content)
	}
	return RecordTypeAny(content)
}

func ContentFromValue(recordType, content string) []string {
	rt := ToRecordType(recordType, content)
	if rt == nil {
		return nil
	}
	return rt.ToContent()
}

// ResourceMeta for ResourceRecords
type ResourceMeta struct {
	name     string
	value    interface{}
	validErr error
}

// Valid error
func (rm ResourceMeta) Valid() error {
	return rm.validErr
}

// NewResourceMetaIP for ip meta
func NewResourceMetaIP(ips ...string) ResourceMeta {
	for _, v := range ips {
		ip := net.ParseIP(v)
		if ip == nil {
			return ResourceMeta{validErr: fmt.Errorf("empty ip")}
		}
	}
	return ResourceMeta{
		name:  "ip",
		value: ips,
	}
}

// NewResourceMetaAsn for asn meta
func NewResourceMetaAsn(asn ...uint64) ResourceMeta {
	return ResourceMeta{
		name:  "asn",
		value: asn,
	}
}

// NewResourceMetaLatLong for lat long meta
func NewResourceMetaLatLong(latlong string) ResourceMeta {
	latlong = strings.TrimLeft(latlong, "(")
	latlong = strings.TrimLeft(latlong, "[")
	latlong = strings.TrimLeft(latlong, "{")
	latlong = strings.TrimRight(latlong, ")")
	latlong = strings.TrimRight(latlong, "]")
	latlong = strings.TrimRight(latlong, "}")
	parts := strings.Split(strings.ReplaceAll(latlong, " ", ""), ",")
	if len(parts) != 2 {
		return ResourceMeta{validErr: fmt.Errorf("latlong invalid format")}
	}
	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return ResourceMeta{validErr: fmt.Errorf("lat is invalid: %w", err)}
	}
	long, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return ResourceMeta{validErr: fmt.Errorf("long is invalid: %w", err)}
	}

	return ResourceMeta{
		name:  "latlong",
		value: []float64{lat, long},
	}
}

// NewResourceMetaNotes for notes meta
func NewResourceMetaNotes(notes ...string) ResourceMeta {
	return ResourceMeta{
		name:  "notes",
		value: notes,
	}
}

// NewResourceMetaCountries for Countries meta
func NewResourceMetaCountries(countries ...string) ResourceMeta {
	return ResourceMeta{
		name:  "countries",
		value: countries,
	}
}

// NewResourceMetaContinents for continents meta
func NewResourceMetaContinents(continents ...string) ResourceMeta {
	return ResourceMeta{
		name:  "continents",
		value: continents,
	}
}

// NewResourceMetaDefault for default meta
func NewResourceMetaDefault() ResourceMeta {
	return ResourceMeta{
		name:  "default",
		value: true,
	}
}

// SetContent to ResourceRecords
func (r *ResourceRecords) SetContent(recordType, val string) *ResourceRecords {
	r.Content = ContentFromValue(recordType, val)
	return r
}

// AddMeta to ResourceRecords
func (r *ResourceRecords) AddMeta(meta ResourceMeta) *ResourceRecords {
	if meta.name == "" || meta.value == "" {
		return r
	}
	if r.Meta == nil {
		r.Meta = map[string]interface{}{}
	}
	r.Meta[meta.name] = meta.value
	return r
}

// ZoneRecord dto describe records in Zone
type ZoneRecord struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	TTL          uint     `json:"ttl"`
	ShortAnswers []string `json:"short_answers"`
}

// APIError customization for API calls
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error,omitempty"`
}

// Error implementation
func (a APIError) Error() string {
	return fmt.Sprintf("%d: %s", a.StatusCode, a.Message)
}
