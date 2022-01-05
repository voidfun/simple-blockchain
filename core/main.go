package main

import (
	"context"
	"github.com/web3coach/the-blockchain-bar/database"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const testKsAndrejAccount = "0x3eb92807f1f91a8d4d85bc908c7f86dcddb1df57"
const testKsBabaYagaAccount = "0x6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8"
const testKsAndrejFile = "test_andrej--3eb92807f1f91a8d4d85bc908c7f86dcddb1df57"
const testKsBabaYagaFile = "test_babayaga--6fdc0d8d15ae6b4ebf45c52fd2aafbcbb19a65c8"
const testKsAccountsPwd = "security123"
const nodeVersion = "0.0.0-alpha 01abcd Test Run"

func main() {
	NodeRun("127.0.0.1", 9527)
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

	n := New(datadir, ip, uint64(port), database.NewAccount(DefaultMiner), PeerNode{}, nodeVersion, 2)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx, true, "")
	if err != nil {
		log.Fatal(err)
	}
}
