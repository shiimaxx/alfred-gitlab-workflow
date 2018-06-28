package workflow

import (
	"encoding/json"
	"os"

	"bufio"
	"io/ioutil"

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

func makeErrorResponse(err error) string {
	var items []Item
	items = append(items, Item{
		UID:   "error",
		Title: err.Error(),
	})
	output, _ := json.Marshal(Output{Items: items})
	return string(output)
}

func Run(endpointURL, token string, refresh bool) string {
	if _, err := os.Stat("./projects"); !os.IsNotExist(err) && !refresh {
		output, err := ioutil.ReadFile("./projects")
		if err != nil {
			return makeErrorResponse(err)
		}
		return string(output)
	}

	c := gitlab.NewClient(nil, endpointURL, token)
	projects, err := c.GetProjects()
	if err != nil {
		return makeErrorResponse(err)
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
		return makeErrorResponse(err)
	}

	f, _ := os.OpenFile("./projects", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	defer f.Close()
	w := bufio.NewWriter(f)
	w.Write(output)

	return string(output)
}
