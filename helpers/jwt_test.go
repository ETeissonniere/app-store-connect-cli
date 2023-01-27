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
	"crypto/elliptic"
	"crypto/rand"
	"reflect"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func TestGenerateJWT(t *testing.T) {
	dummyIssuerId := "dummyIssuerId"
	dummyBundleId := "dummyBundleId"
	dummyKeyId := "dummyKeyId"
	dummyPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	testSigned, err := GenerateJWT(dummyIssuerId, dummyBundleId, dummyKeyId, dummyPrivateKey)
	if err != nil {
		t.Errorf("GenerateJWT() error = %v", err)
	}
	testToken, err := jwt.NewParser().ParseWithClaims(testSigned, &appleClaims{}, func(token *jwt.Token) (interface{}, error) {
		return dummyPrivateKey.Public(), nil
	})
	if err != nil {
		t.Errorf("jwt.Parse() error = %v", err)
	}
	testClaims := testToken.Claims.(*appleClaims)
	testHeaders := testToken.Header

	t.Run("should not live longer than an hour", func(t *testing.T) {
		expiration := testClaims.ExpiresAt.Time
		now := time.Now()

		if expiration.Sub(now) > time.Hour {
			t.Errorf("token expiration is too long: %v", expiration.Sub(now))
		}
	})

	t.Run("headers follow apple format", func(t *testing.T) {
		assertHeaderEqual := func(header string, expected interface{}) {
			if testHeaders[header] != expected {
				t.Fatalf("wrong header %s: expected %v, got %v", header, expected, testHeaders[header])
			}
		}
		assertHeaderEqual("alg", "ES256")
		assertHeaderEqual("kid", dummyKeyId)
		assertHeaderEqual("typ", "JWT")
	})

	t.Run("claims follow apple format", func(t *testing.T) {
		assertClaimEqual := func(claimName string, claim interface{}, expected interface{}) {
			if !reflect.DeepEqual(claim, expected) {
				t.Fatalf("wrong claim %s: expected %v, got %v", claimName, expected, claim)
			}
		}
		assertClaimEqual("iss", testClaims.Issuer, dummyIssuerId)
		assertClaimEqual("aud", testClaims.Audience, appleAudienceClaim)
		assertClaimEqual("bid", testClaims.BundleId, dummyBundleId)
	})
}
