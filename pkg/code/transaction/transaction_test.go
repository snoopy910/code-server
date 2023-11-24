package transaction

import (
	"testing"

	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/code-payments/code-server/pkg/solana"
	"github.com/code-payments/code-server/pkg/solana/system"
	"github.com/code-payments/code-server/pkg/solana/token"
	"github.com/code-payments/code-server/pkg/testutil"
	"github.com/code-payments/code-server/pkg/code/common"
	code_data "github.com/code-payments/code-server/pkg/code/data"
)

func TestTransaction_MakeNoncedTransaction_HappyPath(t *testing.T) {
	subsidizer := testutil.SetupRandomSubsidizer(t, code_data.NewTestDataProvider())

	nonceAccount, err := common.NewAccountFromPublicKeyString("non9MZDuwcTzNYfWFu18XT4MLi3Pf6vscuuMuKTbrTx")
	require.NoError(t, err)

	untypedBlockhash, err := base58.Decode("9eRZTogvYM4WC8PRrw27fpzcZTvEvQuaREQyRETyw46d")
	require.NoError(t, err)
	var typedBlockhash solana.Blockhash
	copy(typedBlockhash[:], untypedBlockhash)

	ixns := []solana.Instruction{
		token.Transfer(
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			1,
		), token.Transfer(
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			2,
		),
		token.Transfer(
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			testutil.NewRandomAccount(t).PublicKey().ToBytes(),
			3,
		),
	}

	txn, err := MakeNoncedTransaction(nonceAccount, typedBlockhash, ixns...)
	require.NoError(t, err)

	assert.Equal(t, typedBlockhash, txn.Message.RecentBlockhash)
	assert.EqualValues(t, txn.Message.Accounts[0], subsidizer.PublicKey().ToBytes())

	require.Len(t, txn.Message.Instructions, 4)

	actual, err := system.DecompileAdvanceNonce(txn.Message, 0)
	require.NoError(t, err)
	assert.EqualValues(t, nonceAccount.PublicKey().ToBytes(), actual.Nonce)
	assert.EqualValues(t, subsidizer.PublicKey().ToBytes(), actual.Authority)

	for i := range ixns {
		actual, err := token.DecompileTransfer(txn.Message, i+1)
		require.NoError(t, err)
		assert.EqualValues(t, i+1, actual.Amount)
	}
}

func TestTransaction_MakeNoncedTransaction_NoInstructions(t *testing.T) {
	nonceAccount, err := common.NewAccountFromPublicKeyString("non9MZDuwcTzNYfWFu18XT4MLi3Pf6vscuuMuKTbrTx")
	require.NoError(t, err)

	untypedBlockhash, err := base58.Decode("9eRZTogvYM4WC8PRrw27fpzcZTvEvQuaREQyRETyw46d")
	require.NoError(t, err)
	var typedBlockhash solana.Blockhash
	copy(typedBlockhash[:], untypedBlockhash)

	_, err = MakeNoncedTransaction(nonceAccount, typedBlockhash)
	assert.Error(t, err)
}
