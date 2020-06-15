pragma solidity 0.4.24;

contract TomoXPrice {
    function GetLastPrice(address base, address quote) public view returns(uint256 p) {
        address[2] memory input;
        input[0] = base;
        input[1] = quote;
        assembly {
            // GetLastPrice precompile!
            if iszero(staticcall(not(0), 0x41, input, 0x28, p, 0x20)) {
                revert(0, 0)
            }
        }
        return p;
    }

    function GetEpochPrice(address base, address quote) public view returns(uint256 p) {
        address[2] memory input;
        input[0] = base;
        input[1] = quote;
        assembly {
            // GetEpochPrice precompile!
            if iszero(staticcall(not(0), 0x42, input, 0x28, p, 0x20)) {
                revert(0, 0)
            }
        }
        return p;
    }
}


