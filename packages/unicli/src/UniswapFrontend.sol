// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./interface/IFrontend.sol";
import "./interface/IUniswap.sol";
import "./interface/VElem.sol";

IUniswapV2Factory constant uniFactory = IUniswapV2Factory(
    0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D
);

contract UniswapFrontend is IFrontend {
    function render(bytes calldata appState)
        external
        pure
        override
        returns (VElem[] memory vdom)
    {
        require(appState.length == 0, "Unexpected state");

        vdom = new VElem[](5);
        vdom[0] = V.Text(1, "HELLO WORLD");
        vdom[1] = V.Amount(2, "Amount in", 18);
        vdom[2] = V.Dropdown(3, "Token in", _tokens());
        vdom[3] = V.Dropdown(3, "Token out", _tokens());
        vdom[4] = V.Button(5, "Swap");
    }

    function _tokens() public pure returns (DropOpt[] memory ret) {
        ret = new DropOpt[](3);
        ret[0] = DropOpt(1, "ETH");
        ret[1] = DropOpt(0x00f80a32a835f79d7787e8a8ee5721d0feafd78108, "DAI");
        ret[2] = DropOpt(0x00c778417e063141139fce010982780140aa0cd5ab, "WETH");
    }

    function act(bytes calldata appState, Action calldata action)
        external
        returns (bytes memory newAppState)
    {}
}
