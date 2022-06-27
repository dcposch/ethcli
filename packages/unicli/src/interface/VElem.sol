// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

uint64 constant TYPE_TEXT = uint64(uint256(keccak256("text")));
uint64 constant TYPE_IN_AMOUNT = uint64(uint256(keccak256("amount")));
uint64 constant TYPE_IN_DROPDOWN = uint64(uint256(keccak256("dropdown")));
uint64 constant TYPE_IN_TEXTBOX = uint64(uint256(keccak256("textbox")));
uint64 constant TYPE_BUTTON = uint64(uint256(keccak256("button")));

struct VElem {
    /** @dev Text field, input, button, etc. */
    uint64 typeHash;
    /** @dev Text for a text field, dropdown options, etc. See ElemAmount etc.*/
    bytes data;
}

/** @dev Virtual DOM helper library. */
library V {
    function Text(uint256 key, string memory text)
        internal
        pure
        returns (VElem memory)
    {
        return VElem(TYPE_TEXT, abi.encode(ElemText(key, text)));
    }

    function Amount(
        uint256 key,
        string memory label,
        uint64 decimals
    ) internal pure returns (VElem memory) {
        return
            VElem(TYPE_IN_AMOUNT, abi.encode(ElemAmount(key, label, decimals)));
    }

    function Dropdown(
        uint256 key,
        string memory label,
        DropOpt[] memory options
    ) internal pure returns (VElem memory) {
        return
            VElem(
                TYPE_IN_DROPDOWN,
                abi.encode(ElemDropdown(key, label, options))
            );
    }

    function Button(uint256 key, string memory text)
        internal
        pure
        returns (VElem memory)
    {
        return VElem(TYPE_BUTTON, abi.encode(ElemButton(key, text)));
    }
}

struct ElemText {
    uint256 key;
    /** @dev UTF-8 text. Line breaks preserved. May auto-wrap at >=80chars. */
    string text;
}

struct ElemAmount {
    uint256 key;
    /** @dev Form input label */
    string label;
    /** @dev Amount input will return fixed-point uint256 to n decimals. */
    uint64 decimals;
}

struct ElemDropdown {
    uint256 key;
    /** @dev Form input label */
    string label;
    /** Options. User must pick one. */
    DropOpt[] options;
}

struct DropOpt {
    /** @dev Dropdown option ID */
    uint256 id;
    /** @dev Dropdown option display string */
    string text;
}

struct ElemButton {
    uint256 key;
    /** Button text */
    string text;
}
