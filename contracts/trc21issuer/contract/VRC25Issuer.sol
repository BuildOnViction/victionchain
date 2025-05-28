pragma solidity ^0.4.24;

import "./TRC21Issuer.sol";

contract VRC25Issuer is TRC21Issuer {
    mapping(address => uint256) private _maxGasPrice; // if maxGasPrice == 0 then no limit

    event Withdraw(address indexed token, address indexed receiver, uint256 value);

    constructor(uint256 value) TRC21Issuer(value) public {
        revert("Constructor is not supported, override code to TRC21Issuer address");
    }

    function updateMaxGasPrice(address token, uint256 maxGasPrice) public {
        AbstractTokenTRC21 t = AbstractTokenTRC21(token);
        require(t.issuer() == msg.sender, "Only issuer can update max gas price");
        _maxGasPrice[token] = maxGasPrice;
    }

    function apply(address token) public payable onlyValidCapacity(token) {
        require(tokensState[token] == 0, "Token already applied"); // cannot apply twice
        super.apply(token);
    }

    function withdraw(address token, address receiver, uint256 amount) public {
        AbstractTokenTRC21 t = AbstractTokenTRC21(token);
        require(t.issuer() == msg.sender, "Only issuer can withdraw");
        require(tokensState[token] >= amount, "Insufficient capacity to withdraw");

        tokensState[token] = tokensState[token].sub(amount);

        receiver.transfer(amount);
        emit Charge(msg.sender, token, amount);
    }

    function maxGasPrice(address token) public view returns (uint256) {
        return _maxGasPrice[token];
    }
}
