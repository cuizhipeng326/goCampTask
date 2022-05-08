package model

// Kratos hello kratos.
const (
	ResourceTypeSensitive = "sensitive"
)

type Kratos struct {
	Hello string
}

type Article struct {
	ID      int64
	Content string
	Author  string
}

type Permission struct {
	PermissionID       string `json:"permissionId" `
	ParentResourceType string `json:"parentResourceType" `
	ResourceID         string `json:"resourceId" `
	Name               string `json:"name" `
	Action             string `json:"action" `
	ParentID           string `json:"parentId" `
	ResourceType       string `json:"resourceType" `
	ParentResourceID   string `json:"parentResourceId" `
	OperationOffset    int    `json:"operationOffset" `
	MemberOfCreator    int    `json:"memberOfCreator"`
}
type Role struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	RoleType string `json:"roleType"`
	Creator  string `json:"creator"`
	Private  int    `json:"private"`
}

type PermissionGroup struct {
	ResourceId      string `json:"resourceId"`
	PermissionId    string `json:"permissionId"`
	ResourceType    string `json:"resourceType"`
	ParentId        string `json:"parentId"`
	Name            string `json:"name"`
	MemberOfCreator int    `json:"memberOfCreator"`
}

type PermissionList struct {
	Permissions      []Permission      `json:"permissions"`
	Roles            []Role            `json:"roles"`
	PermissionGroups []PermissionGroup `json:"permissionGroups"`
}
