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

// package web3ext contains geth specific web3.js extensions.
package web3ext

var Modules = map[string]string{
	"admin":        Admin_JS,
	"chequebook":   Chequebook_JS,
	"clique":       Clique_JS,
	"posv":         Posv_JS,
	"debug":        Debug_JS,
	"eth":          Eth_JS,
	"miner":        Miner_JS,
	"net":          Net_JS,
	"personal":     Personal_JS,
	"rpc":          RPC_JS,
	"shh":          Shh_JS,
	"tomox":        TomoX_JS,
	"tomoxlending": TomoXLending_JS,
	"swarmfs":      SWARMFS_JS,
	"txpool":       TxPool_JS,
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

const Posv_JS = `
web3._extend({
	property: 'posv',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'posv_getSnapshot',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'posv_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getSigners',
			call: 'posv_getSigners',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSignersAtHash',
			call: 'posv_getSignersAtHash',
			params: 1
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'networkInformation',
			getter: 'posv_networkInformation'
		}),
	]
});
`

const Admin_JS = `
web3._extend({
	property: 'admin',
	methods: [
		new web3._extend.Method({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new web3._extend.Method({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new web3._extend.Method({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new web3._extend.Method({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Method({
			name: 'stopWS',
			call: 'admin_stopWS'
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new web3._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		}),
	]
});
`

const Debug_JS = `
web3._extend({
	property: 'debug',
	methods: [
		new web3._extend.Method({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new web3._extend.Method({
			name: 'seedHash',
			call: 'debug_seedHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1
		}),
		new web3._extend.Method({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new web3._extend.Method({
			name: 'metrics',
			call: 'debug_metrics',
			params: 1
		}),
		new web3._extend.Method({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new web3._extend.Method({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new web3._extend.Method({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new web3._extend.Method({
			name: 'freeOSMemory',
			call: 'debug_freeOSMemory',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'setGCPercent',
			call: 'debug_setGCPercent',
			params: 1,
		}),
		new web3._extend.Method({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new web3._extend.Method({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new web3._extend.Method({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new web3._extend.Method({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new web3._extend.Method({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'mutexProfile',
			call: 'debug_mutexProfile',
			params: 2
		}),
		new web3._extend.Method({
			name: 'setMutexProfileRate',
			call: 'debug_setMutexProfileRate',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMutexProfile',
			call: 'debug_writeMutexProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new web3._extend.Method({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockFromFile',
			call: 'debug_traceBlockFromFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Method({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getBadBlocks',
			call: 'debug_getBadBlocks',
			params: 0,
		}),
		new web3._extend.Method({
			name: 'storageRangeAt',
			call: 'debug_storageRangeAt',
			params: 5,
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByNumber',
			call: 'debug_getModifiedAccountsByNumber',
			params: 2,
			inputFormatter: [null, null],
		}),
		new web3._extend.Method({
			name: 'getModifiedAccountsByHash',
			call: 'debug_getModifiedAccountsByHash',
			params: 2,
			inputFormatter:[null, null],
		}),
	],
	properties: []
});
`

const Eth_JS = `
web3._extend({
	property: 'eth',
	methods: [
		new web3._extend.Method({
			name: 'chainId',
			call: 'eth_chainId',
			params: 0
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'eth_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'resend',
			call: 'eth_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'eth_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'submitTransaction',
			call: 'eth_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Method({
			name: 'getRawTransaction',
			call: 'eth_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getRewardByHash',
			call: 'eth_getRewardByHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'eth_getRawTransactionByBlockHashAndIndex' : 'eth_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		}),
		new web3._extend.Method({
			name: 'getOwnerByCoinbase',
			call: 'eth_getOwnerByCoinbase',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Method({
			name: 'getBlockReceipts',
			call: 'eth_getBlockReceipts',
			params: 1,
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'eth_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
	]
});
`

const Miner_JS = `
web3._extend({
	property: 'miner',
	methods: [
		new web3._extend.Method({
			name: 'start',
			call: 'miner_start',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'stop',
			call: 'miner_stop'
		}),
		new web3._extend.Method({
			name: 'setEtherbase',
			call: 'miner_setEtherbase',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter]
		}),
		new web3._extend.Method({
			name: 'setExtra',
			call: 'miner_setExtra',
			params: 1
		}),
		new web3._extend.Method({
			name: 'setGasPrice',
			call: 'miner_setGasPrice',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Method({
			name: 'getHashrate',
			call: 'miner_getHashrate'
		}),
	],
	properties: []
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

const Personal_JS = `
web3._extend({
	property: 'personal',
	methods: [
		new web3._extend.Method({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new web3._extend.Method({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Method({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		}),
		new web3._extend.Method({
			name: 'openWallet',
			call: 'personal_openWallet',
			params: 2
		}),
		new web3._extend.Method({
			name: 'deriveAccount',
			call: 'personal_deriveAccount',
			params: 3
		}),
		new web3._extend.Method({
			name: 'signTransaction',
			call: 'personal_signTransaction',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'listWallets',
			getter: 'personal_listWallets'
		}),
	]
})
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

const Shh_JS = `
web3._extend({
	property: 'shh',
	methods: [
	],
	properties:
	[
		new web3._extend.Property({
			name: 'version',
			getter: 'shh_version',
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Property({
			name: 'info',
			getter: 'shh_info'
		}),
	]
});
`

const TomoX_JS = `
web3._extend({
	property: 'tomox',
	methods: [
		new web3._extend.Method({
			name: 'version',
			call: 'tomox_version',
			params: 0,
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Method({
			name: 'info',
			call: 'tomox_info',
			params: 0
		}),
		new web3._extend.Method({
            name: 'getFeeByEpoch',
            call: 'tomox_getFeeByEpoch',
            params: 1,
            inputFormatter: [null, web3._extend.formatters.inputAddressFormatter]
        }),
		new web3._extend.Method({
            name: 'sendOrderRawTransaction',
            call: 'tomox_sendOrderRawTransaction',
            params: 1
		}),
		new web3._extend.Method({
            name: 'sendLendingRawTransaction',
            call: 'tomox_sendLendingRawTransaction',
            params: 1
		}),

		new web3._extend.Method({
            name: 'sendOrderTransaction',
            call: 'tomox_sendOrder',
            params: 1
		}),
		new web3._extend.Method({
            name: 'sendLendingTransaction',
            call: 'tomox_sendLending',
            params: 1
		}),
		new web3._extend.Method({
            name: 'getOrderTxMatchByHash',
            call: 'tomox_getOrderTxMatchByHash',
            params: 1
		}),
		new web3._extend.Method({
            name: 'getOrderPoolContent',
            call: 'tomox_getOrderPoolContent',
            params: 0
		}),
		new web3._extend.Method({
            name: 'getOrderStats',
            call: 'tomox_getOrderStats',
            params: 0
		}),
		new web3._extend.Method({
            name: 'getOrderCount',
            call: 'tomox_getOrderCount',
            params: 1
        }),
		new web3._extend.Method({
            name: 'getBestBid',
            call: 'tomox_getBestBid',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBestAsk',
            call: 'tomox_getBestAsk',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBidTree',
            call: 'tomox_getBidTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getAskTree',
            call: 'tomox_getAskTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getOrderById',
            call: 'tomox_getOrderById',
            params: 3
		}),
		new web3._extend.Method({
            name: 'getPrice',
            call: 'tomox_getPrice',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLastEpochPrice',
            call: 'tomox_getLastEpochPrice',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getCurrentEpochPrice',
            call: 'tomox_getCurrentEpochPrice',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getTradingOrderBookInfo',
            call: 'tomox_getTradingOrderBookInfo',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLiquidationPriceTree',
            call: 'tomox_getLiquidationPriceTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getInvestingTree',
            call: 'tomox_getInvestingTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBorrowingTree',
            call: 'tomox_getBorrowingTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLendingOrderBookInfo',
            call: 'tomox_getLendingOrderBookInfo',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLendingOrderTree',
            call: 'tomox_getLendingOrderTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLendingTradeTree',
            call: 'tomox_getLendingTradeTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLiquidationTimeTree',
            call: 'tomox_getLiquidationTimeTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLendingOrderCount',
            call: 'tomox_getLendingOrderCount',
            params: 1
        }),
		new web3._extend.Method({
            name: 'getBestInvesting',
            call: 'tomox_getBestInvesting',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBestBorrowing',
            call: 'tomox_getBestBorrowing',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBids',
            call: 'tomox_getBids',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getAsks',
            call: 'tomox_getAsks',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getInvests',
            call: 'tomox_getInvests',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBorrows',
            call: 'tomox_getBorrows',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getLendingTxMatchByHash',
            call: 'tomox_getLendingTxMatchByHash',
            params: 1
		}),
		new web3._extend.Method({
            name: 'getLiquidatedTradesByTxHash',
            call: 'tomox_getLiquidatedTradesByTxHash',
            params: 1
		}),
		new web3._extend.Method({
            name: 'getLendingOrderById',
            call: 'tomox_getLendingOrderById',
            params: 3
		}),
		new web3._extend.Method({
            name: 'getLendingTradeById',
            call: 'tomox_getLendingTradeById',
            params: 3
		}),
	]
});
`

const TomoXLending_JS = `
web3._extend({
	property: 'tomoxlending',
	methods: [
		new web3._extend.Method({
			name: 'version',
			call: 'tomoxlending_version',
			params: 0,
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Method({
			name: 'info',
			call: 'tomoxlending_info',
			params: 0
		}),
		new web3._extend.Method({
            name: 'createOrder',
            call: 'tomoxlending_createOrder',
            params: 1,
            inputFormatter: [null]
        }),
        new web3._extend.Method({
            name: 'cancelOrder',
            call: 'tomoxlending_cancelOrder',
            params: 1,
            inputFormatter: [null]
        }),
        new web3._extend.Method({
            name: 'getOrders',
            call: 'tomoxlending_getOrders',
            params: 1
        }),
		new web3._extend.Method({
            name: 'getOrderNonce',
            call: 'tomoxlending_getOrderNonce',
            params: 1,
            inputFormatter: [web3._extend.formatters.inputAddressFormatter]
		}),
		new web3._extend.Method({
            name: 'getFeeByEpoch',
            call: 'tomoxlending_GetFeeByEpoch',
            params: 1,
            inputFormatter: [null, web3._extend.formatters.inputAddressFormatter]
        }),
		new web3._extend.Method({
            name: 'getPendingOrders',
            call: 'tomoxlending_getPendingOrders',
            params: 1
        }),
		new web3._extend.Method({
            name: 'getAllPendingHashes',
            call: 'tomoxlending_getAllPendingHashes',
            params: 0
        }),
		new web3._extend.Method({
            name: 'sendOrderRawTransaction',
            call: 'tomoxlending_sendOrderRawTransaction',
            params: 1
        }),
		new web3._extend.Method({
            name: 'sendOrderTransaction',
            call: 'tomoxlending_sendOrder',
            params: 1
		}),
		new web3._extend.Method({
            name: 'getOrderCount',
            call: 'tomoxlending_getOrderCount',
            params: 1
        }),
		new web3._extend.Method({
            name: 'getBestBid',
            call: 'tomoxlending_getBestBid',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBestAsk',
            call: 'tomoxlending_getBestAsk',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getBidTree',
            call: 'tomoxlending_getBidTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getAskTree',
            call: 'tomoxlending_getAskTree',
            params: 2
		}),
		new web3._extend.Method({
            name: 'getOrderById',
            call: 'tomoxlending_getOrderById',
            params: 3
		}),
		new web3._extend.Method({
            name: 'getPrice',
            call: 'tomoxlending_getPrice',
            params: 2
		}),
	]
});
`

/*
   var sendOrderRawTransaction = new Method({
       name: 'sendOrderRawTransaction',
       call: 'eth_sendOrderRawTransaction',
       params: 1,
       inputFormatter: [null]
   });

   var sendOrderTransaction = new Method({
       name: 'sendOrder',
       call: 'tomox_sendOrder',
       params: 1,
       inputFormatter: [null]
   });
*/

const SWARMFS_JS = `
web3._extend({
	property: 'swarmfs',
	methods:
	[
		new web3._extend.Method({
			name: 'mount',
			call: 'swarmfs_mount',
			params: 2
		}),
		new web3._extend.Method({
			name: 'unmount',
			call: 'swarmfs_unmount',
			params: 1
		}),
		new web3._extend.Method({
			name: 'listmounts',
			call: 'swarmfs_listmounts',
			params: 0
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
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
	]
});
`
