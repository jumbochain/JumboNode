package block

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	db "jumbochain.org/database"
	"jumbochain.org/transaction"
)

type AllocJosn struct {
	Allocs []Alloc `json:"alloc"`
}

type Alloc struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func InitGenesis() {

	var allocJosn AllocJosn = readGenesis("genesis.json")
	var transaction_hashes []string
	for i := 0; i < len(allocJosn.Allocs); i++ {
		// fmt.Println(allocJosn.Allocs[i].Address)
		// fmt.Println(allocJosn.Allocs[i].Balance)
		to := allocJosn.Allocs[i].Address
		value := allocJosn.Allocs[i].Balance
		hash, trxBody := transaction.SignGenesisTrx(to, value)

		//add Transaction
		db.AddDataToDatabase("database", []byte(hash), []byte(trxBody))

		//update blockstate
		db.AddDataToDatabase("database", []byte(to), []byte(strconv.Itoa(value)))

		// content := db.FeatchFromDatabase("database", []byte(hash))
		// fmt.Println("------------")
		// fmt.Println(string(content))

		transaction_hashes = append(transaction_hashes, hash)

	}

	fmt.Println(transaction_hashes)

	hash, blockBody := GenesisBlockHash(transaction_hashes)

	//update blockchain
	db.AddDataToDatabase("database", []byte(hash), []byte(blockBody))

	blockNumber := "0"
	//update block number
	db.AddDataToDatabase("database", []byte(blockNumber), []byte(hash))

	lastHash := "lh"
	//update lasthash
	db.AddDataToDatabase("database", []byte(lastHash), []byte(hash))

	currentBlockNumber := "currentBlockNumber"
	//update current block number
	db.AddDataToDatabase("database", []byte(currentBlockNumber), []byte(blockNumber))

	// content := db.FeatchFromDatabase("database", []byte(hash))
	// fmt.Println("------------")
	// fmt.Println(content)

	// var trxBody Block

	// json.Unmarshal(content, &trxBody)
	// fmt.Println(trxBody)

	//create Block

	// block := Createblock(0, hashes, prevHash)

	// fmt.Println(block)

}

func readGenesis(filename string) AllocJosn {
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println("error is opening file")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var allocJosn AllocJosn

	json.Unmarshal(byteValue, &allocJosn)

	return allocJosn

}
