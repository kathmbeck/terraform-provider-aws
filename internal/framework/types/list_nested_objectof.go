// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
)

// ListNestedObjectTypeOf is the attribute type of a ListNestedObjectValueOf.
type ListNestedObjectTypeOf[T any] struct {
	basetypes.ListType
}

var _ basetypes.ListTypable = ListNestedObjectTypeOf[struct{}]{}

func NewListNestedObjectTypeOf[T any](ctx context.Context) ListNestedObjectTypeOf[T] {
	return ListNestedObjectTypeOf[T]{basetypes.ListType{ElemType: NewObjectTypeOf[T](ctx)}}
}

func (t ListNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(ListNestedObjectTypeOf[T])

	if !ok {
		return false
	}

	return t.ListType.Equal(other.ListType)
}

func (t ListNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("ListNestedObjectTypeOf[%T]", zero)
}

func (t ListNestedObjectTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewListNestedObjectValueOfNull[T](ctx), nil
	}
	if in.IsUnknown() {
		return NewListNestedObjectValueOfUnknown[T](ctx), nil
	}

	listValue, diags := basetypes.NewListValue(NewObjectTypeOf[T](ctx), in.Elements())

	if diags.HasError() {
		return NewListNestedObjectValueOfUnknown[T](ctx), diags
	}

	value := ListNestedObjectValueOf[T]{
		ListValue: listValue,
	}

	return value, nil
}

func (t ListNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func (t ListNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListNestedObjectValueOf[T]{}
}

// ListNestedObjectValueOf represents a Terraform Plugin Framework List value whose elements are of type ObjectTypeOf.
type ListNestedObjectValueOf[T any] struct {
	basetypes.ListValue
}

var _ basetypes.ListValuable = ListNestedObjectValueOf[struct{}]{}

func (v ListNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListNestedObjectValueOf[T])

	if !ok {
		return false
	}

	return v.ListValue.Equal(other.ListValue)
}

func (v ListNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewListNestedObjectTypeOf[T](ctx)
}

func NewListNestedObjectValueOfNull[T any](ctx context.Context) ListNestedObjectValueOf[T] {
	return ListNestedObjectValueOf[T]{ListValue: basetypes.NewListNull(NewObjectTypeOf[T](ctx))}
}

func NewListNestedObjectValueOfUnknown[T any](ctx context.Context) ListNestedObjectValueOf[T] {
	return ListNestedObjectValueOf[T]{ListValue: basetypes.NewListUnknown(NewObjectTypeOf[T](ctx))}
}

func NewListNestedObjectValueOf[T any](ctx context.Context, t *T) ListNestedObjectValueOf[T] {
	return NewListNestedObjectValueOfSlice(ctx, []*T{t})
}

func NewListNestedObjectValueOfSlice[T any](ctx context.Context, ts []*T) ListNestedObjectValueOf[T] {
	return newListNestedObjectValueOf[T](ctx, ts)
}

func NewListNestedObjectValueOfValueSlice[T any](ctx context.Context, ts []T) ListNestedObjectValueOf[T] {
	return newListNestedObjectValueOf[T](ctx, ts)
}

func newListNestedObjectValueOf[T any](ctx context.Context, elements any) ListNestedObjectValueOf[T] {
	return ListNestedObjectValueOf[T]{ListValue: fwdiag.Must(basetypes.NewListValueFrom(ctx, NewObjectTypeOf[T](ctx), elements))}
}
