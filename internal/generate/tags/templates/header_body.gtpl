// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package {{ .ServicePackage }}

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/YakDriver/smarterr"
	"github.com/aws/aws-sdk-go-v2/aws"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/aws/aws-sdk-go-v2/service/{{ .AWSService }}"
	awstypes "github.com/aws/aws-sdk-go-v2/service/{{ .AWSService }}/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
    "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
    tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
    "github.com/hashicorp/terraform-provider-aws/internal/tfresource"
    "github.com/hashicorp/terraform-provider-aws/internal/types/option"
    "github.com/hashicorp/terraform-provider-aws/names"
)
