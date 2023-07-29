package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"jumbochain.org/accounts"
	"jumbochain.org/block"
	"jumbochain.org/p2p"
	"jumbochain.org/temp"
	"jumbochain.org/transaction"
)

func main() {

	startNode := flag.Int("startNode", 0, "starts node")
	sourcePort := flag.Int("isValidator", 0, "source port number")
	newAccount := flag.Int("newAccount", 0, "creats new Account")
	sendTrx := flag.Int("sendTrx", 0, "sends transaction")
	from := flag.String("from", "", "from address for transaction")
	to := flag.String("to", "", "to address for transaction")
	value := flag.Int("value", 0, "value  for transaction")
	initGenesis := flag.Int("initGenesis", 0, "init genesis command")
	balance := flag.String("balance", "", "get balance of address")
	currentBlockNumber := flag.Int("currentBlockNumber", 0, "current block number")
	getBlockHashByNumber := flag.String("getBlockHashByNumber", "", "current block number")
	getBlockByHash := flag.String("getBlockByHash", "", "current block number")
	getTransactionByHash := flag.String("getTransactionByHash", "", "current block number")

	flag.Parse()

	if *getTransactionByHash != "" {
		blockHash := block.GetBlockByHash(*getTransactionByHash)
		fmt.Println("Transaction is :")
		fmt.Println(blockHash)
	}

	if *getBlockByHash != "" {
		blockHash := block.GetBlockByHash(*getBlockByHash)
		fmt.Println("block is :")
		fmt.Println(blockHash)
	}

	if *getBlockHashByNumber != "" {
		blockHash := block.GetBlockHashByNumber(*getBlockHashByNumber)
		fmt.Println("block hash is :")
		fmt.Println(blockHash)
	}

	if *currentBlockNumber != 0 {
		currentBlockNumber := block.GetCurrentBlockNumber()
		fmt.Println("current block number is :")
		fmt.Println(currentBlockNumber)

	}

	if *balance != "" {
		balanceOfAddress := block.GetBalance(*balance)
		fmt.Println("balance of address is :")
		fmt.Println(balanceOfAddress)
	}

	if *sendTrx == 1 {
		transaction.SendTxx(*from, *to, *value)
	}

	if *newAccount == 1 {
		accounts.NewAccount()
	}

	if *initGenesis == 1 {
		block.InitGenesis()
		fmt.Println("Genesis has been loaded successfully")
	}

	if *startNode == 1 {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r := rand.Reader

		h, err := p2p.MakeHost(*sourcePort, r)
		if err != nil {
			log.Fatal(err)
		}

		if *sourcePort != 0 {
			p2p.StartListener(ctx, h, *sourcePort)

			ticker := time.NewTicker(10 * time.Second)
			quit := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker.C:
						transaction.VerifyTransactions()
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()

			ticker2 := time.NewTicker(3 * time.Second)
			quit2 := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker2.C:
						block.CreateBlock()
					case <-quit2:
						ticker.Stop()
						return
					}
				}
			}()

			<-ctx.Done()
		} else {
			ticker := time.NewTicker(10 * time.Second)
			quit := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker.C:

						peerlist := temp.ReadCsv("peerlist.csv")

						for i := 0; i < len(peerlist); i++ {
							fmt.Println(peerlist[i][0])
							target := peerlist[i][0]
							info := p2p.RunSender(h, target)

							p2p.SendStream(h, info)

						}
						if err := os.Truncate("TrxMemPool.csv", 0); err != nil {
							log.Printf("Failed to truncate: %v", err)
						}
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()
			<-ctx.Done()
		}
	}

}
