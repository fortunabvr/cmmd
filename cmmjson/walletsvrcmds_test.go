// Copyright (c) 2018 The Commercium developers
// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
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

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
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
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cmmjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &cmmjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cmmjson.NewAddMultisigAddressCmd(2, keys, cmmjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &cmmjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   cmmjson.String("test"),
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return cmmjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &cmmjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &cmmjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &cmmjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &cmmjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &cmmjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetBalanceCmd{
				Account: nil,
				MinConf: cmmjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBalanceCmd(cmmjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.GetBalanceCmd{
				Account: cmmjson.String("acct"),
				MinConf: cmmjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetBalanceCmd(cmmjson.String("acct"), cmmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &cmmjson.GetBalanceCmd{
				Account: cmmjson.String("acct"),
				MinConf: cmmjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNewAddressCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetNewAddressCmd{
				Account:   nil,
				GapPolicy: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getnewaddress", "acct", "ignore")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetNewAddressCmd(cmmjson.String("acct"), cmmjson.String("ignore"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct","ignore"],"id":1}`,
			unmarshalled: &cmmjson.GetNewAddressCmd{
				Account:   cmmjson.String("acct"),
				GapPolicy: cmmjson.String("ignore"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &cmmjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetRawChangeAddressCmd(cmmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.GetRawChangeAddressCmd{
				Account: cmmjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: cmmjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetReceivedByAccountCmd("acct", cmmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &cmmjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: cmmjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &cmmjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: cmmjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetReceivedByAddressCmd("1Address", cmmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &cmmjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: cmmjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewGetTransactionCmd("123", cmmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &cmmjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: cmmjson.Bool(true),
			},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewImportPrivKeyCmd("abc", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &cmmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  cmmjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewImportPrivKeyCmd("abc", cmmjson.String("label"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &cmmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   cmmjson.String("label"),
				Rescan:  cmmjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewImportPrivKeyCmd("abc", cmmjson.String("label"), cmmjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &cmmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   cmmjson.String("label"),
				Rescan:  cmmjson.Bool(false),
			},
		},
		{
			name: "importprivkey optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("importprivkey", "abc", "label", false, 12345)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewImportPrivKeyCmd("abc", cmmjson.String("label"), cmmjson.Bool(false), cmmjson.Int(12345))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false,12345],"id":1}`,
			unmarshalled: &cmmjson.ImportPrivKeyCmd{
				PrivKey:  "abc",
				Label:    cmmjson.String("label"),
				Rescan:   cmmjson.Bool(false),
				ScanFrom: cmmjson.Int(12345),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &cmmjson.KeyPoolRefillCmd{
				NewSize: cmmjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewKeyPoolRefillCmd(cmmjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &cmmjson.KeyPoolRefillCmd{
				NewSize: cmmjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListAccountsCmd{
				MinConf: cmmjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListAccountsCmd(cmmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &cmmjson.ListAccountsCmd{
				MinConf: cmmjson.Int(6),
			},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAccountCmd{
				MinConf:          cmmjson.Int(1),
				IncludeEmpty:     cmmjson.Bool(false),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAccountCmd(cmmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAccountCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(false),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAccountCmd(cmmjson.Int(6), cmmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAccountCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(true),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAccountCmd(cmmjson.Int(6), cmmjson.Bool(true), cmmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAccountCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(true),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAddressCmd{
				MinConf:          cmmjson.Int(1),
				IncludeEmpty:     cmmjson.Bool(false),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAddressCmd(cmmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAddressCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(false),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAddressCmd(cmmjson.Int(6), cmmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAddressCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(true),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListReceivedByAddressCmd(cmmjson.Int(6), cmmjson.Bool(true), cmmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &cmmjson.ListReceivedByAddressCmd{
				MinConf:          cmmjson.Int(6),
				IncludeEmpty:     cmmjson.Bool(true),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: cmmjson.Int(1),
				IncludeWatchOnly:    cmmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListSinceBlockCmd(cmmjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &cmmjson.ListSinceBlockCmd{
				BlockHash:           cmmjson.String("123"),
				TargetConfirmations: cmmjson.Int(1),
				IncludeWatchOnly:    cmmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListSinceBlockCmd(cmmjson.String("123"), cmmjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &cmmjson.ListSinceBlockCmd{
				BlockHash:           cmmjson.String("123"),
				TargetConfirmations: cmmjson.Int(6),
				IncludeWatchOnly:    cmmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListSinceBlockCmd(cmmjson.String("123"), cmmjson.Int(6), cmmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &cmmjson.ListSinceBlockCmd{
				BlockHash:           cmmjson.String("123"),
				TargetConfirmations: cmmjson.Int(6),
				IncludeWatchOnly:    cmmjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListTransactionsCmd{
				Account:          nil,
				Count:            cmmjson.Int(10),
				From:             cmmjson.Int(0),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListTransactionsCmd(cmmjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &cmmjson.ListTransactionsCmd{
				Account:          cmmjson.String("acct"),
				Count:            cmmjson.Int(10),
				From:             cmmjson.Int(0),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListTransactionsCmd(cmmjson.String("acct"), cmmjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &cmmjson.ListTransactionsCmd{
				Account:          cmmjson.String("acct"),
				Count:            cmmjson.Int(20),
				From:             cmmjson.Int(0),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListTransactionsCmd(cmmjson.String("acct"), cmmjson.Int(20),
					cmmjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &cmmjson.ListTransactionsCmd{
				Account:          cmmjson.String("acct"),
				Count:            cmmjson.Int(20),
				From:             cmmjson.Int(1),
				IncludeWatchOnly: cmmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListTransactionsCmd(cmmjson.String("acct"), cmmjson.Int(20),
					cmmjson.Int(1), cmmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &cmmjson.ListTransactionsCmd{
				Account:          cmmjson.String("acct"),
				Count:            cmmjson.Int(20),
				From:             cmmjson.Int(1),
				IncludeWatchOnly: cmmjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &cmmjson.ListUnspentCmd{
				MinConf:   cmmjson.Int(1),
				MaxConf:   cmmjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListUnspentCmd(cmmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &cmmjson.ListUnspentCmd{
				MinConf:   cmmjson.Int(6),
				MaxConf:   cmmjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListUnspentCmd(cmmjson.Int(6), cmmjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &cmmjson.ListUnspentCmd{
				MinConf:   cmmjson.Int(6),
				MaxConf:   cmmjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return cmmjson.NewListUnspentCmd(cmmjson.Int(6), cmmjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &cmmjson.ListUnspentCmd{
				MinConf:   cmmjson.Int(6),
				MaxConf:   cmmjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return cmmjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1,"tree":0}]],"id":1}`,
			unmarshalled: &cmmjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []cmmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &cmmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cmmjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendFromCmd("from", "1Address", 0.5, cmmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &cmmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cmmjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendFromCmd("from", "1Address", 0.5, cmmjson.Int(6),
					cmmjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &cmmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cmmjson.Int(6),
				Comment:     cmmjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendFromCmd("from", "1Address", 0.5, cmmjson.Int(6),
					cmmjson.String("comment"), cmmjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &cmmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     cmmjson.Int(6),
				Comment:     cmmjson.String("comment"),
				CommentTo:   cmmjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cmmjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &cmmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cmmjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cmmjson.NewSendManyCmd("from", amounts, cmmjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &cmmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cmmjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return cmmjson.NewSendManyCmd("from", amounts, cmmjson.Int(6), cmmjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &cmmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     cmmjson.Int(6),
				Comment:     cmmjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &cmmjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSendToAddressCmd("1Address", 0.5, cmmjson.String("comment"),
					cmmjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &cmmjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   cmmjson.String("comment"),
				CommentTo: cmmjson.String("commentto"),
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &cmmjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &cmmjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &cmmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    cmmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return cmmjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &cmmjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]cmmjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    cmmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.RawTxInput{}
				privKeys := []string{"abc"}
				return cmmjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &cmmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]cmmjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    cmmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []cmmjson.RawTxInput{}
				privKeys := []string{}
				return cmmjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					cmmjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &cmmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]cmmjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    cmmjson.String("ALL"),
			},
		},
		{
			name: "verifyseed",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifyseed", "abc")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewVerifySeedCmd("abc", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc"],"id":1}`,
			unmarshalled: &cmmjson.VerifySeedCmd{
				Seed:    "abc",
				Account: nil,
			},
		},
		{
			name: "verifyseed optional",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("verifyseed", "abc", 5)
			},
			staticCmd: func() interface{} {
				account := cmmjson.Uint32(5)
				return cmmjson.NewVerifySeedCmd("abc", account)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc",5],"id":1}`,
			unmarshalled: &cmmjson.VerifySeedCmd{
				Seed:    "abc",
				Account: cmmjson.Uint32(5),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &cmmjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return cmmjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &cmmjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return cmmjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return cmmjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &cmmjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
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
