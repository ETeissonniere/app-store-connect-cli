/*
Copyright © 2023 Eliott Teissonniere

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package helpers

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var appleAudienceClaim = jwt.ClaimStrings{"appstoreconnect-v1"}

type appleClaims struct {
	jwt.RegisteredClaims

	BundleId string `json:"bid"`
}

// A helper function for other commands to generare the JWT necessary to interact
// with Apple's API.
func GenerateJWT(issuerId, bundleId, keyId string, privateKey *ecdsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		appleClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    issuerId,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
				Audience:  appleAudienceClaim,
			},
			BundleId: bundleId,
		},
	)
	token.Header["alg"] = "ES256"
	token.Header["typ"] = "JWT"
	token.Header["kid"] = keyId

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign JWT: %w", err)
	}

	return signed, nil
}
