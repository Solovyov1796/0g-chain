package main

import (
	dbm "github.com/cometbft/cometbft-db"
)

func main() {
	// TODO: Implement tests

	err := InitDB(dbm.BackendType("goleveldb"), "/data/0g-home/data")
	if err != nil {
		println("init db failed", err.Error())
	} else {
		_, err := GetTx("1728B211DA7134EDA67910B12EC1FF276E45ED9E3390AD3ABB5BA027ED46B952")
		if err != nil {
			println("get tx failed", err.Error())
		}
	}
}
