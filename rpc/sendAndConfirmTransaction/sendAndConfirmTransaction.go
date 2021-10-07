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
	return SendAndConfirmTransactionWithOpts(
		ctx,
		rpcClient,
		wsClient,
		transaction,
		false,
		rpc.CommitmentFinalized,
	)
}

// Send and wait for confirmation of a transaction.
func SendAndConfirmTransactionWithOpts(
	ctx context.Context,
	rpcClient *rpc.Client,
	wsClient *ws.Client,
	transaction *solana.Transaction,
	skipPreflight bool, // if true, skip the preflight transaction checks (default: false)
	preflightCommitment rpc.CommitmentType, // optional; Commitment level to use for preflight (default: "finalized").
) (signature solana.Signature, err error) {

	sig, err := rpcClient.SendTransactionWithOpts(
		ctx,
		transaction,
		skipPreflight,
		preflightCommitment,
	)
	if err != nil {
		return sig, err
	}

	sub, err := wsClient.SignatureSubscribe(
		sig,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return sig, err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return sig, err
		}
		if got.Value.Err != nil {
			return sig, fmt.Errorf("transaction confirmation failed: %v", got.Value.Err)
		} else {
			return sig, nil
		}
	}
}
