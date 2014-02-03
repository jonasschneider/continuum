from ubuntu:saucy

add setup.sh /tmp/setup.sh
run /tmp/setup.sh
run rm /tmp/setup.sh

add init /init
add hook /hook
add build /build
add hook.ru /hook.ru

expose 9292
cmd /init
