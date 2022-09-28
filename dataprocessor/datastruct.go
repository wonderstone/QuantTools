package dataprocessor

//  this edition is based on rwmutex and map

import (
	"errors"
	"time"
)

// stock and futures market data element to send at some specific timestamp for some target
type BarDE struct {
	BarTime     string             // 时间标签
	InstID      string             // 标的代码
	IndiDataMap map[string]float64 // 所有信息 key :indicator and bar element name , value: indicator and bar element value
}

func NewBarDE(BarTime string, InstID string, IndiDataMap map[string]float64) *BarDE {
	return &BarDE{
		BarTime:     BarTime,
		InstID:      InstID,
		IndiDataMap: IndiDataMap,
	}
}

// BarCombination with stockdata and futuresdata
// by using map with instrument id as its key
type BarC struct {
	Stockdata   map[string]*BarDE // use instID as key
	Futuresdata map[string]*BarDE
}

func (bc *BarC) GetTimeStamp() (string, error) {
	// check if stockdata is empty
	if len(bc.Stockdata) != 0 {
		for _, v := range bc.Stockdata {
			return v.BarTime, nil
		}
	} else if len(bc.Futuresdata) != 0 {
		for _, v := range bc.Futuresdata {
			return v.BarTime, nil
		}
	}
	return "", errors.New("BarC is empty")
}

func NewBarC(targetNum int) *BarC {
	sde := make(map[string]*BarDE, targetNum)
	fde := make(map[string]*BarDE, targetNum)
	return &BarC{
		Stockdata:   sde,
		Futuresdata: fde,
	}
}

type BarCM struct {
	// for further API scenario user defined part
	InstSIDS   []string //instrument stock ID slice
	IndiSNames []string //indicator stock names slice
	InstFIDS   []string //instrument futures ID slice
	IndiFNames []string //indicator futures names slice
	BeginDate  string
	EndDate    string

	BarCMap       map[string]*BarC              //datetime as key
	BarCMapkeydts []string                      //datetime key slice for iteration
	FMTMDataMap   map[string]map[string]float64 // key: datetime, value: map[instrID]MTMprice
}

func NewBarCM(instSIDS []string, IndiSNames []string, instFIDS []string, IndiFNames []string, beginDate string, endDate string) *BarCM {
	bd, error := time.Parse("2006/1/2 15:04", beginDate)
	if error != nil {
		panic("beginDate parse error")
	}
	ed, error := time.Parse("2006/1/2 15:04", endDate)
	if error != nil {
		panic("endDate parse error")
	}

	//依据分钟级别预估一下数据长度
	DurTrain := int(ed.Sub(bd).Hours() / 24) // 估算天数
	// make一次性分配内存可能会提升性能  不够动态扩展会涉及到资源迁移的读写 粗略按照6小时分配
	barCMap := make(map[string]*BarC, (DurTrain)*60*6)
	// nested map
	FMTMDtM := make(map[string]map[string]float64, (DurTrain))
	return &BarCM{
		InstSIDS:    instSIDS, //标的切片
		InstFIDS:    instFIDS, //标的切片
		BeginDate:   beginDate,
		EndDate:     endDate,
		BarCMap:     barCMap, //时间标签作为key的数据map
		FMTMDataMap: FMTMDtM,
	}
}
