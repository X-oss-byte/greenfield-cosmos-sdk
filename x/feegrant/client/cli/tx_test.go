package cli_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/feegrant/client/cli"
	"github.com/cosmos/cosmos-sdk/x/feegrant/module"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	oneYear  = 365 * 24 * 60 * 60
	tenHours = 10 * 60 * 60
	oneHour  = 60 * 60
)

type CLITestSuite struct {
	suite.Suite

	addedGranter sdk.AccAddress
	addedGrantee sdk.AccAddress
	addedGrant   feegrant.Grant

	kr        keyring.Keyring
	baseCtx   client.Context
	encCfg    testutilmod.TestEncodingConfig
	clientCtx client.Context

	accounts []sdk.AccAddress
}

func TestCLITestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.encCfg = testutilmod.MakeTestEncodingConfig(module.AppModuleBasic{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec)
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockTendermintRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID(sdktestutil.DefaultChainId)

	var outBuf bytes.Buffer
	ctxGen := func() client.Context {
		bz, _ := s.encCfg.Codec.Marshal(&sdk.TxResponse{})
		c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
			Value: bz,
		})

		return s.baseCtx.WithClient(c)
	}
	s.clientCtx = ctxGen().WithOutput(&outBuf)

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	accounts := testutil.CreateKeyringAccounts(s.T(), s.kr, 2)

	granter := accounts[0].Address
	grantee := accounts[1].Address

	s.createGrant(granter, grantee)

	grant, err := feegrant.NewGrant(granter, grantee, &feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100))),
	})
	s.Require().NoError(err)

	s.addedGrant = grant
	s.addedGranter = granter
	s.addedGrantee = grantee
	for _, v := range accounts {
		s.accounts = append(s.accounts, v.Address)
	}
	s.accounts[1] = accounts[1].Address
}

// createGrant creates a new basic allowance fee grant from granter to grantee.
func (s *CLITestSuite) createGrant(granter, grantee sdk.Address) {
	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100))).String()),
	}

	fee := sdk.NewCoin("stake", sdk.NewInt(100))

	args := append(
		[]string{
			granter.String(),
			grantee.String(),
			fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, fee.String()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
			fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(oneYear)),
		},
		commonFlags...,
	)

	cmd := cli.NewCmdFeeGrant()
	out, err := clitestutil.ExecTestCLICmd(s.clientCtx, cmd, args)
	s.Require().NoError(err)

	var resp sdk.TxResponse
	s.Require().NoError(s.clientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
	s.Require().Equal(resp.Code, uint32(0))
}

func (s *CLITestSuite) TestNewCmdFeeGrant() {
	granter := s.accounts[0]
	alreadyExistedGrantee := s.addedGrantee
	clientCtx := s.clientCtx

	fromAddr, fromName, _, err := client.GetFromFields(s.baseCtx, s.kr, granter.String())
	s.Require().Equal(fromAddr, granter)
	s.Require().NoError(err)

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
	}

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"wrong granter address",
			append(
				[]string{
					"wrong_granter",
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"wrong grantee address",
			append(
				[]string{
					granter.String(),
					"wrong_grantee",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"wrong granter key name",
			append(
				[]string{
					"invalid_granter",
					"0x15Ec1D4647Ec0485cD21856eB169E95CB595e45D",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"valid basic fee grant",
			append(
				[]string{
					granter.String(),
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid basic fee grant with granter key name",
			append(
				[]string{
					fromName,
					"0x15Ec1D4647Ec0485cD21856eB169E95CB595e45D",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, fromName),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid basic fee grant without spend limit",
			append(
				[]string{
					granter.String(),
					"0x9f1E2e29698b1b3e8bF262cdF3572aaF925Ec28C",
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid basic fee grant without expiration",
			append(
				[]string{
					granter.String(),
					"0x5093a4E223b6B98a78EfBD0c8D0C3eb44EfEB1DF",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid basic fee grant without spend-limit and expiration",
			append(
				[]string{
					granter.String(),
					"0x5d2596027D155ef78046f06e4C7cd02A1CCE8D0d",
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"try to add existed grant",
			append(
				[]string{
					granter.String(),
					alreadyExistedGrantee.String(),
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 18, &sdk.TxResponse{},
		},
		{
			"invalid number of args(periodic fee grant)",
			append(
				[]string{
					granter.String(),
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(tenHours)),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"period mentioned and period limit omitted, invalid periodic grant",
			append(
				[]string{
					granter.String(),
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, tenHours),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(oneHour)),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"period cannot be greater than the actual expiration(periodic fee grant)",
			append(
				[]string{
					granter.String(),
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, tenHours),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(oneHour)),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"valid periodic fee grant",
			append(
				[]string{
					granter.String(),
					"0x911c52f7644d88e99a11E65bDe12a364400B3B0b",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, oneHour),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(tenHours)),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid periodic fee grant without spend-limit",
			append(
				[]string{
					granter.String(),
					"0xfA4D171191cAcc5e12c3eA41ac8445e225C750F8",
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, oneHour),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(tenHours)),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid periodic fee grant without expiration",
			append(
				[]string{
					granter.String(),
					"0x170A9D72eA1aB2F04AcDDd6CddC3A87c45ACee90",
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, oneHour),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"valid periodic fee grant without spend-limit and expiration",
			append(
				[]string{
					granter.String(),
					"0xA442dDB26b1e91E43E6c1966a883Ae4C54589376",
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, oneHour),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
		{
			"invalid expiration",
			append(
				[]string{
					granter.String(),
					"0xfA4D171191cAcc5e12c3eA41ac8445e225C750F8",
					fmt.Sprintf("--%s=%d", cli.FlagPeriod, oneHour),
					fmt.Sprintf("--%s=%s", cli.FlagPeriodLimit, "10stake"),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
					fmt.Sprintf("--%s=%s", cli.FlagExpiration, "invalid"),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCmdFeeGrant()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
			}
		})
	}
}

func (s *CLITestSuite) TestNewCmdRevokeFeegrant() {
	granter := s.addedGranter
	grantee := s.addedGrantee
	clientCtx := s.clientCtx

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
	}

	// Create new fee grant specifically to test amino.
	encodedGrantee := "0x0EA9D661C1C2B5C7660146F1193700A527D81EFb"
	aminoGrantee, err := sdk.AccAddressFromHexUnsafe(encodedGrantee)
	s.Require().NoError(err)
	s.createGrant(granter, aminoGrantee)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"invalid grantee",
			append(
				[]string{
					"wrong_granter",
					grantee.String(),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"invalid grantee",
			append(
				[]string{
					granter.String(),
					"wrong_grantee",
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, 0, nil,
		},
		{
			"Valid revoke",
			append(
				[]string{
					granter.String(),
					grantee.String(),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, 0, &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCmdRevokeFeegrant()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
			}
		})
	}
}

func (s *CLITestSuite) TestTxWithFeeGrant() {
	// s.T().Skip() // TODO to re-enable in #12274

	clientCtx := s.clientCtx
	granter := s.addedGranter

	// creating an account manually (This account won't be exist in state)
	k, _, err := s.baseCtx.Keyring.NewMnemonic("grantee", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	pub, err := k.GetPubKey()
	s.Require().NoError(err)
	grantee := sdk.AccAddress(pub.Address())

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))).String()),
	}

	fee := sdk.NewCoin("stake", sdk.NewInt(100))

	args := append(
		[]string{
			granter.String(),
			grantee.String(),
			fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, fee.String()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
			fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(oneYear)),
		},
		commonFlags...,
	)

	cmd := cli.NewCmdFeeGrant()

	var res sdk.TxResponse
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())

	testcases := []struct {
		name       string
		from       string
		flags      []string
		expErrCode uint32
	}{
		{
			name:  "granted fee allowance for an account which is not in state and creating any tx with it by using --fee-granter shouldn't fail",
			from:  grantee.String(),
			flags: []string{fmt.Sprintf("--%s=%s", flags.FlagFeeGranter, granter.String())},
		},
		{
			name: "use --fee-payer and --fee-granter together works",
			from: grantee.String(),
			flags: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFeePayer, grantee.String()),
				fmt.Sprintf("--%s=%s", flags.FlagFeeGranter, granter.String()),
			},
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			err := s.msgSubmitLegacyProposal(s.baseCtx, tc.from,
				"Text Proposal", "No desc", govv1beta1.ProposalTypeText,
				tc.flags...,
			)
			s.Require().NoError(err)

			var resp sdk.TxResponse
			s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
		})
	}
}

func (s *CLITestSuite) msgSubmitLegacyProposal(clientCtx client.Context, from, title, description, proposalType string, extraArgs ...string) error {
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
	}

	args := append([]string{
		fmt.Sprintf("--%s=%s", govcli.FlagTitle, title),
		fmt.Sprintf("--%s=%s", govcli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", govcli.FlagProposalType, proposalType),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)

	cmd := govcli.NewCmdSubmitLegacyProposal()

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	s.Require().NoError(err)
	var resp sdk.TxResponse
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())

	return err
}

func (s *CLITestSuite) TestFilteredFeeAllowance() {
	granter := s.addedGranter
	k, _, err := s.baseCtx.Keyring.NewMnemonic("grantee1", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)
	pub, err := k.GetPubKey()
	s.Require().NoError(err)
	grantee := sdk.AccAddress(pub.Address())

	clientCtx := s.clientCtx

	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))).String()),
	}
	spendLimit := sdk.NewCoin("stake", sdk.NewInt(1000))

	allowMsgs := strings.Join([]string{sdk.MsgTypeURL(&govv1.MsgSubmitProposal{}), sdk.MsgTypeURL(&govv1.MsgVoteWeighted{})}, ",")

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"invalid granter address",
			append(
				[]string{
					"not an address",
					"0x942742fA24F6000aB091d42e744D66F4bBd8cAc8",
					fmt.Sprintf("--%s=%s", cli.FlagAllowedMsgs, allowMsgs),
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, spendLimit.String()),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, &sdk.TxResponse{}, 0,
		},
		{
			"invalid grantee address",
			append(
				[]string{
					granter.String(),
					"not an address",
					fmt.Sprintf("--%s=%s", cli.FlagAllowedMsgs, allowMsgs),
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, spendLimit.String()),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			true, &sdk.TxResponse{}, 0,
		},
		{
			"valid filter fee grant",
			append(
				[]string{
					granter.String(),
					grantee.String(),
					fmt.Sprintf("--%s=%s", cli.FlagAllowedMsgs, allowMsgs),
					fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, spendLimit.String()),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
				},
				commonFlags...,
			),
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCmdFeeGrant()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())
			}
		})
	}

	// exec filtered fee allowance
	cases := []struct {
		name         string
		malleate     func() error
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid proposal tx",
			func() error {
				return s.msgSubmitLegacyProposal(s.baseCtx, grantee.String(),
					"Text Proposal", "No desc", govv1beta1.ProposalTypeText,
					fmt.Sprintf("--%s=%s", flags.FlagFeeGranter, granter.String()),
					fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))).String()),
				)
			},
			&sdk.TxResponse{},
			0,
		},
		{
			"valid weighted_vote tx",
			func() error {
				return s.msgVote(s.baseCtx, grantee.String(), "0", "yes",
					fmt.Sprintf("--%s=%s", flags.FlagFeeGranter, granter.String()),
					fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))).String()),
				)
			},
			&sdk.TxResponse{},
			2,
		},
		{
			"should fail with unauthorized msgs",
			func() error {
				args := append(
					[]string{
						grantee.String(),
						"0x170A9D72eA1aB2F04AcDDd6CddC3A87c45ACee90",
						fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, "100stake"),
						fmt.Sprintf("--%s=%s", flags.FlagFeeGranter, granter),
					},
					commonFlags...,
				)

				cmd := cli.NewCmdFeeGrant()
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &sdk.TxResponse{}), out.String())

				return err
			},
			&sdk.TxResponse{},
			7,
		},
	}

	for _, tc := range cases {
		tc := tc

		s.Run(tc.name, func() {
			err := tc.malleate()
			s.Require().NoError(err)
		})
	}
}

// msgVote votes for a proposal
func (s *CLITestSuite) msgVote(clientCtx client.Context, from, id, vote string, extraArgs ...string) error {
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
	}
	args := append([]string{
		id,
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)
	cmd := govcli.NewCmdWeightedVote()

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)

	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &sdk.TxResponse{}), out.String())

	return err
}

func getFormattedExpiration(duration int64) string {
	return time.Now().Add(time.Duration(duration) * time.Second).Format(time.RFC3339)
}
