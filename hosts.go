package datadog

import (
	"fmt"
	"net/url"
)

type HostActionResp struct {
	Action   string `json:"action"`
	Hostname string `json:"hostname"`
	Message  string `json:"message,omitempty"`
}

type HostActionMute struct {
	Message  *string `json:"message,omitempty"`
	EndTime  *string `json:"end,omitempty"`
	Override *bool   `json:"override,omitempty"`
}

// MuteHost mutes all monitors for the given host
func (client *Client) MuteHost(host string, action *HostActionMute) (*HostActionResp, error) {
	var out HostActionResp
	uri := "/v1/host/" + host + "/mute"
	if err := client.doJsonRequest("POST", uri, action, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnmuteHost unmutes all monitors for the given host
func (client *Client) UnmuteHost(host string) (*HostActionResp, error) {
	var out HostActionResp
	uri := "/v1/host/" + host + "/unmute"
	if err := client.doJsonRequest("POST", uri, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HostTotalsResp defines response to GET /v1/hosts/totals.
type HostTotalsResp struct {
	TotalUp     *int `json:"total_up"`
	TotalActive *int `json:"total_active"`
}

// GetHostTotals returns number of total active hosts and total up hosts.
// Active means the host has reported in the past hour, and up means it has reported in the past two hours.
func (client *Client) GetHostTotals() (*HostTotalsResp, error) {
	var out HostTotalsResp
	uri := "/v1/hosts/totals"
	if err := client.doJsonRequest("GET", uri, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// QueryHosts defines response to GET /v1/hosts/.
type HostSearchResp struct {
	TotalReturned int             `json:"total_returned"`
	TotalMatching int             `json:"total_matching"`
	HostListEntry []HostListEntry `json:"host_list"`
}

type HostListEntry struct {
	LastReportedTime int                 `json:"last_reported_time"`
	Name             string              `json:"name"`
	IsMuted          bool                `json:"is_muted"`
	MuteTimeout      int                 `json:"mute_timeout"`
	Apps             []string            `json:"apps"`
	TagsBySource     map[string][]string `json:"tags_by_source"`
	Up               bool                `json:"up"`
	Metrics          map[string]float32  `json:"metrics"`
	Source           []string            `json:"sources"`
	HostName         string              `json:"host_name"`
	Id               int                 `json:"id"`
	Aliases          []string            `json:"aliases"`
}

//Search among hosts live within the past 2 hours. Max 100 results at a time.
func (client *Client) QueryHosts(filter, sort_field, sort_dir string, start, count, from int) (HostSearchResp, error) {
	// Since this is a GET request, we need to build a query string.
	vals := url.Values{}

	if filter != "" {
		vals.Add("filter", filter)
	}
	if sort_field != "" {
		vals.Add("priority", sort_field)
	}
	if sort_dir != "" {
		vals.Add("sources", sort_dir)
	}
	if start != 0 {
		vals.Add("start", string(start))
	}
	if count != 100 {
		vals.Add("count", string(count))
	}
	if from != 0 {
		vals.Add("from", string(from))
	}

	// Now the request and response.
	var out HostSearchResp
	if err := client.doJsonRequest("GET",
		fmt.Sprintf("/v1/hosts?%s", vals.Encode()), nil, &out); err != nil {
		return HostSearchResp{}, err
	}
	return out, nil

}
