package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"

	"github.com/rs/zerolog/log"
	"github.com/spruceid/siwe-go"
)

// var sessionStore = sessions.NewCookieStore([]byte("siwe-quickstart-secret"))
// cookies:
//
// siwe, nonce

type httpResponse struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

var (
	sessionName = "xyz.lucidconnect.auth"
)

func (s *Server) GetNonce() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, sessionName)
		session.Values["nonce"] = siwe.GenerateNonce()

		session.Save(r, w)
		fmt.Println("session", session.ID)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(session.Values["nonce"].(string)))
	}
}

func (s *Server) VerifyMerchant() http.HandlerFunc {
	type requestBody struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}
	type responseData struct {
		Valid           bool   `json:"status"`
		Token           string `json:"token"`
		Address         string `json:"address,omitempty"`
		MerchantId      string `json:"merchant_id,omitempty"`
		IsValidMerchant bool   `json:"is_valid_merchant"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := &requestBody{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Err(err).Msg("decoding request failed")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session, _ := s.sessionStore.Get(r, sessionName)
		fmt.Println("session", session.ID)
		message := request.Message
		signature := request.Signature
		nonce, ok := session.Values["nonce"].(string)
		if !ok {
			log.Error().Msg("nonce is empty")
			response := &httpResponse{Status: http.StatusInternalServerError, Error: ""}
			writeJsonResponse(w, response)
			return
		}
		siweObj, err := siwe.ParseMessage(message)
		if err != nil {
			log.Err(err).Msg("parsing siwe message failed")
			response := &httpResponse{Status: http.StatusBadRequest, Error: "parsing siwe message failed"}
			writeJsonResponse(w, response)
			return
		}

		pkey, err := siweObj.Verify(signature, nil, &nonce, nil)
		if err != nil {
			log.Err(err).Msg("invalid signature")
			data := &responseData{Valid: false}
			response := &httpResponse{Status: http.StatusBadRequest, Data: data, Error: "invalid signature"}
			writeJsonResponse(w, response)
			return
		}

		address := crypto.PubkeyToAddress(*pkey)

		// session.Values["siwe"] = siweObj
		// fmt.Println("siwe:", session.Values["siwe"])
		// if err := session.Save(r, w); err != nil {
		// 	log.Err(err).Send()
		// 	return
		// }
		// fmt.Println("session: ", session.Values)
		jwt, err := generateJwt(message)
		if err != nil {
			log.Err(err).Msg("generating jwt failed")
			response := &httpResponse{Status: http.StatusInternalServerError}
			writeJsonResponse(w, response)
			return
		}
		data := &responseData{Valid: true, Address: address.Hex(), Token: jwt, IsValidMerchant: false}

		if merchant, err := s.database.FetchMerchantByAddress(address.Hex()); err == nil {
			data.MerchantId = merchant.ID.String()
			data.IsValidMerchant = true
		}

		response := &httpResponse{Status: http.StatusOK, Data: data}
		w.Header().Set("Content-Type", "text/plain")
		writeJsonResponse(w, response)
	}
}

func writeJsonResponse(w http.ResponseWriter, response *httpResponse) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Err(err).Send()
		w.WriteHeader(response.Status)
		return
	}
}

func generateJwt(siweMsg string) (string, error) {
	var secretKey = os.Getenv("JWT_SECRET")
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		log.Err(err).Msg("decoding secret key failed")
		return "", err
	}
	log.Info().Msg(string(key))

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["siwe"] = siweMsg

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("claims", token.Claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
