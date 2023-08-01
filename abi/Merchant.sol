// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

// Uncomment this line to use console.log
// import "hardhat/console.sol";

contract LucidMerchant {
    struct Merchant {
        bytes32 name;
        address owner;
        address receivingAddress;
    }

    struct MutableMerchantData {
        bytes32 name;
        address receivingAddress;
    }

    mapping(bytes32 => Merchant) public merchants;

    modifier onlyOwner(bytes32 identifier) {
        require(msg.sender == merchants[identifier].owner, "Unauthorized");
        _;
    }

    function createMerchant(
        bytes32 identifier,
        MutableMerchantData memory data
    ) public {
        require(
            merchants[identifier].owner == address(0),
            "Identifier already exists"
        );

        require(data.name != 0, "Invalid name");
        require(
            data.receivingAddress != address(0),
            "Invalid receiving address"
        );

        Merchant memory newMerchant = Merchant({
            name: data.name,
            owner: msg.sender,
            receivingAddress: data.receivingAddress
        });

        merchants[identifier] = newMerchant;
    }

    function deleteMerchant(bytes32 identifier) public onlyOwner(identifier) {
        delete merchants[identifier];
    }

    function updateMerchant(
        bytes32 identifier,
        MutableMerchantData calldata data
    ) public onlyOwner(identifier) {
        Merchant memory merchant = merchants[identifier];

        require(merchant.owner != address(0), "Merchant does not exist");

        require(data.name != 0, "Invalid name");
        require(
            data.receivingAddress != address(0),
            "Invalid receiving address"
        );

        merchants[identifier].name = data.name;
        merchants[identifier].receivingAddress = data.receivingAddress;
    }

    function getMerchant(
        bytes32 identifier
    ) public view returns (Merchant memory) {
        Merchant memory merchant = merchants[identifier];
        require(merchant.owner != address(0), "Merchant does not exist");

        return merchant;
    }
}
