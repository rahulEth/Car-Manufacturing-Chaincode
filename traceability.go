package main

import (
	// "fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	// "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

var  _TraceingLogger = shim.NewLogger("Traceing")

//Events

const _CreateEvent = "PART_MANU"
const _ChangeOwnership = "CHANGE_OWNERSHIP"
const _ChangeCarOwnership = "CHANGECAR_OWNERSHIP"
const _CreateCar = "CREATE_CAR"
// Trace manages all Traceing related transactions
type Trace struct {
}


type Part struct {
	ProductType    string `json:"prod"`   // wheels, engines, transmission
	Manufacturer   string `json:"manu"`
	SerialNumber   string `json:"snumber"`
	CreateTs       string `json:"cts"`
	Owner          string `json:"owner"`
	UpdateTs       string `json:"uts"`
	// //Foremen
	// ForemenUpdate  []ForemenType  `json:"foremenupdate"`
}

type Car struct {
	ProductType    string    `json:"prod"`
	Manufacturer   string    `json:"manu"`
	Owner          string    `json:"owner"`
	Color          string    `json:"color"` 
	SerialNumber   string    `json:"snumber"`
	PartList       []string  `json:"plist"`
	CreateTs       string    `json:"cts"`
	UpdateTs       string    `json:"uts"`

}



// Create the purchase order
func (tr *Trace) PartManufacturing(stub shim.ChaincodeStubInterface) peer.Response {
	_TraceingLogger.Infof("PartManufacturing")
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
       return shim.Error("Invalid number of arguments provided")
	}	
	var partManu Part
	err := json.Unmarshal([]byte(args[0]), &partManu)
	if err != nil {
		return shim.Error("Invalid json provided as input")
	}
	if recordBytes, _ :=stub.GetState(partManu.SerialNumber); len(recordBytes) > 0 {
        return shim.Error("Part already exist please provide unique Serial Number")
	}
	partManu.Manufacturer = "PartFactory"
	orderJson, _ := json.Marshal(partManu)

	_TraceingLogger.Infof("partManu.SerialNumber..........", partManu.SerialNumber)

	erre := stub.PutState(partManu.SerialNumber, orderJson)

	if erre !=nil {
		_TraceingLogger.Errorf("Unable to save with SerialNumber " + partManu.SerialNumber)
		return shim.Error("Unable to save with SerialNumber " + partManu.SerialNumber)
	}
	_TraceingLogger.Infof("PartManufacturing : PutState Success : " + string(orderJson))
	erer := stub.SetEvent(_CreateEvent, orderJson)

	if erer != nil {
		_TraceingLogger.Errorf("Event not generated for event : PART_MANUFACTURING")
		// return shim.Error("{\"error\" : \"Unable to generate Purchase Order Event.\"}")
	}
	resultData := map[string]interface{}{
		"trxID" : stub.GetTxID(),
		"POID"  : partManu.SerialNumber,
		"message" : "Part Manufacturing successful",
		"Order" : partManu,
		"status" : true,
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

// chaange ownership of the parts from partFactory to Car Manufacturer

func (tr *Trace) ChangeOwnership(stub shim.ChaincodeStubInterface) peer.Response {

	_ , args := stub.GetFunctionAndParameters()
	
	if len(args) < 1 {
		_TraceingLogger.Errorf("changeOwnership : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var changeOwner Part
	err := json.Unmarshal([]byte(args[0]), &changeOwner)
	if err != nil {
		_TraceingLogger.Errorf("changeOwnership : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}

	RecordsBytes, _ := stub.GetState(changeOwner.SerialNumber)

	if len(RecordsBytes) <=0 {
		_TraceingLogger.Errorf("Part does not exist with serial number : " + changeOwner.SerialNumber)
		return shim.Error("Part does not exist with serial number : " + changeOwner.SerialNumber)
	}
    
	var part Part
	
	err1 := json.Unmarshal(RecordsBytes, &part)
	if err1 != nil {
		return shim.Error("Unmarshaling Error")
	}
	if part.ProductType != changeOwner.ProductType {
		_TraceingLogger.Errorf("product Type does not match with product serial number : " + changeOwner.SerialNumber)
		return shim.Error("product Type does not match with product serial number")
	}
	part.Owner = changeOwner.Owner
	part.UpdateTs = changeOwner.UpdateTs

	OrderBytes, err2 := json.Marshal(part)

	if err2 !=nil {
		_TraceingLogger.Errorf("changeOwnership : Marshalling Error : " + string(err2.Error()))
		return shim.Error("changeOwnership : Marshalling Error : " + string(err2.Error()))
	}
	_TraceingLogger.Infof("changeOwnership : saving the Part after ownership change : " + changeOwner.SerialNumber)

	errorr :=stub.PutState(changeOwner.SerialNumber, OrderBytes)

	if errorr != nil {
		_TraceingLogger.Errorf("changeOwnership : Put State Failed Error : " + string(errorr.Error())) 
		return shim.Error("Put State Failed Error : " + string(errorr.Error()))
	}

	_TraceingLogger.Infof("changeOwnership : PutState Success : " + string(OrderBytes))
	err3 := stub.SetEvent(_ChangeOwnership, OrderBytes)
	if err3 != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Event not generating for : " + _ChangeOwnership)
	
	}
	
	resultData := map[string]interface{}{
		"trxID" : stub.GetTxID(),
		"POID"  : changeOwner.SerialNumber,
		"message" : "Ownership change successful",
		"Part" : part,
		"status" : true,
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}


func (tr *Trace) CreateCar(stub shim.ChaincodeStubInterface) peer.Response{
	_TraceingLogger.Infof("createCar")
	_, args := stub.GetFunctionAndParameters()
	
	if len (args) < 1 {
		_TraceingLogger.Errorf("createCar : Invalid number of argumnet provided")
		return shim.Error("Invalid number of argumnet provided")
	}
	var createCar Car
	err := json.Unmarshal([]byte(args[0]), &createCar)
	if err != nil {
		_TraceingLogger.Errorf("createCar : invalid json provided as input")
		return shim.Error("invalid json provided as input")
	}
    for index, _ := range createCar.PartList{
		RecordsBytes, _ := stub.GetState(createCar.PartList[index])
		
		if len(RecordsBytes) <=0 {
			_TraceingLogger.Errorf("Part does not exist with serial number : " + createCar.PartList[index])
			return shim.Error("Part does not exist with serial number : " + createCar.PartList[index])
		}
		
		var part Part
		
		err1 := json.Unmarshal(RecordsBytes, &part)
		if err1 != nil {
			return shim.Error("Unmarshaling Error")
		}
		if part.Owner != "carFactory" {
			_TraceingLogger.Errorf("part owner is not carFacotry with serial number : " + createCar.PartList[index])
		    return shim.Error("part owner is not carFacotry with serial number : " + createCar.PartList[index])
		} 
	}

	createCar.Manufacturer =  "carFactory"
	
	OrderByte, err := json.Marshal(createCar)
	
	if err != nil {
		_TraceingLogger.Errorf("createCar : Marshling Error")
		return shim.Error("createCar : Marshling Error : " + string(err.Error()))	
	}

	_TraceingLogger.Infof("saving the manufactured Car : " + createCar.SerialNumber)

	erroor := stub.PutState(createCar.SerialNumber, OrderByte)

	if erroor != nil {
		  _TraceingLogger.Infof("createCar : putState error : " + string(erroor.Error()))
		  return shim.Error("putState Fail : " + string(erroor.Error()))
	}
	_TraceingLogger.Infof("createCar : Put State success : " + string(OrderByte))

	err2 := stub.SetEvent(_CreateCar, OrderByte)
	
	if err2 != nil {
		_TraceingLogger.Infof("createCar Event not generated for : " + _CreateCar)
	}

		
	resultData := map[string]interface{}{
		"trxID" : stub.GetTxID(),
		"POID"  : createCar.SerialNumber,
		"message" : "Car created successful",
		"Car" : createCar,
		"status" : true,
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)

}

//change car ownership from carFactory to dealer

func (tr *Trace) ChangeCarOwnership(stub shim.ChaincodeStubInterface) peer.Response {

	_ , args := stub.GetFunctionAndParameters()
	
	if len(args) < 1 {
		_TraceingLogger.Errorf("changeCarOwnership : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var changeOwner Car
	err := json.Unmarshal([]byte(args[0]), &changeOwner)
	if err != nil {
		_TraceingLogger.Errorf("changeCarOwnership : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}

	RecordsBytes, _ := stub.GetState(changeOwner.SerialNumber)

	if len(RecordsBytes) <=0 {
		_TraceingLogger.Errorf("Car does not exist with serial number : " + changeOwner.SerialNumber)
		return shim.Error("Car does not exist with serial number : " + changeOwner.SerialNumber)
	}
    
	var car Car
	
	err1 := json.Unmarshal(RecordsBytes, &car)
	if err1 != nil {
		return shim.Error("Unmarshaling Error")
	}
	if car.ProductType != changeOwner.ProductType {
		_TraceingLogger.Errorf("product Type does not match with product serial number : " + changeOwner.SerialNumber)
		return shim.Error("product Type does not match with product serial number")
	}
	car.Owner = changeOwner.Owner
	car.UpdateTs = changeOwner.UpdateTs

	OrderBytes, err2 := json.Marshal(car)

	if err2 !=nil {
		_TraceingLogger.Errorf("changeCarOwnership : Marshalling Error : " + string(err2.Error()))
		return shim.Error("changeCarOwnership : Marshalling Error : " + string(err2.Error()))
	}
	_TraceingLogger.Infof("changeCarOwnership : saving the Car after ownership change : " + changeOwner.SerialNumber)

	errorr :=stub.PutState(changeOwner.SerialNumber, OrderBytes)

	if errorr != nil {
		_TraceingLogger.Errorf("changeCarOwnership : Put State Failed Error : " + string(errorr.Error())) 
		return shim.Error("Put State Failed Error : " + string(errorr.Error()))
	}

	_TraceingLogger.Infof("changeCarOwnership : PutState Success : " + string(OrderBytes))
	err3 := stub.SetEvent(_ChangeCarOwnership, OrderBytes)
	if err3 != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Event not generating for : " + _ChangeCarOwnership)
	
	}
	
	resultData := map[string]interface{}{
		"trxID" : stub.GetTxID(),
		"POID"  : changeOwner.SerialNumber,
		"message" : "Ownership change successful",
		"Car" : car,
		"status" : true,
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}
