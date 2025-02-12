package leetcode

import (
	"container/list"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
)

func LRUCache(arr []string) []string {
	cache := list.New()
	memo := make(map[string]*list.Element)
	for _, val := range arr {
		if node, ok := memo[val]; ok {
			cache.MoveToBack(node)
		} else {
			if cache.Len() < 5 {
				memo[val] = cache.PushBack(val)
			} else {
				delete(memo, cache.Front().Value.(string))
				cache.Remove(cache.Front())
				memo[val] = cache.PushBack(val)
			}
		}
	}

	output := make([]string, 0, cache.Len())
	for node := cache.Front(); node != nil; node = node.Next() {
		output = append(output, node.Value.(string))
	}

	return output
}

// #1 Two Sum
func TwoSum(nums []int, target int) []int {
	hashmap := make(map[int]int)
	for i, v := range nums {
		complement := target - nums[i]
		if value, ok := hashmap[complement]; ok {
			return []int{value, i}
		}
		// We don't need to check if the key exists in the map
		// because the guarantee is that there is only one solution
		hashmap[v] = i
	}

	return []int{}
}

// #242 Valid Anagram
func IsAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	// Let's use a hashmap to store the frequency of each character
	hashmap := make(map[rune]int)
	for _, c := range s {
		hashmap[c]++
	}
	for _, c := range t {
		hashmap[c]--

		// If the frequency of a character is negative, one of the characters is not present in the other string
		// or the frequency of the character is different between the two strings
		if hashmap[c] < 0 {
			return false
		}
	}

	return true
}

// #3 Longest Substring Without Repeating Characters
func LengthOfLongestSubstring(s string) int {
	charToNextIndex := make(map[rune]int)

	longest := 0
	left := 0
	for i, r := range s {
		if prevIdx, ok := charToNextIndex[r]; ok {
			left = max(left, prevIdx)
		}
		longest = max(longest, i-left+1)
		charToNextIndex[r] = i + 1
	}

	return longest
}

// #266 Palindrome Permutation
func CanPermutePalindromeHashmap(s string) bool {
	// If the string has an even length, we expect every character to always occur an even number of times
	// e.g. "aabb" => "abba", "baab"
	// If the string has an odd length, we expect every character to occur an even number of times except for one character
	// e.g. "aabbccc" => "abcccba", "cabcbac", "acbcbca", "bacccab", "bcacacb"
	//
	// Therefore, if more than one character occurs an odd number of times, the string cannot be balanced into a palindrome
	// e.g. "aabbcccddd" => "abcdddcbac"
	occurrences := make(map[rune]int)
	for _, r := range s {
		occurrences[r]++
	}

	oddCount := 0
	for _, v := range occurrences {
		if v%2 == 1 {
			oddCount++
		}
	}

	return oddCount <= 1
}

func CanPermutePalindromeArray(s string) bool {
	occurrences := make([]byte, 128)
	for _, r := range s {
		occurrences[r]++
	}

	oddCount := 0
	for _, v := range occurrences {
		if v%2 == 1 {
			oddCount++
		}
	}

	return oddCount <= 1
}

func CanPermutePalindromeSinglePass(s string) bool {
	occurrences := make([]byte, 128)
	oddCount := 0
	for _, r := range s {
		occurrences[r]++
		if occurrences[r]%2 == 1 {
			oddCount++
		} else {
			oddCount--
		}
	}
	return oddCount <= 1
}

func CanPermutePalindromeSet(s string) bool {
	oddSet := mapset.NewSet[rune]()
	for _, r := range s {
		if oddSet.Contains(r) {
			oddSet.Remove(r)
		} else {
			oddSet.Add(r)
		}
	}
	return oddSet.Cardinality() <= 1
}

// #70 Climbing Stairs
// Golden comment: https://leetcode.com/problems/climbing-stairs/editorial/comments/604446
func ClimbStairsRecursive(n int) int {
	var recurse func(i int, n int, memo []int) int
	recurse = func(i int, n int, memo []int) int {
		if i > n {
			return 0
		}
		if i == n {
			return 1
		}
		if memo[i] > 0 {
			return memo[i]
		}
		memo[i] = recurse(i+1, n, memo) + recurse(i+2, n, memo)
		return memo[i]
	}

	memo := make([]int, n+1)
	return recurse(0, n, memo)
}

func ClimbStairsDynamicProgramming(n int) int {
	if n == 1 {
		return 1
	}

	dp := make([]int, n+1)
	dp[1] = 1
	dp[2] = 2
	for i := 3; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

func ClimbStairsFibonacci(n int) int {
	if n == 1 {
		return 1
	}

	first, second := 1, 2
	for i := 3; i <= n; i++ {
		third := first + second
		first = second
		second = third
	}

	return second
}

// #78 Subsets
func SubsetsIterativeCascade(nums []int) [][]int {
	subsets := [][]int{{}} // A slice with an empty slice as the first element

	// Generate subsets containing the current element
	for _, num := range nums {
		// Each new subset is the result of appending the current element
		// to an existing subset

		// This creates a slice of `len(subsets)` number of slices with each slice being a nil slice.
		newSubsets := make([][]int, len(subsets))
		for i, subset := range subsets {
			// We manually create a newSubset having n_subset+1 capacity to avoid reallocation
			// Use len(subset)+1 for both the length and capacity if you need a bit more performance for
			// very large slices, but it might not be worth it so profile/benchmark first.
			newSubset := make([]int, len(subset), len(subset)+1)
			copy(newSubset, subset)
			newSubset = append(newSubset, num)
			newSubsets[i] = newSubset
		}
		// Once new subsets have been generated, add them to the output slice of subsets
		subsets = append(subsets, newSubsets...)
	}

	return subsets
}

func SubsetsBacktracking(nums []int) [][]int {
	subsets := [][]int{}

	var backtrack func(first int, currentSubset []int)
	backtrack = func(first int, currentSubset []int) {
		// We add a copy of the current subset to the output
		// because Go slices behave like a reference type and will be modified
		// in the recursive calls
		subsetCopy := make([]int, len(currentSubset))
		copy(subsetCopy, currentSubset)
		subsets = append(subsets, subsetCopy)

		for i := first; i < len(nums); i++ {
			subsetCopy = append(subsetCopy, nums[i])
			backtrack(i+1, subsetCopy)
			subsetCopy = subsetCopy[:len(subsetCopy)-1]
		}
	}

	backtrack(0, []int{})

	return subsets
}

func SubsetsBitmask(nums []int) [][]int {
	n := len(nums)
	numSubsets := 1 << n // 2^n
	output := make([][]int, 0, numSubsets)

	for bitmask := 0; bitmask < numSubsets; bitmask++ {
		currentSubset := []int{}
		for i := 0; i < n; i++ {
			// Check if the i-th bit in bitmask is set (1)
			if (bitmask>>i)&1 == 1 {
				currentSubset = append(currentSubset, nums[i])
			}
		}
		output = append(output, currentSubset)
	}

	return output
}

// #22 Generate Parentheses
func GenerateParenthesisBruteForce(n int) []string {
	isValidPermutation := func(permutation string) bool {
		leftCount := 0
		for _, c := range permutation {
			if string(c) == "(" {
				leftCount++
			} else {
				leftCount--
			}

			if leftCount < 0 {
				return false
			}
		}

		return leftCount == 0
	}

	output := []string{}
	queue := list.New()
	queue.PushBack("")
	for queue.Len() > 0 {
		head := queue.Front()
		queue.Remove(head)
		headValue := head.Value.(string)
		if len(headValue) == 2*n {
			if isValidPermutation(headValue) {
				output = append(output, headValue)
			}
			continue
		}
		queue.PushBack(headValue + ")")
		queue.PushBack(headValue + "(")
	}
	return output
}

func GenerateParenthesisBacktracking(n int) []string {
	output := []string{}

	var backtrack func(permutation string, leftCount int, rightCount int)
	backtrack = func(permutation string, leftCount int, rightCount int) {
		if len(permutation) == 2*n {
			output = append(output, permutation)
			return
		}
		if leftCount < n {
			backtrack(permutation+"(", leftCount+1, rightCount)
		}
		if rightCount < leftCount {
			backtrack(permutation+")", leftCount, rightCount+1)
		}
	}
	backtrack("", 0, 0)

	return output
}

func GenerateParenthesisBitmasking(n int) []string {
	output := []string{}

	// 3 valid parentheses per string means a max of 6 characters always
	maxPermutationLength := 2 * n

	// `6 max characters` means up to 6 bits could be used positionally (1000000),
	// which means we will find all permutations from 1-000000 to 1-111111
	maxBit := 1 << maxPermutationLength

	// "(" is 0
	// ")" is 1
	for bitmask := 0; bitmask < maxBit; bitmask++ {
		leftCount := 0
		permutation := ""

		// Technically we can use conventional increments here,
		// but assigning the values might be confusing to other people.
		// If we increment, we'll have to reverse the concatenation.
		// Permutations would just flip. That's why technically, it all still works.
		for i := maxPermutationLength - 1; i >= 0; i-- {
			if (bitmask>>i)&1 == 0 {
				leftCount++
				permutation += "("
			} else {
				leftCount--
				permutation += ")"
			}

			if leftCount < 0 || leftCount > n {
				break
			}
		}

		if len(permutation) == maxPermutationLength && leftCount == 0 {
			output = append(output, permutation)
		}
	}

	return output
}

// #160 Intersection of Two Linked Lists
type ListNode struct {
	Val  int
	Next *ListNode
}

func NewListNodeFromInts(values ...int) *ListNode {
	head := &ListNode{}
	current := head
	for _, v := range values {
		current.Next = &ListNode{Val: v}
		current = current.Next
	}
	return head.Next
}

func GetIntersectionNodeHashmap(headA, headB *ListNode) *ListNode {
	bNodes := map[*ListNode]struct{}{}

	for current := headB; current != nil; current = current.Next {
		bNodes[current] = struct{}{}
	}

	for current := headA; current != nil; current = current.Next {
		if _, ok := bNodes[current]; ok {
			return current
		}
	}

	return nil
}

func GetIntersectionNodeTwoPointers(headA, headB *ListNode) *ListNode {
	pointerA, pointerB := headA, headB
	for pointerA != pointerB {
		if pointerA == nil {
			pointerA = headB
		} else {
			pointerA = pointerA.Next
		}

		if pointerB == nil {
			pointerB = headA
		} else {
			pointerB = pointerB.Next
		}
	}

	return pointerA
}

// #202 Happy Number
func IsHappyHashMap(n int) bool {
	seen := map[int]struct{}{}

	var sumOfSquares func(currentN int) int = func(currentN int) int {
		sum := 0
		for currentN > 0 {
			digit := currentN % 10
			sum += digit * digit
			currentN /= 10
		}
		return sum
	}

	for n != 1 {
		// If we've seen the number before, we're in a cycle
		if _, ok := seen[n]; ok {
			return false
		}
		seen[n] = struct{}{}

		n = sumOfSquares(n)
	}

	return true
}

func IsHappyTwoPointers(n int) bool {
	var sumOfSquares func(currentN int) int = func(currentN int) int {
		sum := 0
		for currentN > 0 {
			digit := currentN % 10
			sum += digit * digit
			currentN /= 10
		}

		return sum
	}

	slowPointer, fastPointer := n, sumOfSquares(n)

	for fastPointer != 1 && slowPointer != fastPointer {
		slowPointer = sumOfSquares(slowPointer)
		fastPointer = sumOfSquares(sumOfSquares(fastPointer))
	}

	return fastPointer == 1
}

// #88 Merge Sorted Arrays
func MergeSortedArraysThreePointersBeginning(nums1 []int, m int, nums2 []int, n int) {
	merged := make([]int, len(nums1))
	p1, p2 := 0, 0

	for i := 0; i < len(merged); i++ {
		if p2 >= n || (p1 < m && nums1[p1] <= nums2[p2]) {
			merged[i] = nums1[p1]
			p1++
		} else {
			merged[i] = nums2[p2]
			p2++
		}
	}

	copy(nums1, merged)
}

func MergeSortedArraysThreePointersEnd(nums1 []int, m int, nums2 []int, n int) {
	// start from the last element of nums1 and nums2, excluding the allocated positions for the merged array
	p1, p2 := m-1, n-1

	// set the merged-position pointer to the padded end of nums1
	for merged := m + n - 1; merged >= 0; merged-- {
		if p2 < 0 || (p1 >= 0 && nums1[p1] > nums2[p2]) {
			nums1[merged] = nums1[p1]
			p1--
		} else {
			nums1[merged] = nums2[p2]
			p2--
		}
	}
}

// #33 Search in Rotated Sorted Array
func SearchInRotatedSortedArrayPivotIndex(nums []int, target int) int {
	binarySearch := func(leftBound int, rightBound int, target int) int {
		left, right := leftBound, rightBound
		for left <= right {
			mid := left + (right-left)/2
			if nums[mid] == target {
				return mid
			} else if nums[mid] > target {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		return -1
	}

	n := len(nums)
	left, right := 0, n-1

	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] > nums[n-1] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if targetIndex := binarySearch(0, left-1, target); targetIndex != -1 {
		return targetIndex
	}
	return binarySearch(left, n-1, target)
}

func SearchInRotatedSortedArrayOneBinarySearch(nums []int, target int) int {
	n := len(nums)
	left, right := 0, n-1

	for left <= right {
		mid := left + (right-left)/2
		// case 1: mid is the target
		if target == nums[mid] {
			return mid
		}

		// case 2: left-to-mid subarray is sorted
		if nums[left] < nums[mid] {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// case 3: mid-to-right subarray is sorted
			if nums[mid] <= target && target < nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

// #4 Median of Two Sorted Arrays
func FindMedianSortedArraysMergeSort(nums1 []int, nums2 []int) float64 {
	m, n := len(nums1), len(nums2)

	// We track the middle values for the median
	prevVal, nextVal := 0, 0

	p1, p2 := 0, 0
	for p1+p2 <= (m+n)/2 {
		if p1 < m && p2 < n {
			if nums1[p1] < nums2[p2] {
				prevVal = nextVal
				nextVal = nums1[p1]
				p1++
			} else {
				prevVal = nextVal
				nextVal = nums2[p2]
				p2++
			}
		} else if p1 < m {
			prevVal = nextVal
			nextVal = nums1[p1]
			p1++
		} else {
			prevVal = nextVal
			nextVal = nums2[p2]
			p2++
		}
	}

	if (m+n)%2 == 0 {
		return float64(prevVal+nextVal) / 2
	} else {
		return float64(nextVal)
	}
}

func FindMedianSortedArraysBinarySearchRecursive(nums1 []int, nums2 []int) float64 {
	// not yet implemented
	return 0
}

// #104 Maximum Depth of Binary Tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewBinaryTreeFromSlice(values []interface{}) *TreeNode {
	if len(values) == 0 {
		return nil
	}

	root := &TreeNode{Val: values[0].(int)}
	queue := list.New()
	queue.PushBack(root)

	for i := 1; i < len(values); i += 2 {
		element := queue.Front()
		queue.Remove(element)
		node := element.Value.(*TreeNode)

		if values[i] != nil {
			node.Left = &TreeNode{Val: values[i].(int)}
			queue.PushBack(node.Left)
		}

		if i+1 < len(values) && values[i+1] != nil {
			node.Right = &TreeNode{Val: values[i+1].(int)}
			queue.PushBack(node.Right)
		}
	}

	return root
}

func MaximumDepthOfBinaryTreeDFSRecursive(root *TreeNode) int {
	var _maxDepth func(node *TreeNode, depth int) int
	_maxDepth = func(node *TreeNode, depth int) int {
		if node == nil {
			return depth
		}
		return max(_maxDepth(node.Left, depth+1), _maxDepth(node.Right, depth+1))
	}

	return _maxDepth(root, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// #98 Validate Binary Search Tree
func IsValidBSTRecursiveTraversalWithValidRange(root *TreeNode) bool {
	var validate func(root *TreeNode, low *int, high *int) bool
	validate = func(root *TreeNode, low *int, high *int) bool {
		if root == nil {
			return true
		}

		if (low != nil && root.Val <= *low) || (high != nil && root.Val >= *high) {
			return false
		}

		return validate(root.Left, low, &root.Val) && validate(root.Right, &root.Val, high)
	}

	return validate(root, nil, nil)
}

// #94 Binary Tree Inorder Traversal
func InorderTraversalIterative(root *TreeNode) []int {
	order := []int{}
	stack := list.New()

	curr := root
	for curr != nil || stack.Len() != 0 {
		for curr != nil {
			stack.PushFront(curr)
			curr = curr.Left
		}

		head := stack.Front()
		stack.Remove(head)
		node := head.Value.(*TreeNode)
		order = append(order, node.Val)
		curr = node.Right
	}

	return order
}

func FindMaxConsecutiveOnes(nums []int) int {
	maxCount, count := 0, 0
	for _, num := range nums {
		if num == 1 {
			count++
		} else {
			maxCount = max(maxCount, count)
			count = 0
		}
	}
	return max(maxCount, count)
}

// #1295 Find Numbers with Even Number of Digits
func FindNumbersWithEvenDigitsLogarithm(nums []int) int {
	// Given a positive integer x, the number of digits in x is ⌊log10(x)⌋+1.
	count := 0
	for _, num := range nums {
		if int(math.Log10(float64(num)))%2 == 1 {
			count++
		}
	}
	return count
}

// #977 Squares of a Sorted Array
func SquaresOfSortedArrayTwoPointers(nums []int) []int {
	intAbs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	n := len(nums)
	squares := make([]int, n)

	left, right := 0, n-1
	for i := n - 1; i >= 0; i-- {
		var squareRoot int // sqrt(9) = 3. We are looking for `3`.
		if intAbs(nums[left]) < intAbs(nums[right]) {
			squareRoot = nums[right]
			right--
		} else {
			squareRoot = nums[left]
			left++
		}
		squares[i] = squareRoot * squareRoot
	}

	return squares
}
