// Copyright 2026 Agenova contributors.
// SPDX-License-Identifier: Apache-2.0

// Package authorization evaluates assignment creation before any authority,
// claim, or runtime side effect is allowed to occur.
package authorization

import (
	"fmt"
	"strings"

	v1alpha1 "github.com/wunderforge/agenova/api/v1alpha1"
	"github.com/wunderforge/agenova/internal/policy"
)

// Request is the trusted, backend-neutral context for assignment admission.
// Principal is supplied out-of-band; Action is derived from validated request
// references and does not itself grant authority.
type Request struct {
	RequestRef string
	Principal  v1alpha1.Principal
	Action     v1alpha1.Action
}

// BundleSource exposes the active immutable policy without coupling admission
// to policy loading or administration.
type BundleSource interface {
	Current() (policy.PolicyBundle, bool)
}

// Evaluator allows Gate to enforce the pre-side-effect ordering boundary.
type Evaluator interface {
	Evaluate(Request) (v1alpha1.Decision, error)
}

// Admission is an internal capability proving that Gate admitted one exact
// assignment request. Its binding fields are intentionally private so callers
// cannot manufacture an Allow token from a public Decision value.
type Admission struct {
	request  Request
	decision v1alpha1.Decision
}

// Matches reports whether this admission belongs to the request context that
// a downstream authority resolver is about to process.
func (a Admission) Matches(requestRef, projectRef, templateRef string) bool {
	return a.decision.Result == v1alpha1.DecisionResultAllow &&
		a.request.RequestRef == requestRef &&
		a.request.Action.Name == "claim.create" &&
		a.request.Action.Project == projectRef &&
		a.request.Action.TemplateRef == templateRef
}

// Decision returns the evidence-ready decision that produced this admission.
func (a Admission) Decision() v1alpha1.Decision {
	return a.decision
}

// Authorizer performs exact-match, default-deny assignment admission.
type Authorizer struct {
	Policies BundleSource
}

// Evaluate returns one evidence-ready decision. Invalid trusted/request
// context is rejected before policy evaluation; absence or mismatch is Deny.
func (a Authorizer) Evaluate(input Request) (v1alpha1.Decision, error) {
	if err := validate(input); err != nil {
		return v1alpha1.Decision{}, err
	}

	decision := v1alpha1.Decision{
		ID:           fmt.Sprintf("decision:%s:authorization", input.RequestRef),
		PrincipalRef: input.Principal.Subject,
		Action:       input.Action.Name,
		Result:       v1alpha1.DecisionResultDeny,
	}
	if a.Policies == nil {
		decision.Reason = "no active policy bundle"
		return decision, nil
	}
	bundle, ok := a.Policies.Current()
	if !ok {
		decision.Reason = "no active policy bundle"
		return decision, nil
	}
	decision.PolicyRef = v1alpha1.PolicyReference{ID: bundle.ID, Version: bundle.Version}

	matched := bundle.Allows(policy.Match{
		Team:        input.Principal.Team,
		Action:      input.Action.Name,
		Project:     input.Action.Project,
		TemplateRef: input.Action.TemplateRef,
	})
	if !matched {
		decision.Reason = "no exact policy rule matched the trusted principal and requested assignment"
		return decision, nil
	}

	decision.Result = v1alpha1.DecisionResultAllow
	decision.Reason = "exact policy rule matched the trusted principal and requested assignment"
	return decision, nil
}

// Gate invokes the authorized continuation exactly once only for Allow.
type Gate struct {
	Evaluator Evaluator
}

// Admit evaluates input and stops all downstream work unless the result is
// exactly Allow. ApprovalRequired is intentionally not treated as authority.
func (g Gate) Admit(input Request, onAllowed func(Admission) error) (v1alpha1.Decision, error) {
	if g.Evaluator == nil {
		return v1alpha1.Decision{}, required("evaluator")
	}
	decision, err := g.Evaluator.Evaluate(input)
	if err != nil || decision.Result != v1alpha1.DecisionResultAllow {
		return decision, err
	}
	if onAllowed == nil {
		return decision, required("onAllowed")
	}
	if err := validateAllowedDecision(input, decision); err != nil {
		return decision, err
	}
	if err := onAllowed(Admission{request: input, decision: decision}); err != nil {
		return decision, err
	}
	return decision, nil
}

func validateAllowedDecision(input Request, decision v1alpha1.Decision) error {
	if decision.ID == "" {
		return required("decision.id")
	}
	if decision.PrincipalRef != input.Principal.Subject {
		return invalid("decision.principalRef", "must match the admitted principal")
	}
	if decision.Action != input.Action.Name {
		return invalid("decision.action", "must match the admitted action")
	}
	if decision.PolicyRef.ID == "" || decision.PolicyRef.Version == "" {
		return required("decision.policyRef")
	}
	if decision.Reason == "" {
		return required("decision.reason")
	}
	return nil
}

func validate(input Request) error {
	fields := []struct {
		path  string
		value string
	}{
		{"requestRef", input.RequestRef},
		{"principal.subject", input.Principal.Subject},
		{"principal.team", input.Principal.Team},
		{"principal.authenticationContext", input.Principal.AuthenticationContext},
		{"action.name", input.Action.Name},
		{"action.project", input.Action.Project},
		{"action.templateRef", input.Action.TemplateRef},
	}
	for _, field := range fields {
		if strings.TrimSpace(field.value) == "" {
			return required(field.path)
		}
	}
	return nil
}

func required(path string) *v1alpha1.ValidationError {
	return &v1alpha1.ValidationError{
		Category:  v1alpha1.ValidationCategoryRequiredField,
		FieldPath: path,
		Detail:    "value is required",
	}
}

func invalid(path, detail string) *v1alpha1.ValidationError {
	return &v1alpha1.ValidationError{
		Category:  v1alpha1.ValidationCategoryInvalidValue,
		FieldPath: path,
		Detail:    detail,
	}
}
