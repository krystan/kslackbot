package main

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "sync/atomic"
    "log"
    "net/http"

    "golang.org/x/net/websocket"
)

type rtmStartResponse struct {
  Ok    bool         `json:"ok"`
  Error string       `json:"error"`
  Url   string       `json:"url"`
  Self  responseSelf `json:"self"`
}

type responseSelf struct {
  Id string `json:"id"`
}

func slackStart(token string) (wsurl, id string, err error) {
  url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
  resp, err := http.Get(url)

  if err != nil {
    return
  }

  if resp.StatusCode != 200 {
    err = fmt.Errorf("Request failed with code %d", resp.StatusCode)
    return
  }

  body, err := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  if err != nil {
    return
  }

  var responseBody rtmStartResponse
  err = json.Unmarshal(body, &responseBody)

  if err != nil {
    return
  }

  if !responseBody.Ok {
    err = fmt.Errorf("Slack error: %s", responseBody.Error)
    return
  }

  wsurl = responseBody.Url
  id = responseBody.Self.Id

  return
}


type SlackMessage struct {
  Id      uint64 `json:"id"`
  Type    string `json:"type"`
  Channel string `json:"channel"`
  Text    string `json:"text"`
}

func getMessage(ws *websocket.Conn) (m SlackMessage, err error) {
  err = websocket.JSON.Receive(ws, &m)
  return
}

var messageCount uint64
func sendMessage(ws *websocket.Conn, m SlackMessage) error {
  m.Id = atomic.AddUint64(&messageCount, 1)
  return websocket.JSON.Send(ws, m)
}

func connectToSlack(token string) (*websocket.Conn, string) {
  wsurl, id, err := slackStart(token)
  if err != nil {
    log.Fatal(err)
  }

  ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
  if err != nil {
    log.Fatal(err)
  }

  return ws, id
}
