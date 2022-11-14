package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/notion-gister/pkg/config"
	"github.com/utkarsh-pro/notion-gister/pkg/executor"
	"github.com/utkarsh-pro/notion-gister/pkg/gister"
	"github.com/utkarsh-pro/notion-gister/pkg/mailer"
	"github.com/utkarsh-pro/notion-gister/pkg/notion"
	"github.com/utkarsh-pro/notion-gister/pkg/utils"
)

func init() {
	config.Setup()
}

func main() {
	ctx := context.Background()
	exec := executor.New(viper.GetString("executor"))

	db := viper.Get("db").([]interface{})
	n := notion.New(viper.GetString("apiKey"))

	for iter, v := range db {
		i := iter
		capturedV := v.(map[string]interface{})

		cronSpec := viper.GetString(fmt.Sprintf("db.%d.cron", i))

		exec.AddFunc(cronSpec, func() {

			items, err := n.ReadDBItems(
				ctx,
				capturedV["id"].(string),
				viper.GetString(fmt.Sprintf("db.%d.notion.filterJSON", i)),
				viper.GetString(fmt.Sprintf("db.%d.notion.sortJSON", i)),
			)
			if err != nil {
				logrus.Error("Failed to read items from DB: ", capturedV["id"], err)
				return
			}

			_items := utils.FromT1ToT2(items, func(item notion.DBItem) gister.Item {
				return &item
			})

			gist := gister.Create(_items, gister.Opts{
				IgnoreRule: gister.IgnoreRule{
					IgnoreTags: []string{"gister:ignore"},
				},
			})

			mailer := mailer.New(
				utils.ViperReturnFirstFound[string](
					fmt.Sprintf("db.%d.mail.smtp.host", i),
					"mail.smtp.host",
				),
				utils.ViperReturnFirstFound[string](
					fmt.Sprintf("db.%d.mail.smtp.port", i),
					"mail.smtp.port",
				),
				utils.ViperReturnFirstFound[string](
					fmt.Sprintf("db.%d.mail.smtp.username", i),
					"mail.smtp.username",
				),
				utils.ViperReturnFirstFound[string](
					fmt.Sprintf("db.%d.mail.smtp.password", i),
					"mail.smtp.password",
				),
				utils.ViperReturnFirstFound[string](
					fmt.Sprintf("db.%d.mail.from", i),
					"mail.from",
				),
				utils.FromT1ToT2(utils.ViperReturnFirstFound[[]interface{}](
					fmt.Sprintf("db.%d.mail.to", i),
					"mail.to",
				), func(v1 interface{}) string {
					return v1.(string)
				}),
				map[string]interface{}{
					"dbname": capturedV["name"],
					"date":   time.Now().Format("2006-01-02"),
					"time":   time.Now().Format(time.RFC822),
					"mail": map[string]interface{}{
						"content": gist,
					},
				},
			)

			if err := mailer.
				Mail(
					utils.ViperReturnFirstFound[string](
						fmt.Sprintf("db.%d.mail.subject", i),
						"mail.subject",
					),
					utils.ViperReturnFirstFound[string](
						fmt.Sprintf("db.%d.mail.body", i),
						"mail.body",
					),
				); err != nil {
				logrus.Error("Error while sending mail: ", err)
			}
		})
	}

	exec.Start()
}
