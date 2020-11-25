package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
	"regexp"
	"strings"
)

func TransferDns(name string) string {
	otherReg, _ := regexp.Compile("[^a-z0-9]")
	trimReg, _ := regexp.Compile("(^[0-9-]+)")

	name = strings.ToLower(name)
	name = otherReg.ReplaceAllString(name, "")
	name = trimReg.ReplaceAllString(name, "")
	return name
}

func TransferGitDns(gitUrl string) string {
	removeGit, _ := regexp.Compile("\\.git$")
	gitUrl = removeGit.ReplaceAllString(gitUrl, "")
	u, err := url.Parse(gitUrl)
	if err != nil {
		return ""
	}
	return TransferDns(u.Path)
}

func ResponseOk(response *resty.Response) error {
	if response.StatusCode() >= 300 && response.StatusCode() != 404 {
		msg := &InValidMessage{}
		err := json.Unmarshal(response.Body(), msg)
		if err != nil {
			return errors.New(fmt.Sprintf("Unmarshal error message: %s", response.String()))
		}
		return errors.New(fmt.Sprintf("code: %d, message: %s", msg.Code, msg.Message))
	} else {
		return nil
	}
}
