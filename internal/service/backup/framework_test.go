package backup_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfbackup "github.com/hashicorp/terraform-provider-aws/internal/service/backup"
)

func TestAccBackupFramework_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":                        testAccFramework_basic,
			"disappears":                   testAccFramework_disappears,
			"UpdateTags":                   testAccFramework_updateTags,
			"UpdateControlScope":           testAccFramework_updateControlScope,
			"UpdateControlInputParameters": testAccFramework_updateControlInputParameters,
			"UpdateControls":               testAccFramework_updateControls,
		},
		"DataSource": {
			"basic":           testAccFrameworkDataSource_basic,
			"ControlScopeTag": testAccFrameworkDataSource_controlScopeTag,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}

func testAccFramework_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	originalDescription := "original description"
	updatedDescription := "updated description"
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_basic(rName, originalDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", originalDescription),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_basic(rName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
		},
	})
}

func testAccFramework_updateTags(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	description := "example description"
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_tags(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2a"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_tagsUpdated(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2b"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "Value3"),
				),
			},
		},
	})
}

func testAccFramework_updateControlScope(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	description := "example description"
	originalControlScopeTagValue := "example"
	updatedControlScopeTagValue := ""
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_controlScopeComplianceResourceID(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "control.0.scope.0.compliance_resource_ids.0", "aws_ebs_volume.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_controlScopeTag(rName, description, originalControlScopeTagValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.tags.Name", originalControlScopeTagValue),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_controlScopeTag(rName, description, updatedControlScopeTagValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.tags.Name", updatedControlScopeTagValue),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
		},
	})
}

func testAccFramework_updateControlInputParameters(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	description := "example description"
	originalRequiredRetentionDays := "35"
	updatedRequiredRetentionDays := "34"
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_controlInputParameter(rName, description, originalRequiredRetentionDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_PLAN_MIN_FREQUENCY_AND_MIN_RETENTION_CHECK"),
					resource.TestCheckResourceAttr(resourceName, "control.0.input_parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredFrequencyUnit",
						"value": "hours",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredRetentionDays",
						"value": originalRequiredRetentionDays,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredFrequencyValue",
						"value": "1",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_controlInputParameter(rName, description, updatedRequiredRetentionDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_PLAN_MIN_FREQUENCY_AND_MIN_RETENTION_CHECK"),
					resource.TestCheckResourceAttr(resourceName, "control.0.input_parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredFrequencyUnit",
						"value": "hours",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredRetentionDays",
						"value": updatedRequiredRetentionDays,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.0.input_parameter.*", map[string]string{
						"name":  "requiredFrequencyValue",
						"value": "1",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
		},
	})
}

func testAccFramework_updateControls(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	description := "example description"
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.name", "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "control.0.scope.0.compliance_resource_types.0", "EBS"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFrameworkConfig_controls(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "control.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.*", map[string]string{
						"name":                    "BACKUP_RECOVERY_POINT_MINIMUM_RETENTION_CHECK",
						"input_parameter.#":       "1",
						"input_parameter.0.name":  "requiredRetentionDays",
						"input_parameter.0.value": "35",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.*", map[string]string{
						"name":              "BACKUP_PLAN_MIN_FREQUENCY_AND_MIN_RETENTION_CHECK",
						"input_parameter.#": "3",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.*", map[string]string{
						"name": "BACKUP_RECOVERY_POINT_ENCRYPTED",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.*", map[string]string{
						"name":                                "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN",
						"scope.#":                             "1",
						"scope.0.compliance_resource_types.#": "1",
						"scope.0.compliance_resource_types.0": "EBS",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "control.*", map[string]string{
						"name": "BACKUP_RECOVERY_POINT_MANUAL_DELETION_DISABLED",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_status"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Test Framework"),
				),
			},
		},
	})
}

func testAccFramework_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var framework backup.DescribeFrameworkOutput

	rName := fmt.Sprintf("tf_acc_test_%s", sdkacctest.RandString(7))
	description := "disappears"
	resourceName := "aws_backup_framework.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccFrameworkPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFrameworkDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFrameworkConfig_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFrameworkExists(ctx, resourceName, &framework),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfbackup.ResourceFramework(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFrameworkPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).BackupConn(ctx)

	_, err := conn.ListFrameworksWithContext(ctx, &backup.ListFrameworksInput{})

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckFrameworkDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).BackupConn(ctx)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_backup_framework" {
				continue
			}

			input := &backup.DescribeFrameworkInput{
				FrameworkName: aws.String(rs.Primary.ID),
			}

			resp, err := conn.DescribeFrameworkWithContext(ctx, input)

			if err == nil {
				if aws.StringValue(resp.FrameworkName) == rs.Primary.ID {
					return fmt.Errorf("Backup Framework '%s' was not deleted properly", rs.Primary.ID)
				}
			}
		}

		return nil
	}
}

func testAccCheckFrameworkExists(ctx context.Context, name string, framework *backup.DescribeFrameworkOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).BackupConn(ctx)
		input := &backup.DescribeFrameworkInput{
			FrameworkName: aws.String(rs.Primary.ID),
		}
		resp, err := conn.DescribeFrameworkWithContext(ctx, input)

		if err != nil {
			return err
		}

		*framework = *resp

		return nil
	}
}

func testAccFrameworkConfig_basic(rName, label string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      compliance_resource_types = [
        "EBS"
      ]
    }
  }

  tags = {
    "Name" = "Test Framework"
  }
}
`, rName, label)
}

func testAccFrameworkConfig_tags(rName, label string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      compliance_resource_types = [
        "EBS"
      ]
    }
  }

  tags = {
    "Name" = "Test Framework"
    "Key2" = "Value2a"
  }
}
`, rName, label)
}

func testAccFrameworkConfig_tagsUpdated(rName, label string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      compliance_resource_types = [
        "EBS"
      ]
    }
  }

  tags = {
    "Name" = "Test Framework"
    "Key2" = "Value2b"
    "Key3" = "Value3"
  }
}
`, rName, label)
}

func testAccFrameworkConfig_controlScopeComplianceResourceID(rName, label string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_ebs_volume" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  type              = "gp2"
  size              = 1
}

resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      compliance_resource_ids = [
        aws_ebs_volume.test.id
      ]

      compliance_resource_types = [
        "EBS"
      ]
    }
  }

  tags = {
    "Name" = "Test Framework"
  }
}
`, rName, label)
}

func testAccFrameworkConfig_controlScopeTag(rName, label, controlScopeTagValue string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      tags = {
        "Name" = %[3]q
      }
    }
  }

  tags = {
    "Name" = "Test Framework"
  }
}
`, rName, label, controlScopeTagValue)
}

func testAccFrameworkConfig_controlInputParameter(rName, label, requiredRetentionDaysValue string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_PLAN_MIN_FREQUENCY_AND_MIN_RETENTION_CHECK"

    input_parameter {
      name  = "requiredFrequencyUnit"
      value = "hours"
    }

    input_parameter {
      name  = "requiredRetentionDays"
      value = %[3]q
    }

    input_parameter {
      name  = "requiredFrequencyValue"
      value = "1"
    }
  }

  tags = {
    "Name" = "Test Framework"
  }
}
`, rName, label, requiredRetentionDaysValue)
}

func testAccFrameworkConfig_controls(rName, label string) string {
	return fmt.Sprintf(`
resource "aws_backup_framework" "test" {
  name        = %[1]q
  description = %[2]q

  control {
    name = "BACKUP_RECOVERY_POINT_MINIMUM_RETENTION_CHECK"

    input_parameter {
      name  = "requiredRetentionDays"
      value = "35"
    }
  }

  control {
    name = "BACKUP_PLAN_MIN_FREQUENCY_AND_MIN_RETENTION_CHECK"

    input_parameter {
      name  = "requiredFrequencyUnit"
      value = "hours"
    }

    input_parameter {
      name  = "requiredRetentionDays"
      value = "35"
    }

    input_parameter {
      name  = "requiredFrequencyValue"
      value = "1"
    }
  }

  control {
    name = "BACKUP_RECOVERY_POINT_ENCRYPTED"
  }

  control {
    name = "BACKUP_RESOURCES_PROTECTED_BY_BACKUP_PLAN"

    scope {
      compliance_resource_types = [
        "EBS"
      ]
    }
  }

  control {
    name = "BACKUP_RECOVERY_POINT_MANUAL_DELETION_DISABLED"
  }

  tags = {
    "Name" = "Test Framework"
  }
}
`, rName, label)
}
