package workflow

import (
	"encoding/json"

	"github.com/shiimaxx/alfred-gitlab-workflow/gitlab"
)

type Output struct {
	Items []Item `json:"items"`
}

type Item struct {
	UID          string `json:"uid"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Arg          string `json:"arg"`
	Autocomplete string `json:"autocomplete"`
	Icon         struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"icon"`
}

func Run(endpointURL string) string {
	c := gitlab.NewClient(nil, endpointURL, "")
	projects, err := c.GetProjects()
	if err != nil {
		return err.Error()
	}

	var Items []Item
	for _, p := range projects {
		item := Item{
			UID:      p.Name,
			Title:    p.NameWithNamespace,
			Subtitle: p.WebURL,
			Arg:      p.WebURL,
		}
		Items = append(Items, item)
	}

	output, err := json.Marshal(Output{Items: Items})
	if err != nil {
		return err.Error()
	}

	return string(output)
}
