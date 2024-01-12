1. Create a project `vehicle lifecycle management` with two org
2. Use PDC to share data
3. create pdc data in org1 and access in org2
4. Implement chaincode lifecycle and pdc  


Vehicle Lifetime Management::>
+ Type::>     Vehicle Name        === Name   
+ ID::>       Vehicle number      === RegNum
+ Color::>    Vehicle Company     === Company
+ Size::>     Vehicle Year of Reg === MfgYear
+ Owner::>    Vehicle Owner       === Owner

+ ID::>       Vehicle number      === RegNum
+ AppraisedValue::> Vehicle Life  === Life


---------------------------------
--------------STEPS--------------
---------------------------------

>> ./network.sh up createChannel -ca -s couchdb
>> ./network.sh deployCC -ccn private -ccp ../asset-transfer-pdc-pps/chaincode-go/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ../asset-transfer-pdc-pps/chaincode-go/collections_config.json

>> export PATH=${PWD}/../bin:${PWD}:$PATH
>> export FABRIC_CFG_PATH=$PWD/../config/

>> export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/
>> fabric-ca-client register --caname ca-org1 --id.name owner --id.secret ownerpw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
>> fabric-ca-client enroll -u https://owner:ownerpw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/owner@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
>> cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/owner@org1.example.com/msp/config.yaml"
>> export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org2.example.com/
>> fabric-ca-client register --caname ca-org2 --id.name buyer --id.secret buyerpw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/org2/tls-cert.pem"
>> fabric-ca-client enroll -u https://buyer:buyerpw@localhost:8054 --caname ca-org2 -M "${PWD}/organizations/peerOrganizations/org2.example.com/users/buyer@org2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org2/tls-cert.pem"
>> cp "${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org2.example.com/users/buyer@org2.example.com/msp/config.yaml"

>> export PATH=${PWD}/../bin:$PATH
>> export FABRIC_CFG_PATH=$PWD/../config/
>> export CORE_PEER_TLS_ENABLED=true
>> export CORE_PEER_LOCALMSPID="Org1MSP"
>> export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
>> export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/owner@org1.example.com/msp
>> export CORE_PEER_ADDRESS=localhost:7051

>> export ASSET_PROPERTIES=$(echo -n "{\"vehicleName\":\"Harrier\",\"vehicleNumber\":\"DL00XXXX\",\"vehicleCompany\":\"TATA\",\"vehicleMfgYear\":2020,\"vehicleLife\":14}" | base64 | tr -d \\n)
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n private -c '{"function":"CreateAsset","Args":[]}' --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"

>> peer chaincode query -C mychannel -n private -c '{"function":"ReadAsset","Args":["DL00XXXX"]}'

