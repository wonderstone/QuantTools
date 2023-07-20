package datatunnel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/vdsdata"
	"google.golang.org/protobuf/proto"
)

const debug = false

// DataTunnel is a data tunnel for strategy
type DataTunnel struct {
	DataTimeStamp string
	//基础策略信息汇总
	StgM map[string]StgTargetsInfo
	// InfoFTSet map[InfoFT]void
	//公有标的区 由map存储key为频率value为StgTargets
	//此部分为向VDS订阅数据所需的信息
	TarM map[string]StgTargets
	//公有策略区 由map存储key为指标信息特征结构体value为策略名作为key的map[string]void
	//此部分为获得VDS数据后对应策略的处理
	InfoM map[InfoFT]map[string]void

	// StgM map[string]StgTargetsInfo
	//公有数据分发区 由map存储的datachannel,key为策略名
	SDataM map[string]chan *dataprocessor.BarC
	FDataM map[string]chan *dataprocessor.BarC
	// 私有数据缓冲sync.map
	cachDataSyncMap sync.Map
	// cachDataMap     map[string]dataprocessor.BarC
}

// NewDataTunnel is a constructor for DataTunnel

func NewDataTunnel() *DataTunnel {
	return &DataTunnel{
		StgM: make(map[string]StgTargetsInfo),
		// InfoFTSet: make(map[InfoFT]void),
		TarM:   make(map[string]StgTargets),
		InfoM:  make(map[InfoFT]map[string]void),
		SDataM: make(map[string]chan *dataprocessor.BarC),
		FDataM: make(map[string]chan *dataprocessor.BarC),
		// 私有数据缓冲map, key为策略名
		cachDataSyncMap: sync.Map{},
	}
}

// RegisterSTG：add a strategy info to DataTunnel
func (dt *DataTunnel) RegisterSTG(sti StgTargetsInfo) {
	// +基础策略信息汇总区StgM
	// check if the sti.StgName has already been registered
	if _, ok := dt.StgM[sti.StgName]; ok {
		// if yes, panic
		panic("strategy name already registered")
	} else {
		// if not, add the deepcopy of sti to dt.StgM
		// deepcopy of sti
		tmp := StgTargetsInfo{
			StgName: sti.StgName,
			STES:    make([]StgTargetsElement, len(sti.STES)),
			InfoET: InfoET{
				StgFreq:         sti.StgFreq,
				StgTimeTriggers: sti.StgTimeTriggers,
			},
		}
		for i, v := range sti.STES {
			tmp.STES[i].FreqType = v.FreqType
			tmp.STES[i].STargets = make(map[string]map[string]void)
			tmp.STES[i].FTargets = make(map[string]map[string]void)
			for k, vv := range v.STargets {
				tmp.STES[i].STargets[k] = make(map[string]void)
				tmp.STES[i].STargets[k] = DeepCopyMap(vv)
			}
			for k, vv := range v.FTargets {
				tmp.STES[i].FTargets[k] = make(map[string]void)
				tmp.STES[i].FTargets[k] = DeepCopyMap(vv)
			}
		}
		dt.StgM[sti.StgName] = tmp
	}
	// +公有标的区TarM由map存储key为频率value为StgTargets
	// iter the sti.STES slice
	for _, v := range sti.STES {

		// check if the TarM has the same key as sti.Frequency
		if _, ok := dt.TarM[v.FreqType]; ok {
			// iter the sti.StgTargets.STargets key to check if the tarM has the same key as sti.Frequency
			for k := range v.StgTargets.STargets {
				// make an InfoFT struct instance and check if it is in dt.InfoFTSet
				// if k is not in tarM's STargets, add the StgTargets to tarM
				if _, ok := dt.TarM[v.FreqType].STargets[k]; !ok {
					dt.TarM[v.FreqType].STargets[k] = make(map[string]void)
					dt.TarM[v.FreqType].STargets[k] = v.StgTargets.STargets[k]

				} else {
					// if k is in tarM's STargets, iter its value to check if STargets[k] has it.
					for kk := range v.StgTargets.STargets[k] {
						if _, ok := dt.TarM[v.FreqType].STargets[k][kk]; !ok {
							dt.TarM[v.FreqType].STargets[k][kk] = void{}
						}
					}
				}
			}
			// iter the sti.StgTargets.FTargets key to check if the tarM has the same key as sti.Frequency
			for k := range v.StgTargets.FTargets {
				// make an InfoFT struct instance and check if it is in dt.InfoFTSet
				// if k is not in tarM's FTargets, add the StgTargets to tarM
				if _, ok := dt.TarM[v.FreqType].FTargets[k]; !ok {
					dt.TarM[v.FreqType].FTargets[k] = make(map[string]void)
					dt.TarM[v.FreqType].FTargets[k] = v.StgTargets.FTargets[k]

				} else {
					// if k is in tarM's FTargets, iter its value to check if FTargets[k] has it.
					for kk := range v.StgTargets.FTargets[k] {
						if _, ok := dt.TarM[v.FreqType].FTargets[k][kk]; !ok {
							dt.TarM[v.FreqType].FTargets[k][kk] = void{}
						}
					}
				}

			}
		} else {
			// if not, add the StgTargets to tarM
			dt.TarM[v.FreqType] = v.StgTargets
		}
	}

	// +公有策略区 InfoM的key为指标信息特征结构体value为策略名作为key的map[string]void

	for _, v := range sti.STES {
		for vv := range v.STargets {
			// 构建InfoFT信息结构体 并检查是否InfoM中该key是否已经存在
			// 如果不存在则添加
			if _, ok := dt.InfoM[InfoFT{vv, v.FreqType}]; !ok {
				dt.InfoM[InfoFT{vv, v.FreqType}] = make(map[string]void)
				dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName] = void{}
			} else {
				// 如果存在则检查是否已经存在该策略名
				if _, ok := dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName]; !ok {
					dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName] = void{}
				}

			}

		}
		for vv := range v.FTargets {
			// 构建InfoFT信息结构体 并检查是否InfoM中该key是否已经存在
			// 如果不存在则添加
			if _, ok := dt.InfoM[InfoFT{vv, v.FreqType}]; !ok {
				dt.InfoM[InfoFT{vv, v.FreqType}] = make(map[string]void)
				dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName] = void{}
			} else {
				// 如果存在则检查是否已经存在该策略名
				if _, ok := dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName]; !ok {
					dt.InfoM[InfoFT{vv, v.FreqType}][sti.StgName] = void{}
				}

			}

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
	// +私有数据缓冲map,key is the stgName,value is dataprocessor.BarC

}

// RemoveSTG：remove a strategy info from DataTunnel
func (dt *DataTunnel) RemoveSTG(sti StgTargetsInfo) {
	// -基础策略信息汇总区StgM
	delete(dt.StgM, sti.StgName)

	// -公有标的区 TarM由map存储key为频率value为StgTargets
	// create the map from dt.StgM
	tmpStgTargets := make(map[string]StgTargets)
	for _, v := range dt.StgM {
		// tmpStgTargets‘ key is the frequency of the data
		// check if the tmpStgTargets has the same key as v.STES element's FreqType
		for _, vv := range v.STES {
			if _, ok := tmpStgTargets[vv.FreqType]; !ok {
				tmpStgTargets[vv.FreqType] = vv.StgTargets
			} else {
				// if tmpStgTargets has the same key as v.STES element's FreqType, then iter the v.STES element's StgTargets
				// iter the sti.StgTargets.STargets key to check if the tarM has the same key as sti.Frequency
				for k := range vv.StgTargets.STargets {
					// if k is not in tarM's STargets, add the StgTargets to tarM
					if _, ok := tmpStgTargets[vv.FreqType].STargets[k]; !ok {
						tmpStgTargets[vv.FreqType].STargets[k] = make(map[string]void)
						tmpStgTargets[vv.FreqType].STargets[k] = vv.StgTargets.STargets[k]

					} else {
						// if k is in tarM's STargets, iter its value to check if STargets[k] has it.
						for kk := range vv.StgTargets.STargets[k] {
							if _, ok := tmpStgTargets[vv.FreqType].STargets[k][kk]; !ok {
								tmpStgTargets[vv.FreqType].STargets[k][kk] = void{}
							}
						}
					}
				}
				// iter the sti.StgTargets.FTargets key to check if the tarM has the same key as sti.Frequency
				for k := range vv.StgTargets.FTargets {
					// if k is not in tarM's FTargets, add the StgTargets to tarM
					if _, ok := tmpStgTargets[vv.FreqType].FTargets[k]; !ok {
						tmpStgTargets[vv.FreqType].FTargets[k] = make(map[string]void)
						tmpStgTargets[vv.FreqType].FTargets[k] = vv.StgTargets.FTargets[k]

					} else {
						// if k is in tarM's FTargets, iter its value to check if FTargets[k] has it.
						for kk := range vv.StgTargets.FTargets[k] {
							if _, ok := tmpStgTargets[vv.FreqType].FTargets[k][kk]; !ok {
								tmpStgTargets[vv.FreqType].FTargets[k][kk] = void{}
							}
						}
					}
				}
			}
		}
	}
	dt.TarM = tmpStgTargets

	// -公有策略区 InfoM的key为指标信息特征结构体value为策略名作为key的map[string]void
	// -基于dt.StgM的遍历，重新生成
	dt.InfoM = make(map[InfoFT]map[string]void)
	for _, v := range dt.StgM {
		// iter STES and check if the InfoM has the same key as v.STES element's FreqType
		for _, vv := range v.STES {
			for k := range vv.StgTargets.STargets {
				if _, ok := dt.InfoM[InfoFT{k, vv.FreqType}]; !ok {
					dt.InfoM[InfoFT{k, vv.FreqType}] = make(map[string]void)
					dt.InfoM[InfoFT{k, vv.FreqType}][sti.StgName] = void{}
				}
			}
			for k := range vv.StgTargets.FTargets {
				if _, ok := dt.InfoM[InfoFT{k, vv.FreqType}]; !ok {
					dt.InfoM[InfoFT{k, vv.FreqType}] = make(map[string]void)
					dt.InfoM[InfoFT{k, vv.FreqType}][sti.StgName] = void{}
				}
			}

		}
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
	// -私有数据缓冲map cachDataMap
	// check if the cachDataSyncMap has the same key as sti.StgName
	dt.cachDataSyncMap.Delete(sti.StgName)

}

// process the data in tunnel
func (dt *DataTunnel) ProcessData(ip string, port int, stopsignal chan bool, reqmaps []map[string]string, targetsChan chan []map[string]string) {
	// 0. Use channal to pass the *dataprocessor.BarDE data
	// 0.1 create the channal for BarDEFreq

	BarDEchan := make(chan *BarDEFreq)
	// 1. Get data from VDS TCP server
	go dt.GetTargetsData(ip, port, stopsignal, reqmaps, targetsChan, BarDEchan)
	// 1.1 loop the chan and get the data from BarDEchan
	for {
		tmpBarDEFreq := <-BarDEchan
		// 1.2 combine the TargetName and FreqType to get the struct then
		InfoFTKey := InfoFT{tmpBarDEFreq.InstID, tmpBarDEFreq.Freq}
		// 1.3 check if the InfoM has the same key as InfoFTKey
		if _, ok := dt.InfoM[InfoFTKey]; ok {
			// 1.4 if ok, then iter the InfoM[InfoFTKey] to get the stgName
			for stgN := range dt.InfoM[InfoFTKey] {
				// 2. Prepare Data for each strategy on each frequency
				dt.UpdateCachDataSyncMap(stgN, &tmpBarDEFreq.BarDE)
				// 2.1 only update the TimeStamp when DataCach updated
				dt.DataTimeStamp = tmpBarDEFreq.BarDE.BarTime
				// 3. Send data to indicator chain for each strategy on time trigger

				// 4. Send data to strategy channel
			}
		}
	}

}

// Lazy man! No get set method! Use public field directly!

// TimeTrigger check
func (dt *DataTunnel) TimeTrigger(iet InfoET, rl string) bool {
	// check if the ts is in the TimeTriggerMap
	switch rl {
	case "S", "s", "second", "Sec", "Seconds", "seconds":
		fmt.Println("update the timestamp to Seconds resolution level")
	case "M", "m", "minute", "Min", "Minutes", "minutes":
		fmt.Println("Minutes resolution level")
	case "H", "h", "hour", "Hour", "Hours", "hours":
		fmt.Println("Hours resolution level")
	default:
		fmt.Println("default")
	}
	switch iet.StgFreq {
	case "Daily":
		fmt.Println("Daily")
	case "Min":
		fmt.Println("Min")
	default:
		panic("wrong frequency type")
	}

	return true
}

// 指标链处理函数
func (dt *DataTunnel) IndicatorChainHandler() {

}

// update the data in cachDataSyncMap
func (dt *DataTunnel) UpdateCachDataSyncMap(stgName string, pbde *dataprocessor.BarDE) {
	// check if the stgName is in the cachDataSyncMap
	if tmp, ok := dt.cachDataSyncMap.Load(stgName); !ok {
		dt.cachDataSyncMap.Store(stgName, dataprocessor.NewBarCSimple())
	} else {
		// change the type of tmp to *dataprocessor.BarC
		tmp.(*dataprocessor.BarC).Stockdata[pbde.InstID] = pbde
		dt.cachDataSyncMap.Store(stgName, tmp)
	}
}

// get targets data from VDS TCP server, receive signal from signal channel to close the connection
func (dt *DataTunnel) GetTargetsData(ip string, port int, stopsignal chan bool, reqmaps []map[string]string, targetsChan chan []map[string]string, BarDEchan chan *BarDEFreq) {
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
			case s := <-stopsignal:
				// todo the logic here needs to be adhere to the Terminal Manufactor（say what pause、restart、stop mean）
				if s {
					conn.Close()
					return
				}
				// todo end!
			default:
				SubProcessData(conn, reqmaps, targetsChan, BarDEchan)
			}
		}
	}

}

// turn an int32 into a byte array with 4 bytes

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

// turn a byte array with 4 bytes into an int32
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

// compare two data source speed

// choose the faster one

// compute the indicators

// func to deepcopy a map
func DeepCopyMap(m map[string]void) map[string]void {
	newMap := make(map[string]void)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// func to sub and process the data
func SubProcessData(conn net.Conn, tmpreqMaps []map[string]string, targetsChan chan []map[string]string, BarDEchan chan *BarDEFreq) {
	// 循环发送所有请求reqmap
	Subdata(conn, tmpreqMaps)
	var doublemsg vdsdata.DoubleMsg
	var stringmsg vdsdata.StringMsg
	var int32msg vdsdata.Int32Msg
	var int64msg vdsdata.Int64Msg
	for {
		// 先check是否targetsChan有数据发送过来
		select {
		case tmpreqMaps := <-targetsChan:
			Subdata(conn, tmpreqMaps)
			fmt.Println("now the sub data is ", tmpreqMaps)
		default:
		}
		//读消息头
		datalen := make([]byte, 4)
		_, err := io.ReadFull(conn, datalen)
		if err != nil {
			panic(err)
		}
		// turn datalen into int
		dtlen := BytesToInt(datalen)
		buf := make([]byte, dtlen)
		_, e := io.ReadFull(conn, buf)
		if e != nil {
			panic(e)
		}
		// 数据解析转换
		var s = vdsdata.VDSRsp{}
		err = proto.Unmarshal(buf, &s)
		if err != nil {
			panic(err)
		}
		// new a BarDEFreq
		tmpBarDEFreq := BarDEFreq{
			// dataprocessor.BarDE part
			BarDE: dataprocessor.BarDE{
				// make the map
				IndiDataMap: make(map[string]float64),
			},
		}
		for k, v := range s.RspMap {
			// new a BarDE and update the data
			switch k {
			case "Msgtype":
				err := v.GetValue().UnmarshalTo(&stringmsg)
				if err != nil {
					panic(err)
				} else {
					tmpBarDEFreq.Freq = stringmsg.Data
					// 	// fmt.Println(k, stringmsg.Data)
				}
			case "Symbol":
				err := v.GetValue().UnmarshalTo(&stringmsg)
				if err != nil {
					panic(err)
				} else {
					tmpBarDEFreq.InstID = stringmsg.Data
					// tmpInfoFT.TargetName = stringmsg.Data
					// fmt.Println(k, stringmsg.Data)
				}

			case "Updatetime":
				err := v.GetValue().UnmarshalTo(&int32msg)
				if err != nil {
					panic(err)
				} else {
					tmpBarDEFreq.BarTime = strconv.Itoa(int(int32msg.Data))
					// fmt.Println(k, int32msg.Data)
				}

			case "Volume":
				err := v.GetValue().UnmarshalTo(&int64msg)
				if err != nil {
					panic(err)
				} else {
					tmpBarDEFreq.IndiDataMap["Volume"] = float64(int64msg.Data)
					// fmt.Println(k, int64msg.Data)
				}

			default:
				err := v.GetValue().UnmarshalTo(&doublemsg)
				if err != nil {
					panic(err)
				} else {
					tmpBarDEFreq.IndiDataMap[k] = doublemsg.Data
					// fmt.Println(k, doublemsg.Data)
				}
			}
		}
		// send the BarDE to the channel
		BarDEchan <- &tmpBarDEFreq
		// fmt.Println("++++++")
		// fmt.Println("BarDE is ", tmpBarDE)
		// fmt.Println("InfoFT is ", tmpInfoFT)
		// fmt.Println("++++++")
		// fmt.Println("now: ", time.Now())
		// fmt.Println("*******************")
	}
}

func Subdata(conn net.Conn, reqmaps []map[string]string) {
	// 循环发送所有请求reqmap
	for _, reqmap := range reqmaps {
		sub := &vdsdata.VDSReq{
			ReqMap: reqmap,
		}
		sData, err := proto.Marshal(sub)
		if err != nil {
			panic(err)
		}
		slen := len(sData)
		q := IntToBytes(slen)
		var b bytes.Buffer
		b.Write(q)
		b.Write([]byte(sData))
		conn.Write([]byte(b.Bytes()))
	}
}
