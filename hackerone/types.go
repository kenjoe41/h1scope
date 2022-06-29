package hackerone

type Scope struct {
	Relationships *Relationship `json:"relationships"`
}

type Relationship struct {
	StructuredScopes *StructuredScopes `json:"structured_scopes"`
}

type StructuredScopes struct {
	Data []*Data `json:"data"`
}

type Data struct {
	Attributes *Attributes `json:"attributes"`
}

type Attributes struct {
	AssetType         string `json:"asset_type"`
	Identifier        string `json:"asset_identifier"`
	EligibleForBounty bool   `json:"eligible_for_bounty"`
}
