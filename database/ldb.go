package database

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func AddDataToDatabase(databasename string, key []byte, value []byte) {
	db, err := leveldb.OpenFile(databasename, nil)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	db.Put(key, value, nil)
}

func FeatchFromDatabase(databasename string, key []byte) []byte {
	db, _ := leveldb.OpenFile(databasename, nil)
	defer db.Close()
	data, _ := db.Get(key, nil)
	// fmt.Println("%v \n", data)

	return data

}

func GetDataFromLevelDB(from string) (float64, error) {

	db, err := leveldb.OpenFile("database", nil)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	//time.Sleep(time.Second)

	fmt.Println(from)

	fmt.Println("-----------")

	fromBalanceBytes, err := db.Get([]byte(from), nil)

	fmt.Println(fromBalanceBytes)

	fmt.Println("-----------")

	//var transaction models.Transaction
	var amount float64
	err = json.Unmarshal(fromBalanceBytes, &amount)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(value)3, &transaction", err)
	}
	fmt.Println(amount)

	return amount, nil

}
