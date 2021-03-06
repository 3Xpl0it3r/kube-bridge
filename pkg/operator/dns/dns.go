package dns

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"l0calh0st.cn/k8s-bridge/configure"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
	"net"
)
//

type dnsConfig struct {
	address string
	port int
}


var globalConfig *configure.Config = configure.NewConfig()


// realDnsServer represent a dns server
type realDnsServer struct {
	dnsConf dnsConfig
	server *net.UDPConn
	cache map[string]string
}

func NewRealDnsServer()Operator{
	return &realDnsServer{
		dnsConf: dnsConfig{},
		server:  nil,
		cache: make(map[string]string),
	}
}

//
// mediator connect server(producer) and client(customer)
type mediator struct {
	clientAddr net.Addr
	body []byte
}

// Run is the entrypoint of real dns server
func(op *realDnsServer)Run(ctx context.Context)error{
	//if op.dnsConf == (dnsConfig{}) {
	//	op.dnsConf = *op.getDefaultDnsConfig()
	//}
	var err error
	bindAddr := &net.UDPAddr{
		IP:   net.ParseIP(globalConfig.Dns.Bind),
		Port: globalConfig.Dns.Port,
		Zone: "",
	}
	op.server,err = net.ListenUDP("udp", bindAddr)
	if err != nil {
		logrus.WithError(err).Error("real dns server start filed")
	}

	media := make(chan mediator)
	go func() {
		defer close(media)
		if err = op.runService(ctx, media);err != nil {
			logging.LogDnsServerController().WithError(err).Error("dns server run server module for accept request from client failed")
		}
	}()
	go func() {
		if err = op.runClient(ctx, media);err != nil {
			logging.LogDnsServerController().WithError(err).Error("dns client module run failed")
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}

func(op *realDnsServer)AddZone(object interface{})error{
	// todo
	var record map[string]string
	var ok bool
	if record, ok = object.(map[string]string);!ok{
		return fmt.Errorf("delete zone from dns cahe, record except map[string]string, but got %v\n", object)
	}
	for k,v := range record{
		op.cache[k] = v
	}
	logging.LogDnsServerController().WithField("Type", "AddResult").Warnln(op.cache)
	return nil
}

func(op *realDnsServer)RemoveZone(object interface{})error{
	var record map[string]string
	var ok bool
	if record, ok = object.(map[string]string);!ok{
		return fmt.Errorf("delete zone from dns cahe, record except map[string]string, but got %v\n", object)
	}
	for k,_ := range record{
		delete(op.cache, k)
		if !ok {

		} else {

		}
	}
	logging.LogDnsServerController().WithField("Type", "RemoveResult").Warnln(op.cache)
	return nil
}

func(op *realDnsServer)UpdateZone(object interface{})error{
	var record map[string]string
	var ok bool
	if record,ok = object.(map[string]string);!ok {
		return fmt.Errorf("update zone to dns cache failed, record expcet map[string]string, but got %v", object)
	}
	for k,v := range record {
		op.cache[k] = v
	}
	logging.LogDnsServerController().WithField("Type", "UpdateResult").Warnln(op.cache)
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
	// logfor del

	packet := gopacket.NewPacket(request, layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	replyMess, _ := dnsPacket.(*layers.DNS)


	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	var ok bool
	ip, ok = op.cache[string(replyMess.Questions[0].Name)]
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



func(op *realDnsServer)getDefaultDnsConfig()*dnsConfig{
	return &dnsConfig{
		address: "0.0.0.0",
		port:    1053,
	}
}