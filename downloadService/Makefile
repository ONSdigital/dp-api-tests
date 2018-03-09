VAULT_ADDR?='http://127.0.0.1:8200'

# The following variables are used to generate a vault token for the app. The reason for declaring variables, is that
# its difficult to move the token code in a Makefile action. Doing so makes the Makefile more difficult to
# read and starts introduction if/else statements.
TOKEN_INFO:="$(shell vault token create -address=$(VAULT_ADDR) -policy=write-psk -period=5m -display-name=dp-api-tests)"
APP_TOKEN:="$(shell echo $(TOKEN_INFO) | awk '{print $$6}')"

vault:
	@echo "$(TOKEN_INFO)"
	@echo "$(APP_TOKEN)"

test:
	HUMAN_LOG=1 VAULT_TOKEN=$(APP_TOKEN) VAULT_ADDR=$(VAULT_ADDR) go test ./...

.PHONY: vault test