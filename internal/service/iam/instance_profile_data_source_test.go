package iam_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccIAMInstanceProfileDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_iam_instance_profile.test"
	resourceName := "aws_iam_instance_profile.test"
	roleResourceName := "aws_iam_role.test"
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iam.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProfileDataSourceConfig_basic(rName1, rName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "/testpath/"),
					resource.TestCheckResourceAttrPair(dataSourceName, "role_arn", roleResourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "role_id", roleResourceName, "unique_id"),
					resource.TestCheckResourceAttr(dataSourceName, "role_name", rName1),
				),
			},
		},
	})
}

func testAccInstanceProfileDataSourceConfig_basic(roleName, profileName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = %[1]q
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
  name = %[2]q
  role = aws_iam_role.test.name
  path = "/testpath/"
}

data "aws_iam_instance_profile" "test" {
  name = aws_iam_instance_profile.test.name
}
`, roleName, profileName)
}
