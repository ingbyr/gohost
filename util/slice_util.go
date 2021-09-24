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

func SliceRemove(arr []string, target string) []string {
	size := len(arr)
	for i, nq := 0, 0; i < size && nq < len(arr); i++ {
		if nq <= i {
			nq = i
		}
		if arr[i] == target {
			// find next not target index
			for nq < len(arr) && arr[nq] == target {
				nq++
			}
			// find nothing
			if nq == len(arr) {
				size = i
				break
			}
			arr[i], arr[nq] = arr[nq], arr[i]
			size--
		}
	}
	return arr[:size]
}

func cache(s []string) map[string]struct{} {
	res := make(map[string]struct{}, len(s))
	for _, item := range s {
		res[item] = struct{}{}
	}
	return res
}
