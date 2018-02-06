package constant

const (
	_              = iota
	StarHyperGiant = 1 << iota
	StarSuperGiant
	StarBrightGiant
	StarGiant
	StarSubGiant
	_
	StarDwarf
	StarSubDwarf
	StarWhiteDwarf
	StarRedDwarf
	StarBrownDwarf
)

// Define currencies
type Curr struct {
	Currency string
	Name     string
	Country  string
	Number   int
}

var Currencies = []Curr{
	Curr{"DZD", "Algerian Dinar", "Algeria", 12},
	Curr{"AUD", "Australian Dollar", "Australia", 36},
	Curr{"EUR", "Euro", "Belgium", 978},
	Curr{"CLP", "Chilean Peso", "Chile", 152},
	Curr{"EUR", "Euro", "Greece", 978},
	Curr{"HTG", "Gourde", "Haiti", 332},
	Curr{"HKD", "Hong Kong Dollar", "Hong Koong", 344},
	Curr{"KES", "Kenyan Shilling", "Kenya", 404},
	Curr{"MXN", "Mexican Peso", "Mexico", 484},
	Curr{"USD", "US Dollar", "United States", 840},
}
