pragma solidity ^0.4.24;

contract TOMOXListing {

    address[] _tokens;
    mapping(address => TokenState) tokensState;


    struct TokenState {
        bool isActive;
    }

    modifier onlyValidApplyNewToken(address token){
        require(token != address(0));
        require(tokensState[token].isActive != true);
        _;
    }

    function tokens() public view returns(address[]) {
        return _tokens;
    }

    function getTokenStatus(address token) public view returns(bool) {
        return tokensState[token].isActive;
    }

    function apply(address token) public payable onlyValidApplyNewToken(token){
        require(msg.value >= 100 ether);
        _tokens.push(token);
        tokensState[token] = TokenState({
            isActive: true
        });
    }
}
