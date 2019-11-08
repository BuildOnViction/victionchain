pragma solidity ^0.4.24;

contract RelayerRegistration {
    using SafeMath for uint;

    /// @dev constructor arguments
    address public CONTRACT_OWNER;
    uint public MaximumRelayers;
    uint public MaximumTokenList;

    /// @dev Data types
    struct Relayer {
        uint256 _deposit;
        uint16 _tradeFee;
        address[] _fromTokens;
        address[] _toTokens;
        uint _index;
        address _owner;
    }

    /// @DEV coinbase -> relayer
    mapping(address => Relayer) private RELAYER_LIST;
    /// @dev index -> coinbase
    mapping(uint => address) public RELAYER_COINBASES;
    /// @dev coinbase -> time
    mapping(address => uint) private RESIGN_REQUESTS;
    /// @dev coinbase -> price
    mapping(address => uint256) public RELAYER_ON_SALE_LIST;

    uint public RelayerCount;
    uint256 public MinimumDeposit;

    /// @dev Events
    /// struct-mapping -> values
    event ConfigEvent(uint max_relayer, uint max_token, uint256 min_deposit);
    event RegisterEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens);
    event UpdateEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens);
    event TransferEvent(address owner, uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens);
    event ResignEvent(uint deposit_release_time, uint256 deposit_amount);
    event RefundEvent(bool success, uint remaining_time, uint256 deposit_amount);

    event SellEvent(bool is_on_sale, address coinbase, uint256 price);
    event BuyEvent(bool success, address coinbase, uint256 price);

    constructor (uint maxRelayers, uint maxTokenList, uint minDeposit) public {
        RelayerCount = 0;
        MaximumRelayers = maxRelayers;
        MaximumTokenList = maxTokenList;
        MinimumDeposit = minDeposit;
        CONTRACT_OWNER = msg.sender;
    }


    /// @dev Modifier
    modifier contractOwnerOnly() {
        require(msg.sender == CONTRACT_OWNER, "Contract Owner Only.");
        _;
    }

    modifier relayerOwnerOnly(address coinbase) {
        require(msg.sender == RELAYER_LIST[coinbase]._owner, "Relayer Owner Only.");
        _;
    }

    modifier onlyActiveRelayer(address coinbase) {
        require(RESIGN_REQUESTS[coinbase] == 0, "The relayer has been requested to close.");
        _;
    }

    modifier nonZeroValue() {
        require(msg.value > 0, "Transfer value must be > 0");
        _;
    }

    modifier notForSale(address coinbase) {
        require(RELAYER_ON_SALE_LIST[coinbase] == 0, "The relayer must be not currently for Sale");
        _;
    }


    /// @dev Contract Config Modifications
    function reconfigure(uint maxRelayer, uint maxToken, uint minDeposit) public contractOwnerOnly {
        require(maxRelayer >= RelayerCount);
        require(maxToken > 4 && maxToken < 1001);
        require(minDeposit > 10000);
        MaximumRelayers = maxRelayer;
        MaximumTokenList = maxToken;
        MinimumDeposit = minDeposit;

        emit ConfigEvent(MaximumRelayers,
                         MaximumTokenList,
                         MinimumDeposit);
    }


    /// @dev State-Alter Methods
    function register(address coinbase, uint16 tradeFee, address[] memory fromTokens, address[] memory toTokens) public payable {
        require(msg.sender != CONTRACT_OWNER, "Contract Owner is forbidden to create a Relayer");
        require(msg.sender != coinbase, "Coinbase and RelayerOwner address must not be the same");
        require(coinbase != CONTRACT_OWNER, "Coinbase must not be same as CONTRACT_OWNER");
        require(msg.value >= MinimumDeposit, "Minimum deposit not satisfied.");
        /// @dev valid relayer configuration
        require(tradeFee >= 1 && tradeFee < 10000, "Invalid Maker Fee");
        require(fromTokens.length <= MaximumTokenList, "Exceeding number of trade pairs");
        require(toTokens.length == fromTokens.length, "Not valid number of Pairs");

        require(RELAYER_LIST[coinbase]._deposit == 0, "Coinbase already registered.");
        require(RelayerCount < MaximumRelayers, "Maximum relayers registered");

        /// @notice Do we need to check the duplication of Token trade-pairs?
        Relayer memory relayer = Relayer(msg.value, tradeFee, fromTokens, toTokens, RelayerCount, msg.sender);
        RELAYER_COINBASES[RelayerCount] = coinbase;
        RELAYER_LIST[coinbase] = relayer;

        RelayerCount++;

        emit RegisterEvent(RELAYER_LIST[coinbase]._deposit,
                           RELAYER_LIST[coinbase]._tradeFee,
                           RELAYER_LIST[coinbase]._fromTokens,
                           RELAYER_LIST[coinbase]._toTokens);
    }


    function update(address coinbase, uint16 tradeFee, address[] memory fromTokens, address[] memory toTokens) public relayerOwnerOnly(coinbase) onlyActiveRelayer(coinbase) notForSale(coinbase) {
        require(tradeFee >= 1 && tradeFee < 10000, "Invalid Maker Fee");
        require(fromTokens.length <= MaximumTokenList, "Exceeding number of trade pairs");
        require(toTokens.length == fromTokens.length, "Not valid number of Pairs");

        RELAYER_LIST[coinbase]._tradeFee = tradeFee;
        RELAYER_LIST[coinbase]._fromTokens = fromTokens;
        RELAYER_LIST[coinbase]._toTokens = toTokens;

        emit UpdateEvent(RELAYER_LIST[coinbase]._deposit,
                         RELAYER_LIST[coinbase]._tradeFee,
                         RELAYER_LIST[coinbase]._fromTokens,
                         RELAYER_LIST[coinbase]._toTokens);
    }


    function transfer(address coinbase, address new_owner) public relayerOwnerOnly(coinbase) onlyActiveRelayer(coinbase) notForSale(coinbase) {
        require(new_owner != address(0) && new_owner != msg.sender);
        require(RELAYER_LIST[new_owner]._tradeFee == 0, "Owner address must not be currently used as relayer-coinbase");

        RELAYER_LIST[coinbase]._owner = new_owner;
        emit TransferEvent(RELAYER_LIST[coinbase]._owner,
                           RELAYER_LIST[coinbase]._deposit,
                           RELAYER_LIST[coinbase]._tradeFee,
                           RELAYER_LIST[coinbase]._fromTokens,
                           RELAYER_LIST[coinbase]._toTokens);
    }


    function depositMore(address coinbase) public payable relayerOwnerOnly(coinbase) onlyActiveRelayer(coinbase) notForSale(coinbase) nonZeroValue {
        require(msg.value >= 1 ether, "At least 1 TOMO is required for a deposit request");
        RELAYER_LIST[coinbase]._deposit.add(msg.value);
        emit UpdateEvent(RELAYER_LIST[coinbase]._deposit,
                         RELAYER_LIST[coinbase]._tradeFee,
                         RELAYER_LIST[coinbase]._fromTokens,
                         RELAYER_LIST[coinbase]._toTokens);
    }


    function resign(address coinbase) public relayerOwnerOnly(coinbase) notForSale(coinbase) {
        require(RELAYER_LIST[coinbase]._deposit > 0, "No relayer associated with this address");
        require(RESIGN_REQUESTS[coinbase] == 0, "Request already received");
        /// @notice: for testing contract, change `4 weeks` to 4 seconds only
        RESIGN_REQUESTS[coinbase] = now + 4 weeks;
        emit ResignEvent(RESIGN_REQUESTS[coinbase],
                         RELAYER_LIST[coinbase]._deposit);
    }


    function refund(address coinbase) public relayerOwnerOnly(coinbase) notForSale(coinbase) {
        require(RESIGN_REQUESTS[coinbase] > 0, "Request not found");
        uint256 amount = RELAYER_LIST[coinbase]._deposit;
        uint deleting_index = RELAYER_LIST[coinbase]._index;

        if (RESIGN_REQUESTS[coinbase] < now) {
            delete RELAYER_LIST[coinbase];
            delete RESIGN_REQUESTS[coinbase];
            /// @notice swap last relayer's index with the deleting relayer's index
            address last_coinbase = RELAYER_COINBASES[RelayerCount - 1];
            delete RELAYER_COINBASES[RelayerCount - 1];

            RELAYER_COINBASES[deleting_index] = last_coinbase;
            RELAYER_LIST[last_coinbase]._index = deleting_index;

            RelayerCount--;
            msg.sender.transfer(amount);
            emit RefundEvent(true,
                             0,
                             amount);
        } else {
            /// Not yet, remind user about the remaining time...
            emit RefundEvent(false,
                             RESIGN_REQUESTS[coinbase] - now,
                             amount);
        }
    }

    /// @dev SELLING/BUYING RELAYERS
    // NOTE: Putting a relayer on sale will immediately halt all the state-changing actions on this relayer...
    // ...including Update, Deposit, Resign, Refund and Transfer
    function sellRelayer(address coinbase, uint256 price) public relayerOwnerOnly(coinbase) onlyActiveRelayer(coinbase) {
        // NOTE: calling this function more than once will act like changing Price-tag
        require(price > 0, "Price tag must be different than Zero(0)");
        RELAYER_ON_SALE_LIST[coinbase] = price;
        emit SellEvent(true, coinbase, price);
    }

    function cancelSelling(address coinbase) public relayerOwnerOnly(coinbase) onlyActiveRelayer(coinbase) {
        require(RELAYER_ON_SALE_LIST[coinbase] > 0, "Relayer is not currently for sale");
        delete RELAYER_ON_SALE_LIST[coinbase];
        emit SellEvent(false, coinbase, 0);
    }

    function buyRelayer(address coinbase) public payable onlyActiveRelayer(coinbase) {
        uint256 price = RELAYER_ON_SALE_LIST[coinbase];

        require(price > 0, "Relayer is not currently for sale");
        require(msg.value == price, "Price-tag must be matched");

        address seller = RELAYER_LIST[coinbase]._owner;
        require(msg.sender != address(0) && msg.sender != seller && seller != address(0), "Address not valid");
        RELAYER_LIST[coinbase]._owner = msg.sender;

        delete RELAYER_ON_SALE_LIST[coinbase];
        seller.transfer(price);
        emit BuyEvent(true, coinbase, msg.value);
    }


    function getRelayerByCoinbase(address coinbase) public view returns (uint, address, uint256, uint16, address[] memory, address[] memory) {
        return (RELAYER_LIST[coinbase]._index,
                RELAYER_LIST[coinbase]._owner,
                RELAYER_LIST[coinbase]._deposit,
                RELAYER_LIST[coinbase]._tradeFee,
                RELAYER_LIST[coinbase]._fromTokens,
                RELAYER_LIST[coinbase]._toTokens);
    }

}