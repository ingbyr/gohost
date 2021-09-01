/*
 @Author: ingbyr
*/

package host

func sub(s1 []string, s2 []string) (res, removed []string) {
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

func union(s1 []string, s2 []string) (res []string, add []string) {
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

func cache(s []string) map[string]struct{} {
	res := make(map[string]struct{}, len(s))
	for _, item := range s {
		res[item] = struct{}{}
	}
	return res
}
