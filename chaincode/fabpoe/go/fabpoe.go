/*
 * DocSHARE source chaincode
 */

/******************************** AMENDMENT HISTORY **********************************************/
/*	DATE		Version		Description
/******************************** AMENDMENT HISTORY **********************************************/
/*	Sep 2018	0		Initial Development
/*	
/*	15 Feb 2019	1		Changes wrt fabric new release	
/*
/*
/******************************** AMENDMENT HISTORY **********************************************/



package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
//	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Doc structure, with 4 properties.  Structure tags are used by encoding/json library
type Asset struct {
	AssetId   			string 		`json:"assetId"`	//this is the 40 byte hash of the document
	TmspStart 			string 		`json:"tmspStart"`
	TmspEnd 			string 		`json:"tmspEnd"`
	AccessList			[]AssetACL 	`json:"accessList"`
	//DocHash				string      	`json:"docHash"`
	OwnerId				string      	`json:"ownerId"`	
	DocTitle			string      	`json:"docTitle"`	
}

type AssetACL struct {

	UserId   		string 		`json:"userId"`
	AccessGrantTmsp 	string      	`json:"accessGrantTmsp"`
	UserDescription		string 		`json:"userDescription"` 
	AccessToFields		string 		`json:"accessToFields"` //optional as we have too few fields.map with boolean 
}


type internalResp struct {
		Status 		int
		Message 	string
		Payload		Asset
}


//List of documents the logged in user has access to 
type UserACL struct {

	UserId   		string 		`json:"userId"`
	AssetId			[]string 	`json:"assetId"` 
}




/*
 * The Init method is called when the Smart Contract "fabDoc" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabDoc"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	fmt.Println("Invoke is running: " + function)
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryDoc" {						//query a doc by Id
		return s.queryDoc(APIstub, args)
	} else if function == "initLedger" {			//?
		return s.initLedger(APIstub)
	} else if function == "createDoc" {				//create a new doc
		return s.createDoc(APIstub, args)
	} else if function == "queryAllDocs" {			//query all docs
		return s.queryAllDocs(APIstub)
	} else if function == "changeDocOwner" {		//owner update
		return s.changeDocOwner(APIstub, args)
	} else if function == "setExpiryOnDoc" {				//end timestamp the document. no more valid after this.
		return s.setExpiryOnDoc(APIstub, args)
	} else if function == "querySchema" {				//end timestamp the document. no more valid after this.
		return s.querySchema(APIstub, args)				//returns the empty asset structure to caller 
	} else if function == "grantAccess" {				//end timestamp the document. no more valid after this.
		return s.grantAccess(APIstub, args)				//returns the empty asset structure to caller 
	} else if function == "getDocHistory" {				//end timestamp the document. no more valid after this.
		return s.getDocHistory(APIstub, args)				//returns the empty asset structure to caller 
	} else if function == "executeRichQuery" {				//end timestamp the document. no more valid after this.
		return s.executeRichQuery(APIstub, args)				//returns the empty asset structure to caller 
	} else if function == "queryAllDocsShared" {				//end timestamp the document. no more valid after this.
		return s.queryAllDocsShared(APIstub, args)				//NOTE: the calling userId should be extracted from 
											//private key using cid package 
	} 

	fmt.Println(function + " Not found")
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	DocAsBytes, _ := APIstub.GetState(args[0])
	if DocAsBytes == nil {

		
		fmt.Println("###############DocAsBytes Empty")
		fmt.Println("No record found for assetId=" + args[0])
		return shim.Error("No record found for assetId=" + args[0])
	}

	fmt.Println("queryDoc Success")
	return shim.Success(DocAsBytes)
}

func (s *SmartContract) querySchema(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//return an empty asset structure to the caller
	asset := Asset{}
	acl := AssetACL{}
	asset.AccessList = append(asset.AccessList, acl)
	assetAsBytes, _ := json.Marshal(asset)
		
	fmt.Println("querySchema: printing marshalled schema asset")
	fmt.Println(string(assetAsBytes))

	
	return shim.Success(assetAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	var acl []AssetACL
	
	acl = append(acl, AssetACL{})

	assets := []Asset{
		Asset{AssetId: "ABCDEF0123456789ABCDEF0123456789000000A1", 	TmspStart: "", 	TmspEnd: "", 	DocTitle: "Partnership agreement XTECH", OwnerId: "initLedger", AccessList :acl},
		Asset{AssetId: "ABCDEF0123456789ABCDEF0123456789000000A3", 	TmspStart: "", 	TmspEnd: "", 	DocTitle: "Lease Agreement Suresh", OwnerId: "initLedger", AccessList :acl},
		Asset{AssetId: "ABCDEF0123456789ABCDEF0123456789000000A2", 	TmspStart: "", 	TmspEnd: "", 	DocTitle: "NDA for Docify", OwnerId: "initLedger", AccessList :acl},
	}

	i := 0
	for i < len(assets) {
		fmt.Println("i is ", i)
		assetAsBytes, _ := json.Marshal(assets[i])
		APIstub.PutState(assets[i].AssetId, assetAsBytes)
		fmt.Println("Added", assets[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var docIn Asset
	var err error

	//check number of arguments recvd
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err = json.Unmarshal([]byte(args[0]),&docIn)
	if err != nil {
		return shim.Error("createDoc: Unmarshal 1 failed" + err.Error())	
		fmt.Printf("\nRecvd Input string=%s",args[0])
	}
	
	//vallidate the input
	msg := validateInput(docIn) 
	if msg != "" {
		return shim.Error(msg)	
	}

	// Now, Check if asset already exists
	// return error in case it does
	assetAsBytes, err := APIstub.GetState(docIn.AssetId)
	if assetAsBytes == nil {
	} else {
		fmt.Println("Attempt to create an asset that already exists: " + docIn.AssetId)
		return shim.Error("This asset already exists: " + docIn.AssetId)
	}

	//calculate timestamps 
	tmsp2 := time.Now()
	docIn.TmspStart = tmsp2.String()
	tmsp3 := time.Unix(1<<63-1, 999999999)
	docIn.TmspEnd = tmsp3.String()

	//now marshal the object prepared so far and write to ledger
	DocAsBytes, _ := json.Marshal(docIn)
	APIstub.PutState(docIn.AssetId, DocAsBytes)

	return shim.Success(nil)
}

func validateInput(docIn Asset) string {

	msg := ""
	//validate if the all the fields have been populated properly by front-end/caller
	if docIn.AssetId == ""   {
    	msg += "|Document ID cannot be empty"
	} 

	//if docIn.DocHash == ""  {
    	//msg += "|Document hash cannot be empty"
	//}

	// Add future validations over here
	//	ownerId				string      `json:"ownerId"`	
	//	DocDesc				string      `json:"docDesc"`		


	return msg

} 

func (s *SmartContract) queryAllDocs(APIstub shim.ChaincodeStubInterface) sc.Response {

	//Note: AssetID will be a combination of a prefix + creation timestamp
	startKey := ""
	endKey   := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllDocs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAllDocsShared(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


	if len(args) != 1 {
		return shim.Error("queryAllDocsShared: Incorrect number of arguments. Expecting 1")
	}
	userId := args[0]

	UsrAsBytes, _ := APIstub.GetState(userId)
	if UsrAsBytes == nil {
		
		fmt.Println("###############DocAsBytes Empty")
		fmt.Println("No record found for UserId=" + userId)
		return shim.Error("No record found for userId=" + userId)
	}

	fmt.Println("queryDoc Success")
	return shim.Success(UsrAsBytes)

}


func (s *SmartContract) changeDocOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//Important: Check if hte document has been terminate timestamped already. 
	//If yes, allow no further updates to this document
	resp := s.checkIfTerminated(APIstub,args[0])
	if resp.Status != 0 {  //Note that nil means error
		fmt.Println(resp.Message)
		return shim.Error(resp.Message)
	} 
	fmt.Println("Document is still active")
	Doc := resp.Payload
	Doc.OwnerId = args[1]
	DocAsBytes, err := json.Marshal(Doc)
	if err != nil {
		fmt.Println("changeDocOwner: Marshal error")
		return shim.Error("changeDocOwner: Marshal error :" + err.Error())
	}

	APIstub.PutState(args[0], DocAsBytes)

	return shim.Success(nil)
}

//end date the document passed
func (s *SmartContract) setExpiryOnDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("setExpiryOnDoc: Incorrect number of arguments. Expecting 2")
	}

	//Important: Check if hte document has been terminated timestamped already. 
	//If yes, allow no further updates to this document
	resp := s.checkIfTerminated(APIstub,args[0])
	if resp.Status != 0 {  //Note that nil means error
		fmt.Println(resp.Message)
		return shim.Error(resp.Message)
	} 
	fmt.Println("Document is still active")
	Doc := resp.Payload

	tmspEnd := ""
	//set to current timestamp if the input timestamp is blank
	if args[1] == "" {
		tmsp2 := time.Now()
		tmspEnd = tmsp2.String()
	} else {
		tmspEnd = args[1]
	}

	Doc.TmspEnd = tmspEnd

	DocAsBytes, err := json.Marshal(Doc)
	if err != nil {
		fmt.Println("setExpiryOnDoc Error Marshal")
		return shim.Error("setExpiryOnDoc: Marshal error :" + err.Error())
	}

	APIstub.PutState(Doc.AssetId, DocAsBytes)
	return shim.Success(nil)
}

//add a user to the Access list of this document
func (s *SmartContract) grantAccess(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//Step1: Add userId the asset's entry on ledger
	//Step2: Add this document'assetId to the user list of docs on the ledger
	if len(args) != 3 {
		return shim.Error("grantAccess: Incorrect number of arguments. Expecting 2")
	}	

	if args[0] == "" || args[1] == "" {
		return shim.Error("grantAccess: assetId and userID cannot be empty")	
	}
	assetId := args[0]
	userId := args[1]

	//Important: Check if hte document has been terminated timestamped already. 
	//If yes, allow no further updates to this document
	resp := s.checkIfTerminated(APIstub,args[0])
	if resp.Status != 0 {  //Note that nil means error
		fmt.Println(resp.Message)
		return shim.Error(resp.Message)
	} 
	fmt.Println("Document is still active")
	Doc := resp.Payload

	newAclRow := AssetACL{} 
	newAclRow.UserId = args[1]
	newAclRow.UserDescription = args[2]
	newAclRow.AccessToFields = "*"
	tmsp2 := time.Now()
	newAclRow.AccessGrantTmsp = tmsp2.String()
	
	//add the new acl row onto the asset 
	Doc.AccessList = append(Doc.AccessList,newAclRow)
	
	DocAsBytes, err := json.Marshal(Doc)
	if err != nil {
		return shim.Error("share asset: Marshal error :" + err.Error())
	}
	APIstub.PutState(Doc.AssetId, DocAsBytes)

	var Usr UserACL
	UsrAsBytes, _ := APIstub.GetState(userId)
	if UsrAsBytes == nil {
		Usr = UserACL{}
		Usr.UserId = userId
		Usr.AssetId = append(Usr.AssetId,assetId)		
	} else {
		json.Unmarshal(UsrAsBytes, &Usr)
		Usr.AssetId = append(Usr.AssetId,assetId)		
	}

	//write this record too on the ledger
	UsrAsBytes, err = json.Marshal(Usr)
	if err != nil {
		return shim.Error("User list: Marshal error :" + err.Error())
	}
	APIstub.PutState(Usr.UserId, UsrAsBytes)


	return shim.Success(nil)
}


//List of documents the logged in user has access to 
//type UserACL struct {
//
//	UserId   		string 		`json:"userId"`
//	AssetId			[]string 	`json:"assetId"` 
//}


//=================================================================================================
// Return history by ID
//=================================================================================================
func (s *SmartContract) getDocHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("Inside HISTORY")

	//check number of arguments recvd
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 AssetID")
	}

	//args[0] user1
	//args[1] assetID
	id := args[0]

	// Read history
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error("Not Found")
	}
	defer resultsIterator.Close()

	// Write return buffer
	var buffer bytes.Buffer
	buffer.WriteString("{ \"values\": [")
	for resultsIterator.HasNext() {
		it, _ := resultsIterator.Next()
		if buffer.Len() > 15 {
			buffer.WriteString(",")
		}
		//var doc Doc
		buffer.WriteString("{\"timestamp\":\"")
		buffer.WriteString(time.Unix(it.Timestamp.Seconds, int64(it.Timestamp.Nanos)).Format(time.Stamp))
		buffer.WriteString("\", \"text\":")
		buffer.WriteString(string(it.Value))
		buffer.WriteString("}")
	}
	buffer.WriteString("]}")

	return shim.Success(buffer.Bytes())
}



/*execute rich query on couchdb*/
func (s *SmartContract) executeRichQuery(stub shim.ChaincodeStubInterface, args []string) sc.Response {
    

    if len(args) != 1 { return shim.Error("Incorrect number of arguments. Expecting 1") }
    
    queryString := args[0]
    
    fmt.Println("Inside rQuery, queryString is: %s", queryString)

    
    resultsIterator, err := stub.GetQueryResult(queryString)
    if err != nil { 
    	return shim.Error("rquery#1" + err.Error())
    }

    defer resultsIterator.Close()
    if err != nil {
   		return shim.Error("rquery#2" + err.Error())
    }

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"AssetKey\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"AssetValue\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- rQuery:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())


}


func (s *SmartContract) checkIfTerminated(stub shim.ChaincodeStubInterface, assetId string) internalResp {

	// check if an existing asset is already end timestamped 
	//   If yes, raise an error and return from here 
	resp := internalResp{}
	DocAsBytes, _ := stub.GetState(assetId)
	if DocAsBytes == nil {
		resp.Status = -1
		resp.Message = "checkIfTerminated: Error=> Asset not found assetId=" + assetId
		resp.Payload = Asset{} 	
		fmt.Println(resp.Message)
		return resp
	}
	Doc := Asset{}
	json.Unmarshal(DocAsBytes, &Doc)

	//calculate greatest possible timestamp
	tmsp3 := time.Unix(1<<63-1, 999999999)
	tmspMax := tmsp3.String()

	if Doc.TmspEnd < tmspMax {
		resp.Status = -1
		resp.Message = "checkIfTerminated: Error=> asset not active assetId=" + assetId
		resp.Payload = Asset{} 	
		fmt.Println(resp.Message)
		return resp
	}

	
	resp.Status = 0
	resp.Message = ""
	resp.Payload = Doc 	
	fmt.Println(resp.Message)
	return resp
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

