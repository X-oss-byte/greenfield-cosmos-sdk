package simulation_test

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/simulation"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// TestWeightedOperations tests the weights of the operations.
func (suite *SimTestSuite) TestWeightedOperations() {
	appParams := make(simtypes.AppParams)

	weightesOps := simulation.WeightedOperations(appParams, suite.cdc, suite.accountKeeper,
		suite.bankKeeper, suite.distrKeeper, suite.stakingKeeper)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accs := suite.getTestingAccounts(r, 3)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{simulation.DefaultWeightMsgSetWithdrawAddress, types.ModuleName, types.TypeMsgSetWithdrawAddress},
		{simulation.DefaultWeightMsgWithdrawDelegationReward, types.ModuleName, types.TypeMsgWithdrawDelegatorReward},
		{simulation.DefaultWeightMsgWithdrawValidatorCommission, types.ModuleName, types.TypeMsgWithdrawValidatorCommission},
		{simulation.DefaultWeightMsgFundCommunityPool, types.ModuleName, types.TypeMsgFundCommunityPool},
	}

	for i, w := range weightesOps {
		operationMsg, _, err := w.Op()(r, suite.app.BaseApp, suite.ctx, accs, sdktestutil.DefaultChainId)
		suite.Require().NoError(err)

		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		suite.Require().Equal(expected[i].weight, w.Weight(), "weight should be the same")
		suite.Require().Equal(expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		suite.Require().Equal(expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgSetWithdrawAddress tests the normal scenario of a valid message of type TypeMsgSetWithdrawAddress.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *SimTestSuite) TestSimulateMsgSetWithdrawAddress() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			ChainID: sdktestutil.DefaultChainId,
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		},
	})

	// execute operation
	op := simulation.SimulateMsgSetWithdrawAddress(suite.txConfig, suite.accountKeeper, suite.bankKeeper, suite.distrKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, sdktestutil.DefaultChainId)
	suite.Require().NoError(err)

	var msg types.MsgSetWithdrawAddress
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("0xd4BFb1CB895840ca474b0D15abb11Cf0f26bc88a", msg.DelegatorAddress)
	suite.Require().Equal("0x6b11EA2aF9b83C6E0BBCe6254d776F82BB6b6C13", msg.WithdrawAddress)
	suite.Require().Equal(types.TypeMsgSetWithdrawAddress, msg.Type())
	suite.Require().Equal(types.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// TestSimulateMsgWithdrawDelegatorReward tests the normal scenario of a valid message
// of type TypeMsgWithdrawDelegatorReward.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *SimTestSuite) TestSimulateMsgWithdrawDelegatorReward() {
	// setup 3 accounts
	s := rand.NewSource(4)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// setup accounts[0] as validator
	validator0 := suite.getTestingValidator0(accounts)

	// setup delegation
	delTokens := sdk.TokensFromConsensusPower(2, sdk.DefaultPowerReduction)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	delegator := accounts[1]
	delegation := stakingtypes.NewDelegation(delegator.Address, validator0.GetOperator(), issuedShares)
	suite.stakingKeeper.SetDelegation(suite.ctx, delegation)
	suite.distrKeeper.SetDelegatorStartingInfo(suite.ctx, validator0.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, math.LegacyOneDec(), 200))

	suite.setupValidatorRewards(validator0.GetOperator())

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			ChainID: sdktestutil.DefaultChainId,
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		},
	})

	// execute operation
	op := simulation.SimulateMsgWithdrawDelegatorReward(suite.txConfig, suite.accountKeeper, suite.bankKeeper, suite.distrKeeper, suite.stakingKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, sdktestutil.DefaultChainId)
	suite.Require().NoError(err)

	var msg types.MsgWithdrawDelegatorReward
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("0x0a12a929C9711b1f71d7C2D67Ff5b81B40a7B916", msg.ValidatorAddress)
	suite.Require().Equal("0x5A698439CCb612132474445A62720649A2FB00E7", msg.DelegatorAddress)
	suite.Require().Equal(types.TypeMsgWithdrawDelegatorReward, msg.Type())
	suite.Require().Equal(types.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// TestSimulateMsgWithdrawValidatorCommission tests the normal scenario of a valid message
// of type TypeMsgWithdrawValidatorCommission.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *SimTestSuite) TestSimulateMsgWithdrawValidatorCommission() {
	suite.testSimulateMsgWithdrawValidatorCommission("atoken")
	suite.testSimulateMsgWithdrawValidatorCommission("tokenxxx")
}

// all the checks in this function should not fail if we change the tokenName
func (suite *SimTestSuite) testSimulateMsgWithdrawValidatorCommission(tokenName string) {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// setup accounts[0] as validator
	validator0 := suite.getTestingValidator0(accounts)

	// set module account coins
	distrAcc := suite.distrKeeper.GetDistributionAccount(suite.ctx)
	suite.Require().NoError(banktestutil.FundModuleAccount(suite.bankKeeper, suite.ctx, distrAcc.GetName(), sdk.NewCoins(
		sdk.NewCoin(tokenName, sdk.NewInt(10)),
		sdk.NewCoin("stake", sdk.NewInt(5)),
	)))
	suite.accountKeeper.SetModuleAccount(suite.ctx, distrAcc)

	// set outstanding rewards
	valCommission := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(tokenName, math.LegacyNewDec(5).Quo(math.LegacyNewDec(2))),
		sdk.NewDecCoinFromDec("stake", math.LegacyNewDec(1).Quo(math.LegacyNewDec(1))),
	)

	suite.distrKeeper.SetValidatorOutstandingRewards(suite.ctx, validator0.GetOperator(), types.ValidatorOutstandingRewards{Rewards: valCommission})
	suite.distrKeeper.SetValidatorOutstandingRewards(suite.ctx, suite.genesisVals[0].GetOperator(), types.ValidatorOutstandingRewards{Rewards: valCommission})

	// setup validator accumulated commission
	suite.distrKeeper.SetValidatorAccumulatedCommission(suite.ctx, validator0.GetOperator(), types.ValidatorAccumulatedCommission{Commission: valCommission})
	suite.distrKeeper.SetValidatorAccumulatedCommission(suite.ctx, suite.genesisVals[0].GetOperator(), types.ValidatorAccumulatedCommission{Commission: valCommission})

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			ChainID: sdktestutil.DefaultChainId,
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		},
	})

	// execute operation
	op := simulation.SimulateMsgWithdrawValidatorCommission(suite.txConfig, suite.accountKeeper, suite.bankKeeper, suite.distrKeeper, suite.stakingKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, sdktestutil.DefaultChainId)
	if !operationMsg.OK {
		suite.Require().Equal("could not find account", operationMsg.Comment)
	} else {
		suite.Require().NoError(err)

		var msg types.MsgWithdrawValidatorCommission
		types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

		suite.Require().True(operationMsg.OK)
		suite.Require().Equal("0x520ecc4903A9F355246c1FF384E694b6dFFcE2Ec", msg.ValidatorAddress)
		suite.Require().Equal(types.TypeMsgWithdrawValidatorCommission, msg.Type())
		suite.Require().Equal(types.ModuleName, msg.Route())
		suite.Require().Len(futureOperations, 0)
	}
}

// TestSimulateMsgFundCommunityPool tests the normal scenario of a valid message of type TypeMsgFundCommunityPool.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *SimTestSuite) TestSimulateMsgFundCommunityPool() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			ChainID: sdktestutil.DefaultChainId,
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		},
	})

	// execute operation
	op := simulation.SimulateMsgFundCommunityPool(suite.txConfig, suite.accountKeeper, suite.bankKeeper, suite.distrKeeper, suite.stakingKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, sdktestutil.DefaultChainId)
	suite.Require().NoError(err)

	var msg types.MsgFundCommunityPool
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("4896096stake", msg.Amount.String())
	suite.Require().Equal("0xd4BFb1CB895840ca474b0D15abb11Cf0f26bc88a", msg.Depositor)
	suite.Require().Equal(types.TypeMsgFundCommunityPool, msg.Type())
	suite.Require().Equal(types.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

type SimTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *runtime.App
	genesisVals []stakingtypes.Validator

	txConfig      client.TxConfig
	cdc           codec.Codec
	stakingKeeper *stakingkeeper.Keeper
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	distrKeeper   keeper.Keeper
}

func (suite *SimTestSuite) SetupTest() {
	var (
		appBuilder *runtime.AppBuilder
		err        error
	)
	suite.app, err = simtestutil.Setup(configurator.NewAppConfig(
		configurator.AuthModule(),
		configurator.AuthzModule(),
		configurator.ParamsModule(),
		configurator.BankModule(),
		configurator.StakingModule(),
		configurator.TxModule(),
		configurator.ConsensusModule(),
		configurator.DistributionModule(),
	), &suite.accountKeeper,
		&suite.bankKeeper,
		&suite.cdc,
		&appBuilder,
		&suite.stakingKeeper,
		&suite.distrKeeper,
		&suite.txConfig,
	)

	suite.NoError(err)

	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{ChainID: sdktestutil.DefaultChainId})

	genesisVals := suite.stakingKeeper.GetAllValidators(suite.ctx)
	suite.Require().Len(genesisVals, 1)
	suite.genesisVals = genesisVals
}

func (suite *SimTestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := suite.stakingKeeper.TokensFromConsensusPower(suite.ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := suite.accountKeeper.NewAccountWithAddress(suite.ctx, account.Address)
		suite.accountKeeper.SetAccount(suite.ctx, acc)
		suite.Require().NoError(banktestutil.FundAccount(suite.bankKeeper, suite.ctx, account.Address, initCoins))
	}

	return accounts
}

func (suite *SimTestSuite) getTestingValidator0(accounts []simtypes.Account) stakingtypes.Validator {
	commission0 := stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyOneDec(), math.LegacyOneDec())
	return suite.getTestingValidator(accounts, commission0, 0)
}

func (suite *SimTestSuite) getTestingValidator(accounts []simtypes.Account, commission stakingtypes.Commission, n int) stakingtypes.Validator {
	require := suite.Require()
	account := accounts[n]
	valPubKey := account.PubKey
	valAddr := sdk.AccAddress(account.PubKey.Address().Bytes())
	validator, err := stakingtypes.NewSimpleValidator(valAddr, valPubKey, stakingtypes.Description{})
	require.NoError(err)
	validator, err = validator.SetInitialCommission(commission)
	require.NoError(err)
	validator.DelegatorShares = math.LegacyNewDec(100)
	validator.Tokens = sdk.NewInt(1000000)

	suite.stakingKeeper.SetValidator(suite.ctx, validator)

	return validator
}

func (suite *SimTestSuite) setupValidatorRewards(valAddress sdk.AccAddress) {
	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, math.LegacyOneDec())}
	historicalRewards := distrtypes.NewValidatorHistoricalRewards(decCoins, 2)
	suite.distrKeeper.SetValidatorHistoricalRewards(suite.ctx, valAddress, 2, historicalRewards)
	// setup current rewards
	currentRewards := distrtypes.NewValidatorCurrentRewards(decCoins, 3)
	suite.distrKeeper.SetValidatorCurrentRewards(suite.ctx, valAddress, currentRewards)
}

func TestSimTestSuite(t *testing.T) {
	suite.Run(t, new(SimTestSuite))
}
