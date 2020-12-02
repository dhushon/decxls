package livd

//Loinc supports the typical Loinc information for reporting "logical observation identifiers namesa and codes" - largely a holding structure at this point.
type Loinc struct {
	//LOINC® Code – the appropriate LOINC for the mapping
	LoincCode string `json:"loinc_code" xls:"LOINC Code"`
	//LOINC Long Name – the long name for the LOINC Code
	LoincLongName string `json:"loinc_long_name" xls:"LOINC Long Name"`
	//LOINC Order Code – the appropriate LOINC Order code for the test result
	LoincOrderCode string `json:"loinc_order_code" xls:"LOINC Order Code"`
	//LOINC Order Code Long Name– the long name for the LOINC Order Code
	LoincOrderCodeLongName string `json:"loinc_order_code_longname" xls:"LOINC Order Code Long Name"`
	//Vendor Comment – human-readable text for clarification, such as “This is a STAT (prioritized) version of the test”
	LoincVendorComment string `json:"loinc_vendor_comment" xls:"LOINC Vendor Comment"`
}
