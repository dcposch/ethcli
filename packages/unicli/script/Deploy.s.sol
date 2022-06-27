// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "forge-std/Script.sol";

import "src/UniswapFrontend.sol";

contract Deploy is Script {
    function setUp() public {}

    function run() public {
        vm.broadcast();
        new UniswapFrontend();
    }
}
