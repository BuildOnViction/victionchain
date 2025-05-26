pragma solidity ^0.4.24;

import "./TRC21Issuer.sol";

contract VRC25Issuer is TRC21Issuer {
    event Withdraw(address indexed token, address indexed receiver, uint256 value);

    constructor(uint256 value) TRC21Issuer(value) public {
        revert("Constructor is not supported, override code to TRC21Issuer address");
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
}
