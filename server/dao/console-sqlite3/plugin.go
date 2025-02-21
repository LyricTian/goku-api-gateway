package console_sqlite3

import (
	SQL "database/sql"
	"errors"
	"fmt"
	"strings"

	log "github.com/eolinker/goku-api-gateway/goku-log"

	database2 "github.com/eolinker/goku-api-gateway/common/database"
	entity "github.com/eolinker/goku-api-gateway/server/entity/console-entity"
)

//GetPluginInfo 获取插件配置信息
func GetPluginInfo(pluginName string) (bool, *entity.Plugin, error) {
	db := database2.GetConnection()
	sql := `SELECT pluginID,pluginName,pluginStatus,IFNULL(pluginConfig,""),pluginPriority,isStop,IFNULL(pluginDesc,""),IFNULL(version,""),pluginType FROM goku_plugin WHERE pluginName = ?;`
	plugin := &entity.Plugin{}
	err := db.QueryRow(sql, pluginName).Scan(&plugin.PluginID, &plugin.PluginName, &plugin.PluginStatus, &plugin.PluginConfig, &plugin.PluginIndex, &plugin.IsStop, &plugin.PluginDesc, &plugin.Version, &plugin.PluginType)
	if err != nil {
		return false, &entity.Plugin{}, err
	}
	return true, plugin, nil
}

// GetPluginList 获取插件列表
func GetPluginList(keyword string, condition int) (bool, []*entity.Plugin, error) {
	db := database2.GetConnection()
	rule := make([]string, 0, 2)

	if keyword != "" {
		searchRule := "pluginName LIKE '%" + keyword + "%' OR pluginDesc LIKE '%" + keyword + "%'"
		rule = append(rule, searchRule)
	}
	if condition > 0 {
		rule = append(rule, fmt.Sprintf("pluginType = %d", condition-1))
	}

	ruleStr := ""
	if len(rule) > 0 {
		ruleStr += "WHERE " + strings.Join(rule, " AND ")
	}
	sql := fmt.Sprintf(`SELECT pluginID,IFNULL(chineseName,""),pluginName,pluginStatus,pluginPriority,IFNULL(pluginDesc,""),isStop,pluginType,IFNULL(version,""),isCheck FROM goku_plugin %s ORDER BY pluginPriority DESC;`, ruleStr)
	rows, err := db.Query(sql)
	if err != nil {
		return false, make([]*entity.Plugin, 0), err
	}
	//延时关闭Rows
	defer rows.Close()
	//获取记录列
	pluginList := make([]*entity.Plugin, 0)

	for rows.Next() {
		var plugin entity.Plugin
		err = rows.Scan(&plugin.PluginID, &plugin.ChineseName, &plugin.PluginName, &plugin.PluginStatus, &plugin.PluginIndex, &plugin.PluginDesc, &plugin.IsStop, &plugin.PluginType, &plugin.Version, &plugin.IsCheck)
		if err != nil {
			return false, make([]*entity.Plugin, 0), err
		}
		pluginList = append(pluginList, &plugin)
	}
	// sort.Sort(sort.Reverse(entity.PluginSlice(pluginList)))
	return true, pluginList, nil
}

// GetPluginCount 获取插件数量
func GetPluginCount() int {
	var count int
	sql := "SELECT COUNT(*) FROM goku_plugin;"
	err := database2.GetConnection().QueryRow(sql).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

// AddPlugin 新增插件信息
func AddPlugin(pluginName, pluginConfig, pluginDesc, version string, pluginPriority, isStop, pluginType int) (bool, string, error) {
	db := database2.GetConnection()
	stmt, err := db.Prepare(`INSERT INTO goku_plugin (pluginName,pluginConfig,pluginDesc,version,pluginStatus,pluginPriority,isStop,official,pluginType) VALUES (?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return false, "[ERROR]Illegal SQL statement!", err
	}
	defer stmt.Close()
	_, err = stmt.Exec(pluginName, pluginConfig, pluginDesc, version, 0, pluginPriority, isStop, "false", pluginType)
	if err != nil {
		return false, "[ERROR]Failed to insert data!", err
	}
	return true, "", nil
}

// EditPlugin 修改插件信息
func EditPlugin(pluginName, pluginConfig, pluginDesc, version string, pluginPriority, isStop, pluginType int) (bool, string, error) {
	db := database2.GetConnection()
	// 查询插件是否是官方插件
	var sql string
	sql = "SELECT pluginType,official FROM goku_plugin WHERE pluginName = ?;"
	var official string
	var oldPluginType int
	err := db.QueryRow(sql, pluginName).Scan(&oldPluginType, &official)
	if err != nil {
		return false, "[ERROR]The plugin is not exist!", err
	}
	Tx, _ := db.Begin()
	paramsArray := make([]interface{}, 0)
	if official == "true" {
		sql = `UPDATE goku_plugin SET pluginConfig = ?,pluginDesc = ? WHERE pluginName = ?`
		paramsArray = append(paramsArray, pluginConfig, pluginDesc, pluginName)
	} else {
		sql = `UPDATE goku_plugin SET pluginPriority = ?,pluginConfig = ?,isStop = ?,pluginDesc = ?,version = ?,pluginType = ? WHERE pluginName = ?`
		paramsArray = append(paramsArray, pluginPriority, pluginConfig, isStop, pluginDesc, version, pluginType, pluginName)
	}
	_, err = Tx.Exec(sql, paramsArray...)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Failed to update data!", err
	}
	flag, err := EditPluginCache(pluginName, oldPluginType, pluginType, Tx)
	if !flag {
		Tx.Rollback()
		return false, "[ERROR]Failed to update plugin cache!", err
	}

	Tx.Commit()
	return true, "", nil
}

// DeletePlugin 删除插件信息
func DeletePlugin(pluginName string) (bool, string, error) {
	db := database2.GetConnection()
	var sql string
	sql = "SELECT pluginType,official FROM goku_plugin WHERE pluginName = ?;"
	var official string
	var pluginType int
	err := db.QueryRow(sql, pluginName).Scan(&pluginType, &official)
	if err != nil {
		return false, "[ERROR]The plugin is not exist!", err
	}
	if official == "true" {
		return false, "[ERROR]Can not delete goku plugin!", errors.New("[error]can not delete goku plugin")
	}
	Tx, _ := db.Begin()
	_, err = Tx.Exec(`DELETE FROM goku_plugin WHERE pluginName = ?`, pluginName)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Failed to delete data!", err
	}
	flag, err := EditPluginCache(pluginName, pluginType, 5, Tx)
	if !flag {
		Tx.Rollback()
		return false, "[ERROR]Failed to update plugin cache!", err
	}

	Tx.Commit()
	return true, "", nil
}

//CheckIndexIsExist 判断插件ID是否存在
func CheckIndexIsExist(pluginName string, pluginPriority int) (bool, error) {
	db := database2.GetConnection()
	sql := "SELECT pluginName FROM goku_plugin WHERE pluginPriority = ?;"
	var p string
	err := db.QueryRow(sql, pluginPriority).Scan(&p)
	if err != nil {
		return false, err
	}
	if pluginName == p {
		return false, err
	}
	return true, nil
}

//GetPluginConfig 获取插件配置及插件信息
func GetPluginConfig(pluginName string) (bool, string, error) {
	db := database2.GetConnection()
	sql := `SELECT IFNULL(pluginConfig,"") FROM goku_plugin WHERE pluginName = ?`
	var pluginConfig string
	err := db.QueryRow(sql, pluginName).Scan(&pluginConfig)
	if err != nil {
		return false, "[ERROR]The plugin is not exist!", err
	}
	return true, pluginConfig, nil
}

//CheckNameIsExist 检查插件名称是否存在
func CheckNameIsExist(pluginName string) (bool, error) {
	db := database2.GetConnection()
	sql := "SELECT pluginName FROM goku_plugin WHERE pluginName = ?;"
	var p string
	err := db.QueryRow(sql, pluginName).Scan(&p)
	if err != nil {
		return false, err
	}
	return true, err
}

//EditPluginStatus 修改插件开启状态
func EditPluginStatus(pluginName string, pluginStatus int) (bool, error) {
	db := database2.GetConnection()
	Tx, _ := db.Begin()
	isCheck := 1

	if pluginStatus == 0 && !strings.Contains(pluginName, "goku-") {
		isCheck = 0
	}
	sql := "UPDATE goku_plugin SET pluginStatus = ?,isCheck = ? WHERE pluginName = ?;"
	if pluginStatus == 1 {
		sql = "UPDATE goku_plugin SET pluginStatus = ?,isCheck = ? WHERE pluginName = ? AND isCheck = 1;"
	}
	_, err := Tx.Exec(sql, pluginStatus, isCheck, pluginName)
	if err != nil {
		Tx.Rollback()
		return false, err
	}
	// 获取使用该插件的策略组列表
	Tx.Commit()
	return true, nil
}

//GetPluginListByPluginType 获取不同类型的插件列表
func GetPluginListByPluginType(pluginType int) (bool, []map[string]interface{}, error) {
	db := database2.GetConnection()
	sql := `SELECT pluginID,pluginName,pluginDesc FROM goku_plugin WHERE pluginType = ? AND pluginStatus = 1;`
	rows, err := db.Query(sql, pluginType)
	if err != nil {
		log.Info(err.Error())
		return false, make([]map[string]interface{}, 0), err
	}
	//延时关闭Rows
	defer rows.Close()
	//获取记录列
	pluginList := make([]map[string]interface{}, 0)

	for rows.Next() {
		var pluginID int
		var pluginName, chineseName string
		err = rows.Scan(&pluginID, &pluginName, &chineseName)
		if err != nil {
			return false, make([]map[string]interface{}, 0), err
		}
		plugin := map[string]interface{}{
			"pluginID":    pluginID,
			"pluginName":  pluginName,
			"pluginType":  pluginType,
			"chineseName": chineseName,
		}
		pluginList = append(pluginList, plugin)
	}
	return true, pluginList, nil
}

//BatchStopPlugin 批量关闭插件
func BatchStopPlugin(pluginNameList string) (bool, string, error) {
	db := database2.GetConnection()
	Tx, _ := db.Begin()
	plugin := strings.Split(pluginNameList, ",")
	code := ""
	s := make([]interface{}, 0)
	for i := 0; i < len(plugin); i++ {
		code += "?"
		if i < len(plugin)-1 {
			code += ","
		}
		s = append(s, plugin[i])
	}
	sql := "UPDATE goku_plugin SET pluginStatus = 0,isCheck = (CASE WHEN (official = 'false') THEN 0 ELSE 1 END) WHERE pluginName IN (" + code + ");"
	_, err := Tx.Exec(sql, s...)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to excute SQL statement!", err
	}
	Tx.Commit()
	return true, "", nil
}

//BatchStartPlugin 批量关闭插件
func BatchStartPlugin(pluginNameList string) (bool, string, error) {
	db := database2.GetConnection()
	Tx, _ := db.Begin()
	plugin := strings.Split(pluginNameList, ",")
	code := ""
	s := make([]interface{}, 0)
	for i := 0; i < len(plugin); i++ {
		code += "?"
		if i < len(plugin)-1 {
			code += ","
		}
		s = append(s, plugin[i])
	}
	sql := "UPDATE goku_plugin SET pluginStatus = 1 WHERE pluginName IN (" + code + ") AND isCheck = 1;"
	_, err := Tx.Exec(sql, s...)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to excute SQL statement!", err
	}
	Tx.Commit()
	return true, "", nil
}

//EditPluginCache 将插件配置写进缓存表中
func EditPluginCache(pluginName string, oldPluginType, pluginType int, Tx *SQL.Tx) (bool, error) {
	if oldPluginType == 1 {
		// 获取策略ID列表
		sql := "SELECT strategyID FROM goku_conn_plugin_strategy WHERE pluginName = ?;"
		strategyIDList := make([]string, 0)
		rows, err := Tx.Query(sql, pluginName)
		if err != nil {
			return false, err
		}
		//延时关闭Rows
		defer rows.Close()
		//获取记录列

		for rows.Next() {
			var strategyID string
			err = rows.Scan(&strategyID)
			if err != nil {
				return false, err
			}
			strategyIDList = append(strategyIDList, strategyID)
		}
		if len(strategyIDList) > 0 {
			if oldPluginType != pluginType {
				// 删除策略组插件
				_, err = Tx.Exec("DELETE FROM goku_conn_plugin_strategy WHERE strategyID AND pluginName = ?;", pluginName)
				if err != nil {
					return false, err
				}
			}
		}
		return true, nil
	} else if oldPluginType == 2 {
		// 获取策略ID列表
		sql := "SELECT strategyID,apiID FROM goku_conn_plugin_api WHERE pluginName = ?;"
		connList := make([]map[string]interface{}, 0)
		rows, err := Tx.Query(sql, pluginName)
		if err != nil {
			return false, err
		}
		//延时关闭Rows
		defer rows.Close()
		//获取记录列

		for rows.Next() {
			var strategyID string
			var apiID int
			err = rows.Scan(&strategyID, &apiID)
			if err != nil {
				return false, err
			}
			connList = append(connList, map[string]interface{}{
				"strategyID": strategyID,
				"apiID":      apiID,
			})
		}
		if len(connList) > 0 {
			if oldPluginType != pluginType {
				// 删除接口插件
				_, err = Tx.Exec("DELETE FROM goku_conn_plugin_api WHERE pluginName = ?;", pluginName)
				if err != nil {
					return false, err
				}
			}
		}
		return true, nil
	}
	return true, nil
}

//EditPluginCheckStatus 更新插件检测状态
func EditPluginCheckStatus(pluginName string, isCheck int) (bool, string, error) {
	db := database2.GetConnection()
	sql := "UPDATE goku_plugin SET isCheck = ? WHERE pluginName = ?;"
	_, err := db.Exec(sql, isCheck, pluginName)
	if err != nil {
		return false, "[ERROR]Fail to update data", err
	}
	return true, "", nil
}
