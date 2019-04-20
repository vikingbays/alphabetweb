package toml

func ParseDatas(data string) (m map[string]interface{}, err error) {
	p, err := parse(data)
	if err != nil {
		return nil, err
	} else {
		return p.mapping, nil
	}
}
