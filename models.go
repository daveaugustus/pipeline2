package main

type ActionOps int

const (
	Insert ActionOps = 1 + iota
	Skip
	Delete
	Update
)

//
type Result struct {
	Meta         Meta         `json:"meta"`
	ParsedResult ParsedResult `json:"parsed_result"`
}

type Meta struct {
	StageResults []StageResult `json:"stage_results"`
	ZipFile      string        `json:"zip_file"`
	UnzipFolder  string        `json:"unzip_folder"`
}

// StageResult But we can utilize to add some other info:
// ParseOrgs -- Number of Skipped, Added, Updated, Deleted
// Same for Users, Association
type StageResult struct {
	StageName string `json:"stage_name"`
	IsSuccess bool   `json:"is_success"`
	Failure   error  `json:"failure"`
}

// ParsedResult is basically composed from three fields
type ParsedResult struct {
	// Orgs will be parsed from unzip folder
	Orgs []Org `json:"orgs"`
	// Orgs will be parsed from unzip folder
	Users     []User                  `json:"users"`
	OrgsUsers []OrgsUsersAssociations `json:"orgs_users_associations"`
}

type OrgsUsersAssociations struct {
	OrgName Org               `json:"org_name"`
	Users   []UserAssociation `json:"user_association"`
}
type UserAssociation struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type KeyDump struct {
	ID                            string      `json:"id"`
	AuthzID                       string      `json:"authz_id"`
	Username                      string      `json:"username"`
	Email                         string      `json:"email"`
	PubkeyVersion                 int         `json:"pubkey_version"`
	PublicKey                     interface{} `json:"public_key"`
	SerializedObject              string      `json:"serialized_object"`
	LastUpdatedBy                 string      `json:"last_updated_by"`
	CreatedAt                     string      `json:"created_at"`
	UpdatedAt                     string      `json:"updated_at"`
	ExternalAuthenticationUID     interface{} `json:"external_authentication_uid"`
	RecoveryAuthenticationEnabled interface{} `json:"recovery_authentication_enabled"`
	Admin                         bool        `json:"admin"`
	HashedPassword                interface{} `json:"hashed_password"`
	Salt                          interface{} `json:"salt"`
	HashType                      interface{} `json:"hash_type"`
}

type Org struct {
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	ActionOps ActionOps `json:"guid"`
}

type User struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	DisplayName      string `json:"display_name"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	MiddleName       string `json:"middle_name"`
	AutomateUsername string `json:"automate_username"`
	Connector        string `json:"connector"`
	IsConflicting    bool   `json:"is_conflicting"`
	IsAdmin          bool   `json:"is_admin"`
}
