chcp 65001
set GOHOST_DEBUG=true
dlv debug --headless --listen=:12345 --api-version=2 --accept-multiclient