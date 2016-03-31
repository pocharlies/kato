package main

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"fmt"
	"io/ioutil"
	"os"

	// Local:
	"github.com/h0tbird/kato/udata"
	"github.com/h0tbird/kato/providers/pkt"
	"github.com/h0tbird/kato/providers/ec2"

	// Community:
	"gopkg.in/alecthomas/kingpin.v2"
)

//-----------------------------------------------------------------------------
// Package variable declarations factored into a block:
//-----------------------------------------------------------------------------

var (

	//----------------------------
	// katoctl: top level command
	//----------------------------

	app = kingpin.New("katoctl", "Katoctl defines and deploys CoreOS clusters.")

	flUdataFile = app.Flag("user-data", "Path to file containing user data.").
			PlaceHolder("KATO_USER_DATA").
			OverrideDefaultFromEnvar("KATO_USER_DATA").
			Short('u').String()

	//-----------------------
	// udata: nested command
	//-----------------------

	cmdUdata = app.Command("udata", "Generate CoreOS cloud-config user-data.")

	flHostID = cmdUdata.Flag("hostid", "hostname = role-id").
			Required().PlaceHolder("CS_HOSTID").
			OverrideDefaultFromEnvar("CS_HOSTID").
			Short('i').String()

	flDomain = cmdUdata.Flag("domain", "Domain name as in (hostname -d)").
			Required().PlaceHolder("CS_DOMAIN").
			OverrideDefaultFromEnvar("CS_DOMAIN").
			Short('d').String()

	flRole = cmdUdata.Flag("role", "Choose one of [ master | node | edge ]").
			Required().PlaceHolder("CS_ROLE").
			OverrideDefaultFromEnvar("CS_ROLE").
			Short('r').String()

	flNs1Apikey = cmdUdata.Flag("ns1-api-key", "NS1 private API key.").
			Required().PlaceHolder("CS_NS1_API_KEY").
			OverrideDefaultFromEnvar("CS_NS1_API_KEY").
			Short('k').String()

	flCAcert = cmdUdata.Flag("ca-cert", "Path to CA certificate.").
			PlaceHolder("CS_CA_CERT").
			OverrideDefaultFromEnvar("CS_CA_CERT").
			Short('c').String()

	flEtcdToken = cmdUdata.Flag("etcd-token", "Provide an etcd discovery token.").
			PlaceHolder("CS_ETCD_TOKEN").
			OverrideDefaultFromEnvar("CS_ETCD_TOKEN").
			Short('e').String()

	flFlannelNetwork = cmdUdata.Flag("flannel-network", "Flannel entire overlay network.").
			PlaceHolder("CS_FLANNEL_NETWORK").
			OverrideDefaultFromEnvar("CS_FLANNEL_NETWORK").
			Short('n').String()

	flFlannelSubnetLen = cmdUdata.Flag("flannel-subnet-len", "Subnet len to llocate to each host.").
			PlaceHolder("CS_FLANNEL_SUBNET_LEN").
			OverrideDefaultFromEnvar("CS_FLANNEL_SUBNET_LEN").
			Short('s').String()

	flFlannelSubnetMin = cmdUdata.Flag("flannel-subnet-min", "Minimum subnet IP addresses.").
			PlaceHolder("CS_FLANNEL_SUBNET_MIN").
			OverrideDefaultFromEnvar("CS_FLANNEL_SUBNET_MIN").
			Short('m').String()

	flFlannelSubnetMax = cmdUdata.Flag("flannel-subnet-max", "Maximum subnet IP addresses.").
			PlaceHolder("CS_FLANNEL_SUBNET_MAX").
			OverrideDefaultFromEnvar("CS_FLANNEL_SUBNET_MAX").
			Short('x').String()

	flFlannelBackend = cmdUdata.Flag("flannel-backend", "Flannel backend type: [ udp | vxlan | host-gw | gce | aws-vpc | alloc ]").
			PlaceHolder("CS_FLANNEL_SUBNET_MAX").
			OverrideDefaultFromEnvar("CS_FLANNEL_SUBNET_MAX").
			Short('b').String()

	//------------------------------
	// setup-packet: nested command
	//------------------------------

	cmdSetupPacket = app.Command("setup-packet", "Setup a Packet.net project to be used by katoctl.")

	//----------------------------
	// run-packet: nested command
	//----------------------------

	cmdRunPacket = app.Command("run-packet", "Starts a CoreOS instance on Packet.net.")

	flPktAPIKey = cmdRunPacket.Flag("api-key", "Packet API key.").
			Required().PlaceHolder("PKT_APIKEY").
			OverrideDefaultFromEnvar("PKT_APIKEY").
			Short('k').String()

	flPktHostName = cmdRunPacket.Flag("hostname", "For the Packet.net dashboard.").
			Required().PlaceHolder("PKT_HOSTNAME").
			OverrideDefaultFromEnvar("PKT_HOSTNAME").
			Short('h').String()

	flPktProjID = cmdRunPacket.Flag("project-id", "Format: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee").
			Required().PlaceHolder("PKT_PROJID").
			OverrideDefaultFromEnvar("PKT_PROJID").
			Short('i').String()

	flPktPlan = cmdRunPacket.Flag("plan", "One of [ baremetal_0 | baremetal_1 | baremetal_2 | baremetal_3 ]").
			Required().PlaceHolder("PKT_PLAN").
			OverrideDefaultFromEnvar("PKT_PLAN").
			Short('p').String()

	flPktOsys = cmdRunPacket.Flag("os", "One of [ coreos_stable | coreos_beta | coreos_alpha ]").
			Required().PlaceHolder("PKT_OS").
			OverrideDefaultFromEnvar("PKT_OS").
			Short('o').String()

	flPktFacility = cmdRunPacket.Flag("facility", "One of [ ewr1 | ams1 ]").
			Required().PlaceHolder("PKT_FACILITY").
			OverrideDefaultFromEnvar("PKT_FACILITY").
			Short('f').String()

	flPktBilling = cmdRunPacket.Flag("billing", "One of [ hourly | monthly ]").
			Required().PlaceHolder("PKT_BILLING").
			OverrideDefaultFromEnvar("PKT_BILLING").
			Short('b').String()

	//---------------------------
	// setup-ec2: nested command
	//---------------------------

	cmdSetupEc2 = app.Command("setup-ec2", "Setup an EC2 elastic IP, VPC and firewall rules to be used by katoctl.")

	flSetupEc2Region = cmdSetupEc2.Flag("region", "EC2 region.").
			Required().PlaceHolder("EC2_REGION").
			OverrideDefaultFromEnvar("EC2_REGION").
			Short('r').String()

	//-------------------------
	// run-ec2: nested command
	//-------------------------

	cmdRunEc2 = app.Command("run-ec2", "Starts a CoreOS instance on Amazon EC2.")

	flEc2HostName = cmdRunEc2.Flag("hostname", "For the EC2 dashboard.").
			PlaceHolder("EC2_HOSTNAME").
			OverrideDefaultFromEnvar("EC2_HOSTNAME").
			Short('h').String()

	flEc2Region = cmdRunEc2.Flag("region", "EC2 region.").
			Required().PlaceHolder("EC2_REGION").
			OverrideDefaultFromEnvar("EC2_REGION").
			Short('r').String()

	flEc2ImageID = cmdRunEc2.Flag("image-id", "EC2 image id.").
			Required().PlaceHolder("EC2_IMAGE_ID").
			OverrideDefaultFromEnvar("EC2_IMAGE_ID").
			Short('i').String()

	flEc2InsType = cmdRunEc2.Flag("instance-type", "EC2 instance type.").
			Required().PlaceHolder("EC2_INSTANCE_TYPE").
			OverrideDefaultFromEnvar("EC2_INSTANCE_TYPE").
			Short('t').String()

	flEc2KeyPair = cmdRunEc2.Flag("key-pair", "EC2 key pair.").
			Required().PlaceHolder("EC2_KEY_PAIR").
			OverrideDefaultFromEnvar("EC2_KEY_PAIR").
			Short('k').String()

	flEc2VpcID = cmdRunEc2.Flag("vpc-id", "EC2 VPC id.").
			Required().PlaceHolder("EC2_VPC_ID").
			OverrideDefaultFromEnvar("EC2_VPC_ID").
			Short('v').String()

	flEc2SubnetIds = cmdRunEc2.Flag("subnet-ids", "EC2 subnet ids.").
			Required().PlaceHolder("EC2_SUBNET_ID").
			OverrideDefaultFromEnvar("EC2_SUBNET_ID").
			Short('s').String()

	flEc2ElasticIP = cmdRunEc2.Flag("elastic-ip", "Allocate an elastic IP [ true | false ]").
			Default("false").PlaceHolder("EC2_ELASTIC_IP").
			OverrideDefaultFromEnvar("EC2_ELASTIC_IP").
			Short('e').String()
)

//----------------------------------------------------------------------------
// Entry point:
//----------------------------------------------------------------------------

func main() {

	// Sub-command selector:
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	//---------------
	// katoctl udata
	//---------------

	case cmdUdata.FullCommand():

		udata := udata.Data {
			HostID:           *flHostID,
			Domain:           *flDomain,
			Role:             *flRole,
			Ns1ApiKey:        *flNs1Apikey,
			EtcdToken:        *flEtcdToken,
			FlannelNetwork:   *flFlannelNetwork,
			FlannelSubnetLen: *flFlannelSubnetLen,
			FlannelSubnetMin: *flFlannelSubnetMin,
			FlannelSubnetMax: *flFlannelSubnetMax,
			FlannelBackend:   *flFlannelBackend,
		}

		err := udata.Render()
		checkError(err)

	//----------------------
	// katoctl setup-packet
	//----------------------

	case cmdSetupPacket.FullCommand():

		pkt := pkt.Data {
		}

		err := setup(&pkt)
		checkError(err)

	//--------------------
	// katoctl run-packet
	//--------------------

	case cmdRunPacket.FullCommand():

		pkt := pkt.Data {
			APIKey:   *flPktAPIKey,
			HostName: *flPktHostName,
			Plan:     *flPktPlan,
			Facility: *flPktFacility,
			Osys:     *flPktOsys,
			Billing:  *flPktBilling,
			ProjID:   *flPktProjID,
		}

		err := run(&pkt)
		checkError(err)

	//-------------------
	// katoctl setup-ec2
	//-------------------

	case cmdSetupEc2.FullCommand():

		ec2 := ec2.Data {
			Region: *flSetupEc2Region,
		}

		err := setup(&ec2)
		checkError(err)

	//-----------------
	// katoctl run-ec2
	//-----------------

	case cmdRunEc2.FullCommand():

		ec2 := ec2.Data {
			Region:    *flEc2Region,
			SubnetIds: *flEc2SubnetIds,
			ImageID:   *flEc2ImageID,
			KeyPair:   *flEc2KeyPair,
			InsType:   *flEc2InsType,
			HostName:  *flEc2HostName,
			ElasticIP: *flEc2ElasticIP,
		}

		err := run(&ec2)
		checkError(err)
	}
}

//--------------------------------------------------------------------------
// func: readUdata
//--------------------------------------------------------------------------

func readUdata() ([]byte, error) {

	// Read data from file:
	if *flUdataFile != "" {
		udata, err := ioutil.ReadFile(*flUdataFile)
		return udata, err
	}

	// Read data from stdin:
	udata, err := ioutil.ReadAll(os.Stdin)
	return udata, err
}

//---------------------------------------------------------------------------
// func: checkError
//---------------------------------------------------------------------------

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
