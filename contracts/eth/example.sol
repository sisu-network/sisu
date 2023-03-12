// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./vault.sol";

contract Example is BaseContract {
    uint256 counter;
    event MessageReceived(uint256 id, bytes message);

    function onReceive(Message calldata input) external returns (uint8 code) {
        require(input.message.length > 0 && input.message.length < 10, "invalid message");

        counter++;
        emit MessageReceived(counter, input.message);

        return 0;
    }
}
