// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"context"
	"fmt"

	apiv2 "github.com/operator-framework/api/pkg/operators/v2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// ErrNoOperatorCondition indicates that the operator condition CRD is nil
	ErrNoOperatorCondition = fmt.Errorf("operator Condition CRD is nil")
)

type Condition interface {
	// Get fetches the condition on the operator's
	// OperatorCondition. It returns an error if there are problems getting
	// the OperatorCondition object or if the specific condition type does not
	// exist.
	Get(ctx context.Context) (*metav1.Condition, error)

	// Set sets the specific condition on the operator's
	// OperatorCondition to the provided status. If the condition is not
	// present, it is added to the CR.
	// To set a new condition, the user can call this method and provide optional
	// parameters if required. It returns an error if there are problems getting or
	// updating the OperatorCondition object.
	Set(ctx context.Context, status metav1.ConditionStatus, option ...Option) error
}

// Option is a function that applies a change to a condition.
// This can be used to set optional condition fields, like reasons
// and messages.
type Option func(*metav1.Condition)

// condition is a Condition that gets and sets a specific
// conditionType in the OperatorCondition CR.
type condition struct {
	namespacedName types.NamespacedName
	condType       apiv2.ConditionType
	client         client.Client
}

var _ Condition = &condition{}

// Get implements conditions.Get
func (c *condition) Get(ctx context.Context) (*metav1.Condition, error) {
	operatorCond := &apiv2.OperatorCondition{}
	err := c.client.Get(ctx, c.namespacedName, operatorCond)
	if err != nil {
		return nil, err
	}
	con := meta.FindStatusCondition(operatorCond.Spec.Conditions, string(c.condType))

	if con == nil {
		return nil, fmt.Errorf("conditionType %v not found", c.condType)
	}
	return con, nil
}

// Set implements conditions.Set
func (c *condition) Set(ctx context.Context, status metav1.ConditionStatus, option ...Option) error {
	operatorCond := &apiv2.OperatorCondition{}
	err := c.client.Get(ctx, c.namespacedName, operatorCond)
	if err != nil {
		return err
	}

	newCond := &metav1.Condition{
		Type:   string(c.condType),
		Status: status,
	}

	for _, opt := range option {
		opt(newCond)
	}
	meta.SetStatusCondition(&operatorCond.Spec.Conditions, *newCond)
	return c.client.Update(ctx, operatorCond)
}

func NewCondition(client client.Client, namespacedName types.NamespacedName, condType apiv2.ConditionType) Condition {
	return &condition{
		client:         client,
		namespacedName: namespacedName,
		condType:       condType,
	}
}
