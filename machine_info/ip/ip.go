package ip

import "DistriAI-Node/config"

type InfoIP struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func GetIpInfo() (InfoIP, error) {
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	//     return InfoIP{}, err
	// }

	// for _, addr := range addrs {
	//     if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//         if ipnet.IP.To4() != nil {
	// 			return InfoIP{IP: ipnet.IP.String()}, nil
	//         }
	//     }
	// }
	return InfoIP{
		IP:   config.GlobalConfig.Console.OuterNetIP,
		Port: config.GlobalConfig.Console.OuterNetPort,
	}, nil
}
