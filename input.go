package random

import (
  "fmt"
  "os/exec"
  "os"
)

// Implement this interface
type Input struct {
  binaryPath string
}

func NewInput(binPath string) (*Input, error) {
  if _, err := os.Stat(binPath); err == nil {
    return &Input{binPath}, nil
  } else {
    return nil, fmt.Errorf("input: file not found '%x'\n", binPath)
  }
}

func (i *Input) GetBits(n int)  {
  output, _ := exec.Command(i.binaryPath).Output()
  fmt.Println(output);
}

