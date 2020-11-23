package kubernetes

import (
	go_tool "github.com/maxiloEmmmm/go-tool"
	"net/url"
	"regexp"
	"strings"
)

func TransferDns(name string) string {
	otherReg, _ := regexp.Compile("[^a-z0-9-]")
	trimReg, _ := regexp.Compile("(^[0-9-]+|[-]+$|[-]{2,})")

	name = strings.ToLower(name)
	name = otherReg.ReplaceAllString(name, "-")
	name = trimReg.ReplaceAllString(name, "")
	return name
}

func TransferGitDns(gitUrl string) string {
	u, err := url.Parse(gitUrl)
	if err != nil {
		return ""
	}
	pathReg, _ := regexp.Compile("(^[\\/]+|[\\/]+$)")
	return TransferDns(go_tool.StringJoin(u.Host, ".", strings.Join(strings.Split(pathReg.ReplaceAllString(u.Path, ""), "/"), "-")))
}
