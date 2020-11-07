# Solana library for Go

Go library to interface with Solana nodes's JSON-RPC interface, Solana's SPL tokens and the
[https://dex.projectserum.com](Serum DEX) instructions.  More contracts to come.


# Command-line
```
$ slnc get balance EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v
1461600 lamports

$ slnc get account EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v
{
  "lamports": 1461600,
  "data": [
    "AQAAABzjWe1aAS4E+hQrnHUaHF6Hz9CgFhuchf/TG3jN/Nj2gCa3xLwWAAAGAQEAAAAqnl7btTwEZ5CY/3sSZRcUQ0/AjFYqmjuGEQXmctQicw==",
    "base64"
  ],
  "owner": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
  "executable": false,
  "rentEpoch": 108
}

$ slnc spl get-mint  SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt

Mint Authority Option:  0
Mint Authority:  73uWQpzn1AmUZkZ7MhafSSwQJmNmQ3fN4guANBrXg8uD
Supply:  9999709435300000
Decimals:  6
Is Initialized:  true
Freeze Authority Option:  0
Freeze Authority:  73uWQpzn1AmUZkZ7MhafSSwQJmNmQ3fN4guANBrXg8uD

$ slnc serum markets
...
SUSHI/USDC -> 7LVJtqSrF6RudMaz5rKGTmR3F3V5TKoDcN6bnk68biYZ
SXP/USDC -> 13vjJ8pxDMmzen26bQ5UrouX8dkXYPW1p3VLVDjxXrKR
MSRM/USDC -> AwvPwwSprfDZ86beBJDNH5vocFvuw4ZbVQ6upJDbSCXZ
FTT/USDC -> FfDb3QZUdMW2R2aqJQgzeieys4ETb3rPrFFfPSemzq7R
YFI/USDC -> 4QL5AQvXdMSCVZmnKXiuMMU83Kq3LCwVfU8CyznqZELG
LINK/USDC -> 7JCG9TsCx3AErSV3pvhxiW4AbkKRcJ6ZAveRmJwrgQ16
HGET/USDC -> 3otQFkeQ7GNUKT3i2p3aGTQKS2SAw6NLYPE5qxh3PoqZ
CREAM/USDC -> 2M8EBxFbLANnCoHydypL1jupnRHG782RofnvkatuKyLL
...

$ slnc serum market 7JCG9TsCx3AErSV3pvhxiW4AbkKRcJ6ZAveRmJwrgQ16


```

# Library usage
Loading a Serum market

```golang

import "github.com/dfuse-io/solana-go/rpc"

addr := solana.MustPublicKeyFromBase58("7JCG9TsCx3AErSV3pvhxiW4AbkKRcJ6ZAveRmJwrgQ16")
cli := rpc.NewClient("http://api.mainnet-beta.solana.com/rpc")
acct, err := cli.GetAccountInfo(context.Background(), addr)
// handle `err`

var m serum.MarketV2
err = struc.Unpack(bytes.NewReader(acct.Value.MustDataToBytes()), &m)
// handle `err`

json.NewEncoder(os.Stdout).Encode(m)
// {
//   "AccountFlags": 3,
//   "OwnAddress": "7JCG9TsCx3AErSV3pvhxiW4AbkKRcJ6ZAveRmJwrgQ16",
//   "VaultSignerNonce": 1,
//   "BaseMint": "CWE8jPTUYhdCTZYWPTe1o5DFqfdjzWKc9WKz6rSjQUdG",
//   "QuoteMint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
//   "BaseVault": "BWxmscpG77m1bWH7thDCisR84JJ1e3ho6rMm4udEq2u7",
//   "BaseDepositsTotal": "30750820000",
//   "BaseFeesAccrued": 0,
//   "QuoteVault": "5tFhdTCzTYMvfVTZnczZEL36YjFnkDTSaoQ7XAZvS7LR",
//   "QuoteDepositsTotal": "269134109688",
//   "QuoteFeesAccrued": 476785321,
//   "QuoteDustThreshold": 100,
//   "RequestQueue": "3SKQw5A69B72SEj3EWBe3DymWS9DPJ2mTUzJtV8i8DQ3",
//   "EventQueue": "DNh4agYmXGYv6po3tcCjM5yFQxZHrWnrF3dKhBXY24ZR",
//   "Bids": "627n4UWdaYxpiHrSd3RJrnhkq1zHSNuAgSg88uZXn1KZ",
//   "Asks": "AjgrF2tKVU96fwygfMdsAdBGa5S1Mttkh7XL9P24QWx7",
//   "BaseLotSize": 10000,
//   "QuoteLotSize": 10,
//   "FeeRateBPS": 0,
//   "ReferrerRebatesAccrued": 481639
// }

```
# Examples




# Contributing

Any contributions are welcome, use your standard GitHub-fu to pitch in and improve.


License
-------

Apache-2
