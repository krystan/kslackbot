  package main

  import (
    "os"
    "fmt"
    "log"
    "strings"
    "context"
    
    "github.com/google/go-github/github"
  )

  func main() {
    if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "usage: kslackbot <args>\n")
      return
    }

    // git client
    gitClient := github.NewClient(nil)

    // open connection
    ws, id := connectToSlack(os.Args[1])
    fmt.Print("bot online, hit ^C to exit")

    gitClient = github.NewClient(nil)

    for {
      m, err := getMessage(ws)

      if err != nil {
        log.Fatal(err)
      }

      // check for message id
      if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
        parts := strings.Fields(m.Text)

        if len(parts) == 3 && parts[1] == "commit" {
          go func(m SlackMessage) {
            commit,err := getLastCommit("krystan", parts[2],  gitClient)

            if err != nil {
              log.Fatal("fatal error getting commit:%v", err)
            }

            m.Text = commit

            sendMessage(ws, m)
          }(m)
        } else {
          m.Text = fmt.Sprintf("unexpected response\n")
          sendMessage(ws, m)
        }
      }
    }
  }


  func getLastCommit(owner, repo string, sgc *github.Client) (string, error) {
    ctx := context.Background()
    commits, res, err := sgc.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{})

    if err != nil {
      log.Printf("err: %s res: %s", err, res)
      return "", err
    }

    return *commits[0].SHA, nil
  }
