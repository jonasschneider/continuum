#!/bin/bash
set -eux

# set the locale for the postgres cluster, default POSIX locale fails
update-locale LANG=C.UTF-8

apt-get update
apt-get -y install autoconf bison build-essential libssl-dev libyaml-dev libreadline6 libreadline6-dev zlib1g zlib1g-dev libpq-dev libssl-dev libxml2-dev libxslt-dev nodejs postgresql git curl ruby1.9.3 wget redis-server

# gems for the webface.. ehehe
gem install sinatra bundler

mkdir -p /var/ci/{log,builds}

# install rbenv and ruby-build
git clone https://github.com/sstephenson/rbenv.git /var/ci/rbenv
mkdir /var/ci/rbenv/plugins
git clone https://github.com/sstephenson/ruby-build.git /var/ci/rbenv/plugins/ruby-build

useradd ci
chown -R ci /var/ci
mkdir /home/ci
chown ci /home/ci
/etc/init.d/postgresql start
su postgres -c 'createuser -s ci'
/etc/init.d/postgresql stop

# add github ssh key
mkdir -p /home/ci/.ssh
cat > /home/ci/.ssh/known_hosts <<END
# github.com SSH-2.0-OpenSSH_5.9p1 Debian-5ubuntu1+github5
|1|rSByF9/SmEVjMqVcR1priTzRXV0=|8aBnwIuz7ZI7KDqOiUV+XGKF3Ik= ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
END
chown ci /home/ci/.ssh/known_hosts

# install ruby in the rbenv, activate it, and install bundler
su ci -c /bin/bash <<END
  set -eux
  export PATH="/var/ci/rbenv/bin:$PATH"
  export RBENV_ROOT=/var/ci/rbenv
  eval "\$(rbenv init -)"
  rbenv install 2.1.0
  rbenv global 2.1.0
  gem install bundler
  rbenv rehash
END
