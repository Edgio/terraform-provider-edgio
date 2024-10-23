package utility

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringSliceToTypesList(slice *[]string) types.List {
	if slice == nil {
		list, _ := types.ListValue(types.StringType, nil)
		return list
	}

	elements := make([]attr.Value, len(*slice))
	for i, v := range *slice {
		elements[i] = types.StringValue(v)
	}

	list, _ := types.ListValue(types.StringType, elements)
	return list
}

func TypesListToStringSlice(list types.List) []string {
	var result []string
	list.ElementsAs(context.Background(), &result, false)
	return result
}

func TypesListToStringSlicePointer(list types.List) *[]string {
	result := TypesListToStringSlice(list)
	return &result
}

func IntSliceToTypesList(slice *[]int64) types.List {
	if slice == nil {
		list, _ := types.ListValue(types.Int64Type, nil)
		return list
	}

	elements := make([]attr.Value, len(*slice))
	for i, v := range *slice {
		elements[i] = types.Int64Value(int64(v))
	}

	list, _ := types.ListValue(types.StringType, elements)
	return list
}

func TypesListToIntSlice(list types.List) []int64 {
	if list.IsNull() || list.IsUnknown() {
		return []int64{}
	}

	var result []int64
	list.ElementsAs(context.Background(), &result, false)
	return result
}

func TypesListToIntSlicePointer(list types.List) *[]int64 {
	result := TypesListToIntSlice(list)
	return &result
}

func StringMapToMapValue(m *map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(*m))
	for k, v := range *m {
		elements[k] = types.StringValue(v)
	}

	mapValue, _ := types.MapValue(types.StringType, elements)
	return mapValue
}

func MapValueToStringMap(m types.Map) map[string]string {
	elements := m.Elements()
	result := make(map[string]string, len(elements))

	for key, value := range elements {
		strValue, _ := value.(types.String)
		result[key] = strValue.ValueString()
	}

	return result
}

func MapValueToStringMapPointer(m types.Map) *map[string]string {
	result := MapValueToStringMap(m)
	return &result
}
