package main
import "testing"
func TestCalculateTotalDeduction(t *testing.T) {
	d := Deduction{
		SpouseDeduction:               true,
		PregnancyExpenses:             80000,
		SecondaryCities:               20000,
		ShopdeeMeekhun:                60000,
		HomeLoanInterest:              150000,
		PurchaseOtopProducts:          25000,
		PurchaseCommunityEnterprises:  30000,
		PurchaseFromSocialEnterprises: 25000,
		PurchaseWithVatETax:           40000,
		PurchaseWithEReceipt:          50000,
	}
	expected := 60000 + 60000 + 60000 + 15000 + 50000 + 100000 + 20000 + 20000 + 20000 + 30000 + 30000
	total := calculateTotalDeduction(d)
	if total != expected {
		t.Errorf("Expected total deduction %d, but got %d", expected, total)
	}
}
func TestCalculateTax(t *testing.T) {
	testCases := []struct {
		name     string 
		income   int    
		expected int    
	}{
		{
			name:"Income 450,000 should result in 22,500 tax",
			income: 450000,
			expected: 22500,
		},
		{
			name:"Income 800,000 should result in 75,000 tax",
			income: 800000,
			expected: 75000,
		},
		{
			name:"Income 2,500,000 should result in 515,000 tax",
			income: 2500000,
			expected: 515000,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := calculateTax(tc.income)
			if actual != tc.expected {
				t.Errorf("For income %d, expected tax %d, but got %d", tc.income, tc.expected, actual)
			}
		})
	}
}