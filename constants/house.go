package constants

type SaleStatus = int

const (
	Unknown        SaleStatus = iota
	NotSale                   // 未销（预）售
	Sold                      // 已网上销售
	Limit                     // 限制出售
	Mortgage                  // 已在建工程抵押
	Seized                    // 已查封
	MortgageSeized            // 已在建工程抵押已查封
)
