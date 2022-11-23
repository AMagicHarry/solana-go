// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sendandconfirmtransaction

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// Send and wait for confirmation of a transaction.
func SendAndConfirmTransaction(
	ctx context.Context,
	rpcClient *rpc.Client,
	wsClient *ws.Client,
	transaction *solana.Transaction,
) (signature solana.Signature, err error) {
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	}

	return SendAndConfirmTransactionWithOpts(
		ctx,
		rpcClient,
		wsClient,
		transaction,
		opts,
		nil,
	)
}

func SendAndConfirmTransactionWithTimeout(
	ctx context.Context,
	rpcClient *rpc.Client,
	wsClient *ws.Client,
	transaction *solana.Transaction,
	timeout time.Duration,
) (signature solana.Signature, err error) {
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	}

	return SendAndConfirmTransactionWithOpts(
		ctx,
		rpcClient,
		wsClient,
		transaction,
		opts,
		&timeout,
	)
}

var ErrTimeout = fmt.Errorf("timeout")

// Send and wait for confirmation of a transaction.
func SendAndConfirmTransactionWithOpts(
	ctx context.Context,
	rpcClient *rpc.Client,
	wsClient *ws.Client,
	transaction *solana.Transaction,
	opts rpc.TransactionOpts,
	timeout *time.Duration,
) (sig solana.Signature, err error) {
	sig, err = rpcClient.SendTransactionWithOpts(
		ctx,
		transaction,
		opts,
	)
	if err != nil {
		return sig, err
	}
	_, err = WaitForConfirmation(
		ctx,
		wsClient,
		sig,
		timeout,
	)
	return sig, err
}

// WaitForConfirmation waits for a transaction to be confirmed.
// If the transaction was confirmed, but it failed while executing (one of the instructions failed),
// then this function will return an error (true, error).
// If the transaction was confirmed, and it succeeded, then this function will return nil (true, nil).
func WaitForConfirmation(
	ctx context.Context,
	wsClient *ws.Client,
	sig solana.Signature,
	timeout *time.Duration,
) (confirmed bool, err error) {
	sub, err := wsClient.SignatureSubscribe(
		sig,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return false, err
	}
	defer sub.Unsubscribe()

	if timeout == nil {
		t := 2 * time.Minute // random default timeout
		timeout = &t
	}

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(*timeout):
			return false, ErrTimeout
		case resp, ok := <-sub.Response():
			if !ok {
				return false, fmt.Errorf("subscription closed")
			}
			if resp.Value.Err != nil {
				// The transaction was confirmed, but it failed while executing (one of the instructions failed).
				return true, fmt.Errorf("confirmed transaction with execution error: %v", resp.Value.Err)
			} else {
				// Success! Confirmed! And there was no error while executing the transaction.
				return true, nil
			}
		case err := <-sub.Err():
			return false, err
		}
	}
}
