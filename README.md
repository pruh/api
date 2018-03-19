# api
[![Build Status](https://travis-ci.org/pruh/api.svg?branch=master)](https://travis-ci.org/pruh/api)
[![Coverage Status](https://coveralls.io/repos/github/pruh/api/badge.svg?branch=master)](https://coveralls.io/github/pruh/api?branch=master)
[![GoDoc](https://godoc.org/github.com/pruh/api?status.svg)](http://godoc.org/github.com/pruh/api)

Simple REST API server in Go with optional basic auth

## Usage

* Rename api.env.template to api.env and add correct data to it, such as Telegram BOT token, basic auth credential, etc.
* Run `docker-compose up -d` to start the server

## api.env

It is a simple key-value pairs file which will be used by docker to set environment variables in container.

* `PORT` port which will be used in container. This parameter is mandatory.

* `TELEGRAM_BOT_TOKEN` telegram bot token to use. This parameter is mandatory.

* `TELEGRAM_DEFAULT_CHAT_ID` default telegram chat ID, which will receive messages from bot. This parameter is optinal and one provided to API methods will be used otherwise.

* `API_V1_CREDS` username/password pairs in JSON format which are allowed to access API: `{"username1":"password1", "username2":"password2"}`. This parameter is optional, and no auth will be required if not set.

## List of API methods

* `/api/v1/telegram/messages/send` POST method which passes message to telegram. It accepts JSON in the following format:

```json
{
    "message": "message to send",
    "chat_id": 1234567890,
    "silent": true
}
```

where:<br />
`message` message to be sent.<br />
`chat_id` chat ID to send message to. This parameter is optional and default Telegram chat ID will be used otherwise.<br />
`silent` send message silently. This parameter is optional and will be set to true by default.
