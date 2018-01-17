
/*
Package Eip enables management and retrieval of Eip from the Open Telekom Cloud
EIP service.
Example to List Eips
	listeip:=eip.ListOpts{}
	allEips,err:=eip.List(client,listeip)
	fmt.Println(out)
	if err != nil {
		panic(err)
	}
	for _, eip := range allEips {
		fmt.Printf("%+v\n", eip)
	}
Example to Create a Eip
	createOpts := eip.CreateOpts{
		 name = "bandwidth_test"
  		 size = 10
         share_type = "PER"
         type = "5_bgp"
	}
	eip, err := eip.Create(vpcClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
Example to Update a Eip
	eipID := "f588ccfa-8750-4d7c-bf5d-2ede24414706"
	updateOpts := vpcs.UpdateOpts{
		port_id = "f588ccfa-8750-4d7c-bf5d-2ede24414706"
	}
	eip, err := eip.Update(vpcClient, eipID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}
Example to Delete a Eip
	eipID := "f588ccfa-8750-4d7c-bf5d-2ede24414706"
	err := eip.Delete(vpcClient, eipID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package eip
