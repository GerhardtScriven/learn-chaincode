/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// Executes when each peer deploys its instance of the chaincode. It starts the chaincode and registers it with the peer.
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//==============================================================================================================================

// Init resets all the things
// GJS: THIS is where we will create the user ledger????????
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
        	return nil, errors.New("Incorrect number of arguments. Expecting 1")
    	}
	
	//GJS  "hello_world" is the KEY below.  What was put in args_0 is the corresponding value
    	err := stub.PutState("hello_world", []byte(args[0]))
    	if err != nil {
        	return nil, err
    	}

    	return nil, nil
}

//==============================================================================================================================

// Invoke is our entry point to invoke a chaincode function
// Invoke functions are captured as transactions, which get grouped into blocks for writing to the ledger. 
// Updating the ledger is achieved by invoking your chaincode. 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("invoke is running " + function)

    	// Handle different functions
    	if function == "init" {
        	return t.Init(stub, "init", args)
    	} else if function == "write" {
        	return t.write(stub, args)
    	}
    	fmt.Println("invoke did not find func: " + function)

    	return nil, errors.New("Received unknown function invocation")

}

//WRITE
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string
    	var err error
    	fmt.Println("running write()")

    	if len(args) != 2 {
        	return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
    	}

    	name = args[0]                            
	//rename for fun
    	value = args[1]
	//GJS use args 0 and 1 to store as key value pairs.   For the MyCREDS, we will have many more args
    	err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
    	if err != nil {
        	return nil, err
    	}
    	return nil, nil
}



//==============================================================================================================================

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("query is running " + function)

    	// Handle different functions
    	if function == "read" {                            
		//read a variable
        	return t.read(stub, args)
    	}
    	fmt.Println("query did not find func: " + function)

    	return nil, errors.New("Received unknown function query")

}

//READ
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    	var name, jsonResp string
    	var err error

    	if len(args) != 1 {
        	return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    	}

    	name = args[0]
    	valAsbytes, err := stub.GetState(name)
    	if err != nil {
        	jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        	return nil, errors.New(jsonResp)
    	}

    	return valAsbytes, nil
}

/*
	UPON DEPLOY (this is the response)
    	"message": "ca9e680f9e7a4c7565e2677b481f9205c9b50504692ee1e4ef9b57491cdd2c0bea2cc9291e32a714dc11f11c085fac6b89b58cc9911cc29314ae9cfe61893ed5"
  	I am using Admin
	This is the code I put as part of deploy (REQUEST):
	
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID": {
      "path": "https://github.com/gerhardtscriven/learn-chaincode/finished"
    },
    "ctorMsg": {
      "function": "init",
      "args": [
        "hi there, I think,"
      ]
    },
    "secureContext": "admin"
  },
  "id": 1
}

AND THEN, WHEN I QUERY, THIS IS MY REQUEST:
{
     "jsonrpc": "2.0",
     "method": "query",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "ca9e680f9e7a4c7565e2677b481f9205c9b50504692ee1e4ef9b57491cdd2c0bea2cc9291e32a714dc11f11c085fac6b89b58cc9911cc29314ae9cfe61893ed5"
         },
         "ctorMsg": {
             "function": "read",
             "args": [
                 "hello_world"
             ]
         },
         "secureContext": "admin"
     },
     "id": 2
 }

AND THIS IS THE RESPONSE
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "hi there, I think,"
  },
  "id": 2
}



AND THEN, WHEN I DO AN INVOKE
{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "ca9e680f9e7a4c7565e2677b481f9205c9b50504692ee1e4ef9b57491cdd2c0bea2cc9291e32a714dc11f11c085fac6b89b58cc9911cc29314ae9cfe61893ed5"
         },
         "ctorMsg": {
             "function": "write",
             "args": [
                 "hello_world",
                 "dude, go away"
             ]
         },
         "secureContext": "admin"
     },
     "id": 3
 }

I GET THIS RESPONSE
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "2d668e5c-114b-4c91-99ae-d641e4c5a5e8"
  },
  "id": 3
}


ID
Secret
admin
2132a5b499
WebAppAdmin
d80ff2824f
user_type1_0
9a60b429c9
user_type1_1
6875d561b7
user_type1_2
9072d21213
user_type1_3
3a697ed3f8
user_type1_4
fd2157092e
user_type2_0
d1fb72cff4
user_type2_1
0223f65c8e
user_type2_2
1e673f176b
user_type2_3
addf17ca10
user_type2_4
a2a6591e05
user_type4_0
1528e1a762
user_type4_1
fe28404d9d
user_type4_2
adc0953bcc
user_type4_3
a650a0f4e6
user_type4_4
83729c20f1
user_type8_0
37bf9d99c4
user_type8_1
fde8207719
user_type8_2
3cf20e1ccf
user_type8_3
6e2c827194
user_type8_4
f2c7a7b8ea
*/