package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidator handles JSON schema validation for VSCU requests
type SchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
	logger  Logger
}

// Logger interface for dependency injection
type Logger interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// NewSchemaValidator creates a new schema validator and loads all schemas
func NewSchemaValidator(schemasDir string, logger Logger) (*SchemaValidator, error) {
	validator := &SchemaValidator{
		schemas: make(map[string]*gojsonschema.Schema),
		logger:  logger,
	}

	// Load all schemas
	err := validator.loadSchemas(schemasDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load schemas: %w", err)
	}

	return validator, nil
}

// loadSchemas loads all JSON schema files from the schemas directory
func (v *SchemaValidator) loadSchemas(schemasDir string) error {
	// Walk through all JSON files in the schemas directory
	err := filepath.Walk(schemasDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			// Load schema file
			schemaBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read schema file %s: %w", path, err)
			}

			// Convert example schema to proper JSON schema
			var schemaMap map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schemaMap); err != nil {
				v.logger.Warnf("Failed to parse schema file %s: %v", path, err)
				return nil
			}

			// Create proper JSON schema structure
			schema := map[string]interface{}{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"type":    "object",
				"properties": map[string]interface{}{},
				"required": []string{},
			}

			properties := schema["properties"].(map[string]interface{})
			required := schema["required"].([]string)

			// Convert example values to type definitions
			for key, value := range schemaMap {
				propSchema := map[string]interface{}{}
				
				switch value.(type) {
				case string:
					propSchema["type"] = "string"
				case float64:
					propSchema["type"] = "number"
				case bool:
					propSchema["type"] = "boolean"
				case nil:
					propSchema["type"] = "null"
				default:
					propSchema["type"] = "string" // default to string for unknown types
				}

				properties[key] = propSchema
				if value != nil {
					required = append(required, key)
				}
			}

			schema["required"] = required

			// Create schema loader
			schemaLoader := gojsonschema.NewGoLoader(schema)
			compiledSchema, err := gojsonschema.NewSchema(schemaLoader)
			if err != nil {
				v.logger.Warnf("Failed to compile schema %s: %v", path, err)
				return nil
			}

			// Store schema using relative path as key
			relPath, err := filepath.Rel(schemasDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			v.schemas[relPath] = compiledSchema
		}

		return nil
	})

	return err
}

// ValidateRequest validates a request against its corresponding schema
func (v *SchemaValidator) ValidateRequest(request interface{}, schemaPath string) error {
	// Get schema
	schema, exists := v.schemas[schemaPath]
	if !exists {
		v.logger.Warnf("Schema not found: %s", schemaPath)
		return nil // Don't block API calls if schema is missing
	}

	// Convert request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		v.logger.Errorf("Failed to marshal request: %v", err)
		return nil // Don't block API calls if marshaling fails
	}

	// Validate
	documentLoader := gojsonschema.NewStringLoader(string(requestJSON))
	result, err := schema.Validate(documentLoader)
	if err != nil {
		v.logger.Errorf("Validation error: %v", err)
		return nil // Don't block API calls if validation fails
	}

	// Check validation result
	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		v.logger.Warnf("Schema validation failed: %s", strings.Join(errors, ", "))
		return nil // Don't block API calls if validation fails
	}

	return nil
}

// GetSchemaPath returns the schema path for a given request type
func GetSchemaPath(requestType interface{}) string {
	t := reflect.TypeOf(requestType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Name() {
	case "InitRequest":
		return "VSCU-Initialization/Initialization_Request_request_schema.json"
	case "SalesRequest":
		return "VSCU-Sales_Management/Send_sales_transaction_information_request_schema.json"
	case "BranchCustomerRequest":
		return "VSCU-Branch_information_management/Save_branch_customer_information_request_schema.json"
	case "ItemRequest":
		return "VSCU-Item_Management/Save_item_information_request_schema.json"
	case "StockRequest":
		return "VSCU-Stock_Information_Management/Save_stock_information_request_schema.json"
	case "PurchaseRequest":
		return "VSCU-Purchase_management/Send_purchase_transaction_information_request_schema.json"
	case "ImportedItemSaveRequest":
		return "VSCU-Imports_Item_Management/Save_imported_item_information_request_schema.json"
	case "CodeListRequest":
		return "VSCU-_Basic_Data_Management/Get_code_list_request_schema.json"
	case "ItemClassRequest":
		return "VSCU-_Basic_Data_Management/Get_item_classification_list_request_schema.json"
	case "CustomerRequest":
		return "VSCU-_Basic_Data_Management/Get_customer_information_request_schema.json"
	case "BranchListRequest":
		return "VSCU-Branch_information_management/Get_branch_list_request_schema.json"
	case "NoticeRequest":
		return "VSCU-_Basic_Data_Management/Get_notice_list_request_schema.json"
	case "ItemCompositionRequest":
		return "VSCU-Item_Management/Save_item_composition_request_schema.json"
	case "PurchaseConfirmationRequest":
		return "VSCU-Purchase_management/Send_purchase_transaction_confirmation_request_schema.json"
	default:
		return ""
	}
}
