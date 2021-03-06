// Copyright (c) 2018 The Commercium developers
// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package cmmjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/CommerciumBlockchain/cmmd/cmmjson"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("addnode", "127.0.0.1", cmmjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewAddNodeCmd("127.0.0.1", cmmjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &cmmjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: cmmjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return cmmjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123}],"id":1}`,
			unmarshalled: &cmmjson.CreateRawTransactionCmd{
				Inputs:  []cmmjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1,"tree":0}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return cmmjson.NewCreateRawTransactionCmd(txInputs, amounts, cmmjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &cmmjson.CreateRawTransactionCmd{
				Inputs:   []cmmjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: cmmjson.Int64(12312333333),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &cmmjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &cmmjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetAddedNodeInfoCmd(true, cmmjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &cmmjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: cmmjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   cmmjson.Bool(true),
				VerboseTx: cmmjson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := cmmjson.Bool(true)
				return cmmjson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockCmd("123", cmmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &cmmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   cmmjson.Bool(true),
				VerboseTx: cmmjson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockCmd("123", cmmjson.Bool(true), cmmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &cmmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   cmmjson.Bool(true),
				VerboseTx: cmmjson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &cmmjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: cmmjson.Bool(true),
			},
		},
		{
			name: "getblocksubsidy",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblocksubsidy", 123, 256)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockSubsidyCmd(123, 256)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocksubsidy","params":[123,256],"id":1}`,
			unmarshalled: &cmmjson.GetBlockSubsidyCmd{
				Height: 123,
				Voters: 256,
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return cmmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &cmmjson.GetBlockTemplateCmd{
				Request: &cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return cmmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &cmmjson.GetBlockTemplateCmd{
				Request: &cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return cmmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &cmmjson.GetBlockTemplateCmd{
				Request: &cmmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getcfilter", "123", "extended")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetCFilterCmd("123", "extended")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123","extended"],"id":1}`,
			unmarshalled: &cmmjson.GetCFilterCmd{
				Hash:       "123",
				FilterType: "extended",
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getcfilterheader", "123", "extended")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetCFilterHeaderCmd("123", "extended")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123","extended"],"id":1}`,
			unmarshalled: &cmmjson.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: "extended",
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetInfoCmd{},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetNetworkHashPSCmd{
				Blocks: cmmjson.Int(120),
				Height: cmmjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNetworkHashPSCmd(cmmjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &cmmjson.GetNetworkHashPSCmd{
				Blocks: cmmjson.Int(200),
				Height: cmmjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNetworkHashPSCmd(cmmjson.Int(200), cmmjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &cmmjson.GetNetworkHashPSCmd{
				Blocks: cmmjson.Int(200),
				Height: cmmjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawMempoolCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetRawMempoolCmd{
				Verbose: cmmjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawMempoolCmd(cmmjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &cmmjson.GetRawMempoolCmd{
				Verbose: cmmjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional 2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawmempool", false, "all")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawMempoolCmd(cmmjson.Bool(false), cmmjson.String("all"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false,"all"],"id":1}`,
			unmarshalled: &cmmjson.GetRawMempoolCmd{
				Verbose: cmmjson.Bool(false),
				TxType:  cmmjson.String("all"),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: cmmjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawTransactionCmd("123", cmmjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &cmmjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: cmmjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &cmmjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: cmmjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetTxOutCmd("123", 1, cmmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &cmmjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: cmmjson.Bool(true),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetWorkCmd(cmmjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &cmmjson.GetWorkCmd{
				Data: cmmjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &cmmjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewHelpCmd(cmmjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &cmmjson.HelpCmd{
				Command: cmmjson.String("getblock"),
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &cmmjson.PingCmd{},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(1),
				Skip:        cmmjson.Int(0),
				Count:       cmmjson.Int(100),
				VinExtra:    cmmjson.Int(0),
				Reverse:     cmmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(0),
				Count:       cmmjson.Int(100),
				VinExtra:    cmmjson.Int(0),
				Reverse:     cmmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), cmmjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(5),
				Count:       cmmjson.Int(100),
				VinExtra:    cmmjson.Int(0),
				Reverse:     cmmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), cmmjson.Int(5), cmmjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(5),
				Count:       cmmjson.Int(10),
				VinExtra:    cmmjson.Int(0),
				Reverse:     cmmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), cmmjson.Int(5), cmmjson.Int(10), cmmjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(5),
				Count:       cmmjson.Int(10),
				VinExtra:    cmmjson.Int(1),
				Reverse:     cmmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), cmmjson.Int(5), cmmjson.Int(10),
					cmmjson.Int(1), cmmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(5),
				Count:       cmmjson.Int(10),
				VinExtra:    cmmjson.Int(1),
				Reverse:     cmmjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSearchRawTransactionsCmd("1Address",
					cmmjson.Int(0), cmmjson.Int(5), cmmjson.Int(10),
					cmmjson.Int(1), cmmjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &cmmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     cmmjson.Int(0),
				Skip:        cmmjson.Int(5),
				Count:       cmmjson.Int(10),
				VinExtra:    cmmjson.Int(1),
				Reverse:     cmmjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &cmmjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: cmmjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendRawTransactionCmd("1122", cmmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &cmmjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: cmmjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSetGenerateCmd(true, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &cmmjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: cmmjson.Int(-1),
				MiningAddr:   nil,
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("setgenerate", true, 6, "22tv7nd31sMmD8BpcVRJAWQLqYCjaCuqpWpz")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSetGenerateCmd(true, cmmjson.Int(6), cmmjson.String("22tv7nd31sMmD8BpcVRJAWQLqYCjaCuqpWpz"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6,"22tv7nd31sMmD8BpcVRJAWQLqYCjaCuqpWpz"],"id":1}`,
			unmarshalled: &cmmjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: cmmjson.Int(6),
				MiningAddr:   cmmjson.String("22tv7nd31sMmD8BpcVRJAWQLqYCjaCuqpWpz"),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &cmmjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &cmmjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := cmmjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return cmmjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &cmmjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &cmmjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &cmmjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &cmmjson.VerifyChainCmd{
				CheckLevel: cmmjson.Int64(3),
				CheckDepth: cmmjson.Int64(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewVerifyChainCmd(cmmjson.Int64(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &cmmjson.VerifyChainCmd{
				CheckLevel: cmmjson.Int64(2),
				CheckDepth: cmmjson.Int64(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewVerifyChainCmd(cmmjson.Int64(2), cmmjson.Int64(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &cmmjson.VerifyChainCmd{
				CheckLevel: cmmjson.Int64(2),
				CheckDepth: cmmjson.Int64(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &cmmjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := cmmjson.MarshalCmd("1.0", testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = cmmjson.MarshalCmd("1.0", testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request cmmjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = cmmjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &cmmjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &cmmjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        cmmjson.Error{Code: cmmjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &cmmjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        cmmjson.Error{Code: cmmjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error type - got `%T` (%v), got `%T`",
				i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(cmmjson.Error); ok {
			gotErrorCode := err.(cmmjson.Error).Code
			if gotErrorCode != terr.Code {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.Code)
				continue
			}
		}
	}
}
