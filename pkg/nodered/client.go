package nodered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

type FlowResponseV2 struct {
	Flows       any    `json:"flows"`
	Rev         string `json:"rev,omitempty"`
	Credentials any    `json:"credentials,omitempty"`
}

type FlowEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Flow struct {
	ID       string    `json:"id"`
	Label    string    `json:"label"`
	Disabled bool      `json:"disabled"`
	Info     string    `json:"info"`
	Env      []FlowEnv `json:"env"`
	Revision string    `json:"-"`
}

func (f Flow) GetName() string {
	value := f.Label
	for _, item := range f.Env {
		if item.Name == "MODULE_NAME" {
			value = item.Value
			break
		}
	}
	return value
}

func (f Flow) GetVersion() string {
	value := f.Revision[0:8]
	for _, item := range f.Env {
		if item.Name == "MODULE_VERSION" {
			value = item.Value
			break
		}
	}
	return value
}

type Client struct {
	api     *http.Client
	BaseURL string
}

func IsTab(v string) bool {
	return v == "tab"
}

func NewClient(baseURL string) *Client {
	c := &Client{
		api:     http.DefaultClient,
		BaseURL: baseURL,
	}
	return c
}

func (c *Client) SetBaseURL(u string) *Client {
	c.BaseURL = u
	return c
}

func (c *Client) url(subpath ...string) string {
	return strings.Join(append([]string{c.BaseURL}, subpath...), "/")
}

func (c *Client) prepare(method string, subpath string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.url(subpath), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Node-RED-API-Version", "v2")
	return req, nil
}

func (c *Client) GetFlows() ([]Flow, error) {
	req, err := c.prepare(http.MethodGet, "flows", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.api.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	node := gjson.ParseBytes(b)
	flows := make([]Flow, 0)
	rev := node.Get("rev").String()
	node.Get("flows").ForEach(func(key, value gjson.Result) bool {
		if IsTab(value.Get("type").String()) {
			f := Flow{
				Revision: rev,
			}
			if err := json.Unmarshal([]byte(value.Raw), &f); err != nil {
				return true
			}
			flows = append(flows, f)
		}
		return true
	})
	return flows, nil
}

// Set new flows
// Docs: https://nodered.org/docs/api/admin/methods/post/flows/
func (c *Client) SetFlow(rev string, flowIn any) (*FlowResponseV2, error) {
	requestBody := &FlowResponseV2{
		Flows: flowIn,
		Rev:   rev,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Replace variables
	// TODO: Does this make sense?
	// host.containers.internal
	// host.docker.internal
	// tedge
	// body = bytes.ReplaceAll(body, []byte("${TEDGE_MQTT_HOST}"), []byte("host.docker.internal"))
	// body = bytes.ReplaceAll(body, []byte("${TEDGE_MQTT_PORT}"), []byte("1883"))

	req, err := c.prepare(http.MethodPost, "flows", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// full, nodes, flows or reload
	req.Header.Add("Node-RED-Deployment-Type", "full")
	req.Header.Add("Content-Type", "application/json")

	slog.Info("Sending request.", "body", body)
	slog.Info("Sending request.", "req", req)

	resp, err := c.api.Do(req)
	if err != nil {
		return nil, err
	}
	slog.Info("response.", "statusCode", resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &FlowResponseV2{}
	err = json.Unmarshal(b, &data)
	return data, err
}

// Delete an existing flow
// Docs: https://nodered.org/docs/api/admin/methods/delete/flow/
func (c *Client) DeleteFlow(flowID string) (*http.Response, error) {
	req, err := c.prepare(http.MethodDelete, fmt.Sprintf("flow/%s", flowID), nil)
	if err != nil {
		return nil, err
	}
	return c.api.Do(req)
}
