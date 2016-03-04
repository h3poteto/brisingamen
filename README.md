# Brisingamen
Brisingamen is a simple slack bot for github pull requests.
It posts pull requests's information which belong to some labels in a your repository.

# Example
Brisingamen read `settings.yml` as setting file, so please read `settings.yml.sample` and write your `settings.yml`.

Generate your private access token in github, and add to `settings.yml` :

```yaml
access_token: "hogehoge"
```

Create a new incoming webhook in slack, and add:

```yaml
slack_url: "https://hooks.slack.com/services/hogehoge"
channel: "#general"
```

And, Please set your repository information:

```yaml
repository:
  owner: "h3poteto"
  repo: "asumibot"
labels:
  - "todo"
```

After that, you can use Brisingamen

```
$ go run main.go
```

or, build those codes:

```
$ go build -o brisingamen main.go
$ ./brisingamen
```




# License

Brisingamen is available as open source under the terms of the MIT License.
