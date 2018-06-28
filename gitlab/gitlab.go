package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://gitlab.com/api/v4/"
	userAgent      = "alfred-gitlab-workflow"
)

// Client for GitLab API
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
	Token     string
	Parallel  int
}

// Project represents a GitLab project
type Project struct {
	ID                int       `json:"id"`
	Description       string    `json:"description"`
	WebURL            string    `json:"web_url"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	CreatedAt         time.Time `json:"created_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	CreatorID         int       `json:"creator_id"`
	Archived          bool      `json:"archived"`
	AvatarURL         string    `json:"avatar_url"`
}

// NewClient constructor for Client
func NewClient(httpClient *http.Client, endpointURL, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	var baseURL *url.URL
	if endpointURL == "" {
		baseURL, _ = url.Parse(defaultBaseURL)
	} else {
		baseURL, _ = url.Parse(endpointURL)
	}

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Token:     token,
		Parallel:  5,
	}
	return c
}

// NewRequest create a GitLab API request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Private-Token", c.Token)
	req.WithContext(ctx)
	return req, nil
}

// GetProjects get GitLab project list
func (c *Client) GetProjects() ([]*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := c.NewRequest(ctx, "HEAD", "projects?per_page=100")
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	xTotalPages := res.Header.Get("X-Total-Pages")
	totalPages, err := strconv.Atoi(xTotalPages)
	if err != nil {
		return nil, err
	}

	if totalPages > 30 {
		xTotal := res.Header.Get("X-Total")
		return nil, fmt.Errorf("too many projects: %s", xTotal)
	}

	limit := make(chan struct{}, c.Parallel)

	var projects []*Project
	var wg sync.WaitGroup
	var m sync.Mutex

	for i := 1; i < totalPages+1; i++ {
		wg.Add(1)
		page := i
		go func() {
			limit <- struct{}{}
			defer wg.Done()
			req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("projects?per_page=100&page=%d", page))
			if err != nil {
				<-limit
				return
			}
			res, err := c.client.Do(req)
			if err != nil {
				<-limit
				return
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				<-limit
				return
			}
			defer res.Body.Close()
			var p []*Project
			json.Unmarshal(body, &p)
			m.Lock()
			projects = append(projects, p...)
			m.Unlock()
			<-limit
		}()
	}
	wg.Wait()

	return projects, nil
}
