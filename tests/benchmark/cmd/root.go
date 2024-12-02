package cmd

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"
	"os"
	"time"

	"cosmossdk.io/errors"
	"github.com/0glabs/0g-chain/tests/benchmark/account"
	"github.com/0glabs/0g-chain/tests/benchmark/generator"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/spf13/cobra"

	"github.com/robfig/cron/v3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "0g-chain workload maker",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

const (
	basePrefix        = 6370
	evmFaucetMnemonic = "hundred flash cattle inquiry gorilla quick enact lazy galaxy apple bitter liberty print sun hurdle oak town cash because round chalk marriage response success"
	// "brief similar type month pause march ribbon rocket vanish space walnut father filter similar wet exact biology ugly empower cousin erode lend science crisp"
	//evmFaucetMnemonic = "news tornado sponsor drastic dolphin awful plastic select true lizard width idle ability pigeon runway lift oppose isolate maple aspect safe jungle author hole"
	//"crash sort dwarf disease change advice attract clump avoid mobile clump right junior axis book fresh mask tube front require until face effort vault"
	// "hundred flash cattle inquiry gorilla quick enact lazy galaxy apple bitter liberty print sun hurdle oak town cash because round chalk marriage response success"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "To run the benchmark",
	Long:  "To run the benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpc-url")
		txSendCountPerDay, _ := cmd.Flags().GetInt("tx-cnt")
		userIncrementPerDay, _ := cmd.Flags().GetInt("usr-cnt")

		speed := txSendCountPerDay / int(24*time.Hour.Seconds())

		println("speed = ", speed)
		if speed == 0 {
			log.Fatal("speed should be greater than 0")
		}

		StartDeamon(func() {
			println("start at:", time.Now().UnixNano())
			doing(rpcUrl, userIncrementPerDay, speed)
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("rpc-url", "", "http://127.0.0.1:8545", "RPC url of the chain")
	rootCmd.PersistentFlags().IntP("tx-cnt", "", 2000000, "tx count per day")
	rootCmd.PersistentFlags().IntP("usr-cnt", "", 500, "user increment per day")

	rootCmd.AddCommand(runCmd)
}

func getAccountPrivateKey(mnemonic string) (*ecdsa.PrivateKey, error) {
	hdPath := hd.CreateHDPath(60, 0, 0)
	privKeyBytes, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to derive private key from mnemonic")
	}
	privKey := &ethsecp256k1.PrivKey{Key: privKeyBytes}
	return crypto.HexToECDSA(hex.EncodeToString(privKey.Bytes()))
}

func doing(rpcUrl string, userIncrementPerDay, speed int) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	faucetPk, err := getAccountPrivateKey(evmFaucetMnemonic)
	if err != nil {
		log.Fatalf("Failed to get the faucet private key: %v", err)
	}
	faucetPkStr := hex.EncodeToString(faucetPk.D.Bytes())

	acctMgr, err := account.NewAccountManager(client, uint32(basePrefix), uint32(userIncrementPerDay), uint32(userIncrementPerDay), faucetPkStr)
	if err != nil {
		log.Fatalf("Failed to create the account manager: %v", err.Error())
	}
	generator, err := generator.NewTransferGenerator(0, 500, faucetPkStr, client, acctMgr)
	if err != nil {
		log.Fatalf("Failed to create the generator: %v", err.Error())
	}

	thisSender := &Sender{
		Speed:     speed,
		Generator: generator,
		Client:    client,
		SendCh:    generator.GenerateTransfer(),
	}

	thisSender.SendCh = generator.GenerateTransfer()
	c := cron.New()
	// c.AddFunc("0 0 0 * *", func() {
	// 	acctMgr.Increament()
	// })

	c.AddFunc("* */2 * * *", func() {
		acctMgr.Increament()
	})

	defer func() {
		c.Stop()
		generator.TearDown()
	}()
	c.Start()
	thisSender.Send()
}
