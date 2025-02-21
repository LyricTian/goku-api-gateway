package node

import (
	console_sqlite3 "github.com/eolinker/goku-api-gateway/server/dao/console-sqlite3"
	entity "github.com/eolinker/goku-api-gateway/server/entity/console-entity"
)

//AddNode 新增节点信息
func AddNode(clusterID int, nodeName, nodeIP, nodePort, gatewayPath string, groupID int) (bool, map[string]interface{}, error) {
	return console_sqlite3.AddNode(clusterID, nodeName, nodeIP, nodePort, gatewayPath, groupID)
}

//EditNode 修改节点
func EditNode(nodeName, nodeIP, nodePort, gatewayPath string, nodeID, groupID int) (bool, string, error) {

	return console_sqlite3.EditNode(nodeName, nodeIP, nodePort, gatewayPath, nodeID, groupID)
}

//DeleteNode 删除节点
func DeleteNode(nodeID int) (bool, string, error) {
	return console_sqlite3.DeleteNode(nodeID)
}

//GetNodeInfo 获取节点信息
func GetNodeInfo(nodeID int) (bool, *entity.Node, error) {
	b, node, e := console_sqlite3.GetNodeInfo(nodeID)
	ResetNodeStatus(node)
	return b, node, e
}

//GetNodeInfoByIPPort 获取节点信息
func GetNodeInfoByIPPort(ip string, port int) (bool, *entity.Node, error) {
	b, node, e := console_sqlite3.GetNodeByIPPort(ip, port)
	ResetNodeStatus(node)
	return b, node, e
}

// GetNodeList 获取节点列表
func GetNodeList(clusterID, groupID int, keyword string) (bool, []*entity.Node, error) {
	b, nodes, e := console_sqlite3.GetNodeList(clusterID, groupID, keyword)
	ResetNodeStatus(nodes...)
	return b, nodes, e
}

//CheckIsExistRemoteAddr 节点IP查重
func CheckIsExistRemoteAddr(nodeID int, nodeIP, nodePort string) bool {
	return console_sqlite3.CheckIsExistRemoteAddr(nodeID, nodeIP, nodePort)
}

//BatchDeleteNode 批量删除节点
func BatchDeleteNode(nodeIDList string) (bool, string, error) {
	flag, nodeIDList, err := console_sqlite3.GetAvaliableNodeListFromNodeList(nodeIDList, 0)
	if !flag {
		return false, err.Error(), err
	} else if nodeIDList == "" {
		return false, "230013", err
	}
	return console_sqlite3.BatchDeleteNode(nodeIDList)
}

//BatchEditNodeGroup 批量修改节点分组
func BatchEditNodeGroup(nodeIDList string, groupID int) (bool, string, error) {
	return console_sqlite3.BatchEditNodeGroup(nodeIDList, groupID)
}

//GetNodeIPList 获取节点IP列表
func GetNodeIPList() (bool, []map[string]interface{}, error) {
	return console_sqlite3.GetNodeIPList()
}
