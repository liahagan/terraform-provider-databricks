package secrets

import (
	"context"

	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecretAcl struct {
	// The name of the scope to apply permissions to.
	Scope string `json:"scope" tf:"force_new"`
	// The principal in which the permission is applied.
	Principal string `json:"principal" tf:"force_new"`
	// The permission level applied to the principal.
	Permission ACLPermission `json:"permission" tf:"force_new"`
}

func ResourceSecretACL() *schema.Resource {
	s := common.StructToSchema(SecretAcl{},
		func(m map[string]*schema.Schema) map[string]*schema.Schema {
			m["scope"].ValidateFunc = validScope
			return m
		})
	p := common.NewPairSeparatedID("scope", "principal", "|||")
	return common.Resource{
		Schema: s,
		Create: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			w, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			var putAcl workspace.PutAcl
			common.DataToStructPointer(d, s, &putAcl)
			if err := w.Secrets.PutAcl(ctx, putAcl); err != nil {
				return err
			}
			p.Pack(d)
			return nil
		},
		Read: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			w, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			var getAclRequest workspace.GetAclRequest
			common.DataToStructPointer(d, s, &getAclRequest)
			v, err := w.Secrets.GetAcl(ctx, getAclRequest)
			if err != nil {
				return err
			}
			return common.StructToData(v, s, d)
		},
		Delete: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			w, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			var deleteAcl workspace.DeleteAcl
			common.DataToStructPointer(d, s, &deleteAcl)
			return w.Secrets.DeleteAcl(ctx, deleteAcl)
		},
	}.ToResource()
}
