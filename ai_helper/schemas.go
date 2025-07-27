package ai_helper

import "google.golang.org/genai"

// --- genai.Schema Definitions ---


func medicationTypes() []string {
	return []string {"Antibiotic", "Antiemetic", "Antifungal", "Antihypertensive", "Antineoplastic", "Antipyretic", "Antiviral", "Bronchodilator", "Diuretic", "Immunosuppressant", "Narcotic", "NSAID", "Steroid", "Other"}
}

func routeEnum() []string{
	return []string{"iv", "oral", "sc", "im", "topical", "inhalation", "unknown"}
}

// Schema for ArticleReference
func articleReferenceSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":   {Type: genai.TypeString, Description: "The title of the scientific article."},
			"authors": {Type: genai.TypeString, Description: "The authors of the article (e.g., 'Doe J, Smith A')."},
			"journal": {Type: genai.TypeString, Description: "The journal where the article was published."},
			"year":    {Type: genai.TypeString, Description: "The publication year (e.g., '2023')."},
			"pmid":    {Type: genai.TypeString, Description: "PubMed ID, if available."},
			"doi":     {Type: genai.TypeString, Description: "Digital Object Identifier, if available."},
		},
		Required: []string{"title", "authors", "journal", "year"},
	}
}

// Schema for MedicationModification
func medicationModificationSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"description": {Type: genai.TypeString, Description: "Description of the impairment level (e.g., 'Mild', 'Moderate', 'Severe')."},
			"adjustment":  {Type: genai.TypeString, Description: "Recommended dose adjustment or modification (e.g., 'Reduce dose by 50%')."},
		},
		Required: []string{"description", "adjustment"},
	}
}

// Schema for MedicationModificationCategory
func medicationModificationCategorySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"category": {
				Type:        genai.TypeString,
				Description: "Category of impairment (e.g., 'Hepatic Impairment', 'Renal Impairment').",
				Format:      "enum",
				Enum:        []string{"hepatic_impairment", "renal_impairment"},
			},
			"modifications": {
				Type:        genai.TypeArray,
				Items:       medicationModificationSchema(),
				Description: "Array of specific modifications for this impairment category.",
			},
		},
		Required: []string{"category", "modifications"},
	}
}

// Schema for Physician
func physicianSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"first_name": {Type: genai.TypeString, Description: "The physician's first name."},
			"last_name":  {Type: genai.TypeString, Description: "The physician's last name."},
		},
		Required: []string{"first_name", "last_name"},
	}
}

// Schema for SummaryProtocol
func summaryProtocolSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"tumor_group": {Type: genai.TypeString,
				Description: "The tumor group associated with the protocol.",
				Format:      "enum",
				Enum:        []string{"breast", "lung", "gastrointestinal", "genitourinary", "head_and_neck", "gynecology", "sarcoma", "leukemia", "bmt", "neuro-oncology", "ocular", "skin", "unknown_primary", "lymphoma", "myeloma", "unknown"},
			},
			"code":                {Type: genai.TypeString, Description: "The protocol's unique code."},
			"name":                {Type: genai.TypeString, Description: "The full name of the protocol."},
			"tags":                {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Tags associated with the protocol... e.g., 'B-Cell Lymphoma', 'Follicular Lymphoma', 'DLBCL', 'R-CHOP', 'Rituximab', 'Cyclophosphamide', 'Doxorubicin', 'Vincristine', 'Prednisone',etc."},
			"notes":               {Type: genai.TypeString, Description: "Any additional notes or comments about the protocol."},
			"protocol_url":        {Type: genai.TypeString, Description: "URL to the protocol document or resource."},
			"patient_handout_url": {Type: genai.TypeString, Description: "URL to the patient handout document or resource."},
			"activated_on":        {Type: genai.TypeString, Description: "Date when the protocol was activated, format: YYYY-MMM-DD (e.g., '2023-Jan-15').", Format: "date-time"}, // Use 'date' format hint
			"revised_on":          {Type: genai.TypeString, Description: "Date when the protocol was last revised, format: YYYY-MMM-DD (e.g., '2024-Feb-29').", Format: "date-time"},
		},
		Required: []string{"tumor_group", "code", "name", "tags", "activated_on", "revised_on"},
	}
}

// Schema for ProtocolEligibilityCriterion
func protocolEligibilityCriterionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"type": {
				Type:        genai.TypeString,
				Description: "Type of criterion: 'inclusion', 'exclusion', or 'unknown'. Each bullet point should be a separate object.",
				Format:      "enum",
				Enum:        []string{"inclusion", "exclusion", "unknown"},
			},
			"description": {Type: genai.TypeString, Description: "The detailed description of the eligibility criterion."},
		},
		Required: []string{"type", "description"},
	}
}

// Schema for ProtocolPrecaution
func protocolPrecautionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":       {Type: genai.TypeString, Description: "A concise title for the precaution (e.g., 'Myelosuppression')."},
			"description": {Type: genai.TypeString, Description: "Detailed description of the precaution and management."},
		},
		Required: []string{"title", "description"},
	}
}

// Schema for ProtocolCaution
func protocolCautionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"description": {Type: genai.TypeString, Description: "Detailed description of the caution."},
		},
		Required: []string{"description"},
	}
}

// Schema for CycleTreatment
func cycleTreatmentSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"medication_name":        {Type: genai.TypeString, Description: "Name of the medication used in this treatment."},
			"medication_description": {Type: genai.TypeString, Description: "Description of the medication."},
			"medication_category": {
				Type:        genai.TypeString,
				Description: "Category of the medication (e.g., 'Antineoplastic').",
				Enum:        medicationTypes(),
			},
			"medication_alternate_names": {
				Type:        genai.TypeArray,
				Items:       &genai.Schema{Type: genai.TypeString},
				Description: "Alternate names for the medication, if available.",
			},
			"dose": {Type: genai.TypeString, Description: "Dosage of the medication (e.g., '100 mg/m2')."},
			"route": {
				Type:        genai.TypeString,
				Description: "Route of administration (e.g., 'iv', 'oral', 'sc', 'im', 'topical', 'inhalation', 'unknown').",
				Enum:        routeEnum(),
			},
			"frequency":            {Type: genai.TypeString, Description: "Frequency of administration (e.g., 'Day 1-2 ', 'Day 1, 8, 15 and 22', 'Day 1 to 14')."},
			"duration":             {Type: genai.TypeString, Description: "Duration of administration (e.g., 'every 28 days')."},
			"administration_guide": {Type: genai.TypeString, Description: "Specific administration guidelines."},
		},
		Required: []string{"medication_name", "medication_description", "medication_category", "dose", "route", "frequency", "duration", "administration_guide"},
	}
}

// Schema for Prescriptions
func prescriptionSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"medication_name":        {Type: genai.TypeString, Description: "Name of the medication used in this prescription."},
			"medication_description": {Type: genai.TypeString, Description: "Description of the medication."},
			"medication_category": {
				Type:        genai.TypeString,
				Description: "Category of the medication (e.g., 'Antineoplastic').",
				Enum:        medicationTypes(),
			},
			"medication_alternate_names": {
				Type:        genai.TypeArray,
				Items:       &genai.Schema{Type: genai.TypeString},
				Description: "Alternate names for the medication, commercial or otherwise.",
			},
			"dose": {Type: genai.TypeString, Description: "Dosage of the medication (e.g., '500mg')."},
			"route": {
				Type:        genai.TypeString,
				Description: "Route of administration (e.g., 'iv', 'oral', 'sc', 'im', 'topical', 'inhalation', 'unknown').",
				Enum:        routeEnum(),
			},
			"frequency":            {Type: genai.TypeString, Description: "Frequency of administration (e.g., 'Start 3 days before chemotherapy and continue for 15 days', 'Start day 7 post-chemotherapy and continue daily for 7 days','30 minutes pre-chemotherapy',etc.)"},
			"duration":             {Type: genai.TypeString, Description: "Duration of administration (e.g., '7 doses every 21 days, 30 tabs, 120 tabs,etc.')."},
			"instructions": {Type: genai.TypeString, Description: "Specific instruction regarding medication use (e.g. 'use as necessary if bone pain associated with filgrastim', 'use 2 tabs after loose stools, and 1 tab after each loose stool afterward',etc.)."},
			"renewals": {Type:genai.TypeNumber,Description:"Number of renewals, if unknown give an estimate based on number of cycles and known information."},
		},
		Required: []string{"medication_name", "medication_description", "medication_category", "dose", "route", "frequency", "duration", "instructions","renewals"},
	}
}

// Schema for TestGroup
func protocolTestGroup() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"category": {
				Type:        genai.TypeString,
				Description: "The category defining this grouping of tests in the structure of the protocol (e.g. Required Pre-treatment Tests, Day 1-4 Tests, Subsequent Pre-treatment Tests, If Clinicalled Indicated Tests, etc.).",
			},
			"comments": {
				Type:        genai.TypeString,
				Description: "A brief description of the category if necessary to better understand..",
			},
			"position":{
				Type: genai.TypeInteger,
				Description: "The order that the test groups are displayed/organized on the front end. Should follow chronological order : '1, 2, 3, etc.'",
			},
			"tests": {
				Type:        genai.TypeArray,
				Items:       testDetailSchema(),
				Description: "List of tests that are part of this category.",
			},
		},
		Required: []string{"category", "comments", "tests"},
	}
}

// Schema for TestGroup
func protocolMedGroups() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"category": {
				Type:        genai.TypeString,
				Description: "The category defining this grouping of prescribed medications (NOT protocol treatments. This section is reserved for pre-medication, supportive medication etc.)",
			},
			"comments": {
				Type:        genai.TypeString,
				Description: "A brief description of the category if necessary to better understand.",
			},
			"prescriptions": {
				Type:        genai.TypeArray,
				Items:       prescriptionSchema(),
				Description: "List of prescriptions that are part of this category.",
			},
		},
		Required: []string{"category", "comments", "prescriptions"},
	}
}

// Schema for ProtocolCycle
func protocolCycleSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"cycle": {
				Type:        genai.TypeString,
				Description: "The cycle number (e.g., 'Cycle 1'). If not specified in source, default to 'Cycle 1+'.",
			},
			"cycle_duration": {
				Type:        genai.TypeString,
				Description: "Duration of the cycle (e.g., '28 days'). If blank in source, default to '28 days'.",
			},
			"treatments": {
				Type:        genai.TypeArray,
				Items:       cycleTreatmentSchema(),
				Description: "List of treatments administered during this cycle.",
			},
		},
		Required: []string{"cycle", "cycle_duration", "treatments"},
	}
}

// Schema for TestDetail
func testDetailSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name":        {Type: genai.TypeString, Description: "Name of the test (e.g., 'CBC', 'Creatinine','Electrolytes','Calcium','ALT','AST')."},
			"description": {Type: genai.TypeString, Description: "Brief description or purpose of the test."},
		},
		Required: []string{"name", "description"},
	}
}


// Schema for ToxicityModification
func toxicityModificationSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"grade": {
				Type:        genai.TypeString, // Use string as enum for exact values "1", "2", "3", "4"
				Description: "The CTCAE grade number (1, 2, 3, or 4).",
				Format:      "enum",
				Enum:        []string{"1", "2", "3", "4"},
			},
			"grade_description": {Type: genai.TypeString, Description: "Description of the grade using CTCAE v5 terminology."},
			"adjustment": {
				Type:        genai.TypeString,
				Description: "Recommended adjustment (e.g., 'Dose reduction', 'Delay', 'Discontinuation'). Leave blank if no information.",
			},
		},
		Required: []string{"grade", "grade_description"},
	}
}

// Schema for Toxicity
func toxicitySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type:        genai.TypeString,
				Description: "Title of the toxicity (e.g., 'Neuropathy', 'Thrombopenia', 'Diarrhea').",
			},
			"description": {Type: genai.TypeString, Description: "Detailed description of the toxicity."},
			"category": {
				Type:        genai.TypeString,
				Description: "Category of the toxicity (e.g., 'Hematologic', 'Neurologic', 'Gastrointestinal').",
				Enum:        []string{"Hematologic", "Neurologic", "Gastrointestinal", "Dermatologic", "Hepatic", "Renal", "Pulmonary", "Cardiovascular", "Endocrine", "Metabolic", "Immune", "Other"},
			},
			"modifications": {
				Type:        genai.TypeArray,
				Items:       toxicityModificationSchema(),
				Description: "Array of modifications for each grade (1, 2, 3, and 4). This array must contain exactly 4 objects, one for each grade.",
				MinItems:    func() *int64 { v := int64(4); return &v }(), // Enforce exactly 4 modifications
				MaxItems:    func() *int64 { v := int64(4); return &v }(), // Enforce exactly 4 modifications
			},
		},
		Required: []string{"title", "description", "category", "modifications"},
	}
}

// Root schema for the entire ProtocolData
func protocolDataSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"article_references": {
				Type:        genai.TypeArray,
				Items:       articleReferenceSchema(),
				Description: "References to scientific articles relevant to the protocol.",
			},			
			"physicians": {
				Type:        genai.TypeArray,
				Items:       physicianSchema(),
				Description: "List of physicians associated with the protocol.",
			},
			"summary_protocol": summaryProtocolSchema(),
			"protocol_eligibility_criteria": {
				Type:        genai.TypeArray,
				Items:       protocolEligibilityCriterionSchema(),
				Description: "Inclusion and exclusion criteria for the protocol.",
			},
			"protocol_precautions": {
				Type:        genai.TypeArray,
				Items:       protocolPrecautionSchema(),
				Description: "Precautions to be taken during the protocol.",
			},
			"protocol_cautions": {
				Type:        genai.TypeArray,
				Items:       protocolCautionSchema(),
				Description: "Cautions to be observed during the protocol.",
			},
			"protocol_cycles": {
				Type:        genai.TypeArray,
				Items:       protocolCycleSchema(),
				Description: "Details of treatment cycles within the protocol.",
			},
			"prescription_groups":{
				Type: genai.TypeArray,
				Items: protocolMedGroups(),
				Description:"Categories of supportive or premedication prescriptions used within the protocol.",
			},
			"test_groups": {
				Type: genai.TypeArray,
				Items:protocolTestGroup(),
				Description:"Categories of lab (or other) tests required or suggested for monitoring prior or during the administration of this protocol.",
			},
			"toxicities": {
				Type:        genai.TypeArray,
				Items:       toxicitySchema(),
				Description: "Information on potential toxicities and their management.",
			},
		},
		Required: []string{
			"article_references",			
			"physicians",
			"summary_protocol",
			"protocol_eligibility_criteria",
			"protocol_precautions",
			"protocol_cautions",
			"protocol_cycles",
			"prescription_groups",
			"test_groups",
			"toxicities",
		},
	}
}
