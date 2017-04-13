package random

import (
	"fmt"
	"os"
	"os/exec"
)

type Input struct {
	binaryPath string
	buffer     *[]byte
}

// Builds a new input type relating to the binary at path [binPath]
func NewInput(binPath string) (*Input, error) {
	// check file exists
	if _, err := os.Stat(binPath); err == nil {
		return &Input{binPath, &[]byte{}}, nil
	} else {
		return nil, fmt.Errorf("input: file not found '%x'\n", binPath)
	}
}

// Fetches n bits from the buffer. If the buffer is empty, fetch a new batch of bits first
func (i *Input) GetBits(n int) *BitString {
	if n <= 0 {
		return nil
	}

	// collect bits until we have enough
	for len(*i.buffer) < n {
		// run binary, then collect it's stdout to the buffer
		output, _ := exec.Command(i.binaryPath).Output()
		newBuffer := append(*i.buffer, output...)
		i.buffer = &newBuffer
	}

	// round up n bits to closest byte boundary
	numBytes := ((n - 1) / 8) + 1

	// take those bytes from the buffer
	bytes := (*i.buffer)[:numBytes]
	tmp := (*i.buffer)[numBytes:]
	i.buffer = &tmp
	bs, _ := BitStringFromBytes(&bytes)

	// discard extra bits
	return bs.Substring(0, n)
}
