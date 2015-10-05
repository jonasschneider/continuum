package main

import (
  "net/http"
  "log"
  "os"
)

var GithubApiToken, GithubSharedSecret, RootPath, EntrypointPath, PostmarkApiToken, PostmarkSenderEmail, ExternalHostname string

func main() {
  GithubSharedSecret = ensure("GITHUB_SHARED_SECRET")
  GithubApiToken = ensure("GITHUB_API_TOKEN")
  RootPath = ensure("CI_ROOT")
  EntrypointPath = ensure("CI_ENTRYPOINT")
  PostmarkApiToken = ensure("POSTMARK_API_TOKEN")
  PostmarkSenderEmail = ensure("POSTMARK_SENDER_EMAIL")
  ExternalHostname = ensure("EXTERNAL_HOSTNAME")

  http.HandleFunc("/githubhook", handleGithubPost)
  http.HandleFunc("/builds/", handleBuildLog)
  log.Fatal(http.ListenAndServe(":80", nil))
}

func ensure(name string) string {
  v := os.Getenv(name)
  if v == "" { log.Fatalln("Please set",name) }
  return v
}
