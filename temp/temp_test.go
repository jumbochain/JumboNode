package temp

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

var j = []byte(`{"foo":1,"bar":2,"baz":[3,4]}`)

// func TestYu(t *testing.T) {
// 	// a map container to decode the JSON structure into
// 	c := make(map[string]json.RawMessage)

// 	fmt.Println(c)

// 	// unmarschal JSON
// 	e := json.Unmarshal(j, &c)

// 	fmt.Println(c)
// 	fmt.Println(e)

// 	// panic on error
// 	if e != nil {
// 		panic(e)
// 	}

// 	// a string slice to hold the keys
// 	k := make([]string, len(c))

// 	// iteration counter
// 	i := 0

// 	// copy c's keys into k
// 	for s, _ := range c {
// 		k[i] = s
// 		// fmt.Println(K[i])
// 		i++
// 	}

// 	// output result to STDOUT
// 	fmt.Printf("%#v\n", k)
// 	fmt.Println("")

// }

func TestStart(t *testing.T) {
	tx := Transaction{
		From:  "add1",
		To:    "add2",
		Value: 100,
	}

	by := EncodeToBytes(tx)

	fmt.Println(by)

	p := DecodeToPerson(by)

	fmt.Println(p)

	// reqBodyBytes := new(bytes.Buffer)
	// json.NewEncoder(reqBodyBytes).Encode(tx)

	// by := reqBodyBytes.Bytes() // this is the []byte

	// fmt.Println(by)

}

func TestRead(t *testing.T) {
	trxs_inMemPool := ReadCsv("TrxMemPool.csv")
	fmt.Println(trxs_inMemPool)

}
func TestReadCsv(t *testing.T) {

	// "TrxMemPool.csv"
	f, err := os.Open("TrxMemPool.csv")

	if err != nil {
		fmt.Println("error is opening file")
	}

	defer f.Close()

	r := csv.NewReader(f)

	// r.
	// skip first line

	records, err := r.ReadAll()

	if err != nil {
		fmt.Println("error in reading file")
	}

	fmt.Println(records)
}
