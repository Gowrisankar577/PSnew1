package auth

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"ps_portal/service"
	"ps_portal/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type MicrosoftSignInRequest struct {
	IdToken string `json:"id_token" binding:"required"`
}

// cache verifiers per tenant (tid)
var verifierCache sync.Map // map[string]*oidc.IDTokenVerifier

func getVerifierForTenant(ctx gin.Context, tid, clientID string) (*oidc.IDTokenVerifier, error) {
	if v, ok := verifierCache.Load(tid); ok {
		return v.(*oidc.IDTokenVerifier), nil
	}
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tid)
	provider, err := oidc.NewProvider(ctx.Request.Context(), issuer)
	if err != nil {
		return nil, fmt.Errorf("oidc provider init failed: %w", err)
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID, // issuer check is ON by default (good)
	})
	verifierCache.Store(tid, verifier)
	return verifier, nil
}

// Extract 'tid' from the raw JWT (unverified). We'll verify afterwards.
func extractTid(rawIDToken string) (string, error) {
	parts := strings.Split(rawIDToken, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode payload: %w", err)
	}
	var p struct {
		Tid string `json:"tid"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return "", fmt.Errorf("unmarshal payload: %w", err)
	}
	if p.Tid == "" {
		return "", fmt.Errorf("tid missing in token")
	}
	return p.Tid, nil
}

func tenantAllowed(tid string) bool {
	allow := strings.TrimSpace(os.Getenv("MS_ALLOWED_TENANTS")) // e.g. "aaaaaaaa-...,bbbb-..."
	if allow == "" {
		// org-only frontends will already block consumers; leave empty to allow any org tenant
		return true
	}
	for _, t := range strings.Split(allow, ",") {
		if strings.EqualFold(strings.TrimSpace(t), tid) {
			return true
		}
	}
	return false
}

func MicrosoftLogin(c *gin.Context) {
	var req MicrosoftSignInRequest
	appDomain := os.Getenv("APP_DOMAIN")
	clientID := os.Getenv("MS_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth config error"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// 1) Get tid (unverified) to select the correct issuer
	tid, err := extractTid(req.IdToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Request"})
		return
	}

	// Optional: block personal Microsoft accounts by rejecting the "consumers" tenant.
	// Personal accounts tenant GUID: 9188040d-6c67-4c5b-b112-36a304b66dad
	if strings.EqualFold(tid, "9188040d-6c67-4c5b-b112-36a304b66dad") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Personal accounts not allowed"})
		return
	}

	// 2) Build/cached verifier for that tenant and verify (issuer, audience, signature)
	verifier, err := getVerifierForTenant(*c, tid, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth config error"})
		return
	}
	idt, err := verifier.Verify(c.Request.Context(), req.IdToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Request"})
		return
	}

	// 3) Read claims (from a verified token)
	var claims struct {
		Email             string `json:"email"`
		PreferredUsername string `json:"preferred_username"`
		UPN               string `json:"upn"`
		Name              string `json:"name"`
		Oid               string `json:"oid"`
		Tid               string `json:"tid"`
		Iss               string `json:"iss"`
	}
	if err := idt.Claims(&claims); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to read token claims"})
		return
	}
	if !strings.EqualFold(claims.Tid, tid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Request"})
		return
	}

	// 4) Enforce allow-list if provided
	if !tenantAllowed(claims.Tid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tenant not allowed"})
		return
	}

	// 5) Resolve email
	email := claims.Email
	if email == "" {
		email = claims.PreferredUsername
	}
	if email == "" {
		email = claims.UPN
	}
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	// 6) Your existing login flow
	id, userId, name, dept, year, yearGroup, authToken, role, svcErr := service.LoginService(email, c)
	if svcErr != nil {
		if errors.Is(svcErr, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to login, please try again later"})
		}
		return
	}

	token := ""
	if authToken == "" {
		token, err = utils.GenerateJWT(name, email, id, userId, dept, year, yearGroup, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to login, please try again later"})
			return
		}
	} else {
		token = authToken
	}

	service.UpdateAuthToken(id, token)

	// NOTE: For local HTTP testing, you may need secure=false.
	c.SetCookie("PS", token, 10800, "/", appDomain, true, true)

	c.JSON(http.StatusOK, gin.H{"Message": "Login Success", "name": name})
}
