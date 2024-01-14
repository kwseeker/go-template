package lc1577

/**
给你两个整数数组 nums1 和 nums2 ，请你返回根据以下规则形成的三元组的数目（类型 1 和类型 2 ）：

类型 1：三元组 (i, j, k) ，如果 nums1[i]2 == nums2[j] * nums2[k] 其中 0 <= i < nums1.length 且 0 <= j < k < nums2.length
类型 2：三元组 (i, j, k) ，如果 nums2[i]2 == nums1[j] * nums1[k] 其中 0 <= i < nums2.length 且 0 <= j < k < nums1.length

示例 1：

输入：nums1 = [7,4], nums2 = [5,2,8,9]
输出：1
解释：类型 1：(1,1,2), nums1[1]^2 = nums2[1] * nums2[2] (4^2 = 2 * 8)
示例 2：

输入：nums1 = [1,1], nums2 = [1,1,1]
输出：9
解释：所有三元组都符合题目要求，因为 1^2 = 1 * 1
类型 1：(0,0,1), (0,0,2), (0,1,2), (1,0,1), (1,0,2), (1,1,2), nums1[i]^2 = nums2[j] * nums2[k]
类型 2：(0,0,1), (1,0,1), (2,0,1), nums2[i]^2 = nums1[j] * nums1[k]
示例 3：

输入：nums1 = [7,7,8,3], nums2 = [1,2,9,7]
输出：2
解释：有两个符合题目要求的三元组
类型 1：(3,0,2), nums1[3]^2 = nums2[0] * nums2[2]
类型 2：(3,0,1), nums2[3]^2 = nums1[0] * nums1[1]
示例 4：

输入：nums1 = [4,7,9,11,23], nums2 = [3,5,1024,12,18]
输出：0
解释：不存在符合题目要求的三元组
*/
func numTriplets(nums1 []int, nums2 []int) int {
	//创建两个哈希表
	cnt1 := make(map[int]int)
	cnt2 := make(map[int]int)
	for _, v := range nums1 {
		cnt1[v]++
	}
	for _, v := range nums2 {
		cnt2[v]++
	}
	//获取两个哈希表的key数组
	//keys1 := make([]int, 0, len(cnt1))
	//for k := range cnt1 {
	//    keys1 = append(keys1, k)
	//}
	//keys2 := make([]int, 0, len(cnt2))
	//for k := range cnt2 {
	//    keys2 = append(keys2, k)
	//}

	res := 0
	//最直观的写法
	////类型1 计算结果
	//for _, v1 := range keys1 {
	//	for i := 0; i+1 < len(nums2); i++ {
	//		for j := i + 1; j < len(nums2); j++ {
	//			if v1*v1 == nums2[i]*nums2[j] {
	//				res += cnt1[v1]
	//			}
	//		}
	//	}
	//}
	////类型2 计算结果
	//for _, v1 := range keys2 {
	//	for i := 0; i+1 < len(nums1); i++ {
	//		for j := i + 1; j < len(nums1); j++ {
	//			if v1*v1 == nums1[i]*nums1[j] {
	//				res += cnt2[v1]
	//			}
	//		}
	//	}
	//}
	//优化后
	//类型1 计算结果
	res += calculateTriplets(cnt1, cnt2)
	//类型2 计算结果
	res += calculateTriplets(cnt2, cnt1)

	return res
}

func calculateTriplets(cnt1 map[int]int, cnt2 map[int]int) int {
	keys1 := make([]int, 0, len(cnt1))
	for k := range cnt1 {
		keys1 = append(keys1, k)
	}
	keys2 := make([]int, 0, len(cnt2))
	for k := range cnt2 {
		keys2 = append(keys2, k)
	}

	res := 0
	for _, v1 := range keys1 {
		square := v1 * v1
		//calculatedValues := make(map[int]byte)
		for _, v2 := range keys2 {
			//这里没必要再遍历了,可以直接通过除法去求候选值 （优化点1），其实有漏洞（除数如果为0呢）比如 nums1 = [0,0], nums2 = [0,0,0]，不过官方的测试用例要求数字是1-100000,只是题目上没有明说
			candidate := 0
			if square%v2 != 0 {
				continue
			}
			candidate = square / v2
			//为防止重复计算 比如 2×8 与 8×2 分别就记录一次, 其实没必要用容器类记录（优化点2），由于数值的特殊性，防止重复计算只需要加个if大小判断
			//if _, ok := calculatedValues[candidate]; ok {
			//	continue
			//}
			//calculatedValues[v2] = 0
			//calculatedValues[candidate] = 0
			//这里有两种情况 v2 == candidate 、 v2 != candidate
			if v2 == candidate { //排列组合情况, 比如1个1,2个1,3个1,...,n个1, 分别是0，1，3，n*(n-1)/2种组合
				res += cnt1[v1] * cnt2[v2] * (cnt2[candidate] - 1) / 2
			} else if v2 < candidate { //排列组合情况：比如1个1+2个2,...,m个1+n个2，m*n种组合, v2 < candidate 表示只记录一次 比如 2×8 与 8×2, 只记录2×8
				res += cnt1[v1] * cnt2[v2] * cnt2[candidate]
			}
		}
	}

	return res
}
