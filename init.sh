RELAYER="evmos1jmghmexanv84dj826gp24l7nfhm2zmrd8987cq"
LOCALKEY="evmos1fjx8p8uzx3h5qszqnwvelulzd659j8uafwws7e"
LOCALKEY2="evmos1f3d3t8y604x9ev4dfgf4hx270gdcrfal2m0hr3"
KEY="mykey"
CHAINID="evmos_9000-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
TRACE="--trace"
# TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.evmosd*

make install

~/go/bin/evmosd config keyring-backend $KEYRING
~/go/bin/evmosd config chain-id $CHAINID

# if $KEY exists it should be deleted
# ~/go/bin/evmosd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

~/go/bin/evmosd keys add $KEY --recover --keyring-backend $KEYRING --algo $KEYALGO
# evmos1fjx8p8uzx3h5qszqnwvelulzd659j8uafwws7e
# B0D4244D56065C138CED0F2F4371E66FCB0EEDDCE95AC276AD95B37045960672
# moral depend knock trouble situate credit carry local state meadow they approve desk slender blush lawsuit behave print involve orbit black social cinnamon casino

# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
~/go/bin/evmosd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to aevmos
cat $HOME/.evmosd/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aevmos"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
cat $HOME/.evmosd/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aevmos"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
cat $HOME/.evmosd/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aevmos"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
cat $HOME/.evmosd/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aevmos"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
cat $HOME/.evmosd/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="aevmos"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# increase block time (?)
cat $HOME/.evmosd/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# Set gas limit in genesis
cat $HOME/.evmosd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# Get close date
node_address=$(~/go/bin/evmosd keys list | grep  "address: " | cut -c12-)
current_date=$(date -u +"%Y-%m-%dT%TZ")
cat $HOME/.evmosd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["airdrop_start_time"]=$current_date' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
# Add account to claims
amount_to_claim=10000
cat $HOME/.evmosd/config/genesis.json | jq -r --arg node_address "$node_address" --arg amount_to_claim "$amount_to_claim" '.app_state["claims"]["claims_records"]=[{"initial_claimable_amount":$amount_to_claim, "actions_completed":[false, false, false, false],"address":$node_address}]' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

cat $HOME/.evmosd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claim"]["params"]["duration_of_decay"]="1000000s"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json
cat $HOME/.evmosd/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claim"]["params"]["duration_until_decay"]="100000s"' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# Claim module account:
# 0xA61808Fe40fEb8B3433778BBC2ecECCAA47c8c47 || evmos15cvq3ljql6utxseh0zau9m8ve2j8erz89m5wkz
cat $HOME/.evmosd/config/genesis.json | jq -r --arg amount_to_claim "$amount_to_claim" '.app_state["bank"]["balances"] += [{"address":"evmos15cvq3ljql6utxseh0zau9m8ve2j8erz89m5wkz","coins":[{"denom":"aevmos", "amount":$amount_to_claim}]}]' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.evmosd/config/config.toml
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.evmosd/config/config.toml
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.evmosd/config/config.toml
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.evmosd/config/config.toml
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.evmosd/config/config.toml
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.evmosd/config/config.toml
  fi
fi

sed -i -e 's/\"allow_messages\":.*/\"allow_messages\": [\"\/cosmos.bank.v1beta1.MsgSend\", \"\/cosmos.staking.v1beta1.MsgDelegate\", \"\/ethermint.evm.v1.MsgEthereumTx\", \"\/ethermint.evm.v1.MsgEthereumIcaTx\", \"\/ethermint.evm.v1.AccessListTx\", \"\/ethermint.evm.v1.DynamicFeeTx\", \"\/ethermint.evm.v1.LegacyTx\", \"\/ethermint.evm.v1.Msg\", \"\/ethermint.evm.v1.ExtensionOptionsEthereumTx\", \"\/ethermint.types.v1.ExtensionOptionsWeb3Tx\"]/g' $HOME/.evmosd/config/genesis.json

# sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.evmosd/config/app.toml

# Allocate genesis accounts (cosmos formatted addresses)
# ~/go/bin/evmosd keys delete $RELAYER  --keyring-backend $KEYRING
# ~/go/bin/evmosd keys delete $LOCALKEY  --keyring-backend $KEYRING

~/go/bin/evmosd add-genesis-account $KEY 100000000000000000000000000000000aevmos --keyring-backend $KEYRING
~/go/bin/evmosd add-genesis-account $KEY 200000000stake --keyring-backend $KEYRING

~/go/bin/evmosd add-genesis-account $RELAYER 100000000000000000000000000000000aevmos,200000000stake --keyring-backend $KEYRING
~/go/bin/evmosd add-genesis-account $LOCALKEY 100000000000000000000000000000000aevmos --keyring-backend $KEYRING
~/go/bin/evmosd add-genesis-account $LOCALKEY2 100000000000000000000000000000000aevmos,200000000stake --keyring-backend $KEYRING

# ~/go/bin/evmosd keys add $RELAYER --keyring-backend $KEYRING
# ~/go/bin/evmosd keys add $LOCALKEY --keyring-backend $KEYRING
# ~/go/bin/evmosd tx bank send $KEY $RELAYER 100000000000aevmos --chain-id=evmos_9000-1 --keyring-backend $KEYRING
# ~/go/bin/evmosd tx bank send $KEY $RELAYER 200000stake --chain-id=evmos_9000-1 --keyring-backend $KEYRING

# Update total supply with claim values
validators_supply=$(cat $HOME/.evmosd/config/genesis.json | jq -r '.app_state["bank"]["supply"][0]["amount"]')
# Bc is required to add this big numbers
# total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
total_supply=300000000000000000000000000010000
cat $HOME/.evmosd/config/genesis.json | jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' > $HOME/.evmosd/config/tmp_genesis.json && mv $HOME/.evmosd/config/tmp_genesis.json $HOME/.evmosd/config/genesis.json

# Sign genesis transaction
~/go/bin/evmosd gentx $KEY 1000000000000000000000aevmos --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
~/go/bin/evmosd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
~/go/bin/evmosd validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
~/go/bin/evmosd start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc-web.address 192.168.0.106:9091 --json-rpc.address 192.168.0.106:8545 --json-rpc.ws-address 192.168.0.106:8546 --rpc.laddr tcp://192.168.0.106:26657



# ~/go/bin/evmosd start --pruning=nothing --trace --log_level trace --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable

# ~/go/bin/evmosd start --pruning=nothing --trace --log_level trace --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc-web.address 192.168.0.106:9091 --json-rpc.address 192.168.0.106:8545 --json-rpc.ws-address 192.168.0.106:8546 --rpc.laddr tcp://192.168.0.106:26657
