// Copyright 2018 PalletOne

// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// package web3ext contains gptn specific web3.js extensions.

package web3ext

var Modules = map[string]string{
	"admin":      Admin_JS,
	"chequebook": Chequebook_JS,
	"clique":     Clique_JS,
	"debug":      Debug_JS,
	"ptn":        Ptn_JS,
	"dag":        Dag_JS,
	"net":        Net_JS,
	"personal":   Personal_JS,
	"rpc":        RPC_JS,
	"wallet":     Wallet_JS,
	"txpool":     TxPool_JS,
	"mediator":   Mediator_JS,
	"contract":   Contract_JS,
}

const Chequebook_JS = `
web3._extend({
	property: 'chequebook',
	methods: [
		new web3._extend.Method({
			name: 'deposit',
			call: 'chequebook_deposit',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Property({
			name: 'balance',
			getter: 'chequebook_balance',
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Method({
			name: 'cash',
			call: 'chequebook_cash',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'issue',
			call: 'chequebook_issue',
			params: 2,
			inputFormatter: [null, null]
		}),
	]
});
`

const Clique_JS = `
web3._extend({
	property: 'clique',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'clique_getSnapshot',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'clique_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getSigners',
			call: 'clique_getSigners',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSignersAtHash',
			call: 'clique_getSignersAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'propose',
			call: 'clique_propose',
			params: 2
		}),
		new web3._extend.Method({
			name: 'discard',
			call: 'clique_discard',
			params: 1
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'proposals',
			getter: 'clique_proposals'
		}),
	]
});
`

const Net_JS = `
web3._extend({
	property: 'net',
	methods: [],
	properties: [
		new web3._extend.Property({
			name: 'version',
			getter: 'net_version'
		}),
	]
});
`

const Contract_JS = `
web3._extend({
	property: 'contract',
	methods: [
		new web3._extend.Method({
			name: 'ccinvoke',
			call: 'contract_ccinvoke',
			params: 3,
			inputFormatter: [null,null,null]
		}),
		new web3._extend.Method({
			name: 'ccinstalltx',
        	call: 'contract_ccinstalltx',
        	params: 8, //from, to , daoAmount, daoFee , tplName, path, version
			inputFormatter: [null, null, null,null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'ccdeploytx',
        	call: 'contract_ccdeploytx',
        	params: 6, //from, to , daoAmount, daoFee , templateId , args  
			inputFormatter: [null, null, null,null, null, null]
		}),
		new web3._extend.Method({
			name: 'ccinvoketx',
        	call: 'contract_ccinvoketx',
        	params: 7, //from, to, daoAmount, daoFee , contractAddr, args[]string------>["fun", "key", "value"], certid
			inputFormatter: [null, null, null,null, null, null, null]
		}),
        new web3._extend.Method({
			name: 'ccinvoketxPass',
			call: 'contract_ccinvoketxPass',
			params: 9, //from, to, daoAmount, daoFee , contractAddr, args[]string------>["fun", "key", "value"],passwd,duration, certid
			inputFormatter: [null, null, null,null, null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'ccinvokeToken',
        	call: 'contract_ccinvokeToken',
        	params: 9, //from, to, toToken, daoAmount, daoFee, daoAmountToken, assetToken, contractAddr, args[]string------>["fun", "key", "value"]
			inputFormatter: [null, null, null,null, null, null,null, null, null]
		}),
		new web3._extend.Method({
			name: 'ccquery',
			call: 'contract_ccquery',
			params: 2, //contractAddr,args[]string---->["func","arg1","arg2","..."]
			inputFormatter: [null,null]
		}),
		new web3._extend.Method({
			name: 'ccstoptx',
        	call: 'contract_ccstoptx',
        	params: 6, //from, to, daoAmount, daoFee, contractId, deleteImage
			inputFormatter: [null, null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'depositContractInvoke',
        	call: 'contract_depositContractInvoke',
        	params: 5, //from, to, daoAmount, daoFee,param[]string
			inputFormatter: [null, null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'depositContractQuery',
        	call: 'contract_depositContractQuery',
        	params: 1, //param[]string
			inputFormatter: [null]
		}),
	],
	properties: []
});
`
const RPC_JS = `
web3._extend({
	property: 'rpc',
	methods: [],
	properties: [
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const TxPool_JS = `
web3._extend({
	property: 'txpool',
	methods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new web3._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new web3._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = web3._extend.utils.toDecimal(status.pending);
      			status.orphans = web3._extend.utils.toDecimal(status.orphans);
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
		new web3._extend.Property({
			name: 'pending',
			getter: 'txpool_pending'
		}),
		new web3._extend.Property({
			name: 'queue',
			getter: 'txpool_queue'
		}),
	]
});
`
