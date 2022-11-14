package notion

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// DBItem represents a single item in a Notion database that
// gister cares about.
type DBItem struct {
	id          string
	title       string
	url         string
	tags        []string
	createdTime time.Time
}

func (i DBItem) String() string {
	loc, _ := time.LoadLocation(viper.GetString("timezone"))

	return fmt.Sprintf("%s (%s)", i.title, i.createdTime.In(loc).Format(time.RFC822))
}

func (i DBItem) ID() string {
	return i.id
}

func (i DBItem) Title() string {
	return i.title
}

func (i DBItem) URL() string {
	return i.url
}

func (i DBItem) Tags() []string {
	return i.tags
}

func (i DBItem) CreatedTime() time.Time {
	return i.createdTime
}
