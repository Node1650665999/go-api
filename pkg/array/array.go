package array

import (
	"fmt"
	"math"
	"reflect"
)

//MapToSlice map转数组
func MapToSlice(mp interface{}) interface{} {
	mpType := reflect.TypeOf(mp)
	if mpType.Kind() != reflect.Map {
		panic("the mp type must be map")
	}

	valOf := reflect.ValueOf(mp)
	length := valOf.Len()
	if length == 0 {
		return nil
	}


	//获得slice对应类型
	sliceT := reflect.SliceOf(mpType.Elem())
	sliceV := reflect.MakeSlice(sliceT, 0, length)

	keys := valOf.MapKeys()
	for _, k := range keys {
		value := valOf.MapIndex(k)
		sliceV = reflect.Append(sliceV, value)
	}

	return sliceV.Interface()
}

//MapToSliceWithKeys map转数组并以给定的 keys 排序
func MapToSliceWithKeys(mp interface{}, sortKeys interface{}) interface{} {
	mpType := reflect.TypeOf(mp)
	if mpType.Kind() != reflect.Map {
		panic("the mp type must be map")
	}

	valOf := reflect.ValueOf(mp)
	length := valOf.Len()
	if length == 0 {
		return nil
	}

	//获得slice对应类型
	sliceT := reflect.SliceOf(mpType.Elem())
	sliceV := reflect.MakeSlice(sliceT, 0, length)

	//保持顺序
	sortKeysValOf := reflect.ValueOf(sortKeys)
	sortKeysLen := sortKeysValOf.Len()
	for i := 0; i < sortKeysLen; i++ {
		sortKey := sortKeysValOf.Index(i).Interface()
		value := valOf.MapIndex(reflect.ValueOf(sortKey))
		//判断key是否存在
		if !value.IsValid() {
			continue
		}

		sliceV = reflect.Append(sliceV, value)
	}

	return sliceV.Interface()
}

//SliceToMap 根据传入的key将数组转换为map
func SliceToMap(slice interface{}, keyName string) interface{} {
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		panic("the slice type must be slice")
	}

	valOf := reflect.ValueOf(slice)
	length := valOf.Len()
	if length == 0 {
		return nil
	}

	//检测 keyName 是否存在
	keyValOf := reflect.Indirect(valOf.Index(0)).FieldByName(keyName)
	if !keyValOf.IsValid() {
		panic(fmt.Sprintf("key=%v doesn't exist", keyName))
	}

	// 获得map对应类型
	mapT := reflect.MapOf(keyValOf.Type(), valOf.Index(0).Type())
	mapV := reflect.MakeMap(mapT)
	for i := 0; i < length; i++ {
		value := valOf.Index(i)
		//if reflect.Indirect(value).Kind() != reflect.Struct {
		//	panic("value type must be struct")
		//}
		mapV.SetMapIndex(reflect.Indirect(valOf.Index(i)).FieldByName(keyName), value)
	}
	return mapV.Interface()
}

//SlicePaging 实现数组分页
func SlicePaging(slice interface{}, page, pageSize int) (data interface{}, paging map[string]int) {
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		panic("the slice type must be slice")
	}

	valOf := reflect.ValueOf(slice)
	//总条数
	total := valOf.Len()
	if total == 0 {
		return nil, nil
	}

	//边界处理：第一页
	if page < 0 {
		page = 1
	}

	//边界处理：每页数量
	switch {
	case pageSize > 200:
		pageSize = 200
	case pageSize <= 0:
		pageSize = 20
	}

	//总页数
	pageCount := int(math.Ceil(float64(total) / float64(pageSize)))
	if page >= pageCount {
		page = pageCount
	}

	//下一页,如果没有,直接nextPage = 0
	nextPage := page + 1
	if nextPage > pageCount {
		nextPage = 0
	}

	//offset and limit
	sliceStart := (page - 1) * pageSize
	sliceEnd := sliceStart + pageSize
	if sliceEnd > total {
		sliceEnd = total
	}

	paging = map[string]int{
		"page":       page,
		"next_page":  nextPage,
		"page_size":  pageSize,
		"total":      total,
		"page_count": pageCount,
	}

	newSlice := valOf.Slice(sliceStart, sliceEnd)

	return newSlice.Interface(), paging
}


//SliceColumn 从二维数组中查找某个值
func SliceColumn(arrMap interface{}, field interface{}) (values []string) {
	list := arrMap.([]map[string]interface{})
	for row := range list {
		for key, val := range list[row] {
			if key == field {
				values = append(values, fmt.Sprint(val))
			}
		}
	}
	return values
}

//SliceUnique 数组去重
func SliceUnique(tmp interface{}) interface{} {
	switch arr := tmp.(type) {
	case []string:
		result := make([]string, 0, len(arr))
		temp := map[string]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []int:
		result := make([]int, 0, len(arr))
		temp := map[int]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []int64:
		result := make([]int64, 0, len(arr))
		temp := map[int64]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []int32:
		result := make([]int32, 0, len(arr))
		temp := map[int32]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []float64:
		result := make([]float64, 0, len(arr))
		temp := map[float64]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []float32:
		result := make([]float32, 0, len(arr))
		temp := map[float32]struct{}{}
		for _, item := range arr {
			if _, ok := temp[item]; !ok {
				temp[item] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	}
	return nil
}

//InSlice 数组中是否包含某个元素
func InSlice(need interface{}, haystack interface{}) bool {
	switch key := need.(type) {
	case int:
		for _, item := range haystack.([]int) {
			if item == key {
				return true
			}
		}
	case string:
		for _, item := range haystack.([]string) {
			if item == key {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if item == key {
				return true
			}
		}
	case int32:
		for _, item := range haystack.([]int32) {
			if item == key {
				return true
			}
		}
	case float64:
		for _, item := range haystack.([]float64) {
			if item == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}



