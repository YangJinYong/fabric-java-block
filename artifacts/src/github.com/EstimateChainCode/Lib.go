package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

func getListResult(resultsIterator shim.StateQueryIteratorInterface) ([]byte,error){

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("{\"records\":[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//buffer.WriteString("{")
		//buffer.WriteString("\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", " )
		//buffer.WriteString("\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

//分页返回
func getPaginationQueryResults(resultsIterator shim.StateQueryIteratorInterface, responseMetadata *pb.QueryResponseMetadata) ([]byte,error) {

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("{\"records\":[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//buffer.WriteString("{")
		//buffer.WriteString("\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")
		//
		//buffer.WriteString(", " )
		//buffer.WriteString("\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString(",\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}")


	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func getHistoryListResult(resultsIterator shim.HistoryQueryIteratorInterface) ([]byte,error){

	type EstimateHistory struct {
		TxId    string   `json:"txId"`
		IsDelete string `json:"isDelete"`
		Value   EstimateStruct   `json:"value"`
		Timestamp string `json:"timestamp"`
	}
	var history []EstimateHistory
	var recordStruct EstimateStruct

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	//var buffer bytes.Buffer
	//buffer.WriteString("[")

	//bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		//queryResponse, err := resultsIterator.Next()
		//if err != nil {
		//	return nil, err
		//}
		//// Add a comma before array members, suppress it for the first array member
		//if bArrayMemberAlreadyWritten == true {
		//	buffer.WriteString(",")
		//}
		//item,_:= json.Marshal( queryResponse)
		//buffer.Write(item)
		//bArrayMemberAlreadyWritten = true

		historyData, err := resultsIterator.Next()
		if err != nil {
			return nil,err
		}

		var tx EstimateHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		tx.IsDelete = strconv.FormatBool(historyData.IsDelete)
		tx.Timestamp = time.Unix(historyData.Timestamp.Seconds, int64(historyData.Timestamp.Nanos)).String()
		if historyData.Value == nil {                  //marble has been deleted
			var emptyMarble EstimateStruct
			tx.Value = emptyMarble                 //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &recordStruct) //un stringify it aka JSON.parse()
			tx.Value = recordStruct                      //copy marble over

		}
		history = append(history, tx)              //add this tx to the list
	}
	//buffer.WriteString("]")
	//fmt.Printf("queryResult:\n%s\n", buffer.String())

	historyAsBytes, _ := json.Marshal(history)
	fmt.Printf("queryHistoryResult:\n%s\n", string(historyAsBytes))
	return historyAsBytes,nil
}

