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
	Content []string          `json:"content"`
	Meta    map[string]string `json:"meta"`
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
	value    string
	validErr error
}

// Valid error
func (rm ResourceMeta) Valid() error {
	return rm.validErr
}

// NewResourceMetaIP for ip meta
func NewResourceMetaIP(ipNet string) ResourceMeta {
	_, ip, err := net.ParseCIDR(ipNet)
	if err != nil {
		return ResourceMeta{validErr: fmt.Errorf("ip cidr: %w", err)}
	}
	return ResourceMeta{
		name:  "ip",
		value: ip.Network(),
	}
}

// NewResourceMetaAsn for asn meta
func NewResourceMetaAsn(asn string) ResourceMeta {
	_, err := strconv.ParseUint(asn, 10, 64)
	if err != nil {
		return ResourceMeta{validErr: fmt.Errorf("asn uint: %w", err)}
	}
	return ResourceMeta{
		name:  "asn",
		value: fmt.Sprint(asn),
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

	return ResourceMeta{
		name:  "latlong",
		value: fmt.Sprintf("%s,%s", parts[0], parts[1]),
	}
}

// NewResourceMetaNotes for notes meta
func NewResourceMetaNotes(note string) ResourceMeta {
	return ResourceMeta{
		name:  "notes",
		value: note,
	}
}

// NewResourceMetaCountries for Countries meta
func NewResourceMetaCountries(country string) ResourceMeta {
	return ResourceMeta{
		name:  "countries",
		value: country,
	}
}

// NewResourceMetaContinents for continents meta
func NewResourceMetaContinents(continent string) ResourceMeta {
	return ResourceMeta{
		name:  "continents",
		value: continent,
	}
}

// NewResourceMetaFallback for fallback meta
func NewResourceMetaFallback() ResourceMeta {
	return ResourceMeta{
		name:  "fallback",
		value: "true",
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
		r.Meta = map[string]string{}
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
