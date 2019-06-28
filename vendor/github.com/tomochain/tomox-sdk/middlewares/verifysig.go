package middlewares

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

func VerifySignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Signature"] == nil || r.Header["Hash"] == nil || r.Header["Pubkey"] == nil {
			httputils.WriteError(w, http.StatusUnauthorized, "There is not enough parameters in header")
			return
		}

		hash := common.Hex2Bytes(r.Header["Hash"][0])
		signature := common.Hex2Bytes(r.Header["Signature"][0])
		publicKeyBytes := common.Hex2Bytes(r.Header["Pubkey"][0])

		signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
		verified := crypto.VerifySignature(publicKeyBytes, hash, signatureNoRecoverID)

		if !verified {
			httputils.WriteError(w, http.StatusUnauthorized, "Signature Invalid")
			return
		}

		next.ServeHTTP(w, r)
	})
}
