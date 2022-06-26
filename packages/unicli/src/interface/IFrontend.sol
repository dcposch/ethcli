// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

uint256 constant TYPE_TEXT = uint256(keccak256("text"));

uint256 constant TYPE_IN_AMOUNT = uint256(keccak256("amount"));
uint256 constant TYPE_IN_DROPDOWN = uint256(keccak256("dropdown"));
uint256 constant TYPE_IN_TEXTBOX = uint256(keccak256("textbox"));

uint256 constant TYPE_BUTTON = uint256(keccak256("button"));

struct VdomElem {
    /** @dev Text field, input, button, etc. */
    uint256 typeHash;
    /** @dev Text for a text field, options for a dropdown, etc. */
    bytes data;
}

struct DataAmount {
    string label;
    /** @dev Amount input will return fixed-point uint256 to n decimals. */
    uint256 decimals;
}

struct DataDropdown {
    string label;
    /** Options. User must pick one. */
    DataDropOption[] options;
}

struct DataDropOption {
    /** @dev Dropdown option ID */
    uint256 id;
    /** @dev Dropdown option display string */
    string display;
}

struct Action {
    /** @dev 0 = first button, etc. */
    uint256 buttonId;
    /** @dev Value of each input.  */
    bytes[] inputs;
}

interface IFrontend {
    function render(bytes calldata appState)
        external
        view
        returns (VdomElem[] memory vdom);

    function act(bytes calldata appState, Action calldata action)
        external
        returns (bytes memory newAppState);
}
