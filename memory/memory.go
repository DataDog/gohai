package memory

// Memory holds memory metadata about the host
type Memory struct {
	// Total is the total memory for the host in byte
	Total uint64
	// SwapTotal is the swap memory size in byte (Unix only)
	SwapTotal uint64
}

const name = "memory"

func (self *Memory) Name() string {
	return name
}

func (self *Memory) Collect() (result interface{}, err error) {
	result, err = getMemoryInfo()
	return
}

// Get returns a Memory struct  already initialized
func Get() (*Memory, error) {
	// Legacy code from gohai returns memory in:
	// - byte for Windows
	// - mix of byte and MB for OSX
	// - KB on linux
	//
	// this method being new we can align this behavior to return bytes everywhere without breaking backward
	// compatibility

	mem, swap, err := getMemoryInfoByte()
	if err != nil {
		return nil, err
	}

	return &Memory{
		Total:     mem,
		SwapTotal: swap,
	}, nil
}
