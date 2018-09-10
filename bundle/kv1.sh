#!/bin/sh
export VAULT_ADDR=http://$VAULT_DEV_LISTEN_ADDRESS
export VAULT_TOKEN=$VAULT_DEV_ROOT_TOKEN_ID

sleep 1
vault secrets move secret secret-v2
vault secrets enable -version=1 -path=secret kv
