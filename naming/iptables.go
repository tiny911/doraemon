package naming

var tables = map[string]string{
	// 线下环境
	"172.33.0.14": "10.69.56.55",  //dev
	"172.32.0.7":  "10.69.58.241", //test
	"172.33.0.7":  "10.69.58.18",  //sandbox

	// 预览环境
	"172.31.16.46": "10.69.56.26", //online_pre
}

func switchIP(ip, env string) string {
	if env == "online" {
		return ip
	}

	if addr, exists := tables[ip]; exists {
		return addr
	}

	return ip
}
