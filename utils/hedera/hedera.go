package hedera

import (
	"encoding/hex"
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/warlockz/ase-service/configs"
	"github.com/warlockz/ase-service/types"
)

type SetDataParameter struct {
	file_data string
	user      string
	file_name string
}

var client *hedera.Client

func ConnectHedera() *hedera.Client {
	var err error

	client = hedera.ClientForTestnet()
	if err != nil {
		println(err.Error(), ": error creating client")
	}
	configOperatorID := configs.EnvHederaOperatorID()
	configOperatorKey := configs.EnvHederaOperatorPrivateKey()

	//client.SetOperator(configOperatorID, configOperatorKey)

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
		}
		// fmt.Println(operatorAccountID.Realm)
		// fmt.Println(operatorKey.String())
		client.SetOperator(operatorAccountID, operatorKey)

	}
	fmt.Println("Connected to Hedera")
	return client

}

var hederaClient *hedera.Client = ConnectHedera()
var contractId string = configs.EnvContractID()

func executeTransaction(functionName string, parametersArray []types.ParametersType) *hedera.ContractFunctionResult {
	// functionVariables :=
	FinalFunctionVariables := addFunctionParameters(parametersArray, &*hedera.NewContractFunctionParameters(), 0)
	// fmt.Println(FinalFunctionVariables)
	// functionVariables :=  hedera.NewContractFunctionParameters().AddString("Test doc2").AddString("user1").AddString("test_doc")
	newContractId, err := hedera.ContractIDFromString(contractId)
	//newfunctionVariables := *hedera.ContractFunctionParameters
	transaction := hedera.NewContractExecuteTransaction().
		SetContractID(newContractId).
		SetGas(10000000).
		SetFunction(functionName, FinalFunctionVariables)
	//Sign with the client operator privateclient key to pay for the transaction and submit the query to a Hedera network
	txResponse, err := transaction.Execute(hederaClient)

	if err != nil {
		panic(err)
	}

	// Get Transaction Record
	txRecord, err := txResponse.GetRecord(hederaClient)
	if err != nil {
		panic(err)
	}
	// Get Contract Result
	contractResult, err := txRecord.GetContractExecuteResult()
	if err != nil {
		panic(err)
	}
	fmt.Println("Result Gas used")
	fmt.Println(contractResult.GasUsed)

	//Request the receipt of the transaction
	txReceipt, err := txResponse.GetReceipt(hederaClient)

	if err != nil {
		panic(err)
	}
	//Get the transaction consensus status
	transactionStatus := txReceipt.Status

	fmt.Printf("The transaction consensus status %v\n", transactionStatus)

	return &contractResult

}

func addFunctionParameters(parameterArray []types.ParametersType, functionParameter *hedera.ContractFunctionParameters, index int) *hedera.ContractFunctionParameters {
	//newFunctionParameter *hedera.ContractFunctionParameters := nil
	if index >= len(parameterArray) {
		return functionParameter
	}

	var err error
	switch parameterArray[index].Datatype {
	case "string":
		functionParameter = functionParameter.AddString(parameterArray[index].Value)
		break
	case "int32":
		functionParameter = functionParameter.AddInt32(200)
	case "uint256":
		functionParameter = functionParameter.AddUint256([]byte(parameterArray[index].Value))
	case "stringArray":
		functionParameter = functionParameter.AddStringArray(parameterArray[index].Array)
	case "address":
		operatorAccountID, err := hedera.AccountIDFromString(parameterArray[index].Value)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
		}
		addressString := operatorAccountID.ToSolidityAddress()
		functionParameter, err = functionParameter.AddAddress(addressString)
	case "bytes32":
		functionParameter = functionParameter.AddBytes32(parameterArray[index].Bytes32)
	}

	if err != nil {
		panic(err)
	}

	return addFunctionParameters(parameterArray, functionParameter, index+1)
}

func SetData(parametersArray []types.ParametersType) {
	// Execute Transaction
	executeTransaction("addData", parametersArray)

	// // Get Output Parameters
	// hash := contractResult.GetBytes32(0)
	// hashString := string(hash[:])

	// timestamp := contractResult.GetUint256(1)
	// timestampString := string(timestamp[:])

	// jsonFile, err := os.Open("./notarization.json")

	//  // if we os.Open returns an error then handle it
	// if err != nil {
	//     fmt.Println(err)
	// }

	// defer jsonFile.Close()

	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// var result map[string]interface{}

	// json.Unmarshal([]byte(byteValue), &result)

	// fmt.Println(result)
	//data := *hedera.ContractFunctionResultFromBytes(hash)
	//fmt.Println(data)
	// setDataOutput := types.SetDataOutput{Hash: hashString, Timestamp: timestampString}

	// return setDataOutput
}

func ViewData(parametersArray []types.ParametersType) {
	newContractId, err := hedera.ContractIDFromString(contractId)
	FinalFunctionVariables := addFunctionParameters(parametersArray, &*hedera.NewContractFunctionParameters(), 0)
	// Call the contract to receive the greeting
	callResult, err := hedera.NewContractCallQuery().
		SetContractID(newContractId).
		// The amount of gas to use for the call
		// All of the gas offered will be used and charged a corresponding fee
		SetGas(100000).
		// This query requires payment, depends on gas used
		SetQueryPayment(hedera.NewHbar(1)).
		// Specified which function to call, and the parameters to pass to the function
		SetFunction("viewData", FinalFunctionVariables).
		// This requires payment
		SetMaxQueryPayment(hedera.NewHbar(5)).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing contract call query")
		return
	}

	fmt.Printf("Message: %v\n", callResult.GetString(0))

}

func GenerateTrapdoor(parametersArray []types.ParametersType) ([]byte, string, [32]byte) {

	contractResult := executeTransaction("search_request", parametersArray)
	var data hedera.ContractFunctionResult = *contractResult

	fileHash := data.GetString(0)
	//fmt.Println("fileHash")

	//	fmt.Println(fileHash)

	trapdoorHashData := data.GetBytes32(1)
	var trapdoorHash [32]byte
	copy(trapdoorHash[:], trapdoorHashData)
	//fmt.Println("trapdoorHashData")

	//fmt.Println(hex.EncodeToString(trapdoorHashData))
	// hashString := string(hash[:])

	keyCipherText := data.GetString(2)
	//fmt.Println("keyCipherText")

	// fmt.Println(keyCipherText)
	r, err := hex.DecodeString(keyCipherText)

	if err != nil {
		panic(err)

	}
	// fmt.Printf("Blockchain data  to \n[%x]\n", r)

	// timestamp := contractResult.GetUint256(1)
	// timestampString := string(timestamp[:])

	return r, fileHash, trapdoorHash

}

func VerifyTrapdoor(parametersArray []types.ParametersType) bool {
	contractResult := executeTransaction("verify_trapdoor", parametersArray)
	var data hedera.ContractFunctionResult = *contractResult
	verified := data.GetBool(0)
	return verified
}

func VerifyResult(parametersArray []types.ParametersType) bool {
	contractResult := executeTransaction("verify_result", parametersArray)
	var data hedera.ContractFunctionResult = *contractResult
	verified := data.GetBool(0)
	return verified
}
