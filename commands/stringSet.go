package commands

type stringSet []string

func (ss *stringSet) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

func (ss *stringSet) Get() interface{} {
	return stringSet(*ss)
}

func (ss *stringSet) String() string {
	return "[]"
}
