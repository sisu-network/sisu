// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0;

contract Dummy {
    string name;

    constructor(string memory initName) {
        name = initName;
    }

    function setName(string memory newName) public {
        name = newName;
    }

    function getName() public view returns (string memory) {
        return name;
    }
}
