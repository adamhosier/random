package random

import (
	"fmt"
	"math"
)

type BitString struct {
	length int
	data   []bool
}

// Builds a BitString from a string of zeros and ones e.g. "01101001"
func BitStringFromString(s string) (*BitString, error) {
	n := len(s)
	bs := make([]bool, n)
	for i, c := range s {
		if c == '1' {
			bs[i] = true
		} else if c == '0' {
			bs[i] = false
		} else {
			return nil, fmt.Errorf("bitstring: Invalid character '%v' in BitStringFromString", c)
		}
	}
	return &BitString{n, bs}, nil
}

// Builds a BitString from some string containing UTF-8 encoded binary
func BitStringFromBytes(bytes *[]byte) (*BitString, error) {
	n := len(*bytes) * 8
	bs := make([]bool, n)
	for i, b := range *bytes {
		for j := 0; j < 8; j++ {
			bs[i*8+j] = (b & 0x1) == 0x1
			b = b >> 1
		}
	}
	return &BitString{n, bs}, nil
}

// Gets a slice of all possible bit strings of length [n] (000, 001, 010, ... for n=3)
func BitStringsOfLength(n int) []*BitString {
	// Base case
	if n == 1 {
		return []*BitString{{1, []bool{false}}, {1, []bool{true}}}
	}

	// Recursive case
	bss := make([]*BitString, int(math.Pow(2.0, float64(n))))
	for i, bs := range BitStringsOfLength(n - 1) {
		d := make([]bool, bs.length)
		copy(d, bs.data)
		bss[2*i] = &BitString{n, append(d, false)}
		bss[2*i+1] = &BitString{n, append(d, true)}
	}
	return bss
}

// Gets the value of the bit at position [i] in [bs]
func (bs *BitString) At(i int) bool {
	return bs.data[i]
}

// Adds a value [b] to the end of this [bs]
func (bs *BitString) Add(b bool) {
	bs.length++
	bs.data = append(bs.data, b)
}

// Gets the first [n] bits from [bs]
func (bs *BitString) First(n int) *BitString {
	d := make([]bool, n)
	copy(d, bs.data[:n])
	return &BitString{n, d}
}

// Finds the substring of [bs] starting at [start] of length [len]
func (bs *BitString) Substring(start, len int) *BitString {
	d := make([]bool, len)
	copy(d, bs.data[start:start+len])
	return &BitString{len, d}
}

// Adds [bs1] to the end of [bs] returning a new BitString
func (bs *BitString) Extend(bs1 *BitString) *BitString {
	return &BitString{bs.length + bs1.length, append(bs.data, bs1.data...)}
}

// Partitions [bs] into blocks of length [len] discarding extra bits at the end
func (bs *BitString) Partition(len int) []*BitString {
	bss := make([]*BitString, bs.length/len)
	for i := 0; i < bs.length/len; i++ {
		bss[i] = &BitString{len, bs.data[i*len : (i+1)*len]}
	}
	return bss
}

// Count number of ones in [bs]
func (bs *BitString) Ones() int {
	var sum int
	for _, b := range bs.data {
		if b {
			sum++
		}
	}
	return sum
}

// Calculates the proportion of ones to zeros in [bs]
func (bs *BitString) Proportion() float64 {
	return float64(bs.Ones()) / float64(bs.length)
}

// Tests if [block] matches [template] at position i
func (bs *BitString) HasTemplateAt(template *BitString, i int) bool {
	for j, b := range template.data {
		if bs.At(i+j) != b {
			return false
		}
	}
	return true
}

// Returns an int holding the inner product of this BitString and another
func (bs *BitString) InnerProduct(other *BitString) int {
	res := 0
	for i := 0; i < int(math.Min(float64(bs.length), float64(other.length))); i++ {
		if bs.At(i) && other.At(i) {
			res++
		}
	}
	return res
}

// Converts [bs] to a string of zeros and ones
func (bs *BitString) String() string {
	var s string
	for _, b := range bs.data {
		if b {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

// Converts [bs] to an integer (max length 64)
func (bs *BitString) Int() int {
	var n int
	for i, b := range bs.data {
		if b {
			n |= 0x1 << uint64(bs.length-i-1)
		}
	}
	return n
}

// Compares if two bit strings contain the same bit sequence
func (bs *BitString) Equals(other *BitString) bool {
	if bs.length != other.length {
		return false
	}
	for i, b := range bs.data {
		if other.data[i] != b {
			return false
		}
	}
	return true
}