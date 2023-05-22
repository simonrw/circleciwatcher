package circleci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var baseUrl = &url.URL{Host: "circleci.com", Scheme: "https", Path: "api/v2/"}

func New(token string) *Client {
	httpClient := &http.Client{}
	return &Client{Token: token, httpClient: httpClient}
}

type Client struct {
	Token string

	httpClient *http.Client
}

type State = string

var (
	StateCreated      State = "created"
	StateErrored            = "errored"
	StateSetupPending       = "setup-pending"
	StateSetup              = "setup"
	StatePending            = "pending"
)

type Pipeline struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type Status = string

var (
	StatusSuccess      Status = "success"
	StatusRunning             = "running"
	StatusNotRun              = "not_run"
	StatusFailed              = "failed"
	StatusError               = "error"
	StatusFailing             = "failing"
	StatusOnHold              = "on_hold"
	StatusCanceled            = "canceled"
	StatusUnauthorized        = "unauthorized"
)

type Workflow struct {
	ID         string `json:"id"`
	Status     Status `json:"status"`
	CanceledBy string `json:"canceled_by"`
	ErroredBy  string `json:"errored_by"`
	Tag        string `json:"tag"`
	StartedBy  string `json:"started_by"`
	CreatedAt  string `json:"created_at"`
	StoppedAt  string `json:"stopped_at"`
}

func (c *Client) GetPipeline(scheme, owner, repo string, number int) (*Pipeline, error) {
	method := "GET"
	path := fmt.Sprintf("project/gh/%s/%s/pipeline/%d", owner, repo, number)
	p := Pipeline{}
	err := c.request(method, path, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) GetPipelineWorkflows(token string) ([]Workflow, error) {

	res := struct {
		Items []Workflow `json:"items"`
	}{}
	if err := c.request("GET", fmt.Sprintf("pipeline/%s/workflow", token), &res); err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) GetWorkflow(id string) (*Workflow, error) {
	w := Workflow{}
	if err := c.request("GET", fmt.Sprintf("workflow/%s", id), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *Client) request(method, path string, responseStruct interface{}) error {
	url := baseUrl.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		// error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading error response body: %w", err)
		}

		if len(body) > 0 {
			message := struct {
				Message string `json:"message"`
			}{}
			err = json.Unmarshal(body, &message)
			if err != nil {
				return fmt.Errorf("invalid JSON error body: %w", err)
			}
			return &Error{StatusCode: resp.StatusCode, Message: message.Message}
		}
	}

	if responseStruct != nil {
		err = json.NewDecoder(resp.Body).Decode(responseStruct)
		if err != nil {
			return fmt.Errorf("decoding body: %w", err)
		}
	}
	return nil
}

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("status code %d: %s", e.StatusCode, e.Message)
}
