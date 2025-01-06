/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// Package fst implements a functional toolkit designed to build hierarchical
// and treelike structures in a composable manner.
package fst

type Builder[T any, N Node[N, A], A any] struct {
	Build func(T) N
}

type Node[N any, A any] interface {
	Append(N)
	Annotate(A)
}

type Annotation interface{}

func New[T any, N Node[N, A], A any](build func(T) N) *Builder[T, N, A] {
	return &Builder[T, N, A]{
		build,
	}
}

// Append adds one or more children to the existing tree.
func (b *Builder[T, N, A]) Append(children ...*Builder[T, N, A]) *Builder[T, N, A] {
	return b.Tap(func(x T, parent N) {
		for _, child := range children {
			parent.Append(child.Build(x))
		}
	})
}

// Annotate adds one of more annotations to the current node.
func (b *Builder[T, N, A]) Annotate(annotations ...A) *Builder[T, N, A] {
	return b.Tap(func(x T, parent N) {
		for _, annotation := range annotations {
			parent.Annotate(annotation)
		}
	})
}

// Lift allows for the dynamic insertion nodes into a tree.
func (b *Builder[T, N, A]) Lift(f func(context T) *Builder[T, N, A]) *Builder[T, N, A] {
	return b.Tap(func(x T, parent N) {
		parent.Append(f(x).Build(x))
	})
}

// Tap is a utility method that abstracts common behaviour required by
// [Builder.Append], [Builder.Annotate], and [Builder.Lift].
func (b *Builder[T, N, A]) Tap(f func(context T, parent N)) *Builder[T, N, A] {
	return New[T](func(x T) N {
		result := b.Build(x)
		f(x, result)
		return result
	})
}

// Scope maps the execution context from a parent tree, to be accepted by a
// subtree.
func Scope[T1, T2 any, N Node[N, A], A any](b *Builder[T2, N, A], f func(T1) T2) *Builder[T1, N, A] {
	return New[T1, N, A](func(x T1) N {
		return b.Build(f(x))
	})
}
