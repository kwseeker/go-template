package lc0049

import (
	"sort"
)

/**
49. 字母异位词分组

给你一个字符串数组，请你将 字母异位词 组合在一起。可以按任意顺序返回结果列表。
字母异位词 是由重新排列源单词的所有字母得到的一个新单词。

字母异位词就是组成字符相同，但是顺序不一样

示例 1:
输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]
输出: [["bat"],["nat","tan"],["ate","eat","tea"]]

示例 2:
输入: strs = [""]
输出: [[""]]

示例 3:
输入: strs = ["a"]
输出: [["a"]]

提示：
1 <= strs.length <= 104
0 <= strs[i].length <= 100
strs[i] 仅包含小写字母
*/

//解法1： 排序, 创建一个map，按字典序重新排列的字符串作为key, 字符串先排序然后判断map是否存在，存在则添加到数组，不存在则添加数组
func groupAnagrams(strs []string) [][]string {
	//长度相同、组成字符相同、各种字符数量也要相同

	//1 存储按字典序排列的字母异位词
	group := make(map[string][]string)

	for _, word := range strs {
		//word按字典序重新排序
		//charSlice := strings.Split(word, "")
		//sort.Strings(charSlice)
		//newWord := strings.Join(charSlice, "")
		//性能更好的写法
		s := []byte(word)
		sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
		newWord := string(s)

		//if sa, ok := group[newWord]; ok {
		//	sa = append(sa, word)
		//	group[newWord] = sa
		//} else {
		//	group[newWord] = []string{word}
		//}
		//上面的可以简写为
		group[newWord] = append(group[newWord], word)
	}

	resMatrix := make([][]string, 0)
	for _, sa := range group {
		resMatrix = append(resMatrix, sa)
	}
	return resMatrix
}

//解法1： 排序, 创建一个map，按字典序重新排列的字符串作为key, 字符串先排序然后判断map是否存在，存在则添加到数组，不存在则添加数组
//func groupAnagrams(strs []string) [][]string {
//	//长度相同、组成字符相同、各种字符数量也要相同
//
//	//1 存储按字典序排列的字母异位词
//	group := make(map[string][]string)
//
//	for _, word := range strs {
//		//word按字典序重新排序
//		charSlice := strings.Split(word, "")
//		sort.Strings(charSlice)
//		newWord := strings.Join(charSlice, "")
//
//		if sa, ok := group[newWord]; ok {
//			sa = append(sa, word)
//			group[newWord] = sa
//		} else {
//			group[newWord] = []string{word}
//		}
//	}
//
//	resMatrix := make([][]string, 0)
//	for _, sa := range group {
//		resMatrix = append(resMatrix, sa)
//	}
//	return resMatrix
//}

//解法1： 排序, 创建一个字符串数组，保存按字典序重新排列的字符串
//func groupAnagrams(strs []string) [][]string {
//	//长度相同、组成字符相同、各种字符数量也要相同
//
//	//1 存储按字典序排列的字母异位词
//	var models []string
//	//2 保存分组结果
//	var resMatrix [][]string
//
//	for _, word := range strs {
//		//word按字典序重新排序
//		charSlice := strings.Split(word, "")
//		sort.Strings(charSlice)
//		newWord := strings.Join(charSlice, "")
//		//是否已经存在相同字母异位词的分组
//		isExist := false
//		for i, modelsWord := range models { //这个for可以被优化掉
//			if newWord == modelsWord {
//				//已经存在
//				resMatrix[i] = append(resMatrix[i], word)
//				isExist = true
//				break
//			}
//		}
//		if !isExist { //不存在相同字母异位词的分组，新建一个分组
//			models = append(models, newWord)
//			resMatrix = append(resMatrix, []string{word})
//		}
//	}
//
//	//将所有分组结果放入resMatrix中返回
//	return resMatrix
//}
