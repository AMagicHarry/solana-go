package rpc

import (
	"context"
)

// GetBlockHeight returns the current block height of the node.
func (cl *Client) GetBlockHeight(
	ctx context.Context,
	commitment CommitmentType,
) (out uint64, err error) {
	params := []interface{}{}
	if commitment != "" {
		params = append(params, commitment)
	}
	err = cl.rpcClient.CallFor(&out, "getBlockHeight", params...)
	return
}