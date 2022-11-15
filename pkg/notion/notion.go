package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	gnotion "github.com/dstotijn/go-notion"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/notion-gister/pkg/utils"
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
) ([]map[string]string, error) {
	var items []map[string]string

	filter, err := createFilterFromJSON(filterTemplate)
	if err != nil {
		return nil, err
	}

	sorts, err := createSortFromJSON(sortTemplate)
	if err != nil {
		return nil, err
	}

	result, err := n.QueryDatabaseLoadAll(ctx, id, filter, sorts)
	if err != nil {
		return nil, err
	}

	for _, page := range result {
		props, ok := page.Properties.(gnotion.DatabasePageProperties)
		if !ok {
			logrus.Warnf("invalid properties: %T %v", page.Properties, page.Properties)
			continue
		}

		item := map[string]string{
			"__id":             page.ID,
			"__createdTime":    utils.PrettyTime(page.CreatedTime),
			"__lastEditedTime": utils.PrettyTime(page.LastEditedTime),
			"__url":            page.URL,
		}

		for name, prop := range props {
			item[name] = stringifyNotionProperty(prop)
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

func stringifyNotionProperty(prop gnotion.DatabasePageProperty) string {
	switch prop.Type {
	case gnotion.DBPropTypeTitle:
		var title string
		for _, text := range prop.Title {
			title += text.PlainText
		}

		return title
	case gnotion.DBPropTypeRichText:
		var rtext string
		for _, text := range prop.RichText {
			rtext += text.PlainText
		}

		return rtext
	case gnotion.DBPropTypeNumber:
		if prop.Number == nil {
			return ""
		}

		return strconv.FormatFloat(*prop.Number, 'f', -1, 64)
	case gnotion.DBPropTypeSelect:
		if prop.Select == nil {
			return ""
		}

		return prop.Select.Name
	case gnotion.DBPropTypeMultiSelect:
		var mselect []string
		for _, selectOption := range prop.MultiSelect {
			mselect = append(mselect, selectOption.Name)
		}

		return strings.Join(mselect, ", ")
	case gnotion.DBPropTypeDate:
		if prop.Date == nil {
			return ""
		}

		return stringifyNotionDate(prop.Date)
	case gnotion.DBPropTypePeople:
		var people []string
		for _, person := range prop.People {
			people = append(people, person.Name)
		}

		return strings.Join(people, ", ")
	case gnotion.DBPropTypeFiles:
		// won't support
		return ""
	case gnotion.DBPropTypeCheckbox:
		if prop.Checkbox == nil {
			return ""
		}

		return strconv.FormatBool(*prop.Checkbox)
	case gnotion.DBPropTypeURL:
		if prop.URL == nil {
			return ""
		}

		return *prop.URL
	case gnotion.DBPropTypeEmail:
		if prop.Email == nil {
			return ""
		}

		return *prop.Email
	case gnotion.DBPropTypePhoneNumber:
		if prop.PhoneNumber == nil {
			return ""
		}

		return *prop.PhoneNumber
	case gnotion.DBPropTypeStatus:
		if prop.Status == nil {
			return ""
		}

		return prop.Status.Name
	case gnotion.DBPropTypeFormula:
		if prop.Formula == nil {
			return ""
		}

		switch prop.Formula.Type {
		case gnotion.FormulaResultTypeNumber:
			if prop.Formula.Number == nil {
				return ""
			}

			return strconv.FormatFloat(*prop.Formula.Number, 'f', -1, 64)
		case gnotion.FormulaResultTypeString:
			if prop.Formula.String == nil {
				return ""
			}

			return *prop.Formula.String
		case gnotion.FormulaResultTypeDate:
			if prop.Formula.Date == nil {
				return ""
			}

			return stringifyNotionDate(prop.Formula.Date)
		default:
			return ""
		}
	case gnotion.DBPropTypeRelation:
		// won't support
		return ""
	case gnotion.DBPropTypeRollup:
		switch prop.Rollup.Type {
		case gnotion.RollupResultTypeNumber:
			if prop.Rollup.Number == nil {
				return ""
			}

			return strconv.FormatFloat(*prop.Rollup.Number, 'f', -1, 64)
		case gnotion.RollupResultTypeArray:
			var array []string
			for _, item := range prop.Rollup.Array {
				array = append(array, stringifyNotionProperty(item))
			}

			return strings.Join(array, ",")
		case gnotion.RollupResultTypeDate:
			if prop.Rollup.Date == nil {
				return ""
			}

			return stringifyNotionDate(prop.Rollup.Date)
		case gnotion.RollupResultTypeUnsupported:
			return "[Unsupported]"
		case gnotion.RollupResultTypeIncomplete:
			return "[Incomplete]"
		default:
			return ""
		}
	case gnotion.DBPropTypeCreatedTime:
		if prop.CreatedTime == nil {
			return ""
		}

		return utils.PrettyTime(*prop.CreatedTime)
	case gnotion.DBPropTypeCreatedBy:
		if prop.CreatedBy == nil {
			return ""
		}

		return prop.CreatedBy.Name
	case gnotion.DBPropTypeLastEditedTime:
		if prop.LastEditedTime == nil {
			return ""
		}

		return utils.PrettyTime(*prop.LastEditedTime)
	case gnotion.DBPropTypeLastEditedBy:
		if prop.LastEditedBy == nil {
			return ""
		}

		return prop.LastEditedBy.Name
	default:
		return ""
	}
}

func stringifyNotionDate(date *gnotion.Date) string {
	if date == nil {
		return ""
	}

	if date.TimeZone == nil {
		if date.End == nil {
			return utils.PrettyTime(date.Start.Time)
		}

		return fmt.Sprintf("%s - %s", utils.PrettyTime(date.Start.Time), utils.PrettyTime(date.End.Time))
	}

	return fmt.Sprintf(
		"%s - %s",
		utils.TimeInZone(date.Start.Time, *date.TimeZone, time.RFC1123),
		utils.TimeInZone(date.End.Time, *date.TimeZone, time.RFC1123),
	)
}
