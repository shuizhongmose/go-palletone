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

package configure

// MainnetBootnodes are the pnode URLs of the P2P bootstrap nodes running on
// the main PalletOne network.
var MainnetBootnodes = []string{
	//"pnode://98627f6b5fa0f549ff644e66fc2b9801ea039614784ea65d1e9e942fed573004a65d0b2e85d4a949c28844efe6f357d97c1198e1194eb436fb0a8bbfafb8e6d2@123.126.106.88:30303",
	//"pnode://eaef728cb3bb7d96f5efb377b2a21d134061eb4ae0276a9573055f998b7fde4fece4aac1d93c5d9bb57d8ddf054f922720dd6188f30c8f85cc68ee9e37465958@123.126.106.88:30305",
	//"pnode://d9d50a943c836e1576589e948c5f0b92429173e617906bcecbe16b528f3bd924f831d86b1abdb646c9ecfed466760f18423f35ef26a5ef73ac102418da6dc3c9@123.126.106.88:30306",
	//"pnode://beacea1636d1ee955247199d09cdc39db38db115616e22dd6c50939b44d40a876cbda4f7cedf86a2398ff67f8d67c63fbc88be1e86277377c75250ae023aa15b@123.126.106.89:30307",
	//"pnode://5e93e974f036fa917d66277dafb020c2c705321c82fefde81b186a934a42867324f966a83f51884a0b1e8694bd241405fbdac3131f644de3d1d4b8f88a52886a@123.126.106.89:30308",
}

// TestnetBootnodes are the pnode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"pnode://27108152fecd83368b40c3451b6eb774f829a89273ccdcd2c82887e2dc8afc6264d7a9797bb7abeec6be59b3c9b2f8ef6796f30e5f9b5d27b5880896de39a713@123.126.106.85:30309",
	"pnode://1d27062fc6656f4f7ff4360ffa6464f9e6ebcb27c3c766fba74ec8969a3b2d4c075a2e37d83dd9bb0e28cd761baffe97ee372f3141e5631fa3edd9f0c718633c@123.126.106.85:30310",
}

var DiscoveryV5Bootnodes = []string{}
