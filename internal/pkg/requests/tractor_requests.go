package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

const RootEndpoint = "http://tractor.rodeofx.com/Tractor/"

var Endpoints = map[string]string{
	"tasktree": "monitor?q=jtree&jid=%s",
	"logs":     "monitor?q=tasklogs&owner=%s&jid=%s&tid=%s",
}

var jidPattern = regexp.MustCompile(`http://[A-z\-\.]*/tv/#jid=(?P<jid>[0-9]*)`)

func ExtractJID(arg string) string {
	slog.Info("Extracting jid from url")
	matches := jidPattern.FindStringSubmatch(arg)
	if len(matches) > 1 {
		jid := matches[1]
		slog.Info("JID found ", "jid", jid)
		return jid
	}
	return arg
}

// Get tree data
func GetTaskTree(jid string) (map[string]any, []interface{}, error) {
	url := RootEndpoint + fmt.Sprintf(Endpoints["tasktree"], jid)
	slog.Info("Querying task tree at", "url", url)

	req, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not fetch task tree for %s, error : %w", jid, err)
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status: %d", req.StatusCode)
	}

	//Read body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not read payload for %s. Err: %w", jid, err)
	}

	var jsonObject map[string]any
	err = json.Unmarshal(body, &jsonObject)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not unmarshal payload for %s. Err : %w", jid, err)
	}

	root, ok := jsonObject["users"].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("JSON payload does not have correct data. Err: %w", err)
	}
	var user string
	for u := range root {
		user = u
		break
	}
	userData := root[user].(map[string]any)
	jidkey := userData[fmt.Sprintf("J%s", jid)].(map[string]any)
	data := jidkey["data"].(map[string]any)
	tasksData := jidkey["children"].([]any)

	return data, tasksData, nil
}

func GetTaskLog(owner, jobID, taskID string) (string, error) {
	url := fmt.Sprintf("http://tractor-log-viewer.rodeofx.com/tractor/%s/J%s/%s.log", owner, jobID, taskID)
	slog.Info("Fetching logs from endpoint ", "endpoint", url)

	requestLogPath, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Could not fetch the log. err : %w", err)
	}

	defer requestLogPath.Body.Close()

	if requestLogPath.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", requestLogPath.StatusCode)
	}

	logContent, err := io.ReadAll(requestLogPath.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read the log. err : %w", err)
	}

	return string(logContent), nil
}
