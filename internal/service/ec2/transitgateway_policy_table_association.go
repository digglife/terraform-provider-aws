package ec2

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceTransitGatewayPolicyTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTransitGatewayPolicyTableAssociationCreate,
		Read:   ResourceTransitGatewayPolicyTableAssociationRead,
		Delete: ResourceTransitGatewayPolicyTableAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_gateway_attachment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"transit_gateway_policy_table_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func ResourceTransitGatewayPolicyTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayAttachmentID := d.Get("transit_gateway_attachment_id").(string)
	transitGatewayPolicyTableID := d.Get("transit_gateway_policy_table_id").(string)
	id := TransitGatewayPolicyTableAssociationCreateResourceID(transitGatewayPolicyTableID, transitGatewayAttachmentID)
	input := &ec2.AssociateTransitGatewayPolicyTableInput{
		TransitGatewayAttachmentId:  aws.String(transitGatewayAttachmentID),
		TransitGatewayPolicyTableId: aws.String(transitGatewayPolicyTableID),
	}

	_, err := conn.AssociateTransitGatewayPolicyTable(input)

	if err != nil {
		return fmt.Errorf("creating EC2 Transit Gateway Policy Table Association (%s): %w", id, err)
	}

	d.SetId(id)

	if _, err := WaitTransitGatewayPolicyTableAssociationCreated(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID); err != nil {
		return fmt.Errorf("waiting for EC2 Transit Gateway Policy Table Association (%s) create: %w", d.Id(), err)
	}

	return ResourceTransitGatewayPolicyTableAssociationRead(d, meta)
}

func ResourceTransitGatewayPolicyTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayPolicyTableID, transitGatewayAttachmentID, err := TransitGatewayPolicyTableAssociationParseResourceID(d.Id())

	if err != nil {
		return err
	}

	transitGatewayPolicyTableAssociation, err := FindTransitGatewayPolicyTableAssociationByTwoPartKey(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Transit Gateway Policy Table Association %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 Transit Gateway Policy Table Association (%s): %w", d.Id(), err)
	}

	d.Set("resource_id", transitGatewayPolicyTableAssociation.ResourceId)
	d.Set("resource_type", transitGatewayPolicyTableAssociation.ResourceType)
	d.Set("transit_gateway_attachment_id", transitGatewayPolicyTableAssociation.TransitGatewayAttachmentId)
	d.Set("state", transitGatewayPolicyTableAssociation.State)
	d.Set("transit_gateway_policy_table_id", transitGatewayPolicyTableAssociation)

	return nil
}

func ResourceTransitGatewayPolicyTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	transitGatewayPolicyTableID, transitGatewayAttachmentID, err := TransitGatewayPolicyTableAssociationParseResourceID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting EC2 Transit Gateway Policy Table Association: %s", d.Id())
	_, err = conn.DisassociateTransitGatewayPolicyTable(&ec2.DisassociateTransitGatewayPolicyTableInput{
		TransitGatewayAttachmentId:  aws.String(transitGatewayAttachmentID),
		TransitGatewayPolicyTableId: aws.String(transitGatewayPolicyTableID),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidPolicyTableIDNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting EC2 Transit Gateway Policy Table Association (%s): %w", d.Id(), err)
	}

	if _, err := WaitTransitGatewayPolicyTableAssociationDeleted(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID); err != nil {
		return fmt.Errorf("waiting for EC2 Transit Gateway Policy Table Association (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func transitGatewayPolicyTableAssociationUpdate(conn *ec2.EC2, transitGatewayPolicyTableID, transitGatewayAttachmentID string, associate bool) error {
	id := TransitGatewayPolicyTableAssociationCreateResourceID(transitGatewayPolicyTableID, transitGatewayAttachmentID)
	_, err := FindTransitGatewayPolicyTableAssociationByTwoPartKey(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID)

	if tfresource.NotFound(err) {
		if associate {
			input := &ec2.AssociateTransitGatewayPolicyTableInput{
				TransitGatewayAttachmentId:  aws.String(transitGatewayAttachmentID),
				TransitGatewayPolicyTableId: aws.String(transitGatewayPolicyTableID),
			}

			_, err := conn.AssociateTransitGatewayPolicyTable(input)

			if err != nil {
				return fmt.Errorf("creating EC2 Transit Gateway Policy Table Association (%s): %w", id, err)
			}

			if _, err := WaitTransitGatewayPolicyTableAssociationCreated(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID); err != nil {
				return fmt.Errorf("waiting for EC2 Transit Gateway Policy Table Association (%s) create: %w", id, err)
			}
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("reading EC2 Transit Gateway Policy Table Association (%s): %w", id, err)
	}

	if !associate {
		// Disassociation must be done only on already associated state.
		if _, err := WaitTransitGatewayPolicyTableAssociationCreated(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID); err != nil {
			return fmt.Errorf("waiting for EC2 Transit Gateway Policy Table Association (%s) create: %w", id, err)
		}

		input := &ec2.DisassociateTransitGatewayPolicyTableInput{
			TransitGatewayAttachmentId:  aws.String(transitGatewayAttachmentID),
			TransitGatewayPolicyTableId: aws.String(transitGatewayPolicyTableID),
		}

		if _, err := conn.DisassociateTransitGatewayPolicyTable(input); err != nil {
			return fmt.Errorf("deleting EC2 Transit Gateway Policy Table Association (%s): %w", id, err)
		}

		if _, err := WaitTransitGatewayPolicyTableAssociationDeleted(conn, transitGatewayPolicyTableID, transitGatewayAttachmentID); err != nil {
			return fmt.Errorf("waiting for EC2 Transit Gateway Policy Table Association (%s) delete: %w", id, err)
		}
	}

	return nil
}

const transitGatewayPolicyTableAssociationIDSeparator = "_"

func TransitGatewayPolicyTableAssociationCreateResourceID(transitGatewayPolicyTableID, transitGatewayAttachmentID string) string {
	parts := []string{transitGatewayPolicyTableID, transitGatewayAttachmentID}
	id := strings.Join(parts, transitGatewayPolicyTableAssociationIDSeparator)

	return id
}

func TransitGatewayPolicyTableAssociationParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, transitGatewayPolicyTableAssociationIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected TRANSIT-GATEWAY-POLICY-TABLE-ID%[2]sTRANSIT-GATEWAY-ATTACHMENT-ID", id, transitGatewayPolicyTableAssociationIDSeparator)
}
