#!/bin/bash
# microtick and bitcanna contributed significantly here.
# Pebbledb state sync script.
set -uxe

# Set Golang environment variables.
export GOPATH=~/go
export PATH=$PATH:~/go/bin

# Initialize chain.
furyad init test

# Get Genesis
wget https://download.dimi.sh/furya-phoenix2-genesis.tar.gz
tar -xvf furya-phoenix2-genesis.tar.gz
mv furya-phoenix2-genesis.json "$HOME/.furya/config/genesis.json"




# Get "trust_hash" and "trust_height".
INTERVAL=1000
LATEST_HEIGHT="$(curl -s https://furya-rpc.polkachu.com/block | jq -r .result.block.header.height)"
BLOCK_HEIGHT="$((LATEST_HEIGHT-INTERVAL))"
TRUST_HASH="$(curl -s "https://furya-rpc.polkachu.com/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)"

# Print out block and transaction hash from which to sync state.
echo "trust_height: $BLOCK_HEIGHT"
echo "trust_hash: $TRUST_HASH"

# Export state sync variables.
export FURYAD_STATESYNC_ENABLE=true
export FURYAD_P2P_MAX_NUM_OUTBOUND_PEERS=200
export FURYAD_STATESYNC_RPC_SERVERS="https://rpc-furya-ia.notional.ventures:443,https://furya-rpc.polkachu.com:443"
export FURYAD_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export FURYAD_STATESYNC_TRUST_HASH=$TRUST_HASH

# Fetch and set list of seeds from chain registry.
FURYAD_P2P_SEEDS="$(curl -s https://raw.githubusercontent.com/cosmos/chain-registry/master/furya/chain.json | jq -r '[foreach .peers.seeds[] as $item (""; "\($item.id)@\($item.address)")] | join(",")')"
export FURYAD_P2P_SEEDS

# Start chain.
furyad start --x-crisis-skip-assert-invariants 
