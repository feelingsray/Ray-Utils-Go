package commonIota

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/feelingsray/Ray-Utils-Go/tools/common"
)

// IOTA 数联网接入模块
// 与北大团队数联网进行数据授权接口模块
type IOTA struct {
	Url      string // IOTA节点或平台地址
	Contract string // IOTA所用协议
	Timeout  int    // IOTA平台请求超时时间
	//Db       *mgo.Database // 网关端才会用到这个，平台传nil即可
	//Cached   *ledis.DB
}

// 给数联网开一个接口,输入DOID，DataHash，时间戳，
//返回TransactionHash

/*
 * NewIOTA 数联网实例
 * url: http://IP:Port/SCIDE/SCManager
 * contract: ShanxiProxy
 */

func NewIOTA(url, contract string, timeout int) (*IOTA, error) {
	iotaModel := new(IOTA)
	iotaModel.Url = url
	iotaModel.Contract = contract
	iotaModel.Timeout = timeout
	//if db != nil {
	//	iotaModel.Db = db
	//}
	//if cached != nil {
	//	iotaModel.Cached = cached
	//}
	return iotaModel, nil
}

//func (iota *IOTA) CreateRepo(name, suffix string) (string, string, error) {
//	body := make(map[string]any)
//	body["action"] = "executeContract"
//	body["contractID"] = iota.Contract
//	body["operation"] = "createRepo"
//	arg := make(map[string]any)
//	arg["name"] = name
//	arg["suffix"] = suffix
//	body["arg"] = arg
//	bodyByte, err := common.DumpJson(body)
//	if err != nil {
//		return "", "", err
//	}
//	resp, err := httpDo(iota.Url, "POST", nil, string(bodyByte), iota.Timeout)
//	if err != nil {
//		return "", "", err
//	}
//	token := gjson.Get(resp, "token").String()
//	repoId := gjson.Get(resp, "repoId").String()
//	if token == "" || repoId == "" {
//		return "", "", errors.New("token or repoId is empty")
//	}
//	return token, repoId, nil
//}

// RegistryDOID 注册获取DOID和DataHash
// 网关上传数据时调用
func (iota *IOTA) RegistryDOID(csDataTime int64, repoId, token string, data map[string]any) (*Ledger, error) {
	jsonData, err := common.DumpJson(data)
	if err != nil {
		return nil, err
	}
	body := make(map[string]any)
	body["action"] = "executeContract"
	body["contractID"] = iota.Contract
	body["operation"] = "registryDOID"
	arg := make(map[string]any)
	arg["data"] = string(jsonData)
	arg["repoId"] = repoId
	arg["token"] = token
	body["arg"] = arg
	bodyByte, err := common.DumpJson(body)
	if err != nil {
		return nil, err
	}
	resp, err := httpDo(iota.Url, "POST", nil, string(bodyByte), iota.Timeout)
	if err != nil {
		return nil, err
	}
	respStatus := gjson.Get(resp, "status").String()
	respResult := gjson.Get(resp, "result")
	if respStatus != "Success" {
		return nil, errors.New(fmt.Sprintf("IOTA接口调用错误:%s", respResult.String()))
	}
	doid := respResult.Get("doId").String()
	dataHash := respResult.Get("dataHash").String()
	if doid == "" || dataHash == "" {
		return nil, errors.New(fmt.Sprintf("返回数据格式错误:%s", resp))
	}
	ledger := new(Ledger)
	ledger.DOID = doid
	ledger.DataHash = dataHash
	//ledger.TransactionHash = ""
	ledger.Data = data
	ledger.CsDataTime = csDataTime
	//_, err = iota.Db.C("t_ledger").Upsert(bson.M{"doid": ledger.DOID, "data_hash": ledger.DataHash}, ledger)
	//if err != nil {
	//	return nil, err
	//}
	return ledger, nil
}

// SendTransaction 发送存证(异步调用)返回账本Hash
//func (iota *IOTA) SendTransaction(ledger *Ledger) (*Ledger, error) {
//	body := make(map[string]any)
//	body["action"] = "executeContract"
//	body["contractID"] = iota.Contract
//	body["operation"] = "sendTransaction"
//	arg := make(map[string]any)
//	arg["doId"] = ledger.DOID
//	arg["dataHash"] = ledger.DataHash
//	body["arg"] = arg
//	ok := false
//	if ok {
//		bodyByte, err := common.DumpJson(body)
//		if err != nil {
//			return nil, err
//		}
//		resp, err := httpDo(iota.Url, "POST", nil, string(bodyByte), iota.Timeout)
//		if err != nil {
//			return nil, err
//		}
//		respStatus := gjson.Get(resp, "status").String()
//		respResult := gjson.Get(resp, "result")
//		if respStatus != "Success" {
//			return nil, errors.New(fmt.Sprintf("IOTA接口调用错误:%s", respResult.String()))
//		}
//		transactionHash := respResult.Get("transactionHash").String()
//		if transactionHash == "" {
//			return nil, errors.New(fmt.Sprintf("返回数据格式错误:%s", resp))
//		}
//		ledger.TransactionHash = transactionHash
//	} else {
//		ledger.TransactionHash = "471c95b24dde9332f7db3981bd29201bb535bd0a"
//	}
//	return ledger, nil
//}

func (iota *IOTA) VerifyDataByDoID(ledger *Ledger) (bool, string, error) {
	body := make(map[string]any)
	body["action"] = "executeContract"
	body["contractID"] = iota.Contract
	body["operation"] = "verifyDataAndHash"
	arg := make(map[string]any)
	jsonData, err := common.DumpJson(ledger.Data)
	if err != nil {
		return false, "", err
	}
	arg["data"] = string(jsonData)
	arg["doId"] = ledger.DOID
	body["arg"] = arg
	jsonBodyByte, err := common.DumpJson(body)
	if err != nil {
		return false, "", err
	}
	resp, err := httpDo(iota.Url, "POST", nil, string(jsonBodyByte), iota.Timeout)
	if err != nil {
		return false, "", err
	}
	finalResult := gjson.Get(resp, "finalResult").Bool()
	url := gjson.Get(resp, "url").String()
	if finalResult != true {
		return finalResult, "", errors.New(fmt.Sprintf("IOTA接口调用错误:%s", err.Error()))
	}
	if url == "" {
		return finalResult, url, errors.New(fmt.Sprintf("IOTA接口调用错误:%s", err.Error()))
	}
	return finalResult, url, nil
}

// GetTransactionByHash 通过账本Hash验证数据返回DOID和DataHash
//func (iota *IOTA) GetTransactionByHash(transactionHash string, csDataTime int64) (*RespIOTTransaction, error) {
//	body := make(map[string]any)
//	body["action"] = "executeContract"
//	body["contractID"] = iota.Contract
//	body["operation"] = "getTransactionByHash"
//	arg := make(map[string]any)
//	arg["transactionHash"] = transactionHash
//	body["arg"] = arg
//	jsonBodyByte, err := common.DumpJson(body)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := httpDo(iota.Url, "POST", nil, string(jsonBodyByte), iota.Timeout)
//	if err != nil {
//		return nil, err
//	}
//	respStatus := gjson.Get(resp, "status").String()
//	respResult := gjson.Get(resp, "result")
//	if respStatus != "Success" {
//		return nil, errors.New(fmt.Sprintf("IOTA接口调用错误:%s", respResult.String()))
//	}
//	respBody := new(RespIOTTransaction)
//	respBody.TransactionHash = transactionHash
//	respBody.DOID = respResult.Get("data").Get("doId").String()
//	respBody.DataHash = respResult.Get("data").Get("dataHash").String()
//	respBody.CsDataTime = csDataTime
//	respBody.By = "GetTransactionByHash"
//	return respBody, nil
//}
