Continuum: a minimalist CI server
-------------------------------

    $ go install github.com/jonasschneider/continuum

    # adjust your config and start it (see env.sh.example)
    $ source env.sh
    $ continuum

    # point a github hook at http://$EXTERNAL_HOSTNAME/githubhook?secret=$GITHUB_SHARED_SECRET
    # then watch builds happen!
