
# upload the smart contract, then create a validator. Confirm it works

export FURYAD_NODE="tcp://localhost:26657"
FLAGS="--from=furya1 --gas=2500000 --fees=50000ufury --node=http://localhost:26657 --yes --keyring-backend=test --home $HOME/.furya1 --chain-id=local-1 --output=json"

furyad tx wasm store ./keeper/contract/furya_staking_hooks_example.wasm $FLAGS

sleep 5

txhash=$(furyad tx wasm instantiate 1 '{}' --label=furya_staking --no-admin $FLAGS | jq -r .txhash)
sleep 5
addr=$(furyad q tx $txhash --output=json --node=http://localhost:26657 | jq -r .logs[0].events[2].attributes[0].value) && echo $addr

# register addr to staking
furyad tx cw-hooks register staking $addr $FLAGS
furyad q cw-hooks staking-contracts

# furyad tx cw-hooks unregister staking $addr $FLAGS
# furyad q cw-hooks staking-contracts

# get config
furyad q wasm contract-state smart $addr '{"get_config":{}}' --node=http://localhost:26657

# get last validator
furyad q wasm contract-state smart $addr '{"last_val_change":{}}' --node=http://localhost:26657
furyad q wasm contract-state smart $addr '{"last_delegation_change":{}}' --node=http://localhost:26657

# create validator
furyad tx staking create-validator --amount 1ufury --commission-rate="0.05" --commission-max-rate="1.0" --commission-max-change-rate="1.0" --moniker="test123" --from=furya2 --pubkey=$(furyad tendermint show-validator --home $HOME/.furya) --min-self-delegation="1" --gas=1000000 --fees=50000ufury --node=http://localhost:26657 --yes --keyring-backend=test --home $HOME/.furya1 --chain-id=local-1 --output=json

# furyad export --output-document=$HOME/Desktop/export.json --home=$HOME/.furya1