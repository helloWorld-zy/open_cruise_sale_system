package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

// RBAC represents the RBAC system
type RBAC struct {
	enforcer *casbin.Enforcer
}

// NewRBAC creates a new RBAC instance with Casbin
func NewRBAC() (*RBAC, error) {
	// Define RBAC model
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	// Load default policies
	if err := loadDefaultPolicies(e); err != nil {
		return nil, err
	}

	return &RBAC{enforcer: e}, nil
}

// loadDefaultPolicies sets up default RBAC policies
func loadDefaultPolicies(e *casbin.Enforcer) error {
	// Super Admin - full access
	policies := [][]string{
		{"super_admin", "*", "*"},
		{"operations", "cruises", "read"},
		{"operations", "cruises", "write"},
		{"operations", "cabin-types", "read"},
		{"operations", "cabin-types", "write"},
		{"operations", "facilities", "read"},
		{"operations", "facilities", "write"},
		{"operations", "routes", "read"},
		{"operations", "routes", "write"},
		{"operations", "voyages", "read"},
		{"operations", "voyages", "write"},
		{"operations", "cabins", "read"},
		{"operations", "cabins", "write"},
		{"operations", "orders", "read"},
		{"operations", "analytics", "read"},
		{"finance", "orders", "read"},
		{"finance", "orders", "refund"},
		{"finance", "payments", "read"},
		{"finance", "refund-requests", "read"},
		{"finance", "refund-requests", "approve"},
		{"finance", "analytics", "read"},
		{"finance", "reconciliation", "read"},
		{"customer_service", "orders", "read"},
		{"customer_service", "orders", "update"},
		{"customer_service", "users", "read"},
		{"customer_service", "refund-requests", "read"},
		{"customer_service", "refund-requests", "process"},
	}

	_, err := e.AddPolicies(policies)
	if err != nil {
		return err
	}

	groupingPolicies := [][]string{
		{"operations", "super_admin"},
		{"finance", "super_admin"},
		{"customer_service", "super_admin"},
		{"finance", "operations"},
	}

	_, err = e.AddGroupingPolicies(groupingPolicies)
	return err
}

// CheckPermission checks if a subject has permission
func (r *RBAC) CheckPermission(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}

// AddRoleForUser adds a role to a user
func (r *RBAC) AddRoleForUser(user, role string) (bool, error) {
	return r.enforcer.AddGroupingPolicy(user, role)
}

// GetRolesForUser gets all roles for a user
func (r *RBAC) GetRolesForUser(user string) ([]string, error) {
	return r.enforcer.GetRolesForUser(user)
}
