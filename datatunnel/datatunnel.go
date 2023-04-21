package datatunnel

import (
	"fmt"
	"github.com/wonderstone/QuantTools/vdsdata"
	"net"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/dataprocessor"
)

const debug = false

type void struct{}

// StgTargets struct
type StgTargets struct {
	STargets map[string]void
	FTargets map[string]void
}

// StgTargetsElement struct
type StgTargetsElement struct {
	Frequency string
	StgTargets
}

type StgTargetsInfo struct {
	StgName string
	// StgTargetsElement slice
	STES []StgTargetsElement
}

// DataTunnel is a data tunnel for strategy
type DataTunnel struct {
	//公有标的区 由map存储key为频率value为StgTargets
	TarM map[string]StgTargets
	//公有策略区 由map存储key为策略名value为StgTargetsInfo切片
	StgM map[string]StgTargetsInfo

	// StgM map[string]StgTargetsInfo
	//公有数据分发区 由map存储的datachannel,key为策略名
	SDataM map[string]chan *dataprocessor.BarC
	FDataM map[string]chan *dataprocessor.BarC
}

// NewDataTunnel is a constructor for DataTunnel

func NewDataTunnel() *DataTunnel {
	return &DataTunnel{
		TarM:   make(map[string]StgTargets),
		StgM:   make(map[string]StgTargetsInfo),
		SDataM: make(map[string]chan *dataprocessor.BarC),
		FDataM: make(map[string]chan *dataprocessor.BarC),
	}
}

// RegisterSTG：add a strategy info to DataTunnel
func (dt *DataTunnel) RegisterSTG(sti StgTargetsInfo) {
	// +公有策略区
	// make a deep copy of sti all elements with copy
	stiCopy := StgTargetsInfo{
		StgName: sti.StgName,
		STES:    make([]StgTargetsElement, len(sti.STES)),
	}
	for i, v := range sti.STES {
		stiCopy.STES[i] = StgTargetsElement{
			Frequency: v.Frequency,
			StgTargets: StgTargets{
				STargets: make(map[string]void, len(v.STargets)),
				FTargets: make(map[string]void, len(v.FTargets)),
			},
		}
		for k := range v.STargets {
			stiCopy.STES[i].STargets[k] = void{}
		}
		for k := range v.FTargets {
			stiCopy.STES[i].FTargets[k] = void{}
		}
	}
	dt.StgM[sti.StgName] = stiCopy
	// +公有标的区
	// iter the sti.STES slice
	for _, v := range sti.STES {
		// check if the TarM has the same key as sti.Frequency
		if _, ok := dt.TarM[v.Frequency]; ok {
			// iter the sti.StgTargets.STargets key to check if the tarM has the same key as sti.Frequency
			for k := range v.StgTargets.STargets {
				// if not, add the StgTargets to tarM
				dt.TarM[v.Frequency].STargets[k] = void{}
			}
			// iter the sti.StgTargets.FTargets key to check if the tarM has the same key as sti.Frequency
			for k := range v.StgTargets.FTargets {
				// if not, add the StgTargets to tarM
				dt.TarM[v.Frequency].FTargets[k] = void{}
			}
		} else {
			// if not, add the StgTargets to tarM
			dt.TarM[v.Frequency] = v.StgTargets
		}
	}
	// +公有数据分发区
	// check if the SDataM has the same key as sti.StgName
	if _, ok := dt.SDataM[sti.StgName]; !ok {
		dt.SDataM[sti.StgName] = make(chan *dataprocessor.BarC)
	}
	// check if the FDataM has the same key as sti.StgName
	if _, ok := dt.FDataM[sti.StgName]; !ok {
		dt.FDataM[sti.StgName] = make(chan *dataprocessor.BarC)
	}
	// DCE: debug info
	if debug {
		// this part is for test only
		log.Info().Str("Stg Name", sti.StgName).
			Msg("Stg Registered")
	}
}

// RemoveSTG：remove a strategy info from DataTunnel
func (dt *DataTunnel) RemoveSTG(sti StgTargetsInfo) {
	// -公有策略区
	delete(dt.StgM, sti.StgName)
	// -公有标的区
	// create a temp map to store all sti frequency as key and void as value
	tmpFrequencies := make(map[string]void)
	tmpStgTargets := make(map[string]StgTargets)

	for _, v := range sti.STES {
		tmpFrequencies[v.Frequency] = void{}
		// add a new key to tmpStgTargets with empty StgTargets

		tmpStgTargets[v.Frequency] = StgTargets{
			STargets: make(map[string]void),
			FTargets: make(map[string]void),
		}
	}

	// iter all the StgM.StgTargetsInfo to check if  sti.Frequency is the same as the StgM.StgTargetsInfo.Frequency
	for _, v := range dt.StgM {
		// iter the StgM.StgTargetsInfo.STES slice to check if sti.Frequency is in tmpStgTargets
		for _, ste := range v.STES {
			// if yes
			if _, ok := tmpFrequencies[ste.Frequency]; ok {
				// add all ste.StgTargets.STargets to tmpStgTargets[ste.Frequency].STargets
				for k := range ste.StgTargets.STargets {
					tmpStgTargets[ste.Frequency].STargets[k] = void{}
				}
				// add all ste.StgTargets.FTargets to tmpStgTargets[ste.Frequency].FTargets
				for k := range ste.StgTargets.FTargets {
					tmpStgTargets[ste.Frequency].FTargets[k] = void{}
				}
			}

		}

	}
	// reassigned to tarM
	for k := range tmpFrequencies {
		dt.TarM[k] = tmpStgTargets[k]
	}

	// -公有数据分发区
	// check if the SDataM has the same key as sti.StgName
	delete(dt.SDataM, sti.StgName)
	// check if the FDataM has the same key as sti.StgName
	delete(dt.FDataM, sti.StgName)
	if debug {
		// this part is for test only
		log.Info().Str("Stg Removed", sti.StgName)
	}
}

// Lazy man! No get set method! Use public field directly!

// get targets data from VDS TCP server, receive signal from signal channel to close the connection
func (dt *DataTunnel) GetTargetsData(ip string, port int, signal chan bool) {
	for {
		// create a tcp connection
		// conn, err := net.Dial("tcp", "123.138.216.197:9009")
		conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
		if err != nil {
			fmt.Println("Failed to connect to server:", err)
			time.Sleep(time.Second)
			continue
		}
		// 在这里处理连接成功后的操作
		fmt.Println("Connected to server.")
		// 不断读取服务器发送的数据
		for {
			select {
			case <-signal:
				conn.Close()
				return
			default:
				// 这里不对
				sub := &vdsdata.VDSSub{
					Symbol:  "600000.SH",
					Exch:    vdsdata.Exch_SH,
					Sectype: vdsdata.VDSSecType_STOCK,
					//Subtype: vdsdata.VDSInterfaceType_SubSnapshot,
				}
				pData, err := proto.Marshal(sub)
				fmt.Println(sub, pData)
				if err != nil {
					panic(err)
				}
				tcp := &vdsdata.VDSTcp{
					Itype:    vdsdata.VDSInterfaceType_SubSnapshot,
					Data:     pData,
					Userdata: []byte{},
				}
				sData, err := proto.Marshal(tcp)
				slen := len(sData)
				var d [4]byte
				p := unsafe.Pointer(&slen)
				q := (*[4]byte)(p)
				//fmt.Println(p, q)
				copy(d[0:], (*q)[0:])
				//blen := bytes(slen)
				var b bytes.Buffer
				b.Write([]byte(d[:]))
				b.Write([]byte(sData))
				//fmt.Println(b.Bytes())
				conn.Write([]byte(b.Bytes()))
				buf := make([]byte, 40960, 40960)
				//var data_storage map[string]vdsdata.VDSSnapshot
				for {
					//读消息
					_, err := conn.Read(buf)
					if err != nil {
						panic(err)
					}
					//fmt.Println(buf)
					//buf = buf[4:]
					//fmt.Println("Recv data len:", cnt)
					var s = vdsdata.VDSTcp{}
					err = proto.Unmarshal(buf, &s)
					var snapshot vdsdata.VDSSnapshot
					err = proto.Unmarshal(s.Data, &snapshot)
					//fmt.Println(snapshot.Date, snapshot.Buylevel, snapshot.Close, snapshot.Exch, snapshot.Symbol, snapshot.Open)
					if snapshot.Uptetime != 0 {
						fmt.Println(snapshot.Uptetime, snapshot.Symbol)
					}
					if snapshot.Symbol == "600000" {
						fmt.Println(snapshot.Last, snapshot.Symbol)
					}
				}
			}
		}
	}

}
