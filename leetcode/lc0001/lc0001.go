package lc0001

/**
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。
你可以按任意顺序返回答案。

示例 1：
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

示例 2：
输入：nums = [3,2,4], target = 6
输出：[1,2]

示例 3：
输入：nums = [3,3], target = 6
输出：[0,1]

提示：
2 <= nums.length <= 104
-109 <= nums[i] <= 109
-109 <= target <= 109
只会存在一个有效答案

进阶：你可以想出一个时间复杂度小于 O(n2) 的算法吗？
*/

//优化2(官方题解)： 实际只需要一次遍历
func twoSum(nums []int, target int) []int {
	hashTable := map[int]int{}
	for i, x := range nums {
		if p, ok := hashTable[target-x]; ok { //（优化点）实际只需要一次遍历，hashTable只需要记录遍历过的元素，后续
			return []int{p, i}
		}
		hashTable[x] = i
	}
	return nil
}

//优化1（性能好、内存占用高）： 数据结构优化、重复值提前判断
//func twoSum(nums []int, target int) []int {
//	vm := make(map[int]int) //值可能重复因为解只有一个只需要记录最大的索引
//	for i := 0; i < len(nums); i++ {
//		idx, ok := vm[nums[i]]
//		if ok && nums[i]*2 == target { //重复值可以提前判断
//			return []int{idx, i}
//		}
//		vm[nums[i]] = i
//	}
//
//	for v, i1 := range vm {
//		//因为数字的特殊性，可以用运算代替遍历
//		candidate := target - v
//		if i2, ok := vm[candidate]; ok && i2 != i1 {
//			return []int{i1, i2}
//		}
//	}
//	return nil
//}

//原始解法
//func twoSum(nums []int, target int) []int {
//	vm := make(map[int][]int)
//	for i := 0; i < len(nums); i++ {
//		ids, ok := vm[nums[i]]
//		if ok {
//			vm[nums[i]] = append(ids, i)
//		} else {
//			vm[nums[i]] = []int{i}
//		}
//	}
//
//	for v := range vm {
//		//因为数字的特殊性，可以用运算代替遍历 （优化点）
//		candidate := target - v
//		if ids, ok := vm[candidate]; ok {
//			if v == candidate && len(ids) > 1 {
//				return []int{ids[0], ids[1]}
//			} else if v != candidate {
//				return []int{vm[v][0], vm[candidate][0]}
//			}
//		}
//	}
//	return nil
//}
