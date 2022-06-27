// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./interface/IFrontend.sol";
import "./interface/IUniswap.sol";

IUniswapV2Factory constant uniFactory = IUniswapV2Factory(
    0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D
);

contract UniswapFrontend is IFrontend {
    function render(bytes calldata appState)
        external
        pure
        override
        returns (VdomElem[] memory vdom)
    {
        require(appState.length == 0, "Unexpected state");

        vdom = new VdomElem[](5);
        vdom[0].typeHash = TYPE_TEXT;
        vdom[0].data = bytes("HELLO WORLD");

        vdom[1].typeHash = TYPE_IN_AMOUNT;
        vdom[1].data = abi.encode(DataAmount("Amount in", 18));

        vdom[2].typeHash = TYPE_IN_DROPDOWN;
        vdom[2].data = abi.encode(DataDropdown("Token in", _tokens()));

        vdom[3].typeHash = TYPE_IN_DROPDOWN;
        vdom[3].data = abi.encode(DataDropdown("Token out", _tokens()));

        vdom[4].typeHash = TYPE_BUTTON;
        vdom[4].data = bytes("Swap");

        // for (uint256 i = 0; i < 10; i++) {
        //     IUniswapV2Pair pair = IUniswapV2Pair(uniFactory.allPairs(i));
        //     IERC20Meta t0 = IERC20Meta(pair.token0());
        //     IERC20Meta t1 = IERC20Meta(pair.token1());
        //     // pair.price0CumulativeLast();
        //     // pair.price1CumulativeLast();
        //     // pair.kLast();
        // }
    }

    function _tokens() public pure returns (DataDropOption[] memory ret) {
        ret = new DataDropOption[](2);

        ret[0] = DataDropOption(1, "ETH");
        ret[2] = DataDropOption(
            0x00f80a32a835f79d7787e8a8ee5721d0feafd78108,
            "DAI"
        );
        ret[3] = DataDropOption(
            0x00c778417e063141139fce010982780140aa0cd5ab,
            "WETH"
        );
    }

    function act(bytes calldata appState, Action calldata action)
        external
        returns (bytes memory newAppState)
    {}
}
