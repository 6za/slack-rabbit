# slack-rabbit
Simple sample of a Slack bot using `go` and `rabbitMQ`. 



# How to create your Bot app

1. Create an app:
https://api.slack.com/apps

2. Fill `App-name` and `workspace`.
3. Select `Event Subscriptions` and enable it. 
4. Provide your app URL at  `Request URL`. (FQDN with the valid certificate)
5. Validate URL
...
Missing some steps and screen-shots
... 
9. Fill `.env` file to feed your docker-compose topology
10. Deploy your host topology docker-compose 
> Current server is a simple echo bot, that anything you send to the bot will be sent back to you. Internally, it is more than that. 
> This sample allow you to receive events from `slack`, they are all sent to a rabbitMQ queue, then you can process it, send to another rabbitMQ queue an the writer will send this messages to slack.

## Sample .env

```bash 
BOT_USER_ID=UXXXXXXXXX
SLACK_TOKEN=xxxx-111111111111-111111111111-XxXxXxXxXxXxXxXxXxXxXxXx
QUEUE_HOSTNAME=queuem
QUEUE_USER=guest
QUEUE_PASSWORD=guest
PUBLIC_DOMAIN=yoursite.com
```

## Run Slack-Rabbit

```bash 
git clone git@github.com:6za/slack-rabbit.git
cd slack-rabbit
export BASE_DIR=$PWD
./build.sh
docker-compose up

```

# Reference: 

Subscribe to events: 
- https://api.slack.com/events-api#events_api_request_urls
