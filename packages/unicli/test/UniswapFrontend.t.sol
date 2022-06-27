// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "forge-std/Test.sol";

import "src/UniswapFrontend.sol";

contract ContractTest is Test {
    function setUp() public {}

    function testRender() public {
        UniswapFrontend f = new UniswapFrontend();

        f.render(hex"");
    }
}
