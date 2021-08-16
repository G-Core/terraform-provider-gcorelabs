package dnssdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	defaultBaseURL = "https://dnsapi.gcorelabs.com"
	tokenHeader    = "APIKey"
	defaultTimeOut = 10 * time.Second
)

// Client for DNS API.
type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
	authHeader string
}

// ZonesFilter find zones
type ZonesFilter struct {
	Names []string
}

type authHeader string

// BearerAuth by header
func BearerAuth(token string) func() authHeader {
	return func() authHeader {
		return authHeader(fmt.Sprintf("Bearer %s", token))
	}
}

// PermanentAPIKeyAuth by header
func PermanentAPIKeyAuth(token string) func() authHeader {
	return func() authHeader {
		return authHeader(fmt.Sprintf("%s %s", tokenHeader, token))
	}
}

func (zf ZonesFilter) query() string {
	if len(zf.Names) == 0 {
		return ""
	}
	return url.Values{"name": zf.Names}.Encode()
}

// NewClient constructor of Client.
func NewClient(authorizer func() authHeader, opts ...func(*Client)) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	cl := &Client{
		authHeader: string(authorizer()),
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: defaultTimeOut},
	}
	for _, op := range opts {
		op(cl)
	}
	return cl
}

// CreateZone add new zone.
// https://dnsapi.gcorelabs.com/docs#operation/CreateZone
func (c *Client) CreateZone(ctx context.Context, name string) (uint64, error) {
	res := CreateResponse{}
	params := AddZone{Name: name}
	err := c.do(ctx, http.MethodPost, "/v2/zones", params, &res)
	if err != nil {
		return 0, fmt.Errorf("request: %w", err)
	}
	if res.Error != "" {
		return 0, APIError{StatusCode: http.StatusOK, Message: res.Error}
	}

	return res.ID, nil
}

// Zones gets all zones.
// https://dnsapi.gcorelabs.com/docs#operation/Zones
func (c *Client) Zones(ctx context.Context, filters ...func(zone *ZonesFilter)) ([]Zone, error) {
	res := ListZones{}
	filter := ZonesFilter{}
	for _, op := range filters {
		op(&filter)
	}
	err := c.do(ctx, http.MethodGet, "/v2/zones?"+filter.query(), nil, &res)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return res.Zones, nil
}

// ZonesWithRecords gets all zones with records information.
func (c *Client) ZonesWithRecords(ctx context.Context, filters ...func(zone *ZonesFilter)) ([]Zone, error) {
	zones, err := c.Zones(ctx, filters...)
	if err != nil {
		return nil, fmt.Errorf("all zones: %w", err)
	}
	gr, _ := errgroup.WithContext(ctx)
	for i, z := range zones {
		z := z
		i := i
		gr.Go(func() error {
			zone, errGet := c.Zone(ctx, z.Name)
			if errGet != nil {
				return fmt.Errorf("%s: %w", z.Name, errGet)
			}
			zones[i] = zone
			return nil
		})
	}
	err = gr.Wait()
	if err != nil {
		return nil, fmt.Errorf("zone info: %w", err)
	}

	return zones, nil
}

// DeleteZone gets zone information.
// https://dnsapi.gcorelabs.com/docs#operation/DeleteZone
func (c *Client) DeleteZone(ctx context.Context, name string) error {
	name = strings.Trim(name, ".")
	uri := path.Join("/v2/zones", name)

	err := c.do(ctx, http.MethodDelete, uri, nil, nil)
	if err != nil {
		return fmt.Errorf("request %s: %w", name, err)
	}

	return nil
}

// Zone gets zone information.
// https://dnsapi.gcorelabs.com/docs#operation/Zone
func (c *Client) Zone(ctx context.Context, name string) (Zone, error) {
	name = strings.Trim(name, ".")
	zone := Zone{}
	uri := path.Join("/v2/zones", name)

	err := c.do(ctx, http.MethodGet, uri, nil, &zone)
	if err != nil {
		return Zone{}, fmt.Errorf("get zone %s: %w", name, err)
	}

	return zone, nil
}

// RRSet gets RRSet item.
// https://dnsapi.gcorelabs.com/docs#operation/RRSet
func (c *Client) RRSet(ctx context.Context, zone, name, recordType string) (RRSet, error) {
	zone, name = strings.Trim(zone, "."), strings.Trim(name, ".")
	var result RRSet
	uri := path.Join("/v2/zones", zone, name, recordType)

	err := c.do(ctx, http.MethodGet, uri, nil, &result)
	if err != nil {
		return RRSet{}, fmt.Errorf("request %s -> %s: %w", zone, name, err)
	}

	return result, nil
}

// DeleteRRSet removes RRSet type records.
// https://dnsapi.gcorelabs.com/docs#operation/DeleteRRSet
func (c *Client) DeleteRRSet(ctx context.Context, zone, name, recordType string) error {
	zone, name = strings.Trim(zone, "."), strings.Trim(name, ".")
	uri := path.Join("/v2/zones", zone, name, recordType)

	err := c.do(ctx, http.MethodDelete, uri, nil, nil)
	if err != nil {
		// Support DELETE idempotence https://developer.mozilla.org/en-US/docs/Glossary/Idempotent
		statusErr := new(APIError)
		if errors.As(err, statusErr) && statusErr.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("delete record request: %w", err)
	}

	return nil
}

// DeleteRRSetRecord removes RRSet record.
func (c *Client) DeleteRRSetRecord(ctx context.Context, zone, name, recordType string, contents ...string) error {
	// get current records info
	rrSet, err := c.RRSet(ctx, zone, name, recordType)
	if err != nil {
		errAPI := new(APIError)
		if errors.As(err, errAPI) && errAPI.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("rrset: %w", err)
	}
	if len(rrSet.Records) == 0 {
		return nil
	}
	// setup new records
	newRecords := make([]ResourceRecords, 0, len(rrSet.Records))
	for _, record := range rrSet.Records {
		if len(record.Content) == 0 {
			continue
		}
		newVal := make([]string, 0, len(contents))
		for _, val := range record.Content {
			exist := false
			for _, content := range contents {
				if val == content {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			// keep existing content
			newVal = append(newVal, val)
		}
		if len(newVal) > 0 {
			newRecords = append(newRecords, ResourceRecords{Content: newVal})
		}
	}
	rrSet.Records = newRecords
	// delete on empty content
	if len(rrSet.Records) == 0 {
		err = c.DeleteRRSet(ctx, zone, name, recordType)
		if err != nil {
			err = fmt.Errorf("delete rrset: %w", err)
		}
		return err
	}
	// update with removing deleted content
	err = c.UpdateRRSet(ctx, zone, name, recordType, rrSet)
	if err != nil {
		err = fmt.Errorf("update rrset: %w", err)
	}
	return err
}

// AddZoneRRSet create or extend resource record.
func (c *Client) AddZoneRRSet(ctx context.Context,
	zone, recordName, recordType string,
	values []ResourceRecords, ttl int) error {

	record := RRSet{TTL: ttl, Records: values}

	records, err := c.RRSet(ctx, zone, recordName, recordType)
	if err == nil && len(records.Records) > 0 {
		record.Records = append(record.Records, records.Records...)
		return c.UpdateRRSet(ctx, zone, recordName, recordType, record)
	}

	return c.CreateRRSet(ctx, zone, recordName, recordType, record)
}

// CreateRRSet https://dnsapi.gcorelabs.com/docs#operation/CreateRRSet
func (c *Client) CreateRRSet(ctx context.Context, zone, name, recordType string, record RRSet) error {
	zone, name = strings.Trim(zone, "."), strings.Trim(name, ".")
	uri := path.Join("/v2/zones", zone, name, recordType)

	return c.do(ctx, http.MethodPost, uri, record, nil)
}

// UpdateRRSet https://dnsapi.gcorelabs.com/docs#operation/UpdateRRSet
func (c *Client) UpdateRRSet(ctx context.Context, zone, name, recordType string, record RRSet) error {
	zone, name = strings.Trim(zone, "."), strings.Trim(name, ".")
	uri := path.Join("/v2/zones", zone, name, recordType)

	return c.do(ctx, http.MethodPut, uri, record, nil)
}

func (c *Client) do(ctx context.Context, method, uri string, bodyParams interface{}, dest interface{}) error {
	var bs []byte
	if bodyParams != nil {
		var err error
		bs, err = json.Marshal(bodyParams)
		if err != nil {
			return fmt.Errorf("encode bodyParams: %w", err)
		}
	}

	endpoint, err := c.BaseURL.Parse(path.Join(c.BaseURL.Path, uri))
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	log.Printf("[DEBUG] dns api request: %s %s %s \n", method, uri, bs)

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), strings.NewReader(string(bs)))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authHeader)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusMultipleChoices {
		all, _ := ioutil.ReadAll(resp.Body)
		e := APIError{
			StatusCode: resp.StatusCode,
		}
		err := json.Unmarshal(all, &e)
		if err != nil {
			e.Message = string(all)
		}
		return e
	}

	if dest == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}
