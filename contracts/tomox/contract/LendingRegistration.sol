pragma solidity 0.4.24;

contract AbstractRegistration {
    mapping(address => uint) public RESIGN_REQUESTS;
    function getRelayerByCoinbase(address) public view returns (uint, address, uint256, uint16, address[] memory, address[] memory);
}

contract Lending {
    
    struct LendingRelayer {
        uint16 _tradeFee;
        address[] _baseTokens;
        uint256[] _terms; // seconds
    }
    
    mapping(address => LendingRelayer) public LENDINGRELAYER_LIST;

    mapping(address => uint256) public COLLATERAL_LIST;
    address[] public COLLATERALS;
    
    address[] public BASES;
    
    uint256[] TERMS;

    AbstractRegistration public relayer;
    
    constructor (address r) public {
        relayer = AbstractRegistration(r);
    }
    
    function addCollateral(address token, uint256 depositRate) public {
        COLLATERAL_LIST[token] = depositRate;
        COLLATERALS.push(token);
    }
    
    function addBaseToken(address token) public {
        BASES.push(token);
    }
    
    function addTerm(uint256 term) public {
        TERMS.push(term);
    }
    
    function register(address coinbase, uint16 tradeFee, address[] memory baseTokens, uint256[] memory terms) public {
        (, address owner,,,,) = relayer.getRelayerByCoinbase(coinbase);
        require(owner == msg.sender, "owner required");
        require(relayer.RESIGN_REQUESTS(coinbase) == 0, "relayer required to close");
        
        LENDINGRELAYER_LIST[coinbase] = LendingRelayer({
            _tradeFee: tradeFee,
            _baseTokens: baseTokens,
            _terms: terms
        });
    }
    
    function getLendingRelayerByCoinbase(address coinbase) public view returns (uint16, address[] memory, uint256[] memory) {
        return (LENDINGRELAYER_LIST[coinbase]._tradeFee,
                LENDINGRELAYER_LIST[coinbase]._baseTokens,
                LENDINGRELAYER_LIST[coinbase]._terms);
    }
}
