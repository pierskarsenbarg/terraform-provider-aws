package rum_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchrum"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudwatchrum "github.com/hashicorp/terraform-provider-aws/internal/service/rum"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRUMMetricsDestination_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var dest cloudwatchrum.MetricDestinationSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rum_metrics_destination.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchrum.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricsDestinationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsDestinationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsDestinationExists(ctx, resourceName, &dest),
					resource.TestCheckResourceAttrPair(resourceName, "app_monitor_name", "aws_rum_app_monitor.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "destination", "CloudWatch"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRUMMetricsDestination_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var dest cloudwatchrum.MetricDestinationSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rum_metrics_destination.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchrum.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricsDestinationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsDestinationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsDestinationExists(ctx, resourceName, &dest),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudwatchrum.ResourceMetricsDestination(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudwatchrum.ResourceMetricsDestination(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRUMMetricsDestination_disappears_appMonitor(t *testing.T) {
	ctx := acctest.Context(t)
	var dest cloudwatchrum.MetricDestinationSummary
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_rum_metrics_destination.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchrum.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMetricsDestinationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsDestinationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsDestinationExists(ctx, resourceName, &dest),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudwatchrum.ResourceAppMonitor(), "aws_rum_app_monitor.test"),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudwatchrum.ResourceMetricsDestination(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckMetricsDestinationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RUMConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_rum_metrics_destination" {
				continue
			}

			_, err := tfcloudwatchrum.FindMetricsDestinationByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("CloudWatch RUM Metrics Destination %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckMetricsDestinationExists(ctx context.Context, n string, v *cloudwatchrum.MetricDestinationSummary) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudWatch RUM Metrics Destination ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RUMConn(ctx)

		output, err := tfcloudwatchrum.FindMetricsDestinationByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccMetricsDestinationConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_rum_app_monitor" "test" {
  name   = %[1]q
  domain = "localhost"
}

resource "aws_rum_metrics_destination" "test" {
  app_monitor_name = aws_rum_app_monitor.test.name
  destination      = "CloudWatch"
}
`, rName)
}
