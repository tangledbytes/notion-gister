package gister

import (
	"fmt"
	"time"

	"github.com/utkarsh-pro/notion-gister/pkg/utils"
)

type Item interface {
	ID() string
	Title() string
	URL() string
	Tags() []string
	CreatedTime() time.Time
}

type Opts struct {
	IgnoreRule IgnoreRule
}

type IgnoreRule struct {
	IgnoreTags []string
	IgnoreTime *time.Time
}

func Create(items []Item, opts Opts) string {
	var gist string = "<ol>\n"

OUTER:
	for _, item := range items {
		if opts.IgnoreRule.IgnoreTime != nil && item.CreatedTime().Before(*opts.IgnoreRule.IgnoreTime) {
			continue
		}

		for _, tag := range item.Tags() {
			if utils.Contains(opts.IgnoreRule.IgnoreTags, tag) {
				continue OUTER
			}
		}

		gist += prepareGistForItem(item)
	}

	gist += "</ol>"

	return gist
}

func prepareGistForItem(item Item) string {
	return htmlForItem(item) + "\n"
}

func htmlForItem(item Item) string {
	return fmt.Sprintf(
		"<li><a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\">%s</a></li>",
		item.URL(),
		item,
	)
}
