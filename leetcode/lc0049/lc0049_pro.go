package lc0049

//最高性能的题解, 用的 字符->索引 的哈希表作为key, 比直接使用 string 作为key性能更好，因为string比较没有哈希表比较性能好
func groupAnagramsPro(strs []string) [][]string {
	//长度相同、组成字符相同、各种字符数量也要相同

	m := make(map[[26]byte][]string)
	for i := len(strs) - 1; i >= 0; i-- {
		sum := [26]byte{}
		for j := len(strs[i]) - 1; j >= 0; j-- {
			sum[strs[i][j]-'a']++
		}

		m[sum] = append(m[sum], strs[i])
	}

	ans := make([][]string, 0)
	for _, v := range m {
		ans = append(ans, v)
	}

	return ans
}
