package hackerone

type Programs struct {
	ProgramsData []*ProgramsData `json:"data"`
	Links        *NextPrograms   `json:"links"`
}

type ProgramsData struct {
	ProgramsAttributes *ProgramsAttributes `json:"attributes"`
}

type ProgramsAttributes struct {
	Handle       string `json:"handle"`
	State        string `json:"state"`
	OffersBounty bool   `json:"offers_bounties"`
}

type NextPrograms struct {
	Next string `json:"next"`
}

type Scope struct {
	ProgramData []*ScopeData `json:"data"`
}

// type ProgramData struct {
// 	Relationships *Relationship `json:"relationships"`
// }

// // type Relationship struct {
// // 	StructuredScopes *StructuredScopes `json:"structured_scopes"`
// // }

// // type StructuredScopes struct {
// // 	ScopeData []*ScopeData `json:"data"`
// // }

type ScopeData struct {
	Attributes *ScopeAttributes `json:"attributes"`
}

type ScopeAttributes struct {
	AssetType         string `json:"asset_type"`
	Identifier        string `json:"asset_identifier"`
	EligibleForBounty bool   `json:"eligible_for_bounty"`
}
