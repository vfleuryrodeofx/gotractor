package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

const ROOT_ENDPOINT = "http://tractor.rodeofx.com/Tractor/"

var ENDPOINTS = map[string]string{
	"tasktree": "monitor?q=jtree&jid=%s",
	"logs":     "monitor?q=tasklogs&owner=%s&jid=%s&tid=%s",
}

func ExtractJID(url string) string {
	pattern := regexp.MustCompile(`http://[A-z\-\.]*/tv/#jid=(?P<jid>[0-9]*)`)
	slog.Info("Extracting jid from url")
	matches := pattern.FindStringSubmatch(url)
	if len(matches) > 1 {
		jid := matches[1]
		slog.Info("JID found ", "jid", jid)
		return jid
	}
	return ""
}

// Get tree data
func GetTaskTree(jid string) (map[string]any, []interface{}) {
	url := fmt.Sprintf("http://tractor.rodeofx.com/Tractor/monitor?q=jtree&jid=%s", jid)
	slog.Info("Querying task tree at", "url", url)

	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	//Read body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var jsonObject map[string]any
	err = json.Unmarshal(body, &jsonObject)
	if err != nil {
		panic(err)
	}

	root := jsonObject["users"].(map[string]any)
	var user string
	for u := range root {
		user = u
		break
	}
	userData := root[user].(map[string]any)
	jidkey := userData[fmt.Sprintf("J%s", jid)].(map[string]any)
	data := jidkey["data"].(map[string]any)
	tasksData := jidkey["children"].([]any)

	return data, tasksData
}

func GetTaskLog(owner, jobID, taskID string) string {
	url := fmt.Sprintf("http://tractor-log-viewer.rodeofx.com/tractor/%s/J%s/%s.log", owner, jobID, taskID)
	slog.Info("Fetching logs from endpoint ", "endpoint", url)

	requestLogPath, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer requestLogPath.Body.Close()

	logContent, err := io.ReadAll(requestLogPath.Body)
	if err != nil {
		panic(err)
	}

	return string(logContent)
}
