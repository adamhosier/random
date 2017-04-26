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
