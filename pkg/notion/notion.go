package notion

import (
	"context"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	gnotion "github.com/dstotijn/go-notion"
	"github.com/sirupsen/logrus"
)

// Notion is a wrapper around the go-notion client with some
// additional functionality.
type Notion struct {
	*gnotion.Client
}

// New creates a new Notion client.
func New(apiKey string) *Notion {
	return &Notion{
		Client: gnotion.NewClient(apiKey),
	}
}

// QueryDatabaseLoadAll loads all pages from a database.
//
// NOTE: This will not load more than 1000 pages.
func (n *Notion) QueryDatabaseLoadAll(
	ctx context.Context,
	id string,
	filter *gnotion.DatabaseQueryFilter,
	sorts []gnotion.DatabaseQuerySort,
) ([]gnotion.Page, error) {
	const pageSize = 100
	const maxPages = 10

	var pages []gnotion.Page
	var cursor string = ""

	// Ensure that no more than 1000 items are loaded at once.
	for i := 0; i < pageSize*maxPages; i++ {
		result, err := n.QueryDatabase(ctx, id, &gnotion.DatabaseQuery{
			Filter:      filter,
			Sorts:       sorts,
			StartCursor: cursor,
			PageSize:    pageSize,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, result.Results...)

		if !result.HasMore {
			break
		}

		if result.NextCursor == nil {
			break
		}

		cursor = *result.NextCursor
	}

	return pages, nil
}

// ReadDBItems reads all items from a database.
func (n *Notion) ReadDBItems(
	ctx context.Context,
	id string,
	filterTemplate,
	sortTemplate string,
) ([]DBItem, error) {
	var items []DBItem

	filter, err := createFilterFromJSON(filterTemplate)
	if err != nil {
		return nil, err
	}

	sorts, err := createSortFromJSON(sortTemplate)
	if err != nil {
		return nil, err
	}

	result, err := n.QueryDatabaseLoadAll(
		ctx,
		id,
		filter,
		sorts,
	)
	if err != nil {
		return nil, err
	}

	for _, page := range result {
		props, ok := page.Properties.(gnotion.DatabasePageProperties)
		if !ok {
			logrus.Warnf("invalid properties: %T %v", page.Properties, page.Properties)
			continue
		}

		item := DBItem{
			id:          page.ID,
			createdTime: page.CreatedTime,
			url:         page.URL,
			title:       "",
		}

		for _, title := range props["Name"].Title {
			item.title += title.PlainText
		}

		for _, tags := range props["Tags"].MultiSelect {
			item.tags = append(item.tags, tags.Name)
		}

		items = append(items, item)
	}

	return items, nil
}

func createFilterFromJSON(filterTemplate string) (*gnotion.DatabaseQueryFilter, error) {
	if filterTemplate == "" {
		return nil, nil
	}

	ftemp, err := template.New("filter").Parse(filterTemplate)
	if err != nil {
		return nil, err
	}

	var filter strings.Builder
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	if err := ftemp.Execute(&filter, map[string]interface{}{
		"yesterday": today.AddDate(0, 0, -1).Format(time.RFC3339Nano),
		"today":     today.Format(time.RFC3339Nano),
		"tomorrow":  today.AddDate(0, 0, 1).Format(time.RFC3339Nano),
	}); err != nil {
		return nil, err
	}

	var f gnotion.DatabaseQueryFilter
	if err := json.Unmarshal([]byte(filter.String()), &f); err != nil {
		return nil, err
	}

	return &f, nil
}

func createSortFromJSON(sortTemplate string) ([]gnotion.DatabaseQuerySort, error) {
	if sortTemplate == "" {
		return nil, nil
	}

	stemp, err := template.New("sort").Parse(sortTemplate)
	if err != nil {
		return nil, err
	}

	var sort strings.Builder
	if err := stemp.Execute(&sort, map[string]interface{}{
		"today": time.Now().Format(time.RFC3339Nano),
	}); err != nil {
		return nil, err
	}

	var s []gnotion.DatabaseQuerySort
	if err := json.Unmarshal([]byte(sort.String()), &s); err != nil {
		return nil, err
	}

	return s, nil
}
