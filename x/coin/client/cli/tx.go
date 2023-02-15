package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdCreateCoin(),
		cmdUpdateCoin(),
		cmdSendCoin(),
		cmdMultiSendCoin(),
		cmdBuyCoin(),
		cmdSellCoin(),
		cmdSellAllCoin(),
		cmdBurnCoin(),
		cmdIssueCheck(),
		cmdRedeemCheck(),
	)

	return cmd
}

func cmdCreateCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [denom] [title] [crr] [initReserve] [initVolume] [limitVolume] [minVolume] [identity]",
		Short: "Creates new coin",
		Long: fmt.Sprintf(`Create custom coin. Reserve, volumes must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s create coin1 "title of coin" 20 10000000000 200000000 10000000000000 "coin identity" --from mykey`,
			cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from           = clientCtx.GetFromAddress()
				denom          = args[0]
				title          = args[1]
				initReserve, _ = sdk.NewIntFromString(args[3])
				initVolume, _  = sdk.NewIntFromString(args[4])
				limitVolume, _ = sdk.NewIntFromString(args[5])
				minVolume, _   = sdk.NewIntFromString(args[6])
				identity       = args[7]
			)

			crr, err := strconv.ParseUint(args[2], 10, 8)
			if err != nil {
				return errors.InvalidCRR
			}

			err = existCoinDenom(clientCtx, denom)
			if err == nil {
				return errors.CoinAlreadyExists
			}

			msg := types.NewMsgCreateCoin(from, denom, title, crr, initVolume, initReserve, limitVolume, minVolume, identity)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdUpdateCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [denom] [limitVolume] [minVolume] [identity]",
		Short: "Update custom coin",
		Long: fmt.Sprintf(`Update your custom coin parameters: limit volume and identity.
Limit volume must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s update coin1 10000000 "some identity" --from mykey`, cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from     = clientCtx.GetFromAddress()
				denom    = args[0]
				identity = args[3]
			)

			limitVolume, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid limit volume")
			}

			minVolume, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid min volume")
			}

			// Check if coin does not exist yet
			resp, err := getCoin(clientCtx, denom)
			if err != nil {
				return err
			}

			if resp.Coin.Creator != from.String() {
				return errors.UpdateOnlyForCreator
			}

			msg := types.NewMsgUpdateCoin(from, denom, limitVolume, minVolume, identity)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdSendCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [recipient] [coinAmount] ",
		Short: "Send coin",
		Long: fmt.Sprintf(`Send coin from one account to other.
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s send dx1hs2wdrm87c92rzhq0vgmgrxr6u57xpr2lcygc2 1000del --from mykey`, cmdcfg.AppBinName, types.ModuleName),
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

			msg := types.NewMsgSendCoin(from, address, coins)

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdMultiSendCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisend [recipient1] [coinToSend1] [recipient2] [coinToSend2]...",
		Short: "Multisend coin",
		Long: fmt.Sprintf(`Send coins from one account to others accounts
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s multisend dx1a..a 1000del dx1b..b 1000tony --from mykey
`, cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from    = clientCtx.GetFromAddress()
				argsLen = len(args)
				sends   = make([]types.MultiSendEntry, 0)
				coins   = make([]sdk.Coin, 0)
			)

			if argsLen%2 != 0 {
				return fmt.Errorf(
					"the number of arguments must be even, put either %d or %d",
					argsLen-1, argsLen+1,
				)
			}

			for i := 0; i < argsLen; i += 2 {
				recipient, err := sdk.AccAddressFromBech32(args[i])
				if err != nil {
					return err
				}

				coin, err := parseCoin(clientCtx, args[i+1])
				if err != nil {
					return err
				}

				send := types.MultiSendEntry{
					Recipient: recipient.String(),
					Coin:      coin,
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
				return errors.InsufficientFunds
			}

			msg := types.NewMsgMultiSendCoin(from, sends)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdBuyCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [amountCoinToBuy] [maxAmountCoinToSell]",
		Short: "Buy coin",
		Long: fmt.Sprintf(`Exchange some coin amount from your wallet to another coin.
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s buy 10000000000tony 12000000del --from mykey`, cmdcfg.AppBinName, types.ModuleName),
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

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdSellCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell [coinAmountToSell] [coinMinAmountToBuy]",
		Short: "Sell coin",
		Long: fmt.Sprintf(`Exchange some coin amount from your wallet to another coin.
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s sell 10000000000tony 12000000del --from mykey`, cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(2),
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

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdSellAllCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell-all [coinDenomToSell] [coinMinAmountToBuy]",
		Short: "Sell all amount of coin",
		Long: fmt.Sprintf(`Sell all amount  of coin with a specific denom from your wallet to buy another coin
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s sell-all del 100000000tony --from mykey
`, cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				from               = clientCtx.GetFromAddress()
				coinDenomToSell    = args[0]
				coinMinAmountToBuy = args[1]
			)

			minAmountToBuy, err := parseCoin(clientCtx, coinMinAmountToBuy)
			if err != nil {
				return err
			}

			err = existCoinDenom(clientCtx, coinDenomToSell)
			if err != nil {
				return err
			}

			err = checkBalance(clientCtx, from, sdkmath.NewInt(1), coinDenomToSell)
			if err != nil {
				return err
			}

			msg := types.NewMsgSellAllCoin(from, coinDenomToSell, minAmountToBuy)
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdBurnCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [coinToBurn]",
		Short: "Burn specified amount of coin",
		Long: fmt.Sprintf(`Burn specified amount of coin from your wallet
Coin amount must be with all 0 (1 coin = 10^18)

Example: 	
$ %s tx %s burn 1000del --from mykey
`, cmdcfg.AppBinName, types.ModuleName),
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

			// broadcast tx
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdIssueCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-check [coinAmount] [nonce] [dueBlock] [passphrase]",
		Short: "Issue check",
		Long: fmt.Sprintf(`Create a check for cashing with another account
Coin amount must be with all 0 (1 coin = 10^18).
'nonce' must be any positive integer.
'dueBlock' - height (block number) of blockchain until check is valid.
'passphrase' - secret key. Only recipient must known it.

Example: 	
$ %s tx %s issue-check 1000del 10 235 "some secret" --from mykey
`, cmdcfg.AppBinName, types.ModuleName),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get args
			var (
				coinStr    = args[0]
				nonce, _   = sdk.NewIntFromString(args[1])
				passphrase = args[3]
			)
			dueBlock, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			// parse tokens and check if such a coin exists
			coin, err := parseCoin(clientCtx, coinStr)
			if err != nil {
				return err
			}

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, err := crypto.ToECDSA(passphraseHash[:])
			if err != nil {
				return err
			}

			// Prepare check without lock
			check := &types.Check{
				ChainID:  clientCtx.ChainID,
				Coin:     coin,
				Nonce:    nonce.BigInt().Bytes(),
				DueBlock: dueBlock,
			}

			// Prepare check lock
			checkHash := check.HashWithoutLock()
			lock, err := crypto.Sign(checkHash[:], passphrasePrivKey)
			if err != nil {
				return err
			}

			// Fill check with prepared lock
			check.Lock = lock

			// Sign check by check issuer
			checkHash = check.Hash()
			signature, _, err := clientCtx.Keyring.Sign(clientCtx.FromName, checkHash[:])
			if err != nil {
				panic(errors.UnableSignCheck)
			}

			check.SetSignature(signature)

			checkBytes, err := rlp.EncodeToBytes(check)
			if err != nil {
				return errors.UnableRLPEncodeToBytesCheck
			}

			return clientCtx.PrintString(base58.Encode(checkBytes))
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom) // nolint:errcheck

	return cmd
}

func cmdRedeemCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-check [check] [passphrase]",
		Short: "Redeem check",
		Long: fmt.Sprintf(`Redeem check 

Example: 	
$ %s tx %s redeem-check 3YEtqixL7ccFTZJaMUHx3T...(result of 'issue-check') "some secret" --from mykey 
`, cmdcfg.AppBinName, types.ModuleName),
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
				return errors.UnableDecodeCheckBase58
			}

			// Parse provided check from raw bytes to ensure it is valid
			_, err = types.ParseCheck(checkBytes)
			if err != nil {
				return err
			}

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, err := crypto.ToECDSA(passphraseHash[:])
			if err != nil {
				return errors.InvalidPassphrase
			}

			// Prepare bytes to sign by private key generated from passphrase
			receiverAddressHash := make([]byte, 32)
			hw := sha3.NewLegacyKeccak256()
			err = rlp.Encode(hw, []interface{}{
				clientCtx.GetFromAddress(),
			})
			if err != nil {
				return errors.UnableRLPEncodeAddress
			}
			hw.Sum(receiverAddressHash[:0])

			// Sign receiver address by private key generated from passphrase
			signature, err := crypto.Sign(receiverAddressHash, passphrasePrivKey)
			if err != nil {
				return errors.UnableSignCheck
			}
			proofBase64 := base64.StdEncoding.EncodeToString(signature)

			// Prepare redeem check message
			msg := types.NewMsgRedeemCheck(clientCtx.GetFromAddress(), checkBase58, proofBase64)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
