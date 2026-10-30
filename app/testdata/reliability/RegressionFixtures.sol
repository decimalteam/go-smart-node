// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

// Deliberately small ABI-compatible contracts for fault-injection tests.
// These are not replacements for deployed network contracts.
contract RegistryFixture {
    mapping(string => address) private addresses;
    function set(string calldata slug, address value) external { addresses[slug] = value; }
    function getContractAddress(string calldata slug) external view returns (address) { return addresses[slug]; }
}

contract AutoUnbondFixture {
    uint256 public calls;
    uint256 public failHold;
    mapping(bytes32 => uint256) public counts;
    function setFailHold(uint256 value) external { failHold = value; }
    function autoUnbondEnqueue(address val, address del, uint256 amount, address token, uint256 hold) external {
        require(failHold == 0 || hold != failHold, "injected enqueue failure");
        calls++;
        counts[keccak256(abi.encode(val, del, amount, token, hold))]++;
    }
}

contract ValidatorFixture {
    uint256 public marker;
    event RewardPerBlockUpdated(uint256 rewardPerBlock, bool enabled);
    function bump() external { marker++; }
    function emitGood(uint256 amount, bool enabled) external payable {
        marker++;
        emit RewardPerBlockUpdated(amount, enabled);
    }
    function emitBad(bool emptyPayload) external payable {
        marker++;
        bytes32 topic = keccak256("RewardPerBlockUpdated(uint256,bool)");
        assembly {
            mstore(0, 1)
            switch emptyPayload
            case 0 { log1(0, 1, topic) }
            default { log1(0, 0, topic) }
        }
    }
}
