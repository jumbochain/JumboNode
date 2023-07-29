package block

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	db "jumbochain.org/database"
	"jumbochain.org/temp"
	tr "jumbochain.org/transaction"
	"jumbochain.org/types"
)

// Define the structure of a block
type Block struct {
	BlockNumber  int      `json:"blocknumber"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
	PrevHash     string   `json:"prevHash"`
	Hash         string   `json:"hash"`
}

func GenesisBlockHash(transactions []string) (string, string) {

	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"

	block := &Block{
		BlockNumber:  0,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	BlockInBytes := []byte(fmt.Sprintf("%v", block))
	blockHash := types.HashFromBytes(BlockInBytes)
	blockHash_string := blockHash.String()

	blockAfterHash := &Block{
		BlockNumber:  0,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         blockHash_string,
	}

	fmt.Println(blockAfterHash)

	blockBody, err := json.MarshalIndent(blockAfterHash, "", " ")

	fmt.Println(string(blockBody))

	// content, err := json.Marshal(txAfterHash)
	if err != nil {
		fmt.Println(err)
	}

	return blockHash_string, string(blockBody)

	// content, err := json.Marshal(blockAfterHash)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(content)
	// fmt.Println(blockHash_string)

	// return blockHash_string, content

	// var trxBody Transaction

	// json.Unmarshal(content, &trxBody)
	// fmt.Println(trxBody)
	// return txHash
}

func GetBalance(address string) int {

	balanceInBytes := db.FeatchFromDatabase("database", []byte(address))

	balance, _ := strconv.Atoi(string(balanceInBytes))

	return balance
}

func GetCurrentBlockNumber() string {
	currentBlockNumber := db.FeatchFromDatabase("database", []byte("currentBlockNumber"))

	return string(currentBlockNumber)
}
func GetBlockHashByNumber(blockNumber string) string {
	blockHash := db.FeatchFromDatabase("database", []byte(blockNumber))

	return string(blockHash)
}
func GetBlockByHash(blockNumber string) string {
	blockHash := db.FeatchFromDatabase("database", []byte(blockNumber))

	return string(blockHash)
}

func GetTransactionByHash(hashq string) string {
	transaction := db.FeatchFromDatabase("database", []byte(hashq))

	return string(transaction)
}

func CreateBlock() {
	transactions := temp.ReadCsv("VerifiedTransactions.csv")
	fmt.Println("---=====")
	fmt.Println(transactions)

	var transaction_hashes []string

	for i := 0; i < len(transactions); i++ {
		fmt.Println(transactions[i][0])
		transaction := transactions[i][0]

		fmt.Println("0000000")
		fmt.Println(transaction)

		transactionBody, err := json.MarshalIndent(transaction, "", " ")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(transactionBody))

		var tx tr.Transaction
		json.Unmarshal([]byte(transaction), &tx)

		fmt.Println(tx.Hash)

		//add Transaction
		db.AddDataToDatabase("database", []byte(tx.Hash), []byte(transactionBody))

		from := tx.From
		to := tx.To
		value := tx.Value

		from_balance_string := string(db.FeatchFromDatabase("database", []byte(from)))
		to_balance_string := string(db.FeatchFromDatabase("database", []byte(to)))

		// strconv.ParseInt(from_balance_string)

		from_balance, err := strconv.ParseInt(from_balance_string, 10, 64)

		to_balance, err := strconv.ParseInt(to_balance_string, 10, 64)

		from_balance_i := int(from_balance)
		to_balance_i := int(to_balance)

		from_new_balance := from_balance_i - value
		to_new_balance := to_balance_i + value

		db.AddDataToDatabase("database", []byte(from), []byte(strconv.Itoa(from_new_balance)))
		db.AddDataToDatabase("database", []byte(to), []byte(strconv.Itoa(to_new_balance)))

		//update blockstate
		// db.AddDataToDatabase("database", []byte(blockNumber), []byte(hash))

		transaction_hashes = append(transaction_hashes, tx.Hash)

		// content, err := json.Marshal(txAfterHash)

	}

	fmt.Println(transaction_hashes)

	prevHash := string(db.FeatchFromDatabase("database", []byte("lh")))
	lastBlock_string := string(db.FeatchFromDatabase("database", []byte("currentBlockNumber")))

	lastBlock, err := strconv.ParseInt(lastBlock_string, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	createBlock(prevHash, transaction_hashes, int(lastBlock+1))

	if err := os.Truncate("VerifiedTransactions.csv", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
}

func createBlock(prevHash string, transactions []string, blockNumber int) {
	block := &Block{
		BlockNumber:  blockNumber,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	BlockInBytes := []byte(fmt.Sprintf("%v", block))
	blockHash := types.HashFromBytes(BlockInBytes)
	blockHash_string := blockHash.String()

	blockAfterHash := &Block{
		BlockNumber:  blockNumber,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         blockHash_string,
	}

	fmt.Println(blockAfterHash)

	blockBody, err := json.MarshalIndent(blockAfterHash, "", " ")

	fmt.Println(string(blockBody))

	if err != nil {
		fmt.Println(err)
	}

	//update blockchain
	db.AddDataToDatabase("database", []byte(blockHash_string), []byte(blockBody))

	//update block number
	db.AddDataToDatabase("database", []byte(strconv.Itoa(blockNumber)), []byte(blockHash_string))

	lastHash := "lh"
	//update lasthash
	db.AddDataToDatabase("database", []byte(lastHash), []byte(blockHash_string))

	currentBlockNumber := "currentBlockNumber"
	//update current block number
	db.AddDataToDatabase("database", []byte(currentBlockNumber), []byte(strconv.Itoa(blockNumber)))

}

// Create a new block with a previous hash
// func Createblock(index int, transactions []string, prevHash string) *Block {
// 	block := &Block{
// 		Index:        index,
// 		Timestamp:    time.Now().String(),
// 		Transactions: transactions,
// 		PrevHash:     prevHash,
// 	}
// 	block.Hash = calculateHash(block)
// 	return block
// }

// func AddDataToDatabase(k []byte, v []byte) {
// 	db, err := leveldb.OpenFile("database", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	db.Put(k, v, nil)
// }

// func FeatchFromDatabase(abc []byte) []byte {
// 	db, _ := leveldb.OpenFile("database", nil)
// 	data, _ := db.Get(abc, nil)
// 	fmt.Println("%v \n", data)

// }

// func MakeGenesisBlock() {
// 	utilsP.ReadFromTxData()

// 	// Create the genesis block
// 	genesisblock := &block{
// 		Index:     0,
// 		Timestamp: time.Now().String(),
// 		Data:      "Genesis block",
// 		PrevHash:  "",
// 	}
// 	genesisblock.Hash = calculateHash(genesisblock)
// 	blockchain = append(blockchain, genesisblock)
// 	blockchainBytes, err := json.MarshalIndent(blockchain, "", "    ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(blockchainBytes))
// 	blockNumber := []byte("0")
// 	AddDataToDatabase(blockNumber, blockchainBytes)
// }

// func AddBlockToChain() {
// 	genesisBlockNumber := []byte("0")
// 	genesisBlock := FeatchFromDatabase(genesisBlockNumber)
// 	// if genesisBlock != 0 {
// 	// 	fmt.Println("genesis is not present")
// 	// }
// 	fmt.Println(genesisBlock)

// 	// Add some transactions to the blockchain
// 	data := utils.ReadFromTxData()
// 	addblock(data)
// 	//addblock("Transaction 2")

// 	// Print the blockchain
// 	blockchainBytes, err := json.MarshalIndent(blockchain, "", "    ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(blockchainBytes))
// }

// Calculate the hash of a block
// func calculateHash(block *Block) string {
// 	record := strconv.Itoa(block.Index) + block.Timestamp + block.Transactions[0] + block.PrevHash
// 	h := sha256.New()
// 	h.Write([]byte(record))
// 	hash := h.Sum(nil)
// 	return hex.EncodeToString(hash)
// }

// // // Define the blockchain
// var blockchain []*Block

// // Add a new block to the blockchain
//func addblock(data string) {
// prevblock := blockchain[len(blockchain)-1]
// newblock := createblock(prevblock.Index+1, data, prevblock.Hash)
//blockchain = append(blockchain, newblock)

// func addblock(data ...string) {
// 	prevblock := blockchain[len(blockchain)-1]
// 	newData := ""
// 	for _, d := range data {
// 		newData += d
// 	}
// 	newblock := Createblock(prevblock.Index+1, newData, prevblock.Hash)
// 	blockchain = append(blockchain, newblock)
// }

// // func main() {
// // 	// Create the genesis block
// // 	genesisblock := &block{
// // 		Index:     0,
// // 		Timestamp: time.Now().String(),
// // 		Data:      "Genesis block",
// // 		PrevHash:  "",
// // 	}
// // 	genesisblock.Hash = calculateHash(genesisblock)
// // 	blockchain = append(blockchain, genesisblock)

// // 	// Add some transactions to the blockchain
// // 	addblock("Transaction1", "Transaction2", "Transaction3", "myTransaction")
// // 	//addblock("Transaction 2")

// // 	// Print the blockchain
// // 	blockchainBytes, err := json.MarshalIndent(blockchain, "", "    ")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	fmt.Println(string(blockchainBytes))
// // }
