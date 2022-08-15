package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/version"
	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"golang.org/x/crypto/sha3"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

//TODO: maybe create recover func and do panics instead of return err

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	coinCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinCmd.AddCommand(
		NewCreateCoinCmd(),
		NewUpdateCoinCmd(),
		NewBuyCoinCmd(),
		NewSellCoinCmd(),
		NewSendCoinCmd(),
		NewMultiSendCoinCmd(),
		NewSellAllCoinCmd(),
		NewIssueCheckCmd(),
		NewRedeemCheckCmd(),
	)

	return coinCmd
}

func NewCreateCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume] [identity]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from           = clientCtx.GetFromAddress()
				title          = args[0]
				symbol         = args[1]
				initReserve, _ = sdk.NewIntFromString(args[3])
				initVolume, _  = sdk.NewIntFromString(args[4])
				limitVolume, _ = sdk.NewIntFromString(args[5])
				identity       = args[6]
			)

			crr, err := strconv.ParseUint(args[2], 10, 8)
			if err != nil {
				return types.ErrInvalidCRR(args[2])
			}

			err = existCoinSymbol(clientCtx, symbol)
			if err == nil {
				return types.ErrCoinAlreadyExists(symbol)
			}

			msg := types.NewMsgCreateCoin(from, title, symbol, crr, initVolume, initReserve, limitVolume, identity)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewUpdateCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [symbol] [limitVolume] [identity]",
		Short: "Update custom coin",
		Long: fmt.Sprintf(`update coin your custom coin parametres 

Example: 	
$ %s tx coin update del 10000000 "" --from mykey`, version.AppName),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from     = clientCtx.GetFromAddress()
				symbol   = args[0]
				identity = args[2]
			)

			limitVolume, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid limit volume")
			}

			// Check if coin does not exist yet
			resp, err := getCoin(clientCtx, symbol)
			if err != nil {
				return err
			}

			if resp.Coin.Creator != from.String() {
				return types.ErrUpdateOnlyForCreator()
			}

			msg := types.NewMsgUpdateCoin(from, symbol, limitVolume, identity)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewBuyCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [amountCoinToBuy] [maxAmountCoinToSell]",
		Short: "Buy coin",
		Long: fmt.Sprintf(`change one token from your wallet to another 

Example: 	
$ %s tx coin buy 1000tony 1000del --from mykey
`, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from                = clientCtx.GetFromAddress()
				amountCoinToBuy     = args[0]
				maxAmountCoinToSell = args[1]
			)

			// parse tokens and check if such a coin exists
			coinToBuy, err := parseCoin(clientCtx, amountCoinToBuy)
			if err != nil {
				return err
			}

			maxAmountToSell, err := parseCoin(clientCtx, maxAmountCoinToSell)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, maxAmountToSell.Amount, maxAmountToSell.Denom)
			if err != nil {
				return err
			}

			// create msg
			msg := types.NewMsgBuyCoin(from, coinToBuy, maxAmountToSell)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewSellCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell [coinAmountToSell] [coinMinAmountToBuy]",
		Short: "Sell coin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get from address
			var (
				from               = clientCtx.GetFromAddress()
				coinAmountToSell   = args[0]
				coinMinAmountToBuy = args[1]
			)

			// parse tokens and check if such a coin exists
			coinToSell, err := parseCoin(clientCtx, coinAmountToSell)
			if err != nil {
				return err
			}

			minAmountToBuy, err := parseCoin(clientCtx, coinMinAmountToBuy)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, coinToSell.Amount, coinToSell.Denom)
			if err != nil {
				return err
			}

			msg := types.NewMsgSellCoin(from, coinToSell, minAmountToBuy)
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewSendCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [receiver] [coinAmount] ",
		Short: "Send coin",
		Long: fmt.Sprintf(`send coins from one account to other 

Example: 	
$ %s tx coin send dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2 1000del --from mykey
`, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from       = clientCtx.GetFromAddress()
				addressStr = args[0]
				coinAmount = args[1]
			)

			coins, err := parseCoin(clientCtx, coinAmount)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(addressStr)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, coins.Amount, coins.Denom)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendCoin(from, coins, address)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewMultiSendCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisend [receiver] [coinToSendAmount]...",
		Short: "Multisend coin",
		Long: fmt.Sprintf(`send coins from one account to others accounts

Example: 	
$ %s tx coin multisend dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2 1000del  dx1hs2wdrmrzhq0vgmgrxr87c926u57xpr2lcygc2 1000tony --from mykey
`, version.AppName),
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from    = clientCtx.GetFromAddress()
				argsLen = len(args)
				sends   = make([]types.Send, 0)
				coins   = make([]sdk.Coin, 0)
			)

			if argsLen%2 != 0 {
				return fmt.Errorf(
					"the number of arguments must be even, put either %d or %d",
					argsLen-1, argsLen+1,
				)
			}

			for i := 0; i < argsLen; i += 2 {
				receiver, err := sdk.AccAddressFromBech32(args[i])
				if err != nil {
					return err
				}

				coin, err := parseCoin(clientCtx, args[i+1])
				if err != nil {
					return err
				}

				send := types.Send{
					Receiver: receiver.String(),
					Coin:     coin,
				}

				sends = append(sends, send)
				coins = append(coins, coin)
			}

			// Check if enough balance
			balances, err := getBalances(clientCtx, from, &query.PageRequest{})
			if err != nil {
				return err
			}

			balance := balances.Balances
			if !balance.IsAllGTE(coins) {
				var wantFunds string
				for _, send := range sends {
					wantFunds += send.Coin.String() + ", "
				}
				wantFunds = strings.TrimSuffix(wantFunds, ", ")
				return types.ErrInsufficientFunds(wantFunds, balance.String())
			}

			msg := types.NewMsgMultiSendCoin(from, sends)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewSellAllCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell-all [coinToSellSymbol] [coinMinAmountToBuy]",
		Short: "Sell all coin",
		Long: fmt.Sprintf(`sell all tokens with a specific symbol from your wallet to buy another token

Example: 	
$ %s tx coin sell_all del 100000000tony --from mykey
`, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from               = clientCtx.GetFromAddress()
				coinToSellSymbol   = args[0]
				coinMinAmountToBuy = args[1]
			)

			minAmountToBuy, err := parseCoin(clientCtx, coinMinAmountToBuy)
			if err != nil {
				return err
			}

			err = existCoinSymbol(clientCtx, coinToSellSymbol)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, sdk.NewInt(1), coinToSellSymbol)
			if err != nil {
				return err
			}

			msg := types.NewMsgSellAllCoin(from, sdk.NewCoin(coinToSellSymbol, sdk.NewInt(0)), minAmountToBuy)
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewBurnCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [coinAmount]",
		Short: "Burn coin",
		Long: fmt.Sprintf(`burn coins 

Example: 	
$ %s tx coin burn 1000del --from mykey
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from       = clientCtx.GetFromAddress()
				coinAmount = args[0]
			)

			coins, err := parseCoin(clientCtx, coinAmount)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, coins.Amount, coins.Denom)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnCoin(from, coins)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewIssueCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-check [coinAmount] [nonce] [dueBlock] [passphrase]",
		Short: "Issue check",
		Long: fmt.Sprintf(`Redeem your transaction

Example: 	
$ %s tx coin issue-check 1000del 10 235 123456789 --from mykey
`, version.AppName),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get args
			var (
				coinAmountStr = args[0]
				nonce, _      = sdk.NewIntFromString(args[1])
				passphrase    = args[3]
			)
			dueBlock, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			// parse tokens and check if such a coin exists
			coinAmount, err := parseCoin(clientCtx, coinAmountStr)
			if err != nil {
				return err
			}

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, _ := ethereumCrypto.ToECDSA(passphraseHash[:])

			// Prepare check without lock
			check := &types.Check{
				ChainID:  clientCtx.ChainID,
				Coin:     coinAmount.Denom,
				Amount:   coinAmount.Amount,
				Nonce:    nonce.BigInt().Bytes(),
				DueBlock: dueBlock,
			}

			// Prepare check lock
			checkHash := check.HashWithoutLock()
			lock, _ := ethereumCrypto.Sign(checkHash[:], passphrasePrivKey)

			// Fill check with prepared lock
			check.Lock = lock
			kr := clientCtx.Keyring

			privKeyArmored, err := kr.ExportPrivKeyArmor(clientCtx.FromName, passphrase)
			if err != nil {
				return types.ErrUnableRetrieveArmoredPkey(clientCtx.FromName, err.Error())
			}
			privKey, algo, err := crypto.UnarmorDecryptPrivKey(privKeyArmored, "")
			if err != nil {
				return types.ErrUnableRetrievePkey(clientCtx.FromName, err.Error())
			}
			if algo != ethsecp256k1.KeyType {
				return types.ErrUnableRetrieveSECPPkey(clientCtx.FromName, algo)
			}

			ethPrivKey, ok := privKey.(*ethsecp256k1.PrivKey)
			if !ok {
				return types.ErrInvalidPkey()
			}

			key, err := ethPrivKey.ToECDSA()
			if err != nil {
				return err
			}

			// Sign check by check issuer
			checkHash = check.Hash()
			signature, err := ethereumCrypto.Sign(checkHash[:], key)
			if err != nil {
				panic(err)
			}
			check.SetSignature(signature)

			checkBytes, err := rlp.EncodeToBytes(check)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(base58.Encode(checkBytes))
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewRedeemCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-check [check] [passphrase]",
		Short: "Redeem check",
		Long: fmt.Sprintf(`Check Redeem 

Example: 	
$ %s tx coin redeem-check {check hash} "" --from mykey 
`, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				checkBase58 = args[0]
				passphrase  = args[1] // TODO: Read passphrase by request to avoid saving it in terminal history =
			)

			// Decode provided check from base58 format to raw bytes
			checkBytes := base58.Decode(checkBase58)
			if len(checkBytes) == 0 {
			}

			// Parse provided check from raw bytes to ensure it is valid
			_, err = types.ParseCheck(checkBytes)
			if err != nil {
				return types.ErrInvalidCheck(err.Error())
			}

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, err := ethereumCrypto.ToECDSA(passphraseHash[:])
			if err != nil {
				return types.ErrInvalidPassphrase(err.Error())
			}

			// Prepare bytes to sign by private key generated from passphrase
			receiverAddressHash := make([]byte, 32)
			hw := sha3.NewLegacyKeccak256()
			err = rlp.Encode(hw, []interface{}{
				clientCtx.GetFromAddress(),
			})
			if err != nil {
				return types.ErrUnableRPLEncodeCheck(err.Error())
			}
			hw.Sum(receiverAddressHash[:0])

			// Sign receiver address by private key generated from passphrase
			signature, err := ethereumCrypto.Sign(receiverAddressHash[:], passphrasePrivKey)
			if err != nil {
				return types.ErrUnableSignCheck(err.Error())
			}
			proofBase64 := base64.StdEncoding.EncodeToString(signature)

			// Prepare redeem check message
			msg := types.NewMsgRedeemCheck(clientCtx.GetFromAddress(), checkBase58, proofBase64)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseCoin(clientCtx client.Context, amount string) (sdk.Coin, error) {
	var (
		coin sdk.Coin
		err  error
	)
	coin, err = sdk.ParseCoinNormalized(amount)
	if err != nil {
		return coin, err
	}

	resp, err := getCoin(clientCtx, coin.Denom)
	switch {
	case err != nil:
		return coin, err
	case resp == nil:
		return coin, types.ErrCoinDoesNotExist(coin.Denom)
	default:
		return coin, nil
	}
}

func checkBalance(clientCtx client.Context, address sdk.AccAddress, needAmount sdk.Int, symbol string) error {
	balance, err := getBalanceWithSymbol(clientCtx, address, symbol)
	if err != nil {
		return err
	}

	amountBalance := balance.Balance.Amount

	if amountBalance.LT(needAmount) {
		return types.ErrInsufficientFunds(needAmount.String(), amountBalance.String())
	}

	return nil
}
