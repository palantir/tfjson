package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func resourceComputeDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeDiskCreate,
		Read:   resourceComputeDiskRead,
		Delete: resourceComputeDiskDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"disk_encryption_key_raw": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"disk_encryption_key_sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"image": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"snapshot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	log.Printf("[DEBUG] Loading zone: %s", d.Get("zone").(string))
	zone, err := config.clientCompute.Zones.Get(
		project, d.Get("zone").(string)).Do()
	if err != nil {
		return fmt.Errorf(
			"Error loading zone '%s': %s", d.Get("zone").(string), err)
	}

	// Build the disk parameter
	disk := &compute.Disk{
		Name:   d.Get("name").(string),
		SizeGb: int64(d.Get("size").(int)),
	}

	// If we were given a source image, load that.
	if v, ok := d.GetOk("image"); ok {
		log.Printf("[DEBUG] Resolving image name: %s", v.(string))
		imageUrl, err := resolveImage(config, v.(string))
		if err != nil {
			return fmt.Errorf(
				"Error resolving image name '%s': %s",
				v.(string), err)
		}

		disk.SourceImage = imageUrl
		log.Printf("[DEBUG] Image name resolved to: %s", imageUrl)
	}

	if v, ok := d.GetOk("type"); ok {
		log.Printf("[DEBUG] Loading disk type: %s", v.(string))
		diskType, err := readDiskType(config, zone, v.(string))
		if err != nil {
			return fmt.Errorf(
				"Error loading disk type '%s': %s",
				v.(string), err)
		}

		disk.Type = diskType.SelfLink
	}

	if v, ok := d.GetOk("snapshot"); ok {
		snapshotName := v.(string)
		log.Printf("[DEBUG] Loading snapshot: %s", snapshotName)
		snapshotData, err := config.clientCompute.Snapshots.Get(
			project, snapshotName).Do()

		if err != nil {
			return fmt.Errorf(
				"Error loading snapshot '%s': %s",
				snapshotName, err)
		}

		disk.SourceSnapshot = snapshotData.SelfLink
	}

	if v, ok := d.GetOk("disk_encryption_key_raw"); ok {
		disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{}
		disk.DiskEncryptionKey.RawKey = v.(string)
	}

	op, err := config.clientCompute.Disks.Insert(
		project, d.Get("zone").(string), disk).Do()
	if err != nil {
		return fmt.Errorf("Error creating disk: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(disk.Name)

	err = computeOperationWaitZone(config, op, project, d.Get("zone").(string), "Creating Disk")
	if err != nil {
		return err
	}
	return resourceComputeDiskRead(d, meta)
}

func resourceComputeDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	disk, err := config.clientCompute.Disks.Get(
		project, d.Get("zone").(string), d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Disk %q", d.Get("name").(string)))
	}

	d.Set("self_link", disk.SelfLink)
	if disk.DiskEncryptionKey != nil && disk.DiskEncryptionKey.Sha256 != "" {
		d.Set("disk_encryption_key_sha256", disk.DiskEncryptionKey.Sha256)
	}

	return nil
}

func resourceComputeDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the disk
	op, err := config.clientCompute.Disks.Delete(
		project, d.Get("zone").(string), d.Id()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing Disk %q because it's gone", d.Get("name").(string))
			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting disk: %s", err)
	}

	zone := d.Get("zone").(string)
	err = computeOperationWaitZone(config, op, project, zone, "Deleting Disk")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
