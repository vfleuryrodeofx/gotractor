package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

func ExtractJID(url string) string {
	pattern := regexp.MustCompile(`http://[A-z\-\.]*/tv/#jid=(?P<jid>[0-9]*)`)
	fmt.Println("Extracting jid from url")
	matches := pattern.FindStringSubmatch(url)
	if len(matches) > 1 {
		jid := matches[1]
		fmt.Println("JID found ", jid)
		return jid
	}
	return ""
}

// Get tree data
func GetTaskTree(jid string) map[string]interface{} {
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

	var jsonObject map[string]interface{}
	err = json.Unmarshal(body, &jsonObject)
	if err != nil {
		panic(err)
	}

	root := jsonObject["users"].(map[string]interface{})
	var user string
	for u := range root {
		user = u
		fmt.Println("User is", user)
		break
	}
	userData := root[user].(map[string]interface{})
	jidkey := userData[fmt.Sprintf("J%s", jid)].(map[string]interface{})
	data := jidkey["data"].(map[string]interface{})
	return data
}
