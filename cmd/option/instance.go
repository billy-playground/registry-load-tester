package option

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Instance struct {
	Count         int
	BatchSize     int
	BatchInterval time.Duration

	flag string
}

func (i *Instance) SetFlag(instanceInput string) {
	i.flag = instanceInput
}

func (i *Instance) Parse() error {
	var err error
	var count int
	var size int
	var interval time.Duration

	numInstancesOption, frequencyOption, ok := strings.Cut(i.flag, "=")
	if _, err := fmt.Sscanf(numInstancesOption, "%d", &count); err != nil {
		return fmt.Errorf("Error parsing number of instances from %q: %v\n", numInstancesOption, err)
	}
	if count <= 0 {
		return fmt.Errorf("Number of instances must be greater than 0\n")
	}

	if ok {
		sizeOption, intervalOption, ok := strings.Cut(frequencyOption, "/")
		if !ok {
			return errors.New("Batch size and interval should be in the format <size>/<interval>\n")
		}
		if _, err = fmt.Sscanf(sizeOption, "%d", &size); err != nil {
			return fmt.Errorf("Error parsing batch size fron %q: %v\n", sizeOption, err)
		}
		if interval, err = time.ParseDuration(intervalOption); err != nil {
			return fmt.Errorf("Error parsing interval from %q: %v\n", intervalOption, err)
		}
		if size <= 0 {
			return fmt.Errorf("Batch size must be greater than 0\n")
		}
		if interval <= 0 {
			return fmt.Errorf("Interval must be greater than 0\n")
		}
	}
	i.Count = count
	i.BatchSize = size
	i.BatchInterval = interval

	return nil
}
