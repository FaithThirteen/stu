package main



func BinarySearch(n []int, target int) int {
	length := len(n)
	low := 0
	high := length - 1
	for low <= high {
		mid := (low + high) / 2
		if n[mid] > target {
			high = mid - 1
		} else if n[mid] < target {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}