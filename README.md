# api
![Build](https://github.com/pruh/api/actions/workflows/test_and_start.yml/badge.svg)
[![CodeCov](https://codecov.io/gh/pruh/api/branch/master/graph/badge.svg)](https://codecov.io/gh/pruh/api)
[![GoDoc](https://godoc.org/github.com/pruh/api?status.svg)](http://godoc.org/github.com/pruh/api)

REST API server written in Go. Supports basic HTTP auth.

## Usage

* Rename api.env.template to api.env and add correct its contents, such as Telegram BOT token, basic auth credential, etc.
* Run `docker-compose up -d` to start the server

## api.env

Simple key-value file which will be used by docker to set container environment variables.

* `PORT` mandatory port to use for service.

* `TELEGRAM_BOT_TOKEN` mandatory telegram bot token.

* `TELEGRAM_DEFAULT_CHAT_ID` default telegram chat ID, which will receive messages from the bot. This parameter is optinal.

* `API_V1_CREDS` username/password pairs in JSON format of users who are allowed to access API: `{"username1":"password1", "username2":"password2"}`. This parameter is optional.

## List of API methods

### Messages:

API to send messages.

* `/api/v1/telegram/messages/send` POST method which passes message to telegram. It accepts JSON in the following format:

  ```json
  {
      "message": "message to send",
      "chat_id": 1234567890,
      "silent": true
  }
  ```

  where `chat_id` is telegram chat id and `silent` is a flag indicating if message should be sent silently
