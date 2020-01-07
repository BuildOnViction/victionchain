pragma solidity 0.4.24;

contract AbstractRegistration {
    mapping(address => uint) public RESIGN_REQUESTS;
    function getRelayerByCoinbase(address) public view returns (uint, address, uint256, uint16, address[] memory, address[] memory);
}

contract Lending {
    
    // @dev collateral = 0x0 => get collaterals from COLLATERALS
    struct LendingRelayer {
        uint16 _tradeFee;
        address[] _baseTokens;
        uint256[] _terms; // seconds
        address[] _collaterals;
    }

    struct Collateral {
        uint256 _depositRate;
        uint256 _liquidationRate;
        uint256 _price;
    }
    
    mapping(address => LendingRelayer) public LENDINGRELAYER_LIST;

    mapping(address => Collateral) public COLLATERAL_LIST;
    address[] public COLLATERALS;
    
    address[] public BASES;
    
    uint256[] TERMS;

    address[] public ALL_COLLATERALS;

    AbstractRegistration public relayer;
    
    constructor (address r) public {
        relayer = AbstractRegistration(r);
    }
    
    function addCollateral(address token, uint256 depositRate, uint256 liquidationRate) public {
        COLLATERAL_LIST[token] = Collateral({
            _depositRate: depositRate,
            _liquidationRate: liquidationRate,
            _price: 0
        });
        COLLATERALS.push(token);
        ALL_COLLATERALS.push(token);
    }
    
    function addBaseToken(address token) public {
        BASES.push(token);
    }
    
    function addTerm(uint256 term) public {
        TERMS.push(term);
    }
    
    function update(address coinbase, uint16 tradeFee, address[] memory baseTokens, uint256[] memory terms, address[] memory collaterals) public {
        (, address owner,,,,) = relayer.getRelayerByCoinbase(coinbase);
        require(owner == msg.sender, "relayer owner required");
        require(relayer.RESIGN_REQUESTS(coinbase) == 0, "relayer required to close");
        require(tradeFee >= 1 && tradeFee < 1000, "Invalid Maker Fee"); // 0.01% -> 10%
        require(baseTokens.length == terms.length, "Not valid number of terms");
        require(baseTokens.length == collaterals.length, "Not valid number of collaterals");
        
        LENDINGRELAYER_LIST[coinbase] = LendingRelayer({
            _tradeFee: tradeFee,
            _baseTokens: baseTokens,
            _terms: terms,
            _collaterals: collaterals
        });
    }
    
    function getLendingRelayerByCoinbase(address coinbase) public view returns (uint16, address[] memory, uint256[] memory, address[] memory) {
        return (LENDINGRELAYER_LIST[coinbase]._tradeFee,
                LENDINGRELAYER_LIST[coinbase]._baseTokens,
                LENDINGRELAYER_LIST[coinbase]._terms,
                LENDINGRELAYER_LIST[coinbase]._collaterals);
    }
}
