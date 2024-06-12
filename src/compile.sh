# Install gvm first
#
# https://github.com/moovweb/gvm/issues/188

[[ -s "$GVM_ROOT/scripts/gvm" ]] && source "$GVM_ROOT/scripts/gvm"
TEST_APP=./src
gvm use go1.20
# echo $TEST_APP
CGO_ENABLED=0 GOOS=linux go build -o $TEST_APP/main $TEST_APP/...
