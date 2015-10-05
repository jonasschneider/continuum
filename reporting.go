package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func reportBuildFailure(github_repo_name, ref, rev, email string, res BuildResult) {
	var url, urlstr string
	if res.Name != "" {
		url = fmt.Sprintf("http://%s/builds/%s?secret=%s", ExternalHostname, res.Name, GithubSharedSecret)
		urlstr = url + "\n\n"
	}

	debuglog := "(no output)"
	if res.DiagOut != nil {
		logOut, err := ioutil.ReadAll(res.DiagOut)
		if err != nil {
			log.Println("error while reading log", err)
		}
		debuglog = string(logOut)
	}
	body := "error: " + res.Error.Error() + "\n\n" + urlstr + debuglog
	subj := fmt.Sprintf("[ci: %s] %s failed", strings.Replace(ref, "refs/heads/", "", 1), rev)
	err := sendMail(Mail{TextBody: body, To: email, Subject: subj})
	if err != nil {
		log.Println("error while sending failure notification:", err)
	} else {
		log.Println("Sent failure mail to", email)
	}

	err = updateGithubStatus(github_repo_name, rev, githubStatus{State: "failure", Description: ExternalHostname + " failed :(", Target: url})
	if err != nil {
		log.Println("failed to update github status:", err)
	}
}

func reportBuildSuccess(github_repo_name, ref, rev, email string, res BuildResult) {
	err := updateGithubStatus(github_repo_name, rev, githubStatus{State: "success", Description: ExternalHostname + " passed! :)"})
	if err != nil {
		log.Println("failed to update github status:", err)
	}
}

func reportBuildStart(github_repo_name, ref, rev, email string) {
	err := updateGithubStatus(github_repo_name, rev, githubStatus{State: "pending", Description: "Building on " + ExternalHostname + "..."})
	if err != nil {
		log.Println("failed to update github status:", err)
	}
}

type Mail struct {
	From, To, Subject, TextBody string
}

func sendMail(m Mail) error {
	m.From = PostmarkSenderEmail
	d, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewReader(d))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", PostmarkApiToken)
	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status 200, but was %d", resp.StatusCode)
	}
	return nil
}

type githubStatus struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Target      string `json:"target_url"`
}

func updateGithubStatus(github_repo_name, rev string, status githubStatus) error {
	d, err := json.Marshal(status)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", github_repo_name, rev), bytes.NewReader(d))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token "+GithubApiToken)
	req.Header.Add("Content-Type", "application/json")
	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("expected status 201, but was %d", resp.StatusCode)
	}
	return nil
}
