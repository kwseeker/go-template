package main

import (
	"log"
)

/*
	滑动窗口解决查找最长不重复子串
*/

// 优化后的答案，将用哈希记录字符是否出现改成用哈希记录字符出现的位置
func lengthOfLongestSubstringPro(s string) (string, int) {
	//key: 字符的Ascii码，value:字符上次出现的位置
	m := make(map[byte]int, 128)
	var i byte
	for i = 0; i < 128; i++ {
		m[i] = -1 //-1表示没出现过
	}
	n := len(s)

	iMax, lk, ans := 0, 0, 0 //lk: 左指针
	for i := 0; i < n; i++ { //i: 右指针
		log.Printf("lk=%d, i=%d\n", lk, i)
		//字符ascii码作为索引
		ascii := s[i]
		//m[ascii]+1: 下轮查询起始位置（没有重复值应该返回0（因为正好设置默认值-1,+1正好得到0）, 有重复值应该返回上次出现索引的下一索引）
		lkShift := false
		if lk < m[ascii]+1 {
			lk = m[ascii] + 1
			log.Println("lk=", lk)
			lkShift = true
		}
		if i-lk+1 > ans {
			if lkShift {
				iMax = lk - 1
			} else {
				iMax = lk
			}
			ans = i - lk + 1
		}
		m[ascii] = i
	}
	subStr := s[iMax : iMax+ans]
	return subStr, ans
}

// 官方的答案
// 	有地方重复查询了，比如字符串 "vgsqwxyzwzaw"，
//	第一轮：
//	左指针 i=0, 右指针 rk+1=0/1/2/3/4/5/7/8，查到第8位时发现有重复
//	第二轮：
//	i重新指向重复元素上次位置4（i=4）, 右指针这时是从5开始继续查重复元素，这里其实可以直接从9开始的（4-8）的元素不可能有重复
// 看下面的日志比较清楚：
// 	2022/09/16 17:07:50 i= 0
//	2022/09/16 17:07:50 rk+1= 0
//	2022/09/16 17:07:50 rk+1= 1
//	2022/09/16 17:07:50 rk+1= 2
//	2022/09/16 17:07:50 rk+1= 3
//	2022/09/16 17:07:50 rk+1= 4
//	2022/09/16 17:07:50 rk+1= 5
//	2022/09/16 17:07:50 rk+1= 6
//	2022/09/16 17:07:50 rk+1= 7
//	2022/09/16 17:07:50 i= 1		//接下来这四次循环，是为了在新一轮开始前清除m
//	2022/09/16 17:07:50 i= 2
//	2022/09/16 17:07:50 i= 3
//	2022/09/16 17:07:50 i= 4
//	2022/09/16 17:07:50 i= 5		//上面Pro()可以直接跳过了上面四次循环
//	2022/09/16 17:07:50 rk+1= 8
//	2022/09/16 17:07:50 i= 6
//	2022/09/16 17:07:50 i= 7
//	2022/09/16 17:07:50 i= 8
//	2022/09/16 17:07:50 rk+1= 9
//	2022/09/16 17:07:50 rk+1= 10
//	2022/09/16 17:07:50 i= 9
//	2022/09/16 17:07:50 rk+1= 11
//	2022/09/16 17:07:50 i= 10
//	2022/09/16 17:07:50 i= 11
func lengthOfLongestSubstring(s string) (string, int) {
	// 哈希集合，记录每个字符是否出现过
	m := map[byte]int{}
	n := len(s)
	// 右指针，初始值为 -1，相当于我们在字符串的左边界的左侧，还没有开始移动	rk: 右指针
	iMax, rk, ans := 0, -1, 0
	for i := 0; i < n; i++ { //	i： 左指针
		log.Println("i=", i)
		if i != 0 {
			// 左指针向右移动一格，移除一个字符
			delete(m, s[i-1])
		}
		for rk+1 < n && m[s[rk+1]] == 0 { //这里存在没必要的判断
			// 不断地移动右指针
			log.Println("rk+1=", rk+1)
			m[s[rk+1]]++
			rk++
		}
		// 第 i 到 rk 个字符是一个极长的无重复字符子串
		if ans < rk-i+1 {
			iMax = i
			ans = rk - i + 1
		}
	}
	subStr := s[iMax : iMax+ans]
	return subStr, ans
}

func main() {
	//randomStringGenerator := util.DefaultGenerator()
	//s := randomStringGenerator.RandomString()
	//log.Println("测试字符串: ", s)

	var subStr string
	var length int
	//subStr, length = lengthOfLongestSubstring("vgsqwxyzwzaw")
	subStr, length = lengthOfLongestSubstringPro("vgsqwxyzwzaw")
	log.Printf("最长不重复子串：%s， 长度: %d\n", subStr, length)
}
