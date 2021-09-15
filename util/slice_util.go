/*
 @Author: ingbyr
*/

package util

func SliceSub(s1 []string, s2 []string) (res, removed []string) {
	s2Cache := cache(s2)
	for _, s := range s1 {
		if _, exist := s2Cache[s]; exist {
			removed = append(removed, s)
		} else {
			res = append(res, s)
		}
	}
	return
}

func SliceUnion(s1 []string, s2 []string) (res []string, add []string) {
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
	return
}

func SortUniqueStringSlice(arr []string) []string {
	last := len(arr)
	for i := 0; i < last; i++ {
		m := i
		for j := i; j < last; j++ {
			if arr[j] < arr[m] {
				m = j
			}
		}
		if i > 0 && arr[m] == arr[i-1] {
			last--
			i--
			arr[m], arr[last] = arr[last], arr[m]
		} else if m != i {
			arr[i], arr[m] = arr[m], arr[i]
		}
	}
	return arr[:last]
}

func cache(s []string) map[string]struct{} {
	res := make(map[string]struct{}, len(s))
	for _, item := range s {
		res[item] = struct{}{}
	}
	return res
}
