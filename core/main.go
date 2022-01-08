package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	json "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const testKsAndrejAccount = "0x3eb92807f1f91a8d4d85bc908c7f86dcddb1df57"
const testKsBabaYagaAccount = "0x6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8"
const testKsAndrejFile = "test_andrej--3eb92807f1f91a8d4d85bc908c7f86dcddb1df57"
const testKsBabaYagaFile = "test_babayaga--6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8"
const testKsAccountsPwd = "security123"
const nodeVersion = "0.0.0-alpha 01abcd Test Run"

func main() {
	//NodeRun("127.0.0.1", 9527)
	NodeMining("127.0.0.1", 9527)
}

func NodeRun(ip string, port int) {
	datadir, err := ioutil.TempDir(os.TempDir(), "node_run")
	if err != nil {
		log.Fatal(err)
	}
	err = RemoveDir(datadir)
	if err != nil {
		log.Fatal(err)
	}

	n := New(datadir, ip, uint64(port), NewAccount(DefaultMiner), PeerNode{}, nodeVersion, 2)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx, true, "")
	if err != nil {
		log.Fatal(err)
	}
}

func NodeMining(ip string, port int) {
	dataDir, andrej, babaYaga, err := setupMiningNodeDir(1000000, 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer RemoveDir(dataDir)

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		ip,
		uint64(port),
		false,
		babaYaga,
		true,
		nodeVersion,
	)

	// Construct a new Node instance and configure
	// Andrej as a miner
	n := New(dataDir, nInfo.IP, nInfo.Port, andrej, nInfo, nodeVersion, 2)

	// Allow the mining to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)

	// Schedule a new TX in 3 seconds from now, in a separate thread
	// because the n.Run() few lines below is a blocking call
	go func() {
		time.Sleep(time.Second * miningIntervalSeconds / 3)

		tx := NewBaseTx(andrej, babaYaga, 1, 1, "")
		signedTx, err := SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, GetKeystoreDirPath(dataDir))
		if err != nil {
			log.Println(err)
			return
		}

		_ = n.AddPendingTX(signedTx, nInfo)
	}()

	// Schedule a TX with insufficient funds in 4 seconds validating
	// the AddPendingTX won't add it to the Mempool
	go func() {
		time.Sleep(time.Second*(miningIntervalSeconds/3) + 1)

		tx := NewBaseTx(babaYaga, andrej, 50, 1, "")
		signedTx, err := SignTxWithKeystoreAccount(tx, babaYaga, testKsAccountsPwd, GetKeystoreDirPath(dataDir))
		if err != nil {
			log.Println(err)
			return
		}

		err = n.AddPendingTX(signedTx, nInfo)
		log.Println(err)
		if err == nil {
			log.Printf("TX should not be added to Mempool because BabaYaga doesn't have %d TBB tokens", tx.Value)
			return
		}
	}()

	// Schedule a new TX in 12 seconds from now simulating
	// that it came in - while the first TX is being mined
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds + 2))

		tx := NewBaseTx(andrej, babaYaga, 2, 2, "")
		signedTx, err := SignTxWithKeystoreAccount(tx, andrej, testKsAccountsPwd, GetKeystoreDirPath(dataDir))
		if err != nil {
			log.Println(err)
			return
		}

		err = n.AddPendingTX(signedTx, nInfo)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	go func() {
		// Periodically check if we mined the 2 blocks
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1000 {
					closeNode()
					return
				}
			}
		}
	}()

	// Run the node, mining and everything in a blocking call (hence the go-routines before)
	_ = n.Run(ctx, true, "")

	if n.state.LatestBlock().Header.Number != 1 {
		log.Fatalln("2 pending TX not mined into 2 blocks under 30m")
	}
}

// setupMiningNodeDir creates a default testing node directory with 2 keystore accounts
//
// Remember to remove the dir once test finishes: defer fs.RemoveDir(dataDir)
func setupMiningNodeDir(andrejBalance uint, forkTip1 uint64) (dataDir string, andrej, babaYaga common.Address, err error) {
	babaYaga = NewAccount(testKsBabaYagaAccount)
	andrej = NewAccount(testKsAndrejAccount)

	dataDir, err = ioutil.TempDir(os.TempDir(), "node_mining")
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	genesisBalances := make(map[common.Address]uint)
	genesisBalances[andrej] = andrejBalance
	genesis := Genesis{Balances: genesisBalances, ForkTIP1: forkTip1}
	genesisJson, err := json.Marshal(genesis)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	err = InitDataDirIfNotExists(dataDir, genesisJson)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	err = copyKeystoreFilesIntoTestDataDirPath(dataDir)
	if err != nil {
		return "", common.Address{}, common.Address{}, err
	}

	return dataDir, andrej, babaYaga, nil
}

func copyKeystoreFilesIntoTestDataDirPath(dataDir string) error {
	pwd, _ := os.Getwd()
	pwd = pwd + "/"
	fmt.Println(pwd)
	andrejSrcKs, err := os.Open(pwd + testKsAndrejFile)
	if err != nil {
		return err
	}
	defer andrejSrcKs.Close()

	ksDir := filepath.Join(GetKeystoreDirPath(dataDir))

	err = os.Mkdir(ksDir, 0777)
	if err != nil {
		return err
	}

	andrejDstKs, err := os.Create(filepath.Join(ksDir, testKsAndrejFile))
	if err != nil {
		return err
	}
	defer andrejDstKs.Close()

	_, err = io.Copy(andrejDstKs, andrejSrcKs)
	if err != nil {
		return err
	}

	babayagaSrcKs, err := os.Open(pwd + testKsBabaYagaFile)
	if err != nil {
		return err
	}
	defer babayagaSrcKs.Close()

	babayagaDstKs, err := os.Create(filepath.Join(ksDir, testKsBabaYagaFile))
	if err != nil {
		return err
	}
	defer babayagaDstKs.Close()

	_, err = io.Copy(babayagaDstKs, babayagaSrcKs)
	if err != nil {
		return err
	}

	return nil
}
