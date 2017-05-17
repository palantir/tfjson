package google

/**
 * Note! You must run these tests once at a time. Google Cloud SQL does
 * not allow you to reuse a database for a short time after you reserved it,
 * and for this reason the tests will fail if the same config is used serveral
 * times in short succession.
 */

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/sqladmin/v1beta4"
)

func TestAccGoogleSqlDatabaseInstance_basic(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_basic2(t *testing.T) {
	var instance sqladmin.DatabaseInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleSqlDatabaseInstance_basic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_basic3(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic3, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(
						&instance),
				),
			},
		},
	})
}
func TestAccGoogleSqlDatabaseInstance_settings_basic(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_slave(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()
	slaveID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_slave, masterID, slaveID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance_master", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance_master", &instance),
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance_slave", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance_slave", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_diskspecs(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_diskspecs, masterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_maintenance(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	masterID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_maintenance, masterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_settings_upgrade(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabaseInstance_settings_downgrade(t *testing.T) {
	var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_settings, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_basic, databaseID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseInstanceExists(
						"google_sql_database_instance.instance", &instance),
					testAccCheckGoogleSqlDatabaseInstanceEquals(
						"google_sql_database_instance.instance", &instance),
				),
			},
		},
	})
}

// GH-4222
func TestAccGoogleSqlDatabaseInstance_authNets(t *testing.T) {
	// var instance sqladmin.DatabaseInstance
	databaseID := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step2, databaseID),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabaseInstance_authNets_step1, databaseID),
			},
		},
	})
}

func testAccCheckGoogleSqlDatabaseInstanceEquals(n string,
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		attributes := rs.Primary.Attributes

		server := instance.Name
		local := attributes["name"]
		if server != local {
			return fmt.Errorf("Error name mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.Tier
		local = attributes["settings.0.tier"]
		if server != local {
			return fmt.Errorf("Error settings.tier mismatch, (%s, %s)", server, local)
		}

		server = strings.TrimPrefix(instance.MasterInstanceName, instance.Project+":")
		local = attributes["master_instance_name"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error master_instance_name mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.ActivationPolicy
		local = attributes["settings.0.activation_policy"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.activation_policy mismatch, (%s, %s)", server, local)
		}

		if instance.Settings.BackupConfiguration != nil {
			server = strconv.FormatBool(instance.Settings.BackupConfiguration.BinaryLogEnabled)
			local = attributes["settings.0.backup_configuration.0.binary_log_enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.binary_log_enabled mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatBool(instance.Settings.BackupConfiguration.Enabled)
			local = attributes["settings.0.backup_configuration.0.enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.enabled mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.BackupConfiguration.StartTime
			local = attributes["settings.0.backup_configuration.0.start_time"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.backup_configuration.start_time mismatch, (%s, %s)", server, local)
			}
		}

		server = strconv.FormatBool(instance.Settings.CrashSafeReplicationEnabled)
		local = attributes["settings.0.crash_safe_replication"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.crash_safe_replication mismatch, (%s, %s)", server, local)
		}

		server = strconv.FormatBool(instance.Settings.StorageAutoResize)
		local = attributes["settings.0.disk_autoresize"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.disk_autoresize mismatch, (%s, %s)", server, local)
		}

		server = strconv.FormatInt(instance.Settings.DataDiskSizeGb, 10)
		local = attributes["settings.0.disk_size"]
		if server != local && len(server) > 0 && len(local) > 0 && local != "0" {
			return fmt.Errorf("Error settings.disk_size mismatch, (%s, %s)", server, local)
		}

		server = instance.Settings.DataDiskType
		local = attributes["settings.0.disk_type"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.disk_type mismatch, (%s, %s)", server, local)
		}

		if instance.Settings.IpConfiguration != nil {
			server = strconv.FormatBool(instance.Settings.IpConfiguration.Ipv4Enabled)
			local = attributes["settings.0.ip_configuration.0.ipv4_enabled"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.ip_configuration.ipv4_enabled mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatBool(instance.Settings.IpConfiguration.RequireSsl)
			local = attributes["settings.0.ip_configuration.0.require_ssl"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.ip_configuration.require_ssl mismatch, (%s, %s)", server, local)
			}
		}

		if instance.Settings.LocationPreference != nil {
			server = instance.Settings.LocationPreference.FollowGaeApplication
			local = attributes["settings.0.location_preference.0.follow_gae_application"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.location_preference.follow_gae_application mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.LocationPreference.Zone
			local = attributes["settings.0.location_preference.0.zone"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.location_preference.zone mismatch, (%s, %s)", server, local)
			}
		}

		if instance.Settings.MaintenanceWindow != nil {
			server = strconv.FormatInt(instance.Settings.MaintenanceWindow.Day, 10)
			local = attributes["settings.0.maintenance_window.0.day"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.day mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatInt(instance.Settings.MaintenanceWindow.Hour, 10)
			local = attributes["settings.0.maintenance_window.0.hour"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.hour mismatch, (%s, %s)", server, local)
			}

			server = instance.Settings.MaintenanceWindow.UpdateTrack
			local = attributes["settings.0.maintenance_window.0.update_track"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error settings.maintenance_window.update_track mismatch, (%s, %s)", server, local)
			}
		}

		server = instance.Settings.PricingPlan
		local = attributes["settings.0.pricing_plan"]
		if server != local && len(server) > 0 && len(local) > 0 {
			return fmt.Errorf("Error settings.pricing_plan mismatch, (%s, %s)", server, local)
		}

		if instance.ReplicaConfiguration != nil &&
			instance.ReplicaConfiguration.MysqlReplicaConfiguration != nil {
			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.CaCertificate
			local = attributes["replica_configuration.0.ca_certificate"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.ca_certificate mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.ClientCertificate
			local = attributes["replica_configuration.0.client_certificate"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.client_certificate mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.ClientKey
			local = attributes["replica_configuration.0.client_key"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.client_key mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatInt(instance.ReplicaConfiguration.MysqlReplicaConfiguration.ConnectRetryInterval, 10)
			local = attributes["replica_configuration.0.connect_retry_interval"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.connect_retry_interval mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.DumpFilePath
			local = attributes["replica_configuration.0.dump_file_path"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.dump_file_path mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatInt(instance.ReplicaConfiguration.MysqlReplicaConfiguration.MasterHeartbeatPeriod, 10)
			local = attributes["replica_configuration.0.master_heartbeat_period"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.master_heartbeat_period mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.Password
			local = attributes["replica_configuration.0.password"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.password mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.SslCipher
			local = attributes["replica_configuration.0.ssl_cipher"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.ssl_cipher mismatch, (%s, %s)", server, local)
			}

			server = instance.ReplicaConfiguration.MysqlReplicaConfiguration.Username
			local = attributes["replica_configuration.0.username"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.username mismatch, (%s, %s)", server, local)
			}

			server = strconv.FormatBool(instance.ReplicaConfiguration.MysqlReplicaConfiguration.VerifyServerCertificate)
			local = attributes["replica_configuration.0.verify_server_certificate"]
			if server != local && len(server) > 0 && len(local) > 0 {
				return fmt.Errorf("Error replica_configuration.verify_server_certificate mismatch, (%s, %s)", server, local)
			}
		}

		return nil
	}
}

func testAccCheckGoogleSqlDatabaseInstanceExists(n string,
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		found, err := config.clientSqlAdmin.Instances.Get(config.Project,
			rs.Primary.Attributes["name"]).Do()

		*instance = *found

		if err != nil {
			return fmt.Errorf("Not found: %s", n)
		}

		return nil
	}
}

func testAccGoogleSqlDatabaseInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*Config)
		if rs.Type != "google_sql_database_instance" {
			continue
		}

		_, err := config.clientSqlAdmin.Instances.Get(config.Project,
			rs.Primary.Attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Database Instance still exists")
		}
	}

	return nil
}

func testAccCheckGoogleSqlDatabaseRootUserDoesNotExist(
	instance *sqladmin.DatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		users, err := config.clientSqlAdmin.Users.List(config.Project, instance.Name).Do()

		if err != nil {
			return fmt.Errorf("Could not list database users for %q: %s", instance.Name, err)
		}

		for _, u := range users.Items {
			if u.Name == "root" && u.Host == "%" {
				return fmt.Errorf("%v@%v user still exists", u.Name, u.Host)
			}
		}

		return nil
	}
}

var testGoogleSqlDatabaseInstance_basic = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
	}
}
`

var testGoogleSqlDatabaseInstance_basic2 = `
resource "google_sql_database_instance" "instance" {
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
	}
}
`
var testGoogleSqlDatabaseInstance_basic3 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "db-f1-micro"
	}
}
`

var testGoogleSqlDatabaseInstance_settings = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false
		replication_type = "ASYNCHRONOUS"
		location_preference {
			zone = "us-central1-f"
		}

		ip_configuration {
			ipv4_enabled = "true"
			authorized_networks {
				value = "108.12.12.12"
				name = "misc"
				expiration_time = "2017-11-15T16:19:00.094Z"
			}
		}

		backup_configuration {
			enabled = "true"
			start_time = "19:19"
		}

		activation_policy = "ON_DEMAND"
	}
}
`

// Note - this test is not feasible to run unless we generate
// backups first.
var testGoogleSqlDatabaseInstance_replica = `
resource "google_sql_database_instance" "instance_master" {
	name = "tf-lw-%d"
	database_version = "MYSQL_5_6"
	region = "us-east1"

	settings {
		tier = "D0"
		crash_safe_replication = true

		backup_configuration {
			enabled = true
			start_time = "00:00"
			binary_log_enabled = true
		}
	}
}

resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	database_version = "MYSQL_5_6"
	region = "us-central"

	settings {
		tier = "D0"
	}

	master_instance_name = "${google_sql_database_instance.instance_master.name}"

	replica_configuration {
		ca_certificate = "${file("~/tmp/fake.pem")}"
		client_certificate = "${file("~/tmp/fake.pem")}"
		client_key = "${file("~/tmp/fake.pem")}"
		connect_retry_interval = 100
		master_heartbeat_period = 10000
		password = "password"
		username = "username"
		ssl_cipher = "ALL"
		verify_server_certificate = false
	}
}
`

var testGoogleSqlDatabaseInstance_slave = `
resource "google_sql_database_instance" "instance_master" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"

		backup_configuration {
			enabled = true
			binary_log_enabled = true
		}
	}
}

resource "google_sql_database_instance" "instance_slave" {
	name = "tf-lw-%d"
	region = "us-central1"

	master_instance_name = "${google_sql_database_instance.instance_master.name}"

	settings {
		tier = "db-f1-micro"
	}
}
`

var testGoogleSqlDatabaseInstance_diskspecs = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"
		disk_autoresize = true
		disk_size = 15
		disk_type = "PD_HDD"
	}
}
`

var testGoogleSqlDatabaseInstance_maintenance = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central1"

	settings {
		tier = "db-f1-micro"

		maintenance_window {
		  day  = 7
		  hour = 3
			update_track = "canary"
	  }
	}
}
`

var testGoogleSqlDatabaseInstance_authNets_step1 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false

		ip_configuration {
			ipv4_enabled = "true"
			authorized_networks {
				value = "108.12.12.12"
				name = "misc"
				expiration_time = "2017-11-15T16:19:00.094Z"
			}
		}
	}
}
`

var testGoogleSqlDatabaseInstance_authNets_step2 = `
resource "google_sql_database_instance" "instance" {
	name = "tf-lw-%d"
	region = "us-central"
	settings {
		tier = "D0"
		crash_safe_replication = false

		ip_configuration {
			ipv4_enabled = "true"
		}
	}
}
`
