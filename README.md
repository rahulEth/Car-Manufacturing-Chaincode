## 12-June-2020
### channel list
peer channel list

### channel create
```sh
export CHANNEL_NAME=mychannel && peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/$CHANNEL_NAME.tx --tls --cafile $ORDERER_CA

```

### channel join
```sh
peer channel join -b mychannel.block

```


### chaincode install
```sh

peer chaincode install -n manufacturing -v 1.0 -p github.com/chaincode/manufacturing/

```
### chaincode instantiate
```sh

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

```

### chaincode upgrade 
```sh

peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -v 2.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

```

### check running containers
```sh

docker ps

```
### check running container logs

```sh
docker logs container_name -f --tail 1000
```

### chaincode functions

partManufacturing
```sh
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -c '{"args":["partManufacturing","{\"snumber\":\"w001\",\"prod\":\"wheel\",\"owner\":\"PartFactory\", \"cts\" : \"1534455533\"}"]}'
```
changeOwnership
```sh
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -c '{"args":["changeOwnership","{\"snumber\":\"001\",\"prod\":\"wheel\",\"owner\" :\"carFactory\",\"uts\":\"456787678\"}"]}'
```

createCar
```sh
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -c '{"args":["createCar","{\"snumber\":\"C0001\",\"owner\":\"carFactory\",\"prod\":\"Car\",\"color\":\"Blue\", \"plist\":["W001", "E001", "T001"], \"cts\":\"154234433\", \"uts\" : \"154234433\"}"]}'
```


ChangeCarOwnership
```sh
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n manufacturing -c '{"args":["changeCarOwnership","{\"snumber\":\"C001\",\"owner\" :\"dealer\",\"prod\":\"Car\",\"uts\":\"456787678\"}"]}'
```
