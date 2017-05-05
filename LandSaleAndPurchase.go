package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//const BUYER = "buyer"
//const SELLER = "seller"
//const CORPORATION = "seller"

type SimpleChaincode struct {
}

//Information to be stored about land in blockchain network
type Land struct {
	LandId    string `json:"landId"`
	Plotno    string `json:"plotno"`
	Locality  string `json:"locality"`
	Area      string `json:"area"`
	OwnerName string `json:"OwnerName"`
}

// this Id will be used to check if that Land already is created by corporation and exists in DB or not
//type LAND_Holder struct {
//LANDs []string `json:"lands"`
//}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

//Init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var bytes1 []byte
	var bytes2 []byte
	var bytes3 []byte
	var bytes4 []byte
	var err error

	var land1 Land
	var land2 Land
	var land3 Land
	var land4 Land

	land1.Area = "105SqFt"
	land1.LandId = "L2"
	land1.Locality = "Baner"
	land1.Plotno = "1"
	land1.OwnerName = "Sumit"

	land4.Area = "100SqFt"
	land4.LandId = "L4"
	land4.Locality = "Baner"
	land4.Plotno = "4"
	land4.OwnerName = "Ishani"

	land3.Area = "100SqFt"
	land3.LandId = "L3"
	land3.Locality = "Aundh"
	land3.Plotno = "3"
	land3.OwnerName = "Parth"

	land2.Area = "1056SqFt"
	land2.LandId = "L2"
	land2.Locality = "Aundh"
	land2.Plotno = "2"
	land2.OwnerName = "Shaily"

	bytes1, err = json.Marshal(land1)
	bytes2, err = json.Marshal(land2)
	bytes3, err = json.Marshal(land3)
	bytes4, err = json.Marshal(land4)

	err = stub.PutState(land1.LandId, bytes1)

	err = stub.PutState(land2.LandId, bytes2)
	err = stub.PutState(land3.LandId, bytes3)
	err = stub.PutState(land4.LandId, bytes4)

	//bytes, err := json.Marshal(landIds)

	if err != nil {
		return nil, errors.New("Error creating  records")
	}

	//err = stub.PutState("landIds", bytes)

	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getLandInfo" {

		return getLandInfo(stub, args)
	}

	return nil, nil

}

func getLandInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var land Land

	var landId = args[0]
	bytes, err := stub.GetState(landId)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &land)

	return bytes, nil
}
