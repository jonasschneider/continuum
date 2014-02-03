require 'sinatra'
require 'json'

post '/githubhook' do
  halt 401 unless params['secret'] == ENV["GITHUB_SHARED_SECRET"]
  payload = JSON.parse(params['payload'])
  $stderr.puts payload.inspect
  ref = payload["ref"]
  rev = payload["after"]
  repo = payload["repository"]["url"].sub("https://github.com/","")
  halt 400 unless repo.match(/\A[a-zA-z0-9\-_\/]+\Z/)
  email = payload["head_commit"]["author"]["email"]
  puts "#{rev} #{ref} #{email}"
  pid = fork do
    $stdout.reopen(File.open("/var/ci/log/hook", "a"))
    $stderr.reopen(File.open("/var/ci/log/hook", "a"))
    ENV["REPO"] = repo
    Kernel.exec("setsid", "/hook", rev, ref, email)
  end
  Process.detach(pid)
end

get '/builds/:name' do
  halt 401 unless params['secret'] == ENV["GITHUB_SHARED_SECRET"]
  build=params['name'].to_i
  content_type 'text/plain'
  send_file "/var/ci/log/#{build}", disposition: :inline
end

run Sinatra::Application
