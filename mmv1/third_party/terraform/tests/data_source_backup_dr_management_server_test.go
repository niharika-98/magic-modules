package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBackupDRManagementServer_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBackupDRManagementServer_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_cloud_backup_dr_management_server.foo", "google_cloud_backup_dr_management_server.foo"),
				),
			},
		},
	})
}


func testAccDataSourceGoogleBackupDRManagementServer_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_compute_network" "default" {
		provider = google-beta
		name = "vpc-network"
	  }
	  
	  resource "google_compute_global_address" "private_ip_address" {
		provider = google-beta
		name          = "vpc-network"
		address_type  = "INTERNAL"
		purpose       = "VPC_PEERING"
		prefix_length = 20
		network       = google_compute_network.default.id
	  }
	  
	  resource "google_service_networking_connection" "default" {
		provider = google-beta
		network                 = google_compute_network.default.id
		service                 = "servicenetworking.googleapis.com"
		reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
	  }
	  
	  resource "google_backup_dr_management_server" "foo" {
		provider = google-beta
		location = "us-central1"
		name     = "management_server"
		type     = "BACKUP_RESTORE" 
		networks {
		  network      = google_compute_network.default.id
		  peering_mode = "PRIVATE_SERVICE_ACCESS"
		}
		depends_on = [ google_service_networking_connection.default ]
	  }

data "ggoogle_backup_dr_management_server" "foo" {
  location = google_backup_dr_management_server.foo.location
}
`, context)
}
