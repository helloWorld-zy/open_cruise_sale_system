package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Role constantsconst (
	RoleSuperAdmin       = "super_admin"
	RoleOperations       = "operations"
	RoleFinance          = "finance"
	RoleCustomerService  = "customer_service"
)

// AllRoles returns all available rolesvar AllRoles = []string{
	RoleSuperAdmin,
	RoleOperations,
	RoleFinance,
	RoleCustomerService,
}

// IsValidRole checks if a role is validfunc IsValidRole(role string) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}

// HashPassword creates a bcrypt hash of the passwordfunc HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hashfunc CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateRandomToken generates a cryptographically secure random tokenfunc GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Permission constants for better type safetyconst (
	PermRead   = "read"
	PermWrite  = "write"
	PermDelete = "delete"
	PermUpdate = "update"
	PermRefund = "refund"
	PermApprove = "approve"
	PermProcess = "process"
)

// Resource constantsconst (
	ResourceCruises        = "cruises"
	ResourceCabinTypes     = "cabin-types"
	ResourceFacilities     = "facilities"
	ResourceRoutes         = "routes"
	ResourceVoyages        = "voyages"
	ResourceCabins         = "cabins"
	ResourceOrders         = "orders"
	ResourcePayments       = "payments"
	ResourceRefundRequests = "refund-requests"
	ResourceUsers          = "users"
	ResourceAnalytics      = "analytics"
	ResourceReconciliation = "reconciliation"
)
