package praefect

import (
	"context"
	"crypto/sha1"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/transactions"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func runPraefectServerAndTxMgr(t testing.TB) (*grpc.ClientConn, *transactions.Manager, testhelper.Cleanup) {
	conf := testConfig(1)
	txMgr := defaultTxMgr()
	cc, _, cleanup := runPraefectServer(t, conf, buildOptions{
		withTxMgr:   txMgr,
		withNodeMgr: nullNodeMgr{}, // to suppress node address issues
	})
	return cc, txMgr, cleanup
}

func TestTransactionSucceeds(t *testing.T) {
	cc, txMgr, cleanup := runPraefectServerAndTxMgr(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := gitalypb.NewRefTransactionClient(cc)

	transactionID, cancelTransaction, err := txMgr.RegisterTransaction([]string{"node1"})
	require.NoError(t, err)
	require.NotZero(t, transactionID)
	defer cancelTransaction()

	hash := sha1.Sum([]byte{})

	response, err := client.StartTransaction(ctx, &gitalypb.StartTransactionRequest{
		TransactionId:        transactionID,
		Node:                 "node1",
		ReferenceUpdatesHash: hash[:],
	})
	require.NoError(t, err)
	require.Equal(t, gitalypb.StartTransactionResponse_COMMIT, response.State)
}

func TestTransactionFailsWithMultipleNodes(t *testing.T) {
	_, txMgr, cleanup := runPraefectServerAndTxMgr(t)
	defer cleanup()

	_, _, err := txMgr.RegisterTransaction([]string{"node1", "node2"})
	require.Error(t, err)
}

func TestTransactionFailures(t *testing.T) {
	cc, _, cleanup := runPraefectServerAndTxMgr(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := gitalypb.NewRefTransactionClient(cc)

	hash := sha1.Sum([]byte{})
	_, err := client.StartTransaction(ctx, &gitalypb.StartTransactionRequest{
		TransactionId:        1,
		Node:                 "node1",
		ReferenceUpdatesHash: hash[:],
	})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestTransactionCancellation(t *testing.T) {
	cc, txMgr, cleanup := runPraefectServerAndTxMgr(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := gitalypb.NewRefTransactionClient(cc)

	transactionID, cancelTransaction, err := txMgr.RegisterTransaction([]string{"node1"})
	require.NoError(t, err)
	require.NotZero(t, transactionID)

	cancelTransaction()

	hash := sha1.Sum([]byte{})
	_, err = client.StartTransaction(ctx, &gitalypb.StartTransactionRequest{
		TransactionId:        transactionID,
		Node:                 "node1",
		ReferenceUpdatesHash: hash[:],
	})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}
