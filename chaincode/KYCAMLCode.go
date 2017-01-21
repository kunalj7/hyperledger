/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
                "errors"
                "fmt"
                //"strconv"
                "encoding/json"
                "github.com/hyperledger/fabric/core/chaincode/shim"
                //"github.com/golang/protobuf/ptypes/timestamp"
)

// KYC  AML Chaincode implementation
type KYCAMLcode struct {
}


var kycAMLIndexTxStr = "_kycAMLIndexTxStr"

type KYCDetails struct{
                USER_NAME string `json:"USER_NAME"`
                USER_ID string `json:"USER_ID"`
                NAME_OF_BANK string `json:"NAME_OF_BANK"`
                KYC_DATE string `json:"KYC_DATE"`
                KYC_VALID_TILL_DATE string `json:"KYC_VALID_TILL_DATE"`
                KYC_DOCUMENT string `json:"KYC_DOCUMENT"`
}


func (t *KYCAMLcode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

                var err error
                // Initialize the chaincode

                fmt.Printf("KYC setup is completed\n")

                var emptyKYCDtls []KYCDetails
                jsonAsBytes, _ := json.Marshal(emptyKYCDtls)
                err = stub.PutState(kycAMLIndexTxStr, jsonAsBytes)
                if err != nil {
                                return nil, err
                }
                return nil, nil
}

// Add region data for the policy
func (t *KYCAMLcode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
                if function == kycAMLIndexTxStr {
                                return t.AddKYCAMLData(stub, args)
                }
                return nil, nil
}

func (t *KYCAMLcode)  AddKYCAMLData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

                var KYCDetailsDataObj KYCDetails
                var KYCDetailsList []KYCDetails
                var err error

                if len(args) != 6 {
                                return nil, errors.New("Incorrect number of arguments. Need 6 arguments")
                }

                // Initialize the chaincode
                KYCDetailsDataObj.USER_NAME = args[0]
                KYCDetailsDataObj.USER_ID = args[1]
                KYCDetailsDataObj.NAME_OF_BANK = args[2]
                KYCDetailsDataObj.KYC_DATE = args[3]
                KYCDetailsDataObj.KYC_VALID_TILL_DATE = args[4]
                KYCDetailsDataObj.KYC_DOCUMENT = args[5]

                fmt.Printf("Input from user:%s\n", KYCDetailsDataObj)

                kycAMLDtlsAsBytes, err := stub.GetState(kycAMLIndexTxStr)
                if err != nil {
                                return nil, errors.New("Failed to retrieve KYC AML Dtls")
                }
                json.Unmarshal(kycAMLDtlsAsBytes, &KYCDetailsList)

                KYCDetailsList = append(KYCDetailsList, KYCDetailsDataObj)
                jsonAsBytes, _ := json.Marshal(KYCDetailsList)

                err = stub.PutState(kycAMLIndexTxStr, jsonAsBytes)
                if err != nil {
                                return nil, err
                }
                return nil, nil
}

// Query callback representing the query of a chaincode
func (t *KYCAMLcode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

                var UserId string // Entities
                var err error
                var resAsBytes []byte


                if len(args) != 1 {
                                return nil, errors.New("Incorrect number of arguments. Expecting UserId to query")
                }

                UserId = args[0]

                resAsBytes, err = t.GetKYCAMLDetails(stub, UserId)

                fmt.Printf("Query Response:%s\n", resAsBytes)

                if err != nil {
                                return nil, err
                }

                return resAsBytes, nil
}

func (t *KYCAMLcode)  GetKYCAMLDetails(stub shim.ChaincodeStubInterface, UserId string) ([]byte, error) {

                //var requiredObj RegionData
                var objFound bool
                KYCDtlsAsBytes, err := stub.GetState(kycAMLIndexTxStr)
                if err != nil {
                                return nil, errors.New("Failed to get KYC AML Details for supplied UserId")
                }
                var KYCDetailsObjects []KYCDetails
                var KYCDetailsObjects1 []KYCDetails
                json.Unmarshal(KYCDtlsAsBytes, &KYCDetailsObjects)
                length := len(KYCDetailsObjects)
                fmt.Printf("Output from chaincode: %s\n", KYCDtlsAsBytes)

                if UserId == "" {
                                res, err := json.Marshal(KYCDetailsObjects)
                                if err != nil {
                                return nil, errors.New("Failed to Marshal the required Obj")
                                }
                                return res, nil
                }

                objFound = false
                // iterate
                for i := 0; i < length; i++ {
                                obj := KYCDetailsObjects[i]
                                if UserId == obj.USER_ID {
                                                KYCDetailsObjects1 = append(KYCDetailsObjects1,obj)
                                                //requiredObj = obj
                                                objFound = true
                                }
                }

                if objFound {
                                res, err := json.Marshal(KYCDetailsObjects1)
                                if err != nil {
                                return nil, errors.New("Failed to Marshal the required Obj")
                                }
                                return res, nil
                } else {
                                res, err := json.Marshal("No Data found")
                                if err != nil {
                                return nil, errors.New("Failed to Marshal the required Obj")
                                }
                                return res, nil
                }
}


func main() {
                err := shim.Start(new(KYCAMLcode))
                if err != nil {
                                fmt.Printf("Error starting Simple chaincode: %s", err)
                }
}
