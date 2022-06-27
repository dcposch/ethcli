// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import "./VElem.sol";

interface IFrontend {
    function render(bytes calldata appState)
        external
        view
        returns (VElem[] memory vdom);

    function act(bytes calldata appState, Action calldata action)
        external
        returns (bytes memory newAppState);
}

struct Action {
    /** @dev Which button was pressed. */
    uint256 buttonKey;
    /** @dev ABI serialization of each input.  */
    bytes[] inputs;
}
