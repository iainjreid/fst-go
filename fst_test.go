/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fst_test

import (
	"reflect"
	"testing"

	"github.com/iainjreid/fst-go"
)

type TestNode struct {
	name        string
	children    []*TestNode
	annotations []*TestAnnotation
}

type TestAnnotation struct {
	key   string
	value string
}

func (t *TestNode) Append(node *TestNode) {
	t.children = append(t.children, node)
}

func (t *TestNode) Annotate(annotation *TestAnnotation) {
	t.annotations = append(t.annotations, annotation)
}

var parent = fst.New(func(string) *TestNode {
	return &TestNode{
		name: "parent",
	}
})

// TestAppend calls [Builder.Append], to ensure that it correctly adds a new
// child to the parent node.
func TestAppend(t *testing.T) {
	var subject = parent.Append(fst.New(func(string) *TestNode {
		return &TestNode{
			name: "child",
		}
	}))

	var expected = &TestNode{
		name: "parent",
		children: []*TestNode{
			{
				name: "child",
			},
		},
	}

	if !reflect.DeepEqual(subject.Build("68yvwz"), expected) {
		t.Fatal("result should be equal to expected output")
	}
}

// TestAnnotate calls [Buidler.Annotate], to ensure that it correctly adds an
// annotation to the parent node.
func TestAnnotate(t *testing.T) {
	var subject = parent.Annotate(&TestAnnotation{
		key:   "test_key",
		value: "test_value",
	})

	var expected = &TestNode{
		name: "parent",
		annotations: []*TestAnnotation{
			{
				key:   "test_key",
				value: "test_value",
			},
		},
	}

	if !reflect.DeepEqual(subject.Build("q708uc"), expected) {
		t.Fatal("result should be equal to expected output")
	}
}

// TestLift calls [Builder.Lift], to ensure that dynamic nodes can be added
// using the provided build context.
func TestLift(t *testing.T) {
	var subject = parent.Lift(func(i string) *fst.Builder[string, *TestNode, *TestAnnotation] {
		return fst.New(func(str string) *TestNode {
			return &TestNode{
				name: str,
			}
		})
	})

	var expected = &TestNode{
		name: "parent",
		children: []*TestNode{
			{
				name: "x8azmu",
			},
		},
	}

	if !reflect.DeepEqual(subject.Build("x8azmu"), expected) {
		t.Fatal("result should be equal to expected output")
	}
}
