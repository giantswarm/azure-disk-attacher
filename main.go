package main

import (
	"context"
	"flag"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	ctx := context.Background()

	subscriptionId := flag.String("subscription", "", "Subscription ID")
	resourceGroup := flag.String("resource-group", "", "Resource Group")
	vmssName := flag.String("vmss", "", "VMSS name")
	instanceID := flag.String("instance-id", "", "Instance ID")
	diskID := flag.String("disk-id", "", "Disk ID")
	lun := flag.Int("lun", -1, "LUN to attach the disk to")

	flag.Parse()

	if *subscriptionId == "" {
		log.Fatalf("--subscription can't be empty")
	}
	if *resourceGroup == "" {
		log.Fatalf("--resource-group can't be empty")
	}
	if *vmssName == "" {
		log.Fatalf("--vmss can't be empty")
	}
	if *instanceID == "" {
		log.Fatalf("--instance-id can't be empty")
	}
	if *diskID == "" {
		log.Fatalf("--disk-id can't be empty")
	}
	if *lun < 0 {
		log.Fatalf("--lun can't be empty")
	}

	credentials := auth.NewMSIConfig()

	authorizer, err := credentials.Authorizer()
	if err != nil {
		log.Fatalf("Error loggin in using MSI: %s", err)
	}

	c := compute.NewVirtualMachineScaleSetVMsClient(*subscriptionId)
	c.Authorizer = authorizer

	log.Printf("Getting VMSS VM with ID %s from VMSS %s in group %s", *instanceID, *vmssName, *resourceGroup)

	vm, err := c.Get(ctx, *resourceGroup, *vmssName, *instanceID, "")
	if err != nil {
		log.Fatalf("Error reading VMSS VM: %s", err)
	}

	dd := *vm.StorageProfile.DataDisks
	dd = append(dd, compute.DataDisk{
		Lun:          to.Int32Ptr(int32(*lun)),
		CreateOption: compute.DiskCreateOptionTypesAttach,
		ManagedDisk: &compute.ManagedDiskParameters{
			ID: to.StringPtr(*diskID),
		},
	})

	log.Printf("Attaching disk %s at LUN %d to VMSS VM with ID %s from VMSS %s in group %s", *diskID, *lun, *instanceID, *vmssName, *resourceGroup)

	future, err := c.Update(ctx, *resourceGroup, *vmssName, *instanceID, compute.VirtualMachineScaleSetVM{
		VirtualMachineScaleSetVMProperties: &compute.VirtualMachineScaleSetVMProperties{
			StorageProfile: &compute.StorageProfile{
				DataDisks: &dd,
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating VMSS VM: %s", err)
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		log.Fatalf("Error waiting for VMSS Update call to complete: %s", err)
	}
}
