pragma solidity 0.4.24;

contract LAbstractRegistration {
    mapping(address => uint) public RESIGN_REQUESTS;
    function getRelayerByCoinbase(address) public view returns (uint, address, uint256, uint16, address[] memory, address[] memory);
}

contract LAbstractTOMOXListing {
    function getTokenStatus(address) public view returns (bool);
}

contract LAbstractTokenTRC21 {
    function issuer() public view returns (address);
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
    
    uint256[] public TERMS;

    address[] public ILO_COLLATERALS;

    LAbstractRegistration public Relayer;

    address public CONTRACT_OWNER;

    address constant private tomoNative = 0x0000000000000000000000000000000000000001;

    LAbstractTOMOXListing public TomoXListing;

    modifier contractOwnerOnly() {
        require(msg.sender == CONTRACT_OWNER, "Contract Owner Only.");
        _;
    }

    function indexOf(address[] memory addrs, address target) internal pure returns (bool){
        for (uint i = 0; i < addrs.length; i ++) {
            if (addrs[i] == target) {
                return true;
            }
        }
        return false;
    }

    function termIndexOf(uint256[] memory terms, uint256 target) internal pure returns (bool){
        for (uint i = 0; i < terms.length; i ++) {
            if (terms[i] == target) {
                return true;
            }
        }
        return false;
    }
    
    constructor (address r, address t) public {
        Relayer = LAbstractRegistration(r);
        TomoXListing = LAbstractTOMOXListing(t);
        CONTRACT_OWNER = msg.sender;
    }
    
    // add/update depositRate liquidationRate price for collateral
    function addCollateral(address token, uint256 depositRate, uint256 liquidationRate, uint256 price) public contractOwnerOnly {
        require(depositRate >= 100 && liquidationRate > 100, "Invalid rates");
        require(depositRate > liquidationRate , "Invalid rates");

        bool b = TomoXListing.getTokenStatus(token) || (token == tomoNative);
        require(b, "Invalid collateral");

        COLLATERAL_LIST[token] = Collateral({
            _depositRate: depositRate,
            _liquidationRate: liquidationRate,
            _price: price
        });

        if (!indexOf(COLLATERALS, token)) {
            COLLATERALS.push(token);
        }
    }

    // update price for collateral
    function setCollateralPrice(address token, uint256 price) public {

        bool b = TomoXListing.getTokenStatus(token) || (token == tomoNative);
        require(b, "Invalid collateral");

        require(COLLATERAL_LIST[token]._depositRate >= 100, "Invalid collateral");

        if (indexOf(COLLATERALS, token)) {
            require(msg.sender == CONTRACT_OWNER, "Contract owner required");
        } else {
            LAbstractTokenTRC21 t = LAbstractTokenTRC21(token);
            require(t.issuer() == msg.sender, "Required token issuer");
        }

        COLLATERAL_LIST[token]._price = price;
    }

    function addILOCollateral(address token, uint256 depositRate, uint256 liquidationRate, uint256 price) public {
        require(depositRate >= 100 && liquidationRate > 100, "Invalid rates");
        require(depositRate > liquidationRate , "Invalid rates");

        bool b = TomoXListing.getTokenStatus(token);
        require(b, "Invalid collateral");

        LAbstractTokenTRC21 t = LAbstractTokenTRC21(token);
        require(t.issuer() == msg.sender, "Required token issuer");
        
        COLLATERAL_LIST[token] = Collateral({
            _depositRate: depositRate,
            _liquidationRate: liquidationRate,
            _price: price
        });

        if (!indexOf(ILO_COLLATERALS, token)) {
            ILO_COLLATERALS.push(token);
        }
    }
    
    // lending tokens
    function addBaseToken(address token) public contractOwnerOnly {
        bool b = TomoXListing.getTokenStatus(token) || (token == tomoNative);
        require(b, "Invalid base token");
        if (!indexOf(BASES, token)) {
            BASES.push(token);
        }
    }
    
    // period of loan
    function addTerm(uint256 term) public contractOwnerOnly {
        require(term >= 60, "Invalid term");

        if (!termIndexOf(TERMS, term)) {
            TERMS.push(term);
        }
    }
    
    function update(address coinbase, uint16 tradeFee, address[] memory baseTokens, uint256[] memory terms, address[] memory collaterals) public {
        (, address owner,,,,) = Relayer.getRelayerByCoinbase(coinbase);
        require(owner == msg.sender, "Relayer owner required");
        require(Relayer.RESIGN_REQUESTS(coinbase) == 0, "Relayer required to close");
        require(tradeFee >= 1 && tradeFee < 1000, "Invalid trade Fee"); // 0.01% -> 10%
        require(baseTokens.length == terms.length, "Not valid number of terms");
        require(baseTokens.length == collaterals.length, "Not valid number of collaterals");

        // validate baseTokens
        bool b = false;
        for (uint i = 0; i < baseTokens.length; i++) {
            b = indexOf(BASES, baseTokens[i]);
            require(b == true, "Invalid lending token");
        }

        // validate terms
        for (i = 0; i < terms.length; i++) {
            b = termIndexOf(TERMS, terms[i]);
            require(b == true, "Invalid term");
        }

        // validate collaterals
        for (i = 0; i < collaterals.length; i++) {
            if (collaterals[i] != address(0)) {
                require(indexOf(ILO_COLLATERALS, collaterals[i]), "Invalid collateral");
            }
        }
        
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
