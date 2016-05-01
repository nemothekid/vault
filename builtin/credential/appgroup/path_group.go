package appgroup

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type groupStorageEntry struct {
	Apps               []string      `json:"apps" structs:"apps" mapstructure:"apps"`
	AdditionalPolicies []string      `json:"additional_policies" structs:"additional_policies" mapstructure:"additional_policies"`
	NumUses            int           `json:"num_uses" structs:"num_uses" mapstructure:"num_uses"`
	TTL                time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL             time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Wrapped            time.Duration `json:"wrapped" structs:"wrapped" mapstructure:"wrapped"`
}

func groupPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name"),
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"apps": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "",
					Description: "Comma separated list of Apps belonging to the group",
				},
				"additional_policies": &framework.FieldSchema{
					Type:    framework.TypeString,
					Default: "",
					Description: `Comma separated list of policies for the Group. The UserID created against the Group,
will have access to the union of all the policies of the Apps. In
addition to those, a set of policies can be assigned using this.
`,
				},
				"num-uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Number of times the a UserID can access the Apps represented by the Group.",
				},
				"ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: "Duration in seconds after which a UserID should expire.",
				},
				"max_ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: "MaxTTL of the UserID created on the App.",
				},
				"wrapped": &framework.FieldSchema{
					Type: framework.TypeDurationSecond,
					Description: `Duration in seconds, if specified, enables Cubbyhole mode. In this mode, a
UserID will not be returned. Instead a new token will be returned. This token
will have the UserID stored in its Cubbyhole. The value represented by 'wrapped'
will be the duration after which the returned token expires.
`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathGroupCreateUpdate,
				logical.UpdateOperation: b.pathGroupCreateUpdate,
				logical.ReadOperation:   b.pathGroupRead,
				logical.DeleteOperation: b.pathGroupDelete,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/apps$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"apps": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "",
					Description: "Comma separated list of Apps belonging to the group",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupAppsUpdate,
				logical.ReadOperation:   b.pathGroupAppsRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-apps"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-apps"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/additional-policies$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"additional_policies": &framework.FieldSchema{
					Type:    framework.TypeString,
					Default: "",
					Description: `Comma separated list of policies for the Group. The UserID created against the Group,
will have access to the union of all the policies of the Apps. In
addition to those, a set of policies can be assigned using this.
`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupAdditionalPoliciesUpdate,
				logical.ReadOperation:   b.pathGroupAdditionalPoliciesRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-additional-policies"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-additional-policies"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/num-uses$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"num-uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Number of times the a UserID can access the Apps represented by the Group.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupNumUsesUpdate,
				logical.ReadOperation:   b.pathGroupNumUsesRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-num-uses"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-num-uses"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/ttl$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: "Duration in seconds after which a UserID should expire.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupTTLUpdate,
				logical.ReadOperation:   b.pathGroupTTLRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-ttl"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-ttl"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/max-ttl$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"max_ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: "MaxTTL of the UserID created on the App.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupMaxTTLUpdate,
				logical.ReadOperation:   b.pathGroupMaxTTLRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-max-ttl"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-max-ttl"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/wrapped$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"wrapped": &framework.FieldSchema{
					Type: framework.TypeDurationSecond,
					Description: `Duration in seconds, if specified, enables Cubbyhole mode. In this mode, a
UserID will not be returned. Instead a new token will be returned. This token
will have the UserID stored in its Cubbyhole. The value represented by 'wrapped'
will be the duration after which the returned token expires.
`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupWrappedUpdate,
				logical.ReadOperation:   b.pathGroupWrappedRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-wrapped"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-wrapped"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/creds$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathGroupCredsRead,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-creds"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-creds"][1]),
		},
		&framework.Path{
			Pattern: "group/" + framework.GenericNameRegex("group_name") + "/creds-specific$",
			Fields: map[string]*framework.FieldSchema{
				"group_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the Group.",
				},
				"user_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "",
					Description: "UserID to be attached to the App.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathGroupCredsSpecificUpdate,
			},
			HelpSynopsis:    strings.TrimSpace(groupHelp["group-creds-specific"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-creds-specific"][1]),
		},
	}
}

func (b *backend) pathGroupCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	// Check if there is already an entry. If entry exists, this is an
	// UpdateOperation.
	group, err := groupEntry(req.Storage, groupName)
	if err != nil {
		return nil, err
	}

	// If entry does not exist, this is a CreateOperation. So, create
	// a new object.
	if group == nil {
		group = &groupStorageEntry{}
	}

	// Update only if value is supplied. Defaults to zero.
	if appsRaw, ok := data.GetOk("apps"); ok {
		group.Apps = strings.Split(appsRaw.(string), ",")
	}

	// Update only if value is supplied. Defaults to zero.
	if additionalPoliciesRaw, ok := data.GetOk("additional_policies"); ok {
		group.AdditionalPolicies = policyutil.ParsePolicies(additionalPoliciesRaw.(string))
	}

	// Update only if value is supplied. Defaults to zero.
	if numUsesRaw, ok := data.GetOk("num_uses"); ok {
		group.NumUses = numUsesRaw.(int)
	}

	// If TTL value is not provided either during update or create, don't bother.
	// Core will set the system default value if the policies does not contain
	// "root" and TTL value is zero.
	// Update only if value is supplied. Defaults to zero.
	if ttlRaw, ok := data.GetOk("ttl"); ok {
		group.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}

	// Update only if value is supplied. Defaults to zero.
	if maxTTLRaw, ok := data.GetOk("max_ttl"); ok {
		group.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	}

	// Check that TTL value provided is greater than MaxTTL.
	//
	// Do not sanitize the TTL and MaxTTL now, just store them as-is.
	// Check the System TTL and MaxTTL values at credential issue time
	// and act accordingly.
	if group.TTL > group.MaxTTL {
		return logical.ErrorResponse("ttl should not be greater than max_ttl"), nil
	}

	// Update only if value is supplied. Defaults to zero.
	if wrappedRaw, ok := data.GetOk("wrapped"); ok {
		group.Wrapped = time.Duration(wrappedRaw.(int)) * time.Second
	}

	// Store the entry.
	return nil, setGroupEntry(req.Storage, groupName, group)
}

func setGroupEntry(s logical.Storage, groupName string, group *groupStorageEntry) error {
	if entry, err := logical.StorageEntryJSON("group/"+strings.ToLower(groupName), group); err != nil {
		return err
	} else {
		return s.Put(entry)
	}
}

func groupEntry(s logical.Storage, groupName string) (*groupStorageEntry, error) {
	if groupName == "" {
		return nil, fmt.Errorf("missing group_name")
	}

	var result groupStorageEntry

	if entry, err := s.Get("group/" + strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	} else if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathGroupRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		// Convert the values to second
		group.TTL = group.TTL / time.Second
		group.MaxTTL = group.MaxTTL / time.Second
		group.Wrapped = group.Wrapped / time.Second

		return &logical.Response{
			Data: structs.New(group).Map(),
		}, nil
	}
}

func (b *backend) pathGroupDelete(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	return nil, req.Storage.Delete("group/" + strings.ToLower(groupName))
}

func (b *backend) pathGroupAppsUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if appsRaw, ok := data.GetOk("apps"); ok {
		group.Apps = strings.Split(appsRaw.(string), ",")
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing apps"), nil
	}
}

func (b *backend) pathGroupAppsRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				"apps": group.Apps,
			},
		}, nil
	}
}

func (b *backend) pathGroupAdditionalPoliciesUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if additionalPoliciesRaw, ok := data.GetOk("additional_policies"); ok {
		group.AdditionalPolicies = policyutil.ParsePolicies(additionalPoliciesRaw.(string))
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing additional_policies"), nil
	}
}

func (b *backend) pathGroupAdditionalPoliciesRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				"additional_policies": group.AdditionalPolicies,
			},
		}, nil
	}
}

func (b *backend) pathGroupNumUsesUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if numUsesRaw, ok := data.GetOk("num_uses"); ok {
		group.NumUses = numUsesRaw.(int)
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing num_uses"), nil
	}
}

func (b *backend) pathGroupNumUsesRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				"num_uses": group.NumUses,
			},
		}, nil
	}
}

func (b *backend) pathGroupTTLUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if ttlRaw, ok := data.GetOk("ttl"); ok {
		if group.TTL = time.Duration(ttlRaw.(int)) * time.Second; group.TTL > group.MaxTTL {
			return logical.ErrorResponse("ttl should not be greater than max_ttl"), nil
		}
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing ttl"), nil
	}
}

func (b *backend) pathGroupTTLRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		group.TTL = group.TTL / time.Second
		return &logical.Response{
			Data: map[string]interface{}{
				"ttl": group.TTL,
			},
		}, nil
	}
}

func (b *backend) pathGroupMaxTTLUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if ttlRaw, ok := data.GetOk("max_ttl"); ok {
		if group.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second; group.TTL > group.MaxTTL {
			return logical.ErrorResponse("max_ttl should be greater than ttl"), nil
		}
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing max_ttl"), nil
	}
}

func (b *backend) pathGroupMaxTTLRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		group.MaxTTL = group.MaxTTL / time.Second
		return &logical.Response{
			Data: map[string]interface{}{
				"max_ttl": group.MaxTTL,
			},
		}, nil
	}
}

func (b *backend) pathGroupWrappedUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	group, err := groupEntry(req.Storage, strings.ToLower(groupName))
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	if wrappedRaw, ok := data.GetOk("wrapped"); ok {
		app.Wrapped = time.Duration(wrappedRaw.(int)) * time.Second
		return nil, setGroupEntry(req.Storage, groupName, group)
	} else {
		return logical.ErrorResponse("missing wrapped"), nil
	}
}

func (b *backend) pathGroupWrappedRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	groupName := data.Get("group_name").(string)
	if groupName == "" {
		return logical.ErrorResponse("missing group_name"), nil
	}

	if group, err := groupEntry(req.Storage, strings.ToLower(groupName)); err != nil {
		return nil, err
	} else if group == nil {
		return nil, nil
	} else {
		group.Wrapped = group.Wrapped / time.Second
		return &logical.Response{
			Data: map[string]interface{}{
				"wrapped": group.Wrapped,
			},
		}, nil
	}
}

func (b *backend) pathGroupCredsRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("pathGroupCredsRead entered\n")
	return nil, nil
}

func (b *backend) pathGroupCredsSpecificUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("pathGroupCredsSpecificUpdate entered\n")
	return nil, nil
}

var groupHelp = map[string][2]string{
	"group":                     {"help", "desc"},
	"group-apps":                {"help", "desc"},
	"group-additional-policies": {"help", "desc"},
	"group-num-uses":            {"help", "desc"},
	"group-ttl":                 {"help", "desc"},
	"group-max-ttl":             {"help", "desc"},
	"group-wrgrouped":           {"help", "desc"},
	"group-creds":               {"help", "desc"},
	"group-creds-specific":      {"help", "desc"},
}
