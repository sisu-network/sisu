package types

func (m *Token) GetAddressForChain(c string) string {
	for i, chain := range m.Chains {
		if chain == c {
			return m.Addresses[i]
		}
	}

	return ""
}
