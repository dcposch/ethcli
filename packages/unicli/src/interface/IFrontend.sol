// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

uint64 constant TYPE_TEXT = uint64(uint256(keccak256("text")));
uint64 constant TYPE_IN_AMOUNT = uint64(uint256(keccak256("amount")));
uint64 constant TYPE_IN_DROPDOWN = uint64(uint256(keccak256("dropdown")));
uint64 constant TYPE_IN_TEXTBOX = uint64(uint256(keccak256("textbox")));
uint64 constant TYPE_BUTTON = uint64(uint256(keccak256("button")));

struct VdomElem {
    /** @dev Text field, input, button, etc. */
    uint64 typeHash;
    /** @dev Text for a text field, options for a dropdown, etc. */
    bytes data;
}

struct DataAmount {
    string label;
    /** @dev Amount input will return fixed-point uint256 to n decimals. */
    uint64 decimals;
}

// 0000000000000000000000000000000000000000000000000000000000000020 // h (DataAmount offset)
// 0000000000000000000000000000000000000000000000000000000000000040 // t head (label offset)
// 0000000000000000000000000000000000000000000000000000000000000012 // t head (decimals) = 18
// 0000000000000000000000000000000000000000000000000000000000000009 // t tail (label) string length
// 416d6f756e7420696e0000000000000000000000000000000000000000000000 // t tail (label) string value

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

// 0000000000000000000000000000000000000000000000000000000000000020
// 0000000000000000000000000000000000000000000000000000000000000020
// 0000000000000000000000000000000000000000000000000000000000000004
// 5377617000000000000000000000000000000000000000000000000000000000
struct DataButton {
    string text;
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
