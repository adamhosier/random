package bitstring

import "testing"

type bitStringTest struct {
	strInput  string
	byteInput *[]byte
	index     int
	wantBit   bool
	wantError bool
	want      BitString
}

func TestNewBitString(t *testing.T) {
	bs := NewBitString()
	if bs.Length != 0 {
		t.Error("NewBitString has non-zero length")
	}
}

func TestBitStringFromString(t *testing.T) {
	cases := []bitStringTest{
		{strInput: "00101101001", wantError: false},
		{strInput: "00101?10101", wantError: true},
		{strInput: "0 1 0 1 0 0", wantError: true},
		{strInput: "     010010", wantError: true},
		{strInput: "001101     ", wantError: true},
	}
	for _, c := range cases {
		_, err := BitStringFromString(c.strInput)
		if err == nil {
			if c.wantError {
				t.Errorf("BitStringFromString(\"%v\") expected an error to be thrown", c.strInput)
			}
		} else {
			if !c.wantError {
				t.Errorf("BitStringFromString(\"%v\") threw an error which wasnt expected", c.strInput)
			}
		}
	}
}

func TestBitStringFromBytes(t *testing.T) {
	cases := []bitStringTest{
		{byteInput: &[]byte{41}, wantError: false},
	}
	for _, c := range cases {
		_, err := BitStringFromBytes(c.byteInput)
		if err == nil {
			if c.wantError {
				t.Errorf("BitStringFromUTF8(\"%v\") expected an error to be thrown", c.strInput)
			} else {

			}
		} else {
			if !c.wantError {
				t.Errorf("BitStringFromUTF8(\"%v\") threw an error which wasnt expected", c.strInput)
			}
		}
	}
}

func TestBitStringFromInt(t *testing.T) {
	got := BitStringFromInt(3, 5)
	want, _ := BitStringFromString("101")
	if !got.Equals(want) {
		t.Errorf("BitStringFromInt(3, 5) == %q, want %q", got, want)
	}

	got = BitStringFromInt(20, 1023)
	want, _ = BitStringFromString("00000000001111111111")
	if !got.Equals(want) {
		t.Errorf("BitStringFromInt(20, 1023) == %q, want %q", got, want)
	}
}

func TestBitStringsOfLength(t *testing.T) {
	n := 2
	bitStrings := BitStringsOfLength(n)

	// make a map containing if a candidate has been seen or not, initially all false
	seen := make(map[string]bool)
	expected := []string{"00", "01", "10", "11"}
	for _, e := range expected {
		seen[e] = false
	}

	// explore the candidates, setting them to true in the map when seen
	for _, bs := range bitStrings {
		seen[bs.String()] = true
	}

	// validate everything has been seen
	for _, e := range expected {
		if !seen[e] {
			t.Errorf("BitStringsOfLength(%d) is missing %s", n, e)
		}
	}
}

func TestBitString_At(t *testing.T) {
	bs, _ := BitStringFromString("0010010010")
	cases := []bitStringTest{
		{index: 0, wantBit: false},
		{index: 1, wantBit: false},
		{index: 2, wantBit: true},
		{index: 3, wantBit: false},
		{index: 4, wantBit: false},
		{index: 5, wantBit: true},
		{index: 6, wantBit: false},
		{index: 7, wantBit: false},
		{index: 8, wantBit: true},
		{index: 9, wantBit: false},
	}
	for _, c := range cases {
		got := bs.At(c.index)
		if got != c.wantBit {
			t.Errorf("BitString.At(%d) == %b, expected %b", c.index, got, c.wantBit)
		}
	}
}

func TestBitString_Add(t *testing.T) {
	bs, _ := BitStringFromString("01")
	bs.Add(false)
	if bs.At(2) {
		t.Error("BitString.Add(0) == 1, expected 0")
	}
	bs.Add(true)
	if !bs.At(3) {
		t.Error("BitString.Add(1) == 0, expected 1")
	}
}

func TestBitString_First(t *testing.T) {
	bs1, _ := BitStringFromString("10010")
	bs2, _ := BitStringFromString("100")

	n := 3
	got := bs1.First(n)
	if !got.Equals(bs2) {
		t.Errorf("BitString.First(%d) == %q, expected %q", n, got, bs2)
	}

	n = 5
	got = bs1.First(n)
	if !got.Equals(bs1) {
		t.Errorf("BitString.First(%d) == %q, expected %q", n, got, bs1)
	}
}

func TestBitString_Invert(t *testing.T) {
	bs, _ := BitStringFromString("1111100000")
	bs.Invert(0)
	want, _ := BitStringFromString("0111100000")
	if !bs.Equals(want) {
		t.Errorf("BitString.Invert(%d) == %q, expected %q", 0, bs, want)
	}

	bs.Invert(9)
	want, _ = BitStringFromString("0111100001")
	if !bs.Equals(want) {
		t.Errorf("BitString.Invert(%d) == %q, expected %q", 9, bs, want)
	}
}

func TestBitString_Substring(t *testing.T) {
	bs1, _ := BitStringFromString("10101")
	bs2, _ := BitStringFromString("010")
	start := 1
	length := 3
	got := bs1.Substring(start, length)
	if !got.Equals(bs2) {
		t.Errorf("BitString.Substring(%d, %d) == %q, expected %q", start, length, got, bs2)
	}
}

func TestBitString_Extend(t *testing.T) {
	bs1, _ := BitStringFromString("00")
	bs2, _ := BitStringFromString("11")
	bs3, _ := BitStringFromString("0011")
	got := bs1.Extend(bs2)
	if !got.Equals(bs3) {
		t.Errorf("BitString.Extend() == %q, expected %q", got, bs3)
	}
}

func TestBitString_Partition(t *testing.T) {
	bs, _ := BitStringFromString("00110")
	bs1, _ := BitStringFromString("00")
	bs2, _ := BitStringFromString("11")
	bss := []*BitString{bs1, bs2}

	got := bs.Partition(2)
	for i, bs := range got {
		if !bs.Equals(bss[i]) {
			t.Error("BitString.Partition produced an invalid partition")
		}
	}
}

func TestBitString_Ones(t *testing.T) {
	bs, _ := BitStringFromString("0011000110101110101101")
	want := 12
	got := bs.Ones()
	if got != want {
		t.Errorf("BitString.Ones() == %d, expected %d", got, want)
	}
}

func TestBitString_Proportion(t *testing.T) {
	bs, _ := BitStringFromString("0011000110101110101101")
	want := 12.0 / 22.0
	got := bs.Proportion()
	if got != want {
		t.Errorf("BitString.Proportion() == %f, expected %f", got, want)
	}
}

func TestBitString_HasTemplateAt(t *testing.T) {
	bs, _ := BitStringFromString("00110")
	bs1, _ := BitStringFromString("11")

	index := 2
	want := true
	got := bs.HasTemplateAt(bs1, index)
	if got != want {
		t.Errorf("BitString.HasTemplateAt(%q, %d) == %b, expected %b", bs1, index, got, want)
	}

	index = 3
	want = false
	got = bs.HasTemplateAt(bs1, index)
	if got != want {
		t.Errorf("BitString.HasTemplateAt(%q, %d) == %b, expected %b", bs1, index, got, want)
	}
}

func TestBitString_BinaryAdd(t *testing.T) {
	bs1, _ := BitStringFromString("00001")
	bs2, _ := BitStringFromString("00010")
	bs3, _ := BitStringFromString("00011")
	bs4, _ := BitStringFromString("00111")
	bs5, _ := BitStringFromString("01111")
	bs6, _ := BitStringFromString("11111")

	got := bs1.BinaryAdd(bs2)
	want, _ := BitStringFromString("00011")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs1, bs2, got, want)
	}

	got = bs2.BinaryAdd(bs3)
	want, _ = BitStringFromString("00101")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs2, bs3, got, want)
	}

	got = bs3.BinaryAdd(bs3)
	want, _ = BitStringFromString("00110")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs3, bs3, got, want)
	}

	got = bs4.BinaryAdd(bs4)
	want, _ = BitStringFromString("01110")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs4, bs4, got, want)
	}

	got = bs5.BinaryAdd(bs5)
	want, _ = BitStringFromString("11110")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs5, bs5, got, want)
	}

	got = bs6.BinaryAdd(bs2)
	want = bs1
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryAdd (%q + %q) == %q, expected %q", bs6, bs2, got, want)
	}
}

func TestBitString_BinaryMul(t *testing.T) {
	bs1, _ := BitStringFromString("00001")
	bs2, _ := BitStringFromString("00010")
	bs3, _ := BitStringFromString("00011")
	bs4, _ := BitStringFromString("00111")
	bs5, _ := BitStringFromString("01111")
	bs6, _ := BitStringFromString("11111")

	got := bs1.BinaryMul(bs2)
	want, _ := BitStringFromString("00010")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryMul (%q * %q) == %q, expected %q", bs1, bs2, got, want)
	}

	got = bs3.BinaryMul(bs3)
	want, _ = BitStringFromString("01001")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryMul (%q * %q) == %q, expected %q", bs3, bs3, got, want)
	}

	got = bs4.BinaryMul(bs3)
	want, _ = BitStringFromString("10101")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryMul (%q * %q) == %q, expected %q", bs4, bs3, got, want)
	}

	got = bs5.BinaryMul(bs5)
	want, _ = BitStringFromString("00001")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryMul (%q * %q) == %q, expected %q", bs5, bs5, got, want)
	}

	got = bs6.BinaryMul(bs2)
	want, _ = BitStringFromString("11110")
	if !got.Equals(want) {
		t.Errorf("BitString.BinaryMul (%q * %q) == %q, expected %q", bs6, bs2, got, want)
	}
}

func TestBitString_Equals(t *testing.T) {
	bs1, _ := BitStringFromString("0010")
	bs2, _ := BitStringFromString("0010")
	want := true
	got := bs1.Equals(bs2)
	if got != want {
		t.Errorf("BitString.Equals (%q Equals %q) == %t, want %t", bs1, bs2, got, want)
	}

	bs1, _ = BitStringFromString("0011")
	bs2, _ = BitStringFromString("001")
	want = false
	got = bs1.Equals(bs2)
	if got != want {
		t.Errorf("BitString.Equals (%q Equals %q) == %t, want %t", bs1, bs2, got, want)
	}

	bs3, _ := BitStringFromString("1")
	bs4 := bs2.Extend(bs3)
	want = true
	got = bs1.Equals(bs4)
	if got != want {
		t.Errorf("BitString.Equals (%q Equals %q) == %t, want %t", bs1, bs4, got, want)
	}
}

func TestBitString_Int(t *testing.T) {
	bs, _ := BitStringFromString("00100101")
	want := 37
	got := bs.Int()
	if got != want {
		t.Errorf("BitString.Int %q == %d, want %d", bs, got, want)
	}

	bs, _ = BitStringFromString("000000")
	want = 0
	got = bs.Int()
	if got != want {
		t.Errorf("BitString.Int %q == %d, want %d", bs, got, want)
	}

	bs, _ = BitStringFromString("11111")
	want = 31
	got = bs.Int()
	if got != want {
		t.Errorf("BitString.Int %q == %d, want %d", bs, got, want)
	}
}

func TestBitString_Bytes(t *testing.T) {
	want := []byte{0xDD, 0xF5, 0x9E, 0x00, 0xFF, 0x53}
	bs, _ := BitStringFromBytes(&want)
	got := bs.Bytes()

	if len(got) != len(want) {
		t.Errorf("BitString.Bytes() == 0x%X, want 0x%X", got, want)
	} else {
		failed := false
		for i, v := range want {
			if got[i] != v {
				failed = true
			}
		}
		if failed {
			t.Errorf("BitString.Bytes() == 0x%X, want 0x%X", got, want)
		}
	}

}

func TestBitString_Compare(t *testing.T) {
	bs1, _ := BitStringFromString("0001")
	bs2, _ := BitStringFromString("00001")
	bs3, _ := BitStringFromString("010100")

	want := 0
	got := bs1.Compare(bs2)
	if got != want {
		t.Errorf("BitString.Compare(%q, %q) == %d, want %d", bs1, bs2, got, want)
	}

	want = 1
	got = bs3.Compare(bs1)
	if got != want {
		t.Errorf("BitString.Compare(%q, %q) == %d, want %d", bs3, bs1, got, want)
	}

	want = -1
	got = bs2.Compare(bs3)
	if got != want {
		t.Errorf("BitString.Compare(%q, %q) == %d, want %d", bs2, bs3, got, want)
	}
}
