package nodered

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
)

type BadRequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("%s. %s", e.Code, e.Message)
}

var ErrNotFound = errors.New("resource not found")

type ServerError struct {
	Err error
}

func (e *ServerError) Error() string {
	return e.Err.Error()
}

func (e *ServerError) Unwrap() error {
	return e.Err
}

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
	api     *resty.Client
	BaseURL string
}

func IsTab(v string) bool {
	return v == "tab"
}

func NewClient(baseURL string) *Client {
	c := &Client{
		api: resty.NewWithClient(http.DefaultClient),
	}
	c.api.Debug = false
	c.api.EnableTrace()
	c.api.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.StatusCode() > 399 || r.StatusCode() < 200 {
			var err error
			switch r.StatusCode() {
			case http.StatusNotFound:
				err = ErrNotFound
			case http.StatusBadRequest:
				badRequestErr := &BadRequestError{}
				if marshalErr := json.Unmarshal(r.Body(), &badRequestErr); marshalErr == nil {
					err = badRequestErr
				} else {
					err = fmt.Errorf("invalid request. %w", err)
				}

			default:
				err = fmt.Errorf("api error")
			}
			return &ServerError{
				Err: err,
			}
		}
		return nil
	})
	c.api.
		SetBaseURL(baseURL).
		SetHeader("Node-RED-API-Version", "v2").
		SetHeader("Content-Type", "application/json")

	return c
}

func (c *Client) SetBaseURL(u string) *Client {
	c.api.SetBaseURL(u)
	return c
}

func (c *Client) SetDebug(v bool) *Client {
	c.api.Debug = v
	return c
}

//
// Flows
//

func (c *Client) GetFlows() ([]Flow, error) {
	resp, err := c.api.R().Get("flows")
	if err != nil {
		return nil, err
	}

	node := gjson.ParseBytes(resp.Body())
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

	// Replace variables
	// TODO: Does this make sense?
	// host.containers.internal
	// host.docker.internal
	// tedge
	// body = bytes.ReplaceAll(body, []byte("${TEDGE_MQTT_HOST}"), []byte("host.docker.internal"))
	// body = bytes.ReplaceAll(body, []byte("${TEDGE_MQTT_PORT}"), []byte("1883"))

	data := &FlowResponseV2{}
	_, err := c.api.R().
		SetHeader("Node-RED-Deployment-Type", "full").
		SetResult(&data).
		SetBody(requestBody).
		Post("flows")

	return data, err
}

// Delete an existing flow
// Docs: https://nodered.org/docs/api/admin/methods/delete/flow/
func (c *Client) DeleteFlow(flowID string) error {
	_, err := c.api.R().Delete("flow/" + flowID)
	return err
}

//
// Projects
//

type ProjectStatus struct {
	Files    map[string]any    `json:"files,omitempty"`
	Commits  map[string]any    `json:"commits,omitempty"`
	Branches map[string]string `json:"branches,omitempty"`
}

type Commit struct {
	Sha     string `json:"sha,omitempty"`
	Subject string `json:"subject,omitempty"`
}

type BranchStatus struct {
	Ahead  int64 `json:"ahead"`
	Behind int64 `json:"behind"`
}

type Branch struct {
	Name    string        `json:"name,omitempty"`
	Remote  string        `json:"remote,omitempty"`
	Status  *BranchStatus `json:"status,omitempty"`
	Commit  *Commit       `json:"commit,omitempty"`
	Current bool          `json:"current"`
}

type Branches struct {
	Branches []Branch `json:"branches,omitempty"`
}

type ProjectsResponse struct {
	Projects []string `json:"projects"`
	Active   string   `json:"active"`
}

type Repository struct {
	URL      string `json:"url,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GitConfig struct {
	Remotes  map[string]Repository `json:"remotes,omitempty"`
	Branches map[string]string     `json:"branches,omitempty"`
}

type Project struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`

	Git               *GitConfig `json:"git,omitempty"`
	CredentialsSecret string     `json:"credentialSecret"`
}

func (c *Client) ProjectList() (*ProjectsResponse, error) {
	data := &ProjectsResponse{}
	_, err := c.api.R().SetResult(data).Get("projects")
	return data, err
}

func (c *Client) ProjectGet(name string) (*Project, error) {
	data := &Project{}
	_, err := c.api.R().SetResult(data).Get("projects/" + name)
	return data, err
}

func (c *Client) ProjectDelete(name string) error {
	_, err := c.api.R().Delete("projects/" + name)
	return err
}

func (c *Client) ProjectClone(name string, url string) (*Project, error) {
	project := Project{
		Name: name,
		Git: &GitConfig{
			Remotes: map[string]Repository{
				"origin": {
					URL: url,
				},
			},
		},
	}

	data := &Project{}
	_, err := c.api.R().
		SetBody(project).
		SetResult(data).
		Post("projects")
	return data, err
}

func (c *Client) ProjectPull(name string) (*Project, error) {
	data := &Project{}
	_, err := c.api.R().
		SetBody("{}").
		SetResult(data).
		Post("projects/" + name + "/pull")
	return data, err
}

func (c *Client) ProjectStatus(name string, clearContext bool) (*ProjectStatus, error) {
	data := &ProjectStatus{}
	_, err := c.api.R().
		SetResult(data).
		Get("projects/" + name + "/status")
	return data, err
}

func (c *Client) ProjectSetActive(name string, clearContext bool) (*Project, error) {
	data := &Project{}
	_, err := c.api.R().
		SetBody(map[string]any{
			"active":       true,
			"clearContext": clearContext,
		}).
		SetResult(data).
		Put("projects/" + name)
	return data, err
}

func (c *Client) ProjectBranches(name string) (*Branches, error) {
	data := &Branches{}
	_, err := c.api.R().SetResult(data).Get("projects/" + name + "/branches")
	return data, err
}
