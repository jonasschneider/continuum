mkdir -p /var/ci/{log,builds}

# install rbenv and ruby-build
git clone https://github.com/sstephenson/rbenv.git /var/ci/rbenv
mkdir /var/ci/rbenv/plugins
git clone https://github.com/sstephenson/ruby-build.git /var/ci/rbenv/plugins/ruby-build

chown -R ci /var/ci

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
