package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtestutil "github.com/cosmos/cosmos-sdk/x/auth/client/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov"
)

type CLITestSuite struct {
	suite.Suite

	kr        keyring.Keyring
	encCfg    testutilmod.TestEncodingConfig
	baseCtx   client.Context
	clientCtx client.Context
	val       sdk.AccAddress
	val1      sdk.AccAddress
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.encCfg = testutilmod.MakeTestEncodingConfig(auth.AppModuleBasic{}, bank.AppModuleBasic{}, gov.AppModuleBasic{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec, keyring.ETHAlgoOption())
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockTendermintRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID(testutil.DefaultChainId)

	var outBuf bytes.Buffer
	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
			Value: bz,
		})
		return s.baseCtx.WithClient(c)
	}
	s.clientCtx = ctxGen().WithOutput(&outBuf)

	kb := s.clientCtx.Keyring
	valAcc, _, err := kb.NewMnemonic("newAccount", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.EthSecp256k1)
	s.Require().NoError(err)
	s.val, err = valAcc.GetAddress()
	s.Require().NoError(err)

	account1, _, err := kb.NewMnemonic("newAccount1", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.EthSecp256k1)
	s.Require().NoError(err)
	s.val1, err = account1.GetAddress()
	s.Require().NoError(err)

	_, _, err = kb.NewMnemonic("newAccount2", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.EthSecp256k1)
	s.Require().NoError(err)

	// Create a dummy account for testing purpose
	_, _, err = kb.NewMnemonic("dummyAccount", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.EthSecp256k1)
	s.Require().NoError(err)
}

func (s *CLITestSuite) TestCLIValidateSignatures() {
	sendTokens := sdk.NewCoins(
		sdk.NewCoin("testtoken", sdk.NewInt(10)),
		sdk.NewCoin("stake", sdk.NewInt(10)))

	res, err := s.createBankMsg(s.clientCtx, s.val, sendTokens,
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly))
	s.Require().NoError(err)

	// write  unsigned tx to file
	unsignedTx := testutil.WriteToNewTempFile(s.T(), res.String())
	defer unsignedTx.Close()

	res, err = authtestutil.TxSignExec(s.clientCtx, s.val, unsignedTx.Name())
	s.Require().NoError(err)
	signedTx, err := s.clientCtx.TxConfig.TxJSONDecoder()(res.Bytes())
	s.Require().NoError(err)

	signedTxFile := testutil.WriteToNewTempFile(s.T(), res.String())
	defer signedTxFile.Close()
	txBuilder, err := s.clientCtx.TxConfig.WrapTxBuilder(signedTx)
	s.Require().NoError(err)
	_, err = authtestutil.TxValidateSignaturesExec(s.clientCtx, signedTxFile.Name())
	s.Require().NoError(err)

	txBuilder.SetMemo("MODIFIED TX")
	bz, err := s.clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)

	modifiedTxFile := testutil.WriteToNewTempFile(s.T(), string(bz))
	defer modifiedTxFile.Close()

	_, err = authtestutil.TxValidateSignaturesExec(s.clientCtx, modifiedTxFile.Name())
	s.Require().EqualError(err, "signatures validation failed")
}

func (s *CLITestSuite) TestCLISignBatch() {
	sendTokens := sdk.NewCoins(
		sdk.NewCoin("testtoken", sdk.NewInt(10)),
		sdk.NewCoin("stake", sdk.NewInt(10)),
	)

	generatedStd, err := s.createBankMsg(s.clientCtx, s.val,
		sendTokens, fmt.Sprintf("--%s=true", flags.FlagGenerateOnly))
	s.Require().NoError(err)

	outputFile := testutil.WriteToNewTempFile(s.T(), strings.Repeat(generatedStd.String(), 3))
	defer outputFile.Close()
	s.clientCtx.HomeDir = strings.Replace(s.clientCtx.HomeDir, "simd", "simcli", 1)

	// sign-batch file - offline is set but account-number and sequence are not
	_, err = authtestutil.TxSignBatchExec(s.clientCtx, s.val, outputFile.Name(), fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID), "--offline")
	s.Require().EqualError(err, "required flag(s) \"account-number\", \"sequence\" not set")

	// sign-batch file - offline and sequence is set but account-number is not set
	_, err = authtestutil.TxSignBatchExec(s.clientCtx, s.val, outputFile.Name(), fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID), fmt.Sprintf("--%s=%s", flags.FlagSequence, "1"), "--offline")
	s.Require().EqualError(err, "required flag(s) \"account-number\" not set")

	// sign-batch file - offline and account-number is set but sequence is not set
	_, err = authtestutil.TxSignBatchExec(s.clientCtx, s.val, outputFile.Name(), fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID), fmt.Sprintf("--%s=%s", flags.FlagAccountNumber, "1"), "--offline")
	s.Require().EqualError(err, "required flag(s) \"sequence\" not set")
}

func (s *CLITestSuite) TestCLIQueryTxCmdByHash() {
	sendTokens := sdk.NewInt64Coin("stake", 10)

	// Send coins.
	out, err := s.createBankMsg(
		s.clientCtx, s.val,
		sdk.NewCoins(sendTokens),
	)
	s.Require().NoError(err)

	var txRes sdk.TxResponse
	s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes))

	testCases := []struct {
		name         string
		args         []string
		expCmdOutput string
	}{
		{
			"not enough args",
			[]string{},
			"",
		},
		{
			"with invalid hash",
			[]string{"somethinginvalid", fmt.Sprintf("--%s=json", flags.FlagOutput)},
			`[somethinginvalid --output=json]`,
		},
		{
			"with valid and not existing hash",
			[]string{"C7E7D3A86A17AB3A321172239F3B61357937AF0F25D9FA4D2F4DCCAD9B0D7747", fmt.Sprintf("--%s=json", flags.FlagOutput)},
			`[C7E7D3A86A17AB3A321172239F3B61357937AF0F25D9FA4D2F4DCCAD9B0D7747 --output=json`,
		},
		{
			"happy case",
			[]string{txRes.TxHash, fmt.Sprintf("--%s=json", flags.FlagOutput)},
			fmt.Sprintf("%s --%s=json", txRes.TxHash, flags.FlagOutput),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			cmd := authcli.QueryTxCmd()
			cmd.SetArgs(tc.args)

			if len(tc.args) != 0 {
				s.Require().Contains(fmt.Sprint(cmd), tc.expCmdOutput)
			}
		})
	}
}

func (s *CLITestSuite) TestCLIQueryTxCmdByEvents() {
	testCases := []struct {
		name         string
		args         []string
		expCmdOutput string
	}{
		{
			"invalid --type",
			[]string{
				fmt.Sprintf("--type=%s", "foo"),
				"bar",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"--type=foo bar --output=json",
		},
		{
			"--type=acc_seq with no addr+seq",
			[]string{
				"--type=acc_seq",
				"",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"--type=acc_seq  --output=json",
		},
		{
			"non-existing addr+seq combo",
			[]string{
				"--type=acc_seq",
				"foobar",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"--type=acc_seq foobar --output=json",
		},
		{
			"--type=signature with no signature",
			[]string{
				"--type=signature",
				"",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"--type=signature  --output=json",
		},
		{
			"non-existing signatures",
			[]string{
				"--type=signature",
				"foo",
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"--type=signature foo --output=json",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			cmd := authcli.QueryTxCmd()
			cmd.SetArgs(tc.args)

			if len(tc.args) != 0 {
				s.Require().Contains(fmt.Sprint(cmd), tc.expCmdOutput)
			}
		})
	}
}

func (s *CLITestSuite) TestCLIQueryTxsCmdByEvents() {
	testCases := []struct {
		name         string
		args         []string
		expCmdOutput string
	}{
		{
			"fee event happy case",
			[]string{
				fmt.Sprintf("--events=tx.fee=%s",
					sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"",
		},
		{
			"no matching fee event",
			[]string{
				fmt.Sprintf("--events=tx.fee=%s",
					sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(0))).String()),
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			"",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			cmd := authcli.QueryTxsByEventsCmd()

			if len(tc.args) != 0 {
				s.Require().Contains(fmt.Sprint(cmd), tc.expCmdOutput)
			}
		})
	}
}

func (s *CLITestSuite) TestCLISendGenerateSignAndBroadcast() {
	sendTokens := sdk.NewCoin("stake", sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	normalGeneratedTx, err := s.createBankMsg(s.clientCtx, s.val,
		sdk.NewCoins(sendTokens), fmt.Sprintf("--%s=true", flags.FlagGenerateOnly))
	s.Require().NoError(err)

	txCfg := s.clientCtx.TxConfig

	normalGeneratedStdTx, err := txCfg.TxJSONDecoder()(normalGeneratedTx.Bytes())
	s.Require().NoError(err)
	txBuilder, err := txCfg.WrapTxBuilder(normalGeneratedStdTx)
	s.Require().NoError(err)
	s.Require().Equal(txBuilder.GetTx().GetGas(), uint64(flags.DefaultGasLimit))
	s.Require().Equal(len(txBuilder.GetTx().GetMsgs()), 1)
	sigs, err := txBuilder.GetTx().GetSignaturesV2()
	s.Require().NoError(err)
	s.Require().Equal(0, len(sigs))

	// Test generate sendTx with --gas=$amount
	limitedGasGeneratedTx, err := s.createBankMsg(s.clientCtx, s.val,
		sdk.NewCoins(sendTokens), fmt.Sprintf("--gas=%d", 100),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
	)
	s.Require().NoError(err)

	limitedGasStdTx, err := txCfg.TxJSONDecoder()(limitedGasGeneratedTx.Bytes())
	s.Require().NoError(err)
	txBuilder, err = txCfg.WrapTxBuilder(limitedGasStdTx)
	s.Require().NoError(err)
	s.Require().Equal(txBuilder.GetTx().GetGas(), uint64(100))
	s.Require().Equal(len(txBuilder.GetTx().GetMsgs()), 1)
	sigs, err = txBuilder.GetTx().GetSignaturesV2()
	s.Require().NoError(err)
	s.Require().Equal(0, len(sigs))

	// Test generate sendTx, estimate gas
	finalGeneratedTx, err := s.createBankMsg(s.clientCtx, s.val,
		sdk.NewCoins(sendTokens), fmt.Sprintf("--gas=%d", flags.DefaultGasLimit),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly))
	s.Require().NoError(err)

	finalStdTx, err := txCfg.TxJSONDecoder()(finalGeneratedTx.Bytes())
	s.Require().NoError(err)
	txBuilder, err = txCfg.WrapTxBuilder(finalStdTx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(flags.DefaultGasLimit), txBuilder.GetTx().GetGas())
	s.Require().Equal(len(finalStdTx.GetMsgs()), 1)

	// Write the output to disk
	unsignedTxFile := testutil.WriteToNewTempFile(s.T(), finalGeneratedTx.String())
	defer unsignedTxFile.Close()

	// Test validate-signatures
	res, err := authtestutil.TxValidateSignaturesExec(s.clientCtx, unsignedTxFile.Name())
	s.Require().EqualError(err, "signatures validation failed")
	s.Require().True(strings.Contains(res.String(), fmt.Sprintf("Signers:\n  0: %v\n\nSignatures:\n\n", s.val.String())))

	// Test sign

	// Does not work in offline mode
	_, err = authtestutil.TxSignExec(s.clientCtx, s.val, unsignedTxFile.Name(), "--offline")
	s.Require().EqualError(err, "required flag(s) \"account-number\", \"sequence\" not set")

	// But works offline if we set account number and sequence
	s.clientCtx.HomeDir = strings.Replace(s.clientCtx.HomeDir, "simd", "simcli", 1)
	_, err = authtestutil.TxSignExec(s.clientCtx, s.val, unsignedTxFile.Name(), "--offline", "--account-number", "1", "--sequence", "1")
	s.Require().NoError(err)

	// Sign transaction
	signedTx, err := authtestutil.TxSignExec(s.clientCtx, s.val, unsignedTxFile.Name())
	s.Require().NoError(err)
	signedFinalTx, err := txCfg.TxJSONDecoder()(signedTx.Bytes())
	s.Require().NoError(err)
	txBuilder, err = s.clientCtx.TxConfig.WrapTxBuilder(signedFinalTx)
	s.Require().NoError(err)
	s.Require().Equal(len(txBuilder.GetTx().GetMsgs()), 1)
	sigs, err = txBuilder.GetTx().GetSignaturesV2()
	s.Require().NoError(err)
	s.Require().Equal(1, len(sigs))
	s.Require().Equal(s.val.String(), txBuilder.GetTx().GetSigners()[0].String())

	// Write the output to disk
	signedTxFile := testutil.WriteToNewTempFile(s.T(), signedTx.String())
	defer signedTxFile.Close()

	// validate Signature
	res, err = authtestutil.TxValidateSignaturesExec(s.clientCtx, signedTxFile.Name())
	s.Require().NoError(err)
	s.Require().True(strings.Contains(res.String(), "[OK]"))

	// Test broadcast

	// Does not work in offline mode
	_, err = authtestutil.TxBroadcastExec(s.clientCtx, signedTxFile.Name(), "--offline")
	s.Require().EqualError(err, "cannot broadcast tx during offline mode")

	// Broadcast correct transaction.
	s.clientCtx.BroadcastMode = flags.BroadcastSync
	_, err = authtestutil.TxBroadcastExec(s.clientCtx, signedTxFile.Name())
	s.Require().NoError(err)
}

func (s *CLITestSuite) TestCLIEncode() {
	sendTokens := sdk.NewCoin("stake", sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	normalGeneratedTx, err := s.createBankMsg(
		s.clientCtx, s.val,
		sdk.NewCoins(sendTokens),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=deadbeef", flags.FlagNote),
	)
	s.Require().NoError(err)
	savedTxFile := testutil.WriteToNewTempFile(s.T(), normalGeneratedTx.String())
	defer savedTxFile.Close()

	// Encode
	encodeExec, err := authtestutil.TxEncodeExec(s.clientCtx, savedTxFile.Name())
	s.Require().NoError(err)
	trimmedBase64 := strings.Trim(encodeExec.String(), "\"\n")

	// Check that the transaction decodes as expected
	decodedTx, err := authtestutil.TxDecodeExec(s.clientCtx, trimmedBase64)
	s.Require().NoError(err)

	txCfg := s.clientCtx.TxConfig
	theTx, err := txCfg.TxJSONDecoder()(decodedTx.Bytes())
	s.Require().NoError(err)
	txBuilder, err := s.clientCtx.TxConfig.WrapTxBuilder(theTx)
	s.Require().NoError(err)
	s.Require().Equal("deadbeef", txBuilder.GetTx().GetMemo())
}

func (s *CLITestSuite) TestSignWithMultisig() {
	// Generate a account for signing.
	account1, err := s.clientCtx.Keyring.Key("newAccount1")
	s.Require().NoError(err)

	addr1, err := account1.GetAddress()
	s.Require().NoError(err)

	// Create an address that is not in the keyring, will be used to simulate `--multisig`
	multisig := "0x8ff4ddbDd327DB98AE56baF3Bf05Be0aD6ad7EeF"
	multisigAddr, err := sdk.AccAddressFromHexUnsafe(multisig)
	s.Require().NoError(err)

	// Generate a transaction for testing --multisig with an address not in the keyring.
	multisigTx, err := clitestutil.MsgSendExec(
		s.clientCtx,
		s.val,
		s.val,
		sdk.NewCoins(
			sdk.NewInt64Coin("stake", 5),
		),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
	)
	s.Require().NoError(err)

	// Save multi tx to file
	multiGeneratedTx2File := testutil.WriteToNewTempFile(s.T(), multisigTx.String())
	defer multiGeneratedTx2File.Close()

	// Sign using multisig. We're signing a tx on behalf of the multisig address,
	// even though the tx signer is NOT the multisig address. This is fine though,
	// as the main point of this test is to test the `--multisig` flag with an address
	// that is not in the keyring.
	_, err = authtestutil.TxSignExec(s.clientCtx, addr1, multiGeneratedTx2File.Name(), "--multisig", multisigAddr.String())
	s.Require().Contains(err.Error(), "error getting account from keybase")
}

func (s *CLITestSuite) TestGetBroadcastCommandOfflineFlag() {
	cmd := authcli.GetBroadcastCommand()
	_ = testutil.ApplyMockIODiscardOutErr(cmd)
	cmd.SetArgs([]string{fmt.Sprintf("--%s=true", flags.FlagOffline), ""})

	s.Require().EqualError(cmd.Execute(), "cannot broadcast tx during offline mode")
}

func (s *CLITestSuite) TestGetBroadcastCommandWithoutOfflineFlag() {
	txCfg := s.clientCtx.TxConfig
	clientCtx := client.Context{}
	clientCtx = clientCtx.WithTxConfig(txCfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd := authcli.GetBroadcastCommand()
	_, out := testutil.ApplyMockIO(cmd)

	// Create new file with tx
	builder := txCfg.NewTxBuilder()
	builder.SetGasLimit(200000)
	from, err := sdk.AccAddressFromHexUnsafe("0xF489272cfE70678A431bD850c5CAeb5026903a46")
	s.Require().NoError(err)
	to, err := sdk.AccAddressFromHexUnsafe("0xF489272cfE70678A431bD850c5CAeb5026903a46")
	s.Require().NoError(err)
	err = builder.SetMsgs(banktypes.NewMsgSend(from, to, sdk.Coins{sdk.NewInt64Coin("stake", 10000)}))
	s.Require().NoError(err)
	txContents, err := txCfg.TxJSONEncoder()(builder.GetTx())
	s.Require().NoError(err)
	txFile := testutil.WriteToNewTempFile(s.T(), string(txContents))
	defer txFile.Close()

	cmd.SetArgs([]string{txFile.Name()})
	err = cmd.ExecuteContext(ctx)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "connect: connection refused")
	s.Require().Contains(out.String(), "connect: connection refused")
}

func (s *CLITestSuite) TestQueryParamsCmd() {
	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"happy case",
			[]string{fmt.Sprintf("--%s=json", flags.FlagOutput)},
			false,
		},
		{
			"with specific height",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			cmd := authcli.QueryParamsCmd()
			clientCtx := s.clientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().NotEqual("internal", err.Error())
			} else {
				var authParams authtypes.Params
				s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), &authParams))
				s.Require().NotNil(authParams.MaxMemoCharacters)
			}
		})
	}
}

// TestTxWithoutPublicKey makes sure sending a proto tx message without the
// public key doesn't cause any error in the RPC layer (broadcast).
// See https://github.com/cosmos/cosmos-sdk/issues/7585 for more details.
func (s *CLITestSuite) TestTxWithoutPublicKey() {
	txCfg := s.clientCtx.TxConfig

	// Create a txBuilder with an unsigned tx.
	txBuilder := txCfg.NewTxBuilder()
	msg := banktypes.NewMsgSend(s.val, s.val, sdk.NewCoins(
		sdk.NewCoin("Stake", sdk.NewInt(10)),
	))
	err := txBuilder.SetMsgs(msg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("Stake", sdk.NewInt(150))))
	txBuilder.SetGasLimit(testdata.NewTestGasLimit())

	// Create a file with the unsigned tx.
	txJSON, err := txCfg.TxJSONEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)
	unsignedTxFile := testutil.WriteToNewTempFile(s.T(), string(txJSON))
	defer unsignedTxFile.Close()

	// Sign the file with the unsignedTx.
	signedTx, err := authtestutil.TxSignExec(s.clientCtx, s.val, unsignedTxFile.Name(), fmt.Sprintf("--%s=true", cli.FlagOverwrite))
	s.Require().NoError(err)

	// Remove the signerInfo's `public_key` field manually from the signedTx.
	// Note: this method is only used for test purposes! In general, one should
	// use txBuilder and TxEncoder/TxDecoder to manipulate txs.
	var tx tx.Tx
	err = s.clientCtx.Codec.UnmarshalJSON(signedTx.Bytes(), &tx)
	s.Require().NoError(err)
	tx.AuthInfo.SignerInfos[0].PublicKey = nil
	// Re-encode the tx again, to another file.
	txJSON, err = s.clientCtx.Codec.MarshalJSON(&tx)
	s.Require().NoError(err)
	signedTxFile := testutil.WriteToNewTempFile(s.T(), string(txJSON))
	defer signedTxFile.Close()
	s.Require().True(strings.Contains(string(txJSON), "\"public_key\":null"))

	// Broadcast tx, test that it shouldn't panic.
	s.clientCtx.BroadcastMode = flags.BroadcastSync
	out, err := authtestutil.TxBroadcastExec(s.clientCtx, signedTxFile.Name())
	s.Require().NoError(err)
	var res sdk.TxResponse
	s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
	s.Require().NotEqual(0, res.Code)
}

func (s *CLITestSuite) createBankMsg(clientCtx client.Context, toAddr sdk.AccAddress, amount sdk.Coins, extraFlags ...string) (testutil.BufferWriter, error) {
	flags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
	}

	flags = append(flags, extraFlags...)
	return clitestutil.MsgSendExec(clientCtx, s.val, toAddr, amount, flags...)
}
