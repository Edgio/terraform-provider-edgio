package utility

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringSliceToTypesList(slice []string) types.List {
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
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

func IntSliceToTypesList(slice []int) types.List {
	elements := make([]attr.Value, len(slice))
	for i, v := range slice {
		elements[i] = types.Int64Value(int64(v))
	}

	list, _ := types.ListValue(types.StringType, elements)
	return list
}

func TypesListToIntSlice(list types.List) []int {
	if list.IsNull() || list.IsUnknown() {
		return []int{}
	}

	var result []int
	list.ElementsAs(context.Background(), &result, false)
	return result
}

func StringMapToMapValue(m map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
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
