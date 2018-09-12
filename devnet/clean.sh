rm -Rf node*/tomo*
rm -Rf logs/*.txt
../build/bin/tomo --datadir node1/ init genesis.json
../build/bin/tomo --datadir node2/ init genesis.json
../build/bin/tomo --datadir node3/ init genesis.json
../build/bin/tomo --datadir node4/ init genesis.json

