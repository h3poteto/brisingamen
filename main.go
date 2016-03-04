package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type SlackMessage struct {
	Text      string `json:"text"`       //投稿内容
	Username  string `json:"username"`   //投稿者名 or Bot名（存在しなくてOK）
	IconEmoji string `json:"icon_emoji"` //アイコン絵文字
	IconURL   string `json:"icon_url"`   //アイコンURL（icon_emojiが存在する場合は、適応されない）
	Channel   string `json:"channel"`    //#部屋名
}

func main() {
	client, err := initClient()
	if err != nil {
		panic(err)
	}

	labels := []string{}
	for _, l := range Config("labels").([]interface{}) {
		labels = append(labels, l.(string))
	}

	listOptions := &github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
	}

	repo := Config("repository").(map[interface{}]interface{})
	issues, resp, err := client.Issues.ListByRepo(repo["owner"].(string), repo["repo"].(string), listOptions)
	if err != nil {
		panic(err)
	}

	log.Println(resp)

	// issueだとpullRequestも含めて全て取れるので、pullRequestに限定する
	pullRequests := []github.Issue{}
	for _, issue := range issues {
		if isPullRequest(issue) {
			pullRequests = append(pullRequests, issue)
		}
	}

	slackPost(pullRequests)
}

func initClient() (c *github.Client, e error) {
	token := Config("access_token")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.(string)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	return client, nil
}

func Config(elem string) interface{} {
	path := filepath.Join("./", "settings.yml")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	return m[elem]
}

func isPullRequest(issue github.Issue) bool {
	if issue.PullRequestLinks == nil {
		return false
	}
	return true
}

func slackPost(pullRequest []github.Issue) {
	slackURL := Config("slack_url")
	channelName := Config("channel")

	text := ""
	for _, pr := range pullRequest {
		text = text + *pr.Title + ": " + *pr.HTMLURL + "\n"
	}

	params, _ := json.Marshal(SlackMessage{
		text,
		"PleaseReview",
		":ghost:",
		"",
		channelName.(string)})

	resp, _ := http.PostForm(
		slackURL.(string),
		url.Values{"payload": {string(params)}},
	)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	log.Println(string(body))

}
