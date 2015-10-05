package main

import (
  "net/http"
  "io/ioutil"
  "os"
  "encoding/json"
  "encoding/base64"
  "strings"
  "os/exec"
  "io"
  "log"
  "time"
  "fmt"
)

type BuildResult struct {
  Name string
  Error error
  DiagOut io.Reader
  PrettyOut io.Reader
}

func fetchScript(github_repo_name, rev string) (error, io.Reader) {
  fetchUrl := "https://api.github.com/repos/"+github_repo_name+"/contents/"+EntrypointPath+"?ref="+rev
  var c http.Client
  req, err := http.NewRequest("GET", fetchUrl, nil)
  if err != nil { return err, nil }
  req.Header.Add("Authorization", "token "+GithubApiToken)
  resp, err := c.Do(req)
  if err != nil { return err, nil }

  var f githubFile

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil { return err, nil }
  err = json.Unmarshal(data, &f)
  if err != nil { return err, nil }

  return nil, base64.NewDecoder(base64.StdEncoding, strings.NewReader(f.Content))
}

func runBuild(github_repo_name, rev string) BuildResult {
  // fail immediately if duplicate
  build_name := fmt.Sprintf("%d-%s", time.Now().Unix(), rev)
  diagLogfile, err := os.Create(RootPath+"/log/"+build_name+".log")
  if err != nil { return BuildResult{Error: err} }

  err, scriptReader := fetchScript(github_repo_name, rev)
  if err != nil { return BuildResult{Error: err} }

  cmd := exec.Command("bash")
  cmd.Dir = "/tmp"
  cmd.Stdin = scriptReader
  cmd.Env = os.Environ()
  cmd.Env = append(cmd.Env, "REPO=git@github.com:"+github_repo_name+".git")
  cmd.Env = append(cmd.Env, "REFSPEC="+rev)
  cmd.Env = append(cmd.Env, "CACHE="+RootPath+"/cache")

  log.Println("Starting",EntrypointPath)

  stdout, err := cmd.StdoutPipe()
  if err != nil { return BuildResult{Error: err} }

  stderr, err := cmd.StderrPipe()
  if err != nil { return BuildResult{Error: err} }

  output := io.MultiReader(stdout, stderr)
  teedOut := io.TeeReader(output, os.Stderr)

  err = cmd.Start()
  if err != nil { return BuildResult{Error: err} }

  err = os.MkdirAll(RootPath+"/log", 0777)
  if err != nil { return BuildResult{Error: err} }

  diagLogfileRead, err := os.Open(diagLogfile.Name())
  if err != nil { return BuildResult{Error: err} }

  _, err = io.Copy(diagLogfile, teedOut)
  if err != nil { return BuildResult{Error: err, DiagOut: diagLogfileRead, Name: build_name} }

  err = cmd.Wait()
  if err != nil { return BuildResult{Error: err, DiagOut: diagLogfileRead, Name: build_name} }

  return BuildResult{Error: nil, DiagOut: diagLogfileRead}
}

func runAndReportBuild(repo_url, ref, rev, email string) {
  github_repo_name := strings.Replace(repo_url, "https://github.com/", "", 1)
  reportBuildStart(github_repo_name, ref, rev, email)
  buildResult := runBuild(github_repo_name, rev)
  if buildResult.Error != nil {
    log.Println("error or failure while running build of",rev,":",buildResult.Error)
    reportBuildFailure(github_repo_name, ref, rev, email, buildResult)
  } else {
    log.Println("build of",rev,"succeeded!")
    reportBuildSuccess(github_repo_name, ref, rev, email, buildResult)
  }
}

type githubFile struct {
  Content string `json:"content"`
}
