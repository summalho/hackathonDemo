package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"
	"strconv"
	"time"
)

type SimpleChaincode struct {
}

type OWNER_ID_Holder struct {
	OWNER_IDs []string `json:"OWNER_IDs"`
}
type PROPERTY_ID_Holder struct {
	PROPERTY_IDs []string `json:"PROPERTY_IDs"`
}
type PROPERTY_HISTORY struct {
	PROPERTY_HISTORY_IDs []string `json:"PROPERTY_IDs"`
}

//Information to be stored about land in blockchain network
type Property struct {
	OwnerId    string `json:"ownerId"`
	PropertyId int    `json:"propertyId"`
	Plotno     string `json:"plotno"`
	City       string `json:"City"`
	Area       string `json:"area"`
	Pincode    string `json:"Pincode"`
	Latitude   string `json:"Latitude"`
	Longitude  string `json:"Longitude"`
}

type Owner struct {
	Id           int    `json:"id"` //generated by blockchain
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AadharNumber string `json:"aadharNumber"`
}

type PropertyHistory struct {
	HistoryId       int    `json:"historyId"`
	PropertyId      string `json:"propertyId"` //generated by blockchain
	OwnerId         string `json:"ownerId"`
	AgreementDate   string `json:"agreementDate"`
	AgreementAmount string `json:"agreementAmount"`
}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

//Init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	var ownerBytes []byte
	var properyBytes []byte
	var propertyHistoryBytes []byte

	var ownerIds OWNER_ID_Holder
	var propertyIds PROPERTY_ID_Holder
	var propertyHistoryHolder PROPERTY_HISTORY

	ownerBytes, err = json.Marshal(ownerIds)
	properyBytes, err = json.Marshal(propertyIds)
	propertyHistoryBytes, err = json.Marshal(propertyHistoryHolder)

	if err != nil {
		return nil, errors.New("Error creating OWNER_ID_Holder record")
	}

	err = stub.PutState("owner_Ids", ownerBytes)
	err = stub.PutState("property_Ids", properyBytes)
	err = stub.PutState("property_history_Holder", propertyHistoryBytes)

	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createOwner" {

		return createOwner(stub, args)
	}
	if function == "createProperty" {

		return createProperty(stub, args)
	}
	if function == "transferOwnerShip" {

		return transferOwnerShip(stub, args)
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "listRegisteredOwners" {

		return listRegisteredOwners(stub, args)
	}
	if function == "listRegisteredProperties" {

		return listRegisteredProperties(stub, args)
	}
	if function == "listRegisteredPropertiesByOwnwer" {

		return listRegisteredPropertiesByOwnwer(stub, args)
	}

	if function == "getOwnerById" {

		return getOwnerById(stub, args)
	}

	if function == "getIds" {

		return getIds(stub, args)
	}

	if function == "getPropertyIds" {

		return getPropertyIds(stub, args)
	}

	return nil, nil

}

func transferOwnerShip(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var propertyfetched Property
	var propbytes []byte

	propertyId := args[0]
	agreementAmount := args[2]

	bytes, err := stub.GetState(propertyId)
	if err != nil {
		return nil, errors.New("Get state using Property Id failed")
	}

	err = json.Unmarshal(bytes, &propertyfetched)

	if err != nil {
		return nil, errors.New("Unmarshalling error")
	}

	propertyfetched.OwnerId = args[3]

	propbytes, err = json.Marshal(propertyfetched)

	stub.PutState(propertyId, propbytes)
	// transaction to change the owner of the property completes here

	savePropertyInHistory(stub, propertyfetched, agreementAmount)

	return propbytes, nil

}

func getOwnerById(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	bytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, err
	}
	return bytes, nil

}

func getIds(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	bytes, err := stub.GetState("owner_Ids")

	if err != nil {
		return nil, err
	}
	return bytes, nil

}

func getPropertyIds(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	bytes, err := stub.GetState("property_Ids")

	if err != nil {
		return nil, err
	}
	return bytes, nil

}

func listRegisteredOwners(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside list generated Users")

	bytes, err := stub.GetState("owner_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var ownerIDHolder OWNER_ID_Holder
	err = json.Unmarshal(bytes, &ownerIDHolder)

	result := "["

	var temp []byte
	var o Owner

	for _, own := range ownerIDHolder.OWNER_IDs {

		fmt.Println("Inside for loop for getting Owner. Owner Id is  ", own)

		o, err = retrieveOwner(stub, own)

		temp, err = json.Marshal(o)

		if err == nil {
			result += string(temp) + ","
		}

	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

func listRegisteredPropertiesByOwnwer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	ownerId := args[0]

	bytes, err := stub.GetState("property_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var propertyIdHolder PROPERTY_ID_Holder
	err = json.Unmarshal(bytes, &propertyIdHolder)

	result := "["

	var temp []byte
	var p Property

	for _, pro := range propertyIdHolder.PROPERTY_IDs {

		fmt.Println("Inside for loop for getting Property. Property Id is  ", pro)

		p, err = retrieveProperty(stub, pro)

		if p.OwnerId == ownerId {
			temp, err = json.Marshal(p)

			if err == nil {
				result += string(temp) + ","
			}

		}

	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil

}

func listRegisteredProperties(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside list registered Properties")

	bytes, err := stub.GetState("property_Ids")

	fmt.Println("Ids recieved", string(bytes))
	var propertyIdHolder PROPERTY_ID_Holder
	err = json.Unmarshal(bytes, &propertyIdHolder)

	result := "["

	var temp []byte
	var p Property

	for _, pro := range propertyIdHolder.PROPERTY_IDs {

		fmt.Println("Inside for loop for getting Property. Property Id is  ", pro)

		p, err = retrieveProperty(stub, pro)

		temp, err = json.Marshal(p)

		if err == nil {
			result += string(temp) + ","
		}

	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

func retrieveProperty(stub shim.ChaincodeStubInterface, propertyId string) (Property, error) {

	fmt.Println("Inside retrieve Property")

	var p Property

	bytes, err := stub.GetState(propertyId)

	fmt.Println("Owner Id is ", propertyId, "and owner details are ", string(bytes))

	if err != nil {
		return p, errors.New("Owner not found")
	}

	err = json.Unmarshal(bytes, &p)

	return p, nil

}

func retrieveOwner(stub shim.ChaincodeStubInterface, ownerId string) (Owner, error) {

	fmt.Println("Inside retrieve Owner")

	var o Owner

	bytes, err := stub.GetState(ownerId)

	fmt.Println("Owner Id is ", ownerId, "and owner details are ", string(bytes))

	if err != nil {
		return o, errors.New("Owner not found")
	}

	err = json.Unmarshal(bytes, &o)

	return o, nil

}

func createOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Inside create user")

	ownerDetails := Owner{}

	ownerDetails.FirstName = args[0]
	ownerDetails.LastName = args[1]
	ownerDetails.AadharNumber = args[2]

	ownerDetails.Id = generateUserId()

	ownerDetailBytes, err := json.Marshal(ownerDetails)

	fmt.Println("Owner Details are : ", string(ownerDetailBytes))

	if err != nil {
		return nil, errors.New("Problem while saving Owner Details in BlockChain Network")

	}

	var id = strconv.Itoa(ownerDetails.Id)

	err = stub.PutState(id, ownerDetailBytes)

	//now owner has been added to block chain network, now we have to save the  Id as well

	bytes, err := stub.GetState("owner_Ids")

	var newOwnerId OWNER_ID_Holder

	err = json.Unmarshal(bytes, &newOwnerId)

	if err != nil {
		return nil, errors.New("error unmarshalling new Owner Id")
	}

	newOwnerId.OWNER_IDs = append(newOwnerId.OWNER_IDs, id)

	bytes, err = json.Marshal(newOwnerId)

	if err != nil {

		return nil, errors.New("error marshalling new Owner Id")
	}

	err = stub.PutState("owner_Ids", bytes)
	fmt.Println("Owner Id Saved is ", string(bytes))

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil
}

func generateUserId() int {

	min := 1
	max := 100000

	rand.Seed(time.Now().Unix())
	return (rand.Intn(max-min) + min)

}

func generatePropertyHistoryId() int {

	min := 100001
	max := 999999

	rand.Seed(time.Now().Unix())
	return (rand.Intn(max-min) + min)

}

func savePropertyInHistory(stub shim.ChaincodeStubInterface, propertyDetails Property, agreementAmount string) {

	var propHis PropertyHistory

	propId := propertyDetails.PropertyId
	propIdSTR := strconv.Itoa(propId)

	t := time.Now()
	timestr := t.Format("2006-01-02 15:04:05")

	propHis.PropertyId = propIdSTR
	propHis.AgreementDate = timestr
	propHis.AgreementAmount = agreementAmount

	propHis.OwnerId = propertyDetails.OwnerId

	propHis.HistoryId = generatePropertyHistoryId()

	hsitoryBytes, err := json.Marshal(propHis)

	var id = strconv.Itoa(propHis.HistoryId)

	stub.PutState(id, hsitoryBytes)
	//history added

	bytes, err := stub.GetState("property_history_Holder")
	var newPropertyhistory PROPERTY_HISTORY

	err = json.Unmarshal(bytes, &newPropertyhistory)

	newPropertyhistory.PROPERTY_HISTORY_IDs = append(newPropertyhistory.PROPERTY_HISTORY_IDs, id)

	bytes, err = json.Marshal(newPropertyhistory)

	err = stub.PutState("property_history_Holder", bytes)
	fmt.Println("Address Id Saved is ", string(bytes))

	if err != nil {
		fmt.Println(err.Error())
	}

}

func createProperty(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Inside create property")

	propertyDetails := Property{}

	propertyDetails.OwnerId = args[0]
	propertyDetails.Area = args[1]
	propertyDetails.City = args[2]
	propertyDetails.Pincode = args[3]
	propertyDetails.Plotno = args[4]
	propertyDetails.Longitude = args[5]
	propertyDetails.Latitude = args[6]
	propertyDetails.PropertyId = generateUserId()

	propertyDetailsBytes, err := json.Marshal(propertyDetails)

	fmt.Println("Property Details are : ", string(propertyDetailsBytes))

	if err != nil {
		return nil, errors.New("Problem while saving Owner Details in BlockChain Network")

	}

	var id = strconv.Itoa(propertyDetails.PropertyId)

	err = stub.PutState(id, propertyDetailsBytes)

	savePropertyInHistory(stub, propertyDetails, "0")

	//now property Details has been added to block chain network, now we have to save the  Id as well

	bytes, err := stub.GetState("property_Ids")

	var newPropertyId PROPERTY_ID_Holder

	err = json.Unmarshal(bytes, &newPropertyId)

	if err != nil {
		return nil, errors.New("error unmarshalling new Property Address")
	}

	newPropertyId.PROPERTY_IDs = append(newPropertyId.PROPERTY_IDs, id)

	bytes, err = json.Marshal(newPropertyId)

	if err != nil {

		return nil, errors.New("error marshalling new Property Address")
	}

	err = stub.PutState("property_Ids", bytes)
	fmt.Println("Address Id Saved is ", string(bytes))

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil

}
