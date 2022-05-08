package args

import "flag"

type RunParam struct {
	Version      string
	MasterIp     string
	LocalIp      string
	InstanceUuid string
}

var RunArgs RunParam

func init() {
	flag.StringVar(&RunArgs.Version, "version", "", "version")
	flag.StringVar(&RunArgs.MasterIp, "masterIp", "", "masterIp")
	flag.StringVar(&RunArgs.LocalIp, "localIp", "", "localIp")
	flag.StringVar(&RunArgs.InstanceUuid, "instanceUuid", "", "instanceUuid")
}
