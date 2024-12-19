package testmanager

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

// CommonValidator provides reusable validation functions
type CommonValidator struct{}

// ValidateCodeData validates code lookup data
func (v *CommonValidator) ValidateCodeData(data interface{}) error {
	code, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid code data format")
	}

	// Required fields validation
	requiredFields := []string{"tin", "bhfId", "lastReqDt"}
	for _, field := range requiredFields {
		if _, exists := code[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate codes array
	if codes, exists := code["codes"].([]interface{}); exists {
		for i, c := range codes {
			codeMap, ok := c.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid code format at index %d", i)
			}
			if err := v.validateCodeFields(codeMap); err != nil {
				return fmt.Errorf("code validation failed at index %d: %s", i, err)
			}
		}
	}

	return nil
}

// validateCodeFields validates individual code fields
func (v *CommonValidator) validateCodeFields(code map[string]interface{}) error {
	requiredFields := []string{"groupCd", "groupNm", "codeCd", "codeNm", "useYn"}
	for _, field := range requiredFields {
		if _, exists := code[field]; !exists {
			return fmt.Errorf("missing required field in code: %s", field)
		}
	}

	// Validate useYn is either "Y" or "N"
	if useYn, ok := code["useYn"].(string); ok {
		if useYn != "Y" && useYn != "N" {
			return fmt.Errorf("invalid useYn value: %s", useYn)
		}
	}

	return nil
}

// ValidateItemData validates item data
func (v *CommonValidator) ValidateItemData(data interface{}) error {
	item, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item data format")
	}

	// Required fields validation
	requiredFields := []string{"tin", "bhfId", "items"}
	for _, field := range requiredFields {
		if _, exists := item[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate items array
	if items, exists := item["items"].([]interface{}); exists {
		for i, itm := range items {
			itemMap, ok := itm.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid item format at index %d", i)
			}
			if err := v.validateItemFields(itemMap); err != nil {
				return fmt.Errorf("item validation failed at index %d: %s", i, err)
			}
		}
	}

	return nil
}

// validateItemFields validates individual item fields
func (v *CommonValidator) validateItemFields(item map[string]interface{}) error {
	requiredFields := []string{
		"itemCd", "itemClsCd", "itemTyCd", "itemNm", "orgnNatCd",
		"pkgUnitCd", "qtyUnitCd", "taxTyCd",
	}
	for _, field := range requiredFields {
		if _, exists := item[field]; !exists {
			return fmt.Errorf("missing required field in item: %s", field)
		}
	}

	// Validate item code format (e.g., "KE1NTXU0000001")
	if itemCd, ok := item["itemCd"].(string); ok {
		if !regexp.MustCompile(`^KE1NTXU\d{7}$`).MatchString(itemCd) {
			return fmt.Errorf("invalid item code format: %s", itemCd)
		}
	}

	return nil
}

// ValidatePurchaseData validates purchase transaction data
func (v *CommonValidator) ValidatePurchaseData(data interface{}) error {
	purchase, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid purchase data format")
	}

	// Required fields validation
	requiredFields := []string{
		"tin", "bhfId", "invcNo", "spplrTin", "pchsDt",
		"totItemCnt", "totTaxblAmt", "totTaxAmt", "totAmt",
	}
	for _, field := range requiredFields {
		if _, exists := purchase[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate date formats
	dateFields := []string{"pchsDt", "wrhsDt"}
	for _, field := range dateFields {
		if dt, ok := purchase[field].(string); ok {
			if _, err := time.Parse("20060102150405", dt); err != nil {
				return fmt.Errorf("invalid date format for %s: %s", field, err)
			}
		}
	}

	// Validate amounts
	if err := v.validateAmounts(purchase); err != nil {
		return err
	}

	// Validate items
	if items, exists := purchase["items"].([]interface{}); exists {
		for i, item := range items {
			if err := v.validatePurchaseItem(item); err != nil {
				return fmt.Errorf("invalid purchase item at index %d: %s", i, err)
			}
		}
	}

	return nil
}

// ValidateStockData validates stock data
func (v *CommonValidator) ValidateStockData(data map[string]interface{}) error {
	required := []string{"tin", "bhfId", "items"}
	for _, field := range required {
		if _, ok := data[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	items, ok := data["items"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid items format")
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid item format")
		}

		itemRequired := []string{"itemCd", "qty", "prc", "splyAmt"}
		for _, field := range itemRequired {
			if _, ok := itemMap[field]; !ok {
				return fmt.Errorf("missing required field in item: %s", field)
			}
		}
	}

	return nil
}

// validateAmounts checks the consistency of amount fields
func (v *CommonValidator) validateAmounts(data map[string]interface{}) error {
	// Validate total amounts
	totTaxblAmt, _ := data["totTaxblAmt"].(float64)
	totTaxAmt, _ := data["totTaxAmt"].(float64)
	totAmt, _ := data["totAmt"].(float64)

	// Basic amount validation
	if totTaxblAmt < 0 || totTaxAmt < 0 || totAmt < 0 {
		return fmt.Errorf("amounts cannot be negative")
	}

	// Validate total amount calculation
	expectedTotal := totTaxblAmt + totTaxAmt
	if expectedTotal != totAmt {
		return fmt.Errorf("total amount mismatch: expected %.2f, got %.2f", expectedTotal, totAmt)
	}

	return nil
}

// validatePurchaseItem validates individual purchase item data
func (v *CommonValidator) validatePurchaseItem(item interface{}) error {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item format")
	}

	requiredFields := []string{
		"itemSeq", "itemCd", "itemClsCd", "itemNm",
		"qty", "prc", "splyAmt", "taxblAmt", "taxAmt", "totAmt",
	}
	for _, field := range requiredFields {
		if _, exists := itemMap[field]; !exists {
			return fmt.Errorf("missing required field in item: %s", field)
		}
	}

	// Validate amounts
	splyAmt, _ := itemMap["splyAmt"].(float64)
	taxAmt, _ := itemMap["taxAmt"].(float64)
	totAmt, _ := itemMap["totAmt"].(float64)

	if splyAmt+taxAmt != totAmt {
		return fmt.Errorf("item amount mismatch: supply amount + tax amount != total amount")
	}

	return nil
}

// ValidateDataConsistency checks if the data is consistent with previous state
func (v *CommonValidator) ValidateDataConsistency(current, previous interface{}) error {
	if previous == nil {
		return nil
	}

	currentVal := reflect.ValueOf(current)
	previousVal := reflect.ValueOf(previous)

	if currentVal.Kind() != previousVal.Kind() {
		return fmt.Errorf("data type mismatch between current and previous state")
	}

	// Add specific consistency checks based on business rules
	return nil
}

// ValidateStockMovement validates a stock movement transaction
func (v *CommonValidator) ValidateStockMovement(ctx *ProcessContext) error {
	// Validate required fields
	required := []string{
		"tin", "bhfId", "sarNo", "regTyCd", "sarTyCd",
		"ocrnDt", "stockRlsDt", "totItemCnt",
	}
	for _, field := range required {
		if _, ok := ctx.Data[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	
	// Validate items
	items, ok := ctx.Data["itemList"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid itemList format")
	}
	
	if len(items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	
	for _, item := range items {
		if itemCd, ok := item["itemCd"].(string); !ok || itemCd == "" {
			return fmt.Errorf("ItemCd is required for all items")
		}
		if qty, ok := item["qty"].(float64); !ok || qty <= 0 {
			return fmt.Errorf("Qty must be positive for item %s", item["itemCd"].(string))
		}
		if prc, ok := item["prc"].(float64); !ok || prc < 0 {
			return fmt.Errorf("Prc cannot be negative for item %s", item["itemCd"].(string))
		}
	}
	
	return nil
}
