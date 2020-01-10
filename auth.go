package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type IDTokenClaims struct {
	Acr               string           `json:"acr"`
	AllowedOrigins    []string         `json:"allowed-origins"`
	Aud               string           `json:"aud"`
	AuthTime          int              `json:"auth_time"`
	Azp               string           `json:"azp"`
	Email             string           `json:"email"`
	Exp               int              `json:"exp"`
	FamilyName        string           `json:"family_name"`
	GivenName         string           `json:"given_name"`
	Iat               int              `json:"iat"`
	Iss               string           `json:"iss"`
	Jti               string           `json:"jti"`
	Name              string           `json:"name"`
	Nbf               int              `json:"nbf"`
	Nonce             string           `json:"nonce"`
	PreferredUsername string           `json:"preferred_username"`
	RealmAccess       Roles            `json:"realm_access"`
	ResourceAccess    map[string]Roles `json:"resource_access"`
	SessionState      string           `json:"session_state"`
	Sub               string           `json:"sub"`
	Typ               string           `json:"typ"`
}

type ipRange struct {
	start net.IP
	end   net.IP
}

var allowedRanges = []ipRange{
	ipRange{
		start: net.ParseIP("xx.xx.xx.xx"),
		end:   net.ParseIP("xx.xx.xx.xx"),
	},
	ipRange{
		start: net.ParseIP("xx.xx.xx.xx"),
		end:   net.ParseIP("xx.xx.xx.xx"),
	},
	ipRange{
		start: net.ParseIP("xx.xx.xx.xx"),
		end:   net.ParseIP("xx.xx.xx.xx"),
	},
}

func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

func isAllowedRange(ipAddress net.IP) bool {
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		for _, r := range allowedRanges {
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func (a *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Detect client IP
		ip := getRealIP(r)

		// Check if IP is private
		//private, err := isPrivateIP(ip)
		//if err != nil {
		//	respondWithError(w, http.StatusBadRequest, err.Error())
		//	return
		//}

		ip = strings.TrimSpace(ip)
		chip := net.ParseIP(ip)

		if isAllowedRange(chip) {
			next.ServeHTTP(w, r)
		} else {
			authHeader := strings.Split(strings.TrimSpace(r.Header.Get("Authorization")), " ")
			if len(authHeader) == 2 && strings.ToLower(authHeader[0]) == "bearer" && len(authHeader[1]) > 0 {

				token, err := a.tokenVerifier.Verify(context.TODO(), authHeader[1])
				if err != nil {
					respondWithError(w, http.StatusUnauthorized, err.Error())
					return
				}

				// parse claims
				var claims IDTokenClaims
				if err := token.Claims(&claims); err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				// Check permission
				if !checkPermission(claims.RealmAccess.Roles) {
					respondWithError(w, http.StatusForbidden, "Access denied")
					return
				}

				next.ServeHTTP(w, r)
			} else {
				respondWithError(w, http.StatusBadRequest, "Token not found")
				return
			}
		}

	})
}

func checkPermission(roles []string) bool {
	if roles != nil {
		for _, r := range roles {
			if r == "bb_user" {
				return true
			}
		}
	}
	return false
}

func getRealIP(r *http.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}

func isPrivateIP(ip string) (bool, error) {
	var err error
	private := false
	IP := net.ParseIP(ip)
	if IP == nil {
		err = errors.New("Invalid IP")
	} else {
		_, isLoopback, _ := net.ParseCIDR("127.0.0.0/8")
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP) || isLoopback.Contains(IP)
	}
	return private, err
}
