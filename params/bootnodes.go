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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// Chancoin Bootnodes Mainnet
	"enode://cea67e9b8f9393c3df8a7843a2552872dd33a045c87b4c8ad2a0d98a4a64e817b9d23f12230fa24e14e0f856ad0052063dbb17d4faaedb45577070303e582ea8@boot.chancoin.org:30301",
	"enode://9de543a66de8503c682ab1f7f242854069c58fc4f9e0769020ad1ea317d0394f6876e3d5d1c4d1c19f8410cbb6fffb5cc8e6ba45bf42d0dc83c75cff18dd72fe@boot.chancoin.moe:30301",
	"enode://a1964d52bb7e7de39081ec6cd88a3628f12a50530cc6671d1c3cf68eeeb584f407a252dc04433e8f1d532238376fedb6d6d8ccea927c678532f206f50fc8a027@boot.signal2noi.se:30301",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	// Chancoin Bootnodes Testnet
	"enode://0bc411381dfe91955d10b8e3d6d210d3a5b056b1504de1a891e503bae6075364dd53130517dcf39071f6f404b7a3b2523c130309a25abd24ebc0cf65d2bfda61@testboot.chancoin.org:30303",
	"enode://b041725145911485b129682f45e451c9eb13a5e7195dd2bcdfd9040437425fa33d2a54500881ae3443a4092316ae574a0bdc859fc6a33455aaa973cba4cd616a@testboot.chancoin.moe:30303",
	"enode://cfa068d689f7e8d6fd9b71466c9b5ade3b61fb43e49b3d5aa9f7a9a52a0a6ad4ee36e0b2ba4522e6e753d43686befb40be0922ea36d781cdddec020152d9d3f5@testboot.signal2noi.se:30303",
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://06051a5573c81934c9554ef2898eb13b33a34b94cf36b202b69fde139ca17a85051979867720d4bdae4323d4943ddf9aeeb6643633aa656e0be843659795007a@35.177.226.168:30303",
	"enode://0cc5f5ffb5d9098c8b8c62325f3797f56509bff942704687b6530992ac706e2cb946b90a34f1f19548cd3c7baccbcaea354531e5983c7d1bc0dee16ce4b6440b@40.118.3.223:30304",
	"enode://1c7a64d76c0334b0418c004af2f67c50e36a3be60b5e4790bdac0439d21603469a85fad36f2473c9a80eb043ae60936df905fa28f1ff614c3e5dc34f15dcd2dc@40.118.3.223:30306",
	"enode://85c85d7143ae8bb96924f2b54f1b3e70d8c4d367af305325d30a61385a432f247d2c75c45c6b4a60335060d072d7f5b35dd1d4c45f76941f62a4f83b6e75daaf@40.118.3.223:30307",
}
