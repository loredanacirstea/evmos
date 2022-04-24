RELAYER="evmos1jmghmexanv84dj826gp24l7nfhm2zmrd8987cq"
LOCALKEY="evmos1fjx8p8uzx3h5qszqnwvelulzd659j8uafwws7e"
LOCALKEY2="evmos1f3d3t8y604x9ev4dfgf4hx270gdcrfal2m0hr3"
KEY="mykey"
CHAINID="evmos_9000-2"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
TRACE="--trace"
# TRACE=""
CHAINDIR=$HOME/.evmosd2

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf $CHAINDIR*

make install

~/go/bin/evmosd config keyring-backend $KEYRING --home $CHAINDIR
~/go/bin/evmosd config chain-id $CHAINID --home $CHAINDIR

# if $KEY exists it should be deleted
# ~/go/bin/evmosd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home $CHAINDIR

~/go/bin/evmosd keys add $KEY --recover --keyring-backend $KEYRING --algo $KEYALGO --home $CHAINDIR

# ~/go/bin/evmosd keys add alice --recover --keyring-backend $KEYRING --algo $KEYALGO
# evmos1fjx8p8uzx3h5qszqnwvelulzd659j8uafwws7e
# B0D4244D56065C138CED0F2F4371E66FCB0EEDDCE95AC276AD95B37045960672
# moral depend knock trouble situate credit carry local state meadow they approve desk slender blush lawsuit behave print involve orbit black social cinnamon casino

# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
~/go/bin/evmosd init $MONIKER --chain-id $CHAINID --home $CHAINDIR

# Change parameter token denominations to aevmos
cat $CHAINDIR/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aevmos"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
cat $CHAINDIR/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aevmos"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
cat $CHAINDIR/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aevmos"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
cat $CHAINDIR/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aevmos"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
cat $CHAINDIR/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="aevmos"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# increase block time (?)
cat $CHAINDIR/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# Set gas limit in genesis
cat $CHAINDIR/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# Get close date
node_address=$(~/go/bin/evmosd keys list --home $CHAINDIR | grep  "address: " | cut -c12-)
current_date=$(date -u +"%Y-%m-%dT%TZ")
cat $CHAINDIR/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["airdrop_start_time"]=$current_date' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
# Add account to claims
amount_to_claim=10000
cat $CHAINDIR/config/genesis.json | jq -r --arg node_address "$node_address" --arg amount_to_claim "$amount_to_claim" '.app_state["claims"]["claims_records"]=[{"initial_claimable_amount":$amount_to_claim, "actions_completed":[false, false, false, false],"address":$node_address}]' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

cat $CHAINDIR/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claim"]["params"]["duration_of_decay"]="1000000s"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json
cat $CHAINDIR/config/genesis.json | jq -r --arg current_date "$current_date" '.app_state["claim"]["params"]["duration_until_decay"]="100000s"' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# Claim module account:
# 0xA61808Fe40fEb8B3433778BBC2ecECCAA47c8c47 || evmos15cvq3ljql6utxseh0zau9m8ve2j8erz89m5wkz
cat $CHAINDIR/config/genesis.json | jq -r --arg amount_to_claim "$amount_to_claim" '.app_state["bank"]["balances"] += [{"address":"evmos15cvq3ljql6utxseh0zau9m8ve2j8erz89m5wkz","coins":[{"denom":"aevmos", "amount":$amount_to_claim}]}]' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $CHAINDIR/config/config.toml
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $CHAINDIR/config/config.toml
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $CHAINDIR/config/config.toml
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $CHAINDIR/config/config.toml
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $CHAINDIR/config/config.toml
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $CHAINDIR/config/config.toml
  fi
fi

# sed -i.bak -e "s%^proxy_app = \"tcp://127.0.0.1:26658\"%proxy_app = \"tcp://127.0.0.1:26653\"%; s%^laddr = \"tcp://127.0.0.1:26657\"%laddr = \"tcp://127.0.0.1:26652\"%; s%^pprof_laddr = \"localhost:6060\"%pprof_laddr = \"localhost:6061\"%; s%^laddr = \"tcp://0.0.0.0:26656\"%laddr = \"tcp://0.0.0.0:26651\"%; s%^prometheus_listen_addr = \":26660\"%prometheus_listen_addr = \":26655\"%" $CHAINDIR/config/config.toml sed -i.bak -e "s%^address = \"0.0.0.0:1317\"%address = \"0.0.0.0:1318\"%; s%^address = \"0.0.0.0:9090\"%address = \"0.0.0.0:9092\"%; s%^address = \"0.0.0.0:9091\"%address = \"0.0.0.0:9093\"%" $CHAINDIR/config/app.toml

# sed -i.bak -e "s%^address = \"0.0.0.0:1317\"%address = \"192.168.0.106:1318\"%; s%^address = \"0.0.0.0:9090\"%address = \"0.0.0.0:9092\"%; s%^address = \"0.0.0.0:9091\"%address = \"0.0.0.0:9093\"%" $CHAINDIR/config/app.toml

sed -i.bak -e "s%^address = \"tcp://0.0.0.0:1317\"%address = \"tcp://192.168.0.106:1318\"%" $CHAINDIR/config/app.toml

sed -i -e 's/\"allow_messages\":.*/\"allow_messages\": [\"\/cosmos.bank.v1beta1.MsgSend\", \"\/cosmos.staking.v1beta1.MsgDelegate\", \"\/ethermint.evm.v1.MsgEthereumTx\", \"\/ethermint.evm.v1.MsgEthereumIcaTx\", \"\/ethermint.evm.v1.AccessListTx\", \"\/ethermint.evm.v1.DynamicFeeTx\", \"\/ethermint.evm.v1.LegacyTx\", \"\/ethermint.evm.v1.Msg\", \"\/ethermint.evm.v1.ExtensionOptionsEthereumTx\", \"\/ethermint.types.v1.ExtensionOptionsWeb3Tx\"]/g' $CHAINDIR/config/genesis.json


# sed -i.bak -e "s%^proxy_app = \"tcp://127.0.0.1:26658\"%proxy_app = \"tcp://127.0.0.1:26653\"%; s%^laddr = \"tcp://127.0.0.1:26657\"%laddr = \"tcp://127.0.0.1:26652\"%; s%^pprof_laddr = \"localhost:6060\"%pprof_laddr = \"localhost:6061\"%; s%^laddr = \"tcp://0.0.0.0:26656\"%laddr = \"tcp://0.0.0.0:26651\"%; s%^prometheus_listen_addr = \":26660\"%prometheus_listen_addr = \":26655\"%" config.toml sed -i.bak -e "s%^address = \"0.0.0.0:9090\"%address = \"0.0.0.0:9092\"%; s%^address = \"0.0.0.0:9091\"%address = \"0.0.0.0:9093\"%" app.toml

# ~/go/bin/evmosd start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home $CHAINDIR --grpc-web.address 192.168.0.106:39091 --json-rpc.address 192.168.0.106:38545 --json-rpc.ws-address 192.168.0.106:38546 --rpc.laddr tcp://192.168.0.106:36657 --rpc.pprof_laddr tcp://192.168.0.106:36060 --rpc.grpc_laddr tcp://192.168.0.106:39092 --address tcp://0.0.0.0:36658 --p2p.laddr tcp://0.0.0.0:36656 --grpc.address 0.0.0.0:39090



# sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $CHAINDIR/config/app.toml

# Allocate genesis accounts (cosmos formatted addresses)
# ~/go/bin/evmosd keys delete $RELAYER  --keyring-backend $KEYRING
# ~/go/bin/evmosd keys delete $LOCALKEY  --keyring-backend $KEYRING

~/go/bin/evmosd add-genesis-account $KEY 100000000000000000000000000000000aevmos --keyring-backend $KEYRING --home $CHAINDIR
~/go/bin/evmosd add-genesis-account $KEY 200000000stake --keyring-backend $KEYRING --home $CHAINDIR

~/go/bin/evmosd add-genesis-account $RELAYER 100000000000000000000000000000000aevmos,200000000stake --keyring-backend $KEYRING --home $CHAINDIR
~/go/bin/evmosd add-genesis-account $LOCALKEY 100000000000000000000000000000000aevmos --keyring-backend $KEYRING --home $CHAINDIR
~/go/bin/evmosd add-genesis-account $LOCALKEY2 100000000000000000000000000000000aevmos,200000000stake --keyring-backend $KEYRING --home $CHAINDIR

# ~/go/bin/evmosd keys add $RELAYER --keyring-backend $KEYRING
# ~/go/bin/evmosd keys add $LOCALKEY --keyring-backend $KEYRING
# ~/go/bin/evmosd tx bank send $KEY $RELAYER 100000000000aevmos --chain-id=evmos_9000-1 --keyring-backend $KEYRING
# ~/go/bin/evmosd tx bank send $KEY $RELAYER 200000stake --chain-id=evmos_9000-1 --keyring-backend $KEYRING

# Update total supply with claim values
validators_supply=$(cat $CHAINDIR/config/genesis.json | jq -r '.app_state["bank"]["supply"][0]["amount"]')
# Bc is required to add this big numbers
# total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
total_supply=300000000000000000000000000010000
cat $CHAINDIR/config/genesis.json | jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' > $CHAINDIR/config/tmp_genesis.json && mv $CHAINDIR/config/tmp_genesis.json $CHAINDIR/config/genesis.json

# Sign genesis transaction
~/go/bin/evmosd gentx $KEY 1000000000000000000000aevmos --keyring-backend $KEYRING --chain-id $CHAINID --home $CHAINDIR

# Collect genesis tx
~/go/bin/evmosd collect-gentxs --home $CHAINDIR

# Run this to ensure everything worked and that the genesis file is setup correctly
~/go/bin/evmosd validate-genesis --home $CHAINDIR

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
~/go/bin/evmosd start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home $CHAINDIR --grpc-web.address 192.168.0.106:39091 --json-rpc.address 192.168.0.106:38545 --json-rpc.ws-address 192.168.0.106:38546 --rpc.laddr tcp://192.168.0.106:36657 --rpc.grpc_laddr tcp://192.168.0.106:39092 --address tcp://0.0.0.0:36658 --p2p.laddr tcp://0.0.0.0:36656 --grpc.address 0.0.0.0:39090


# --api.laddr tcp://192.168.0.106:3317
# --rpc.pprof_laddr tcp://192.168.0.106:36060

# --priv_validator_laddr
# --p2p.laddr tcp://0.0.0.0:26656
# --grpc.address 0.0.0.0:9090

# Error: failed to listen on 0.0.0.0:1317: listen tcp 0.0.0.0:1317: bind: address already in use


# ~/go/bin/evmosd start --pruning=nothing --trace --log_level trace --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable

# ~/go/bin/evmosd start --pruning=nothing --trace --log_level trace --minimum-gas-prices=0.0001aevmos --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --grpc-web.address 192.168.0.106:9091 --json-rpc.address 192.168.0.106:8545 --json-rpc.ws-address 192.168.0.106:8546 --rpc.laddr tcp://192.168.0.106:26657
