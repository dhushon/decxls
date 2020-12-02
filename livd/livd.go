package livd

import (
	"time"
)

// AcronymDict provides a LUT for common abbreviations/acronyms
// TODO: shift to database table?
var AcronymDict map[string]string = map[string]string{"Ab": "Antibody",
	"ACnc":              "Arbitrary concentration",
	"Ag":                "Antigen",
	"APHL":              "Association of Public Health Laboratories",
	"ART":               "Antiretroviral therapy",
	"ARUP":              "ARUP Laboratories (name of a compnay)",
	"Bld":               "Blood",
	"CDC":               "US Centers for Disease Control and Prevention",
	"CHIV":              "HIV Ag/Ab combo",
	"ChLIA":             "Chemiluminescent immunoassay",
	"CIA":               "Chemiluminescent immunoassay",
	"CLIA":              "Clinical Laboratory Improvement Amendments",
	"CMIA":              "Chemi-luminometric immunoassay",
	"dB":                "Database",
	"EIA":               "Enzyme immunoassay",
	"Env":               "Envelope",
	"FDA":               "Food and Drug Administration",
	"GAG":               "HIV-1 GAG protein",
	"gp":                "Glycoprotein",
	"IA":                "Immunoassay",
	"IA.rapid":          "Rapid immunoassay",
	"IB":                "Immunoblot",
	"ID":                "Identifier",
	"IF":                "Immunofluorescence",
	"IFA":               "Immunofluorescence assay",
	"gG":                "Immunoglobulin G",
	"IgM":               "Immunoglobulin M",
	"IICC":              "IVD Industry Connectivity Consortium",
	"IU":                "International units",
	"IVD":               "In vitro diagnostic",
	"LnCnc":             "Log number concentration",
	"log#":              "Log number",
	"LOINC":             "Logical observation identifiers Names and Codes",
	"mL":                "Milliliter",
	"NAA":               "Nucleic acid amplification",
	"NCnc":              "Number concentration (count/vol)",
	"NFr":               "Number fraction",
	"Nom":               "Nominal",
	"Non-probe.amp.tar": "Target amplification followed by non-probe based detection",
	"NYS":               "New York state",
	"Ord":               "Ordinal",
	"PAML":              "Pathology Associates Medical Laboratories",
	"PCR":               "Polymerase chain reaction",
	"PHA":               "Public health agency",
	"Plas":              "Plasma",
	"Prid":              "Presence or identity",
	"Probe.amp.tar":     "Target amplified probe",
	"PrThr":             "Presence / threshold",
	"Pt":                "Point in time",
	"Qn":                "Quantitative",
	"qual":              "Qualitative",
	"RNA":               "Ribonucleic acid",
	"Seq":               "Nucleotide sequence",
	"Ser":               "Serum",
	"SerPl":             "Serum plasma",
	"Susc":              "Susceptibility",
	"TMA":               "Transcription-mediated amplification of nucleic acid",
	"ug":                "Microgram",
	"UID":               "Universal identifier",
	"uL":                "Microliter",
	"Vol":               "Volume",
	"WB":                "Westernblot",
	"XXX":               "Not specified",
	"NP":                "Nasopharyngeal (usuallys swab, but can also be washings/aspirate)",
	"OP":                "Oropharyngeal (usuallys swab)",
	"NMT":               "Nasal mid-turbinate swab",
	"AN":                "Anterior nares swab",
	"NS":                "Nasal Swab",
	"BAL":               "Broncho-alvelar Lavage",
	"SARS":              "Severe Acute Respiratory Syndrome",
	"CoV":               "Coronavirus",
	"COVID-19":          "Coronavirus Disease"}

// ResultDescription provides a mapping between LIVD/language and SnoMedCT codes
type ResultDescription struct {
	Result string   `json:"result,omitempty"`
	SnoMedCT SnoMedCT `json:"snomedct"`
}
// SpecimenDescription provides a mapping between LIVD/language and SnoMedCT codes
type SpecimenDescription struct {
	SpecimenType string   `json:"specimen_type,omitempty"`
	Units    string   `json:"specimen_units,omitempty"`
	SnowMedCT SnoMedCT `json:"snomedct"`
}

//LTime supports the conversion of time should we run into problems with json parsing or expectations
type LTime time.Time

//Livd or LOINC In Vitro Diagnostic (LIVD) Test Code Mapping for SARS-CoV-2 Test results
type Livd struct {
	//Manufacturer - the name of the Instrument manufacturer or Testkit manufacturer
	Manufacturer string `json:"manufacturer" xls:"Manufacturer"`
	//Model – the name of the Instrument or the Testkit
	Model string `json:"model" xls:"Model"`
	//Vendor Analyte Name – the human-readable text the vendor uses to identify the analyte. The text might be displayed by the instrument or could be used within an assay insert or instructions for use.
	VAnalyteName string `json:"vendor_analyte_name" xls:"Vendor Analyte Name"`
	//Vendor Specimen Description – the human-readable text that provides information about the specimen used for the test, such as “Serum or Plasma.” The field documents the vendor description of the specimen used for the IVD test. The field may contain multiple specimen types. For this catalog, the SNOMED-CT mapping for each Specimen is provided in parenthesis
	//e.g. o   Serum (119364003^Serum specimen^SCT)
	//e.g. o   Plasma (119361006^Plasma specimen^SCT)
	VSSpecimenDescription string `json:"summary_vendor_specimen_description" xls:"Vendor Specimen Description"`// TODO: build TypeXLSUnmarshaller and Marshaller to convert
	VSpecimentDescription []*SpecimenDescription `json:"vendor_specimen_descriptions"` 
	//Vendor Result Description – the human-readable text that provides information about the result that is produced.
	//o   For non-numeric results, this field contains the possible result values, along with the SNOMED-CT mapping
	//	§  Positive (10828004^Positive^SCT)
	//	§  Negative (260385009^Negative^SCT)
	//o   For numeric results and associated units of measure this field describes the result by including a representative unit of measure, preferably represented as a UCUM unit.
	VSResultDescription string `json:"vendor_result_description" xls:"Vendor Result Description"` // TODO: build TypeXLSUnmarshaller and Marshaller to convert
	VResultDescription []*ResultDescription `json:"vendor_result_descriptions"`

	Loinc              *Loinc          `json:"loinc"`
	//Testkit Name ID – the Testkit Device Identifier if known, otherwise the concatenation of the <Diagnostic (Letter of Authorization)>_<Manufacturer> fields from the FDA EUA Website
	TestkitID string `json:"testkit_id" xls:"Testkit Name ID"`
	//Testkit Name ID Type – GTIN when the Device Identifier is known, otherwise EUA when the Emergency Use Authorization identifiers are used
	TestkitIDType string `json:"testkit_id_type" xls:"Testkit Name ID Type"`
	//Vendor Analyte Code – one of two possible values
	//  o   For an automated test, it contains Vendor Transmission Code used by the instrument when sending the test result to a health information system, such as an LIS.
	//  o   For a manual test, it is the Vendor Analyte Identifier for the test result produced by the Test Kit.
	VAnalyteCode string `json:"vendor_analyte_code" xls:"Vendor Analyte Code"`
	//Vendor Reference ID – a vendor identifier, such as an identifier that can be used to locate the associated assay insert published by the vendor. This attribute may contain the material number used to order the product from the manufacturer.
	VReferenceID string `json:"vendor_reference_id" xls:"Vendor Reference ID"`

	//Equipment UID – when available, the Device Identifier of the instrument. If not known, then the concatenation of the Model and Manufacturer,  <Model>_<Manufacturer>
	EquipmentUID string `json:"equipment_uid" xls:"Equipent UID"`
	//Equipment UID Type – GTIN when the Device Identifier is provided, otherwise MNI when model and manufacturer are provided
	EquipmentUIDType string `json:"equipment_uid_type" xls:"Equipment UID Type"`
	//Component – the LOINC Component/analyte
	Component string `json:"component" xls:"Component"`
	//Property – the LOINC Kind of Property
	Property string `json:"property" xls:"Property"`
	//Time – the LOINC Time Aspect
	Time string `json:"time" xls:"Time"`
	//System – The LOINC System (Sample)
	System string `json:"system" xls:"System"`
	//Scale – the LOINC Type of Scale
	Scale string `json:"scale" xls:"Scale"`
	//Method – the LOINC Type of Method
	Method string `json:"method" xls:"Method"`
	//Publication Version ID – human-readable information used to differentiate mapping publication versions
	Publication string `json:"publication_version_id" xls:"Publication Version ID"`
	//LOINC Version ID – the version of LOINC used to establish the LOINC mapping
	Version string `json:"loinc_version" xls:"LOINC Version ID"`
}
