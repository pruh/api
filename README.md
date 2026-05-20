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

### Providers:

API to store and retrive providers.
Providers are stored in-memory and are reset when the API process restarts.

The following HTTP methods are supported:

* `/api/v1/providers/` HTTP GET method to return all providers.
  
  `/api/v1/providers/{UUID}` HTTP GET method to return a single provider by UUID.

  Methods return providers in the following format:
  ```json
  {
      "_id": "c146d6f1-8992-4010-85da-80459bb55d10",
      "type": "NJTransit",
      "njtransit": {
        "orig_station_code":"AA",
        "dest_station_code":"AB"
      }
  }
  ```

  based on `type` different optional data will present in json, like `njtransit` in the example above

* `/api/v1/providers/` HTTP POST method to save a provider.
  Method accepts JSON in the following format:

  ```json
  {
      "type": "NJTransit",
      "njtransit": {
        "orig_station_code":"AA",
        "dest_station_code":"AB"
      }
  }
  ```

* `/api/v1/providers/{UUID}` HTTP DELETE method to delete a previously saved provider by UUID.
