package dns

import (
	"context"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
	"net"
)
//

type dnsConfig struct {
	address string
	port int
}

// mediator connect server(producer) and client(customer)
type mediator struct {
	clientAddr net.Addr
	body []byte
}

// realDnsServer represent a dns server
type realDnsServer struct {
	dnsConf dnsConfig
	server *net.UDPConn
}

func NewRealDnsServer()Operator{
	return &realDnsServer{
		dnsConf: dnsConfig{},
		server:  nil,
	}
}


func(op *realDnsServer)Run(ctx context.Context)error{
	var err error
	bindAddr := &net.UDPAddr{
		IP:   net.ParseIP(op.dnsConf.address),
		Port: op.dnsConf.port,
		Zone: "",
	}
	op.server,err = net.ListenUDP("udp", bindAddr)

	media := make(chan mediator)
	defer close(media)
	go func() {
		if err = op.runService(ctx, media);err != nil {
			logging.LogDnsServerController().WithError(err).Error("dns server run server module for accept request from client failed")
		}
	}()
	go func() {
		if err = op.runClient(ctx, media);err != nil {
			logging.LogDnsServerController().WithError(err).Error("dns client module run failed")
		}
	}()
	return err
}

func(op *realDnsServer)AddZone(object interface{})error{
	// todo
	return nil
}

func(op *realDnsServer)RemoveZone(object interface{})error{
	//todo
	return nil
}

func(op *realDnsServer)UpdateZone(object interface{})error{
	//todo
	return nil
}
func(op* realDnsServer)runService(ctx context.Context,object chan <- mediator)error{
	//todo
	var buffer []byte = make([]byte, 512)
	for{
		_, clientAddr, err := op.server.ReadFrom(buffer)
		if err != nil{
			return nil
		}
		select {
		case object<- mediator{
			clientAddr: clientAddr,
			body:       buffer,
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func(op *realDnsServer)runClient(ctx context.Context,object <- chan mediator)error{
	for ; ; {
		select {
		case tmpmedia := <- object:
			go op.handleRequestAndResponse(tmpmedia.clientAddr, tmpmedia.body)
		case <- ctx.Done():
			return ctx.Err()
		}
	}
}

func(op* realDnsServer)handleRequestAndResponse(clientAddr net.Addr, request []byte){
	packet := gopacket.NewPacket(request, layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	replyMess, _ := dnsPacket.(*layers.DNS)

	records := map[string]string{"google.com":"192.168.1.1", }

	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	var ok bool
	ip, ok = records[string(replyMess.Questions[0].Name)]
	if !ok {
		ip = "114.114.114.114"
	}
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(replyMess.Questions[0].Name)

	dnsAnswer.Class = layers.DNSClassIN
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	op.server.WriteTo(buf.Bytes(), clientAddr)
}