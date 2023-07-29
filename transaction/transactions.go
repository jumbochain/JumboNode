package transaction

import (
	"encoding/json"
	"fmt"
	"os"

	ldb "jumbochain.org/database"
	"jumbochain.org/temp"
	"jumbochain.org/types"
)

type Transaction struct {
	To         string `json:"to"`
	From       string `json:"from"`
	Value      int    `json:"value"`
	Hash       string `json:"hash"`
	IsVerified bool   `json:"verified"`
}

func SignGenesisTrx(_to string, _value int) (string, string) {

	_from := "dh00000000000000000GENESIS0000000000000000"
	tx := &Transaction{
		To:    _to,
		From:  _from,
		Value: _value,
	}
	transactionInBytes := []byte(fmt.Sprintf("%v", tx))
	txHash := types.HashFromBytes(transactionInBytes)
	txxHash := txHash.String()

	txAfterHash := &Transaction{
		To:    _to,
		From:  _from,
		Value: _value,
		Hash:  txxHash,
	}

	// var trx Transaction
	transactionBody, err := json.MarshalIndent(txAfterHash, "", " ")

	fmt.Println(string(transactionBody))

	// content, err := json.Marshal(txAfterHash)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(content)
	// fmt.Println(txxHash)

	// var trxBody Transaction

	// json.Unmarshal(transactionBody, &trxBody)
	// fmt.Println(trxBody)

	return txxHash, string(transactionBody)

	//

	// return txHash
}

func SendTxx(_from string, _to string, _value int) (types.Hash, error) {
	tx := &Transaction{
		To:    _to,
		From:  _from,
		Value: _value,
	}

	transactionInBytes := []byte(fmt.Sprintf("%v", tx))
	txHash := types.HashFromBytes(transactionInBytes)
	txxHash := txHash.String()

	// fmt.Println("111111111111", txHash)
	DefaultStatus := false

	// fmt.Println("TO", tx.To)
	// fmt.Println("FROM", tx.From)
	// fmt.Println("Value", tx.Value)
	// fmt.Println("================================================")

	txAfterHash := &Transaction{
		To:         _to,
		From:       _from,
		Value:      _value,
		Hash:       txxHash,
		IsVerified: DefaultStatus,
	}
	// fmt.Println("TO", txAfterHash.To)
	// fmt.Println("FROM", txAfterHash.From)
	// fmt.Println("Value", txAfterHash.Value)
	// fmt.Println("Hash", txAfterHash.Hash)
	// fmt.Println("================================================")

	// fmt.Println("22222222222222222", txAfterHash.Hash)
	notHavingFile := checkFile("TrxMemPool.csv")
	if notHavingFile != nil {
		fmt.Println("Unable to create file", notHavingFile)
	}

	f, err := os.OpenFile("TrxMemPool.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	content, err := json.Marshal(txAfterHash)
	// fmt.Println("content", content)

	n, err := f.Write(content)
	if err != nil {
		fmt.Println(n, err)
	}

	if n, err = f.WriteString("\n"); err != nil {
		fmt.Println(n, err)
	}
	// fmt.Println("TRX:-", txHash)

	return txHash, nil
}

func MakeGenesis(_to string, _from string, _value int) (types.Hash, error) {
	tx := &Transaction{
		To:    _to,
		From:  _from,
		Value: _value,
	}

	transactionInBytes := []byte(fmt.Sprintf("%v", tx))
	txHash := types.HashFromBytes(transactionInBytes)
	// txxHash := txHash.String()

	// DefaultStatus := false

	// txAfterHash := &Transaction{
	// 	To:         _to,
	// 	From:       _from,
	// 	Value:      _value,
	// 	Hash:       txxHash,
	// 	IsVerified: DefaultStatus,
	// }

	notHavingFile := checkFile("txData.csv")
	if notHavingFile != nil {
		fmt.Println("Unable to create file", notHavingFile)
	}

	// f, err := os.OpenFile("txData.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer f.Close()

	// content, err := json.Marshal(txAfterHash)
	// fmt.Println("content", content)

	// n, err := f.Write(content)
	// if err != nil {
	// 	fmt.Println(n, err)
	// }

	// if n, err = f.WriteString("\n"); err != nil {
	// 	fmt.Println(n, err)
	// }
	fmt.Println("TRX:-", txHash)

	return txHash, nil
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

//shiv

var hold string

func VerifyTransactions() (bool, error) {
	var transaction Transaction
	csvData := temp.ReadCsv("TrxMemPoolValidator.csv")
	CSVMap := make(map[int]Transaction, 0)
	VerifiedTransactionsMap := make(map[int]Transaction, 0)
	//fmt.Println(csvData)

	for i := 0; i < len(csvData); i++ {
		fmt.Println(csvData[i][0])
	}

	for i := 0; i < len(csvData); i++ {
		//fmt.Println(csvData[i][0])
		err := json.Unmarshal([]byte(csvData[i][0]), &transaction)
		if err != nil {
			fmt.Println("json.Unmarshal([]byte(value)1, &transaction", err)
		}

		fmt.Println(transaction.From)

		//Open the LevelDB file for reading
		fromValue, err := ldb.GetDataFromLevelDB(transaction.From)

		if err != nil {
			fmt.Println("json.Unmarshal([]byte(value)2, &transaction", err)
		}

		if transaction.From != "" && fromValue >= float64(transaction.Value) {
			// delete(CSVMap, i)
			VerifiedTransactionsMap[i] = transaction
		} else {
			CSVMap[i] = transaction
			// delete(VerifiedTransactionsMap, i)
		}
		//transaction
	}
	temp.TruncateCSVFile()
	for _, value := range CSVMap {

		//fmt.Println("csv file data ", CSVMap[i])
		//fmt.Println("verified transaction data ", VerifiedTransactionsMap[i])
		marshData, err := json.Marshal(value)
		if err != nil {
			fmt.Println("json.Unmarshal([]byte(value)1, &transaction", err)
		}
		temp.AddDataToCSVFile("TrxMemPoolValidator.csv", string(marshData))

	}
	for _, value := range VerifiedTransactionsMap {

		//fmt.Println("csv file data ", CSVMap[i])
		//fmt.Println("verified transaction data ", VerifiedTransactionsMap[i])
		marshData, err := json.Marshal(value)
		if err != nil {
			fmt.Println("json.Unmarshal([]byte(value)1, &transaction", err)
		}
		temp.AddDataToCSVFile("VerifiedTransactions.csv", string(marshData))

	}

	return false, nil

}
