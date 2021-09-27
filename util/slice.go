/*
 @Author: ingbyr
*/

package util

func SliceSub(s1 []string, s2 []string) ([]string, []string) {
	res, removed := make([]string, 0), make([]string, 0)
	s2Cache := cache(s2)
	for _, s := range s1 {
		if _, exist := s2Cache[s]; exist {
			removed = append(removed, s)
		} else {
			res = append(res, s)
		}
	}
	return res, removed
}

func SliceUnion(s1 []string, s2 []string) ([]string, []string) {
	res, add := make([]string, 0), make([]string, 0)
	s1Cache := cache(s1)
	for _, s1i := range s1 {
		res = append(res, s1i)
	}
	for _, s2i := range s2 {
		if _, exist := s1Cache[s2i]; !exist {
			res = append(res, s2i)
			add = append(add, s2i)
		}
	}
	return res, add
}

func SortUniqueStringSlice(arr []string) []string {
	last := len(arr)
	for i := 0; i < last; i++ {
		min := i
		for j := i; j < last; j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		if i > 0 && arr[min] == arr[i-1] {
			last--
			// resort index at i
			i--
			arr[min], arr[last] = arr[last], arr[min]
		} else if min != i {
			arr[i], arr[min] = arr[min], arr[i]
		}
	}
	return arr[:last]
}

func SliceRemove(arr []string, target string) ([]string, bool) {
	i := 0
	for j :=  0; i < len(arr) && j < len(arr); i++ {
		if j < i {
			j = i
		}
		if arr[i] == target {
			// find next not target index
			for j < len(arr) && arr[j] == target {
				j++
			}
			// find nothing
			if j == len(arr) {
				break
			}
			arr[i], arr[j] = arr[j], arr[i]
			j++
		}
	}
	return arr[:i], i < len(arr)
}

func cache(s []string) map[string]struct{} {
	res := make(map[string]struct{}, len(s))
	for _, item := range s {
		res[item] = struct{}{}
	}
	return res
}
