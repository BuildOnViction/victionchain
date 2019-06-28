package services

import (
	"github.com/tomochain/tomox-sdk/errors"

	"github.com/tomochain/tomox-sdk/utils"
)

var logger = utils.Logger

var ErrPairExists = errors.New("Pairs already exists")
var ErrPairNotFound = errors.New("Pair not found")
var ErrBaseTokenNotFound = errors.New("BaseToken not found")
var ErrQuoteTokenNotFound = errors.New("QuoteToken not found")
var ErrQuoteTokenInvalid = errors.New("Quote Token Invalid (not a quote)")
var ErrTokenExists = errors.New("Token already exists")
var ErrAccountNotFound = errors.New("Account not found")
var ErrAccountExists = errors.New("Account already Exists")
var ErrNoContractCode = errors.New("Contract not found at given address")
