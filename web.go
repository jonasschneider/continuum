package main

import (
  "net/http"
  "log"
  "encoding/json"
  "strings"
  "fmt"
  "os"
  "io"
)

func handleGithubPost(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  defer func() {
    if r := recover(); r != nil {
      log.Println(r)
      w.WriteHeader(400)
    }
  }()

  if r.Form.Get("secret") != GithubSharedSecret {
    w.WriteHeader(401)
    return
  }

  var f interface{}
  err := json.Unmarshal([]byte(r.Form.Get("payload")), &f)
  if err != nil { log.Println(err); w.WriteHeader(400); return; }

  ref := f.(map[string]interface{})["ref"].(string)
  rev := f.(map[string]interface{})["after"].(string)

  repo_obj := f.(map[string]interface{})["repository"].(map[string]interface{})
  repo_url := repo_obj["url"].(string)

  email := f.(map[string]interface{})["head_commit"].(map[string]interface{})["author"].(map[string]interface{})["email"].(string)

  log.Println("Starting build of",repo_url,"ref",ref,"at",rev,"by",email)
  go runAndReportBuild(repo_url, ref, rev, email)

  w.WriteHeader(200)
}

func handleBuildLog(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  if r.Form.Get("secret") != GithubSharedSecret {
    w.WriteHeader(401)
    return
  }

  parts := strings.Split(r.URL.Path, "/")
  name := parts[len(parts)-1]
  if strings.Index(name,"/") != -1 { // at least so *something*
    w.WriteHeader(400)
    return
  }
  path := fmt.Sprintf(RootPath+"/log/"+name+".log") // unsafe, but behind auth.. meeh
  f, err := os.Open(path)
  if err != nil {
    log.Println(err)
    w.WriteHeader(500)
    return
  }
  w.Header().Add("Content-Type", "text/plain")
  _, err = io.Copy(w, f)
  if err != nil {
    log.Println(err)
  }
}
