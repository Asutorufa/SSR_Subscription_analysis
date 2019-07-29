package dns

import (
	"log"
	"net"
	"strconv"
	"strings"
)

// DNSv4 <-- dns for ipv4
func DNSv4(DNSServer, domain string) (domainIP string, success bool) {
	// +------------------------------+
	// |             id               |  16bit
	// +------------------------------+
	// |qr|opcpde|aa|tc|rd|ra|z|rcode |
	// +------------------------------+
	// |          QDCOUNT             |
	// +------------------------------+
	// |          ancount             |
	// +------------------------------+
	// |          nscount             |
	// +------------------------------+
	// |          arcount             |
	// +------------------------------+

	// • ID：这是由生成DNS查询的程序指定的16位的标志符。该标志符也被随后的应答报文所用，申请者利用这个标志将应答和原来的请求对应起来。

	// • QR：该字段占1位，用以指明DNS报文是请求（0）还是应答（1）。
	// • OPCODE：该字段占4位，用于指定查询的类型。值为0表示标准查询，值为1表示逆向查询，值为2表示查询服务器状态，值为3保留，值为4表示通知，值为5表示更新报文，值6～15的留为新增操作用。
	// • AA：该字段占1位，仅当应答时才设置。值为1，即意味着正应答的域名服务器是所查询域名的管理机构或者说是被授权的域名服务器。
	// • TC：该字段占1位，代表截断标志。如果报文长度比传输通道所允许的长而被分段，该位被设为1。
	// • RD：该字段占1位，是可选项，表示要求递归与否。如果为1，即意味 DNS解释器要求DNS服务器使用递归查询。

	// • RA：该字段占1位，代表正在应答的域名服务器可以执行递归查询，该字段与查询段无关。
	// • Z：该字段占3位，保留字段，其值在查询和应答时必须为0。
	// • RCODE：该字段占4位，该字段仅在DNS应答时才设置。用以指明是否发生了错误。
	// 允许取值范围及意义如下：
	// 0：无错误情况，DNS应答表现为无错误。
	// 1：格式错误，DNS服务器不能解释应答。
	// 2：严重失败，因为名字服务器上发生了一个错误，DNS服务器不能处理查询。
	// 3：名字错误，如果DNS应答来自于授权的域名服务器，意味着DNS请求中提到的名字不存在。
	// 4：没有实现。DNS服务器不支持这种DNS请求报文。
	// 5：拒绝，由于安全或策略上的设置问题，DNS名字服务器拒绝处理请求。
	// 6 ～15 ：留为后用。

	// • QDCOUNT：该字段占16位，指明DNS查询段中的查询问题的数量。
	// • ANCOUNT：该字段占16位，指明DNS应答段中返回的资源记录的数量，在查询段中该值为0。
	// • NSCOUNT：该字段占16位，指明DNS应答段中所包括的授权域名服务器的资源记录的数量，在查询段中该值为0。
	// • ARCOUNT：该字段占16位，指明附加段里所含资源记录的数量，在查询段中该值为0。
	// (2）DNS正文段
	// 在DNS报文中，其正文段封装在图7-42所示的DNS报文头内。DNS有四类正文段：查询段、应答段、授权段和附加段。

	// id -> 2byte
	// qr -> 1bit opcpde -> 4bit aa -> 1bit tc -> 1bit rd -> 1bit  => sum 1byte
	//  ra -> 1bit z -> 3bit rcode -> 4bit => sum 1byte
	// QDCOUNT -> 2byte
	// ANCOUNT -> 2byte
	// NSCOUNT -> 2byte
	// ARCOUNT -> 2byte
	header := []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var domainSet []byte
	for _, domain := range strings.Split(domain, ".") {
		domainSet = append(domainSet, byte(len(domain)))
		domainSet = append(domainSet, []byte(domain)...)
	}

	// append domain and qType And QClass
	// domain []byte(domainSet), 0x00
	// qType 0x00,0x01
	// qClass 0x00,0x01
	domainSetAndQTypeAndQClass := append([]byte(domainSet), 0x00, 0x00, 0x01, 0x00, 0x01)

	all := append(header, domainSetAndQTypeAndQClass...)

	conn, err := net.Dial("udp", DNSServer)
	if err != nil {
		log.Println(err)
		return "", false
	}
	defer conn.Close()

	// log.Println(all, len(all))

	var b [1024]byte
	conn.Write(all)
	// log.Println("write")
	// var b [1024]byte
	n, _ := conn.Read(b[:])
	// log.Println("header", b[0:12], "qr+opcode+aa+tc+rd:", b[2:3], "ra+z+rcode:", b[3], "rcode:", b[3]&1, "....", b[3]&2, b[3]&4, b[3]&8)
	if b[3]&1 != 0 {
		// log.Println("no such name")
		return "", false
	}
	// log.Println(b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n])
	// ip := b[bytes.Index(b[:n], []byte{192, 12})+2+2+2+4 : n]
	// log.Println("ip:", strconv.Itoa(int(ip[2]))+"."+strconv.Itoa(int(ip[3]))+"."+strconv.Itoa(int(ip[4]))+"."+strconv.Itoa(int(ip[5])))
	// log.Println(strconv.Itoa(int(b[n-4])) + "." + strconv.Itoa(int(b[n-3])) + "." + strconv.Itoa(int(b[n-2])) + "." + strconv.Itoa(int(b[n-1])))
	return strconv.Itoa(int(b[n-4])) + "." + strconv.Itoa(int(b[n-3])) + "." + strconv.Itoa(int(b[n-2])) + "." + strconv.Itoa(int(b[n-1])), true
}