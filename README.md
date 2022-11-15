# Notion Gister

I have been a notion user since 2018. I really like it and I manage almost everything in it like:
1. Expenses
2. Scheduling
3. Tasks
4. Notes
5. Articles that I find on the internet
6. Planning
And many other things.

It works well! The only issue is that if I dump something in the notion then the odds of me taking a look at it again are **very** less. Basically, it is built on **pull** model, i.e. I have to pull the data from notion manually and this isn't very efficient. 

This is the problem Notion Gister tries to solve. Notion Gister if integrated with Notion Databases and configured, it can create a "gist" of the items and send a mail to you at the time you want (like a cron job).

Just like anything I build, it is very barebones at the expense of extensibility. Almost evert aspect of gister can be configured via `.gister.yaml`.

## Gister Config Guide
Gister config has to be named `.gister.yaml` and needs to be in the same directory where gister is running. 

```yaml
# apiKey is the Notion API Key (Can be given via GISTER_APIKEY as well)
apiKey: sadfsdfwefrwerf

# executor can be either "cron" or "none". If nothing is provided, "none" is assumed
# "none" executor will not respect `db.$i.cron` values but is useful if Gister is running in
# an environment which already has cron implemented (like Linux Server or Github Action)
executor: cron

# timezone is used to render the date and time in the given timezone. If no value is give, 
# UTC is assumed
timezone: Asia/Kolkata

# Default Mail Config
mail:
  # Default SMTP related config
  smtp:
    username: asbc # Username for SMTP relay (Alternatively can be set via env var GISTER_SMTP_USERNAME)
    password: absdf # Password for SMPT relay (Alternatively can be set via env var GISTER_SMTP_PASSWORD)
    host: smpt.example.com # SMTP Host (Alternatively can be set via env var GISTER_SMPT_HOST)
    port: "25" # SMTP Port - Has to be a string (Alternatively can be set via env var GISTER_SMPT_PORT) 
  from: <example@example.com> # Default mail sender ID

  # Default mail recivers
  to:
  - test@test.com
  - abc@test.com
  
  # Default Subject of the mail, this is a Go template and receives values like
  # dbname - Name of the database
  # time - Time of generation
  # date - Date of generation
  # notion - An array of map[string]string (In map, key is the name of the notion database
  # property and value is the property value in string format, additional values are present 
  # like __id - page id, __url - page url, __createdTime - Time of creation)
  subject: Hello world for "{{.dbname}}"

  # Default body of the mail, this is a Go template and receives values like
  # dbname - Name of the database
  # time - Time of generation
  # date - Date of generation
  # notion - An array of map[string]string (In map, key is the name of the notion database
  # property and value is the property value in string format, additional values are present 
  # like __id - page id, __url - page url, __createdTime - Time of creation)
  body: |
    Hello this is the mail body generated at {{.time}}

# DB is array of Notion Databases which needs to be scraped by Gister
db:
  - id: abc1234 # ID of the notion database - REQUIRED
    name: abcwe # Name of the database (could be anything) - REQUIRED
    cron: 0 12 * * * #  The times when the gister should kick in - REQUIRED
    notion:
      # filterJSON is also Go Template (refer - https://developers.notion.com/reference/post-database-query-filter)
      # Template receives values like "yesterday", "today", "tomorrow"
      filterJSON:
      # sortJSON is also Go Template (refer - https://developers.notion.com/reference/post-database-query-sort)
      # Template receives values like "today"
      sortJSON:
    # This mail has the same structure as the global mail but is specific to this DB,
    # therefore mail subject, body, senders, receiver, smpt relay everything can be overridden here
    mail:
```

Apart from the config, environment variables can be used as well with prefix GISTER.