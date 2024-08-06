package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"holidays/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var data = make(map[string]models.YearData)

// loadJSONData 加载指定年份的 JSON 数据
func loadJSONData(year string) error {
	// 构建 JSON 文件名和路径
	filename := fmt.Sprintf("%s.json", year)
	filePath := filepath.Join("data", filename)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解析 JSON 文件
	var yearData models.YearData
	if err := json.NewDecoder(file).Decode(&yearData); err != nil {
		return err
	}

	// 存储解析后的数据
	data[year] = yearData
	return nil
}

// GetHolidaysByYear 获取指定年份的所有节假日
func GetHolidaysByYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]

	// 如果数据尚未加载，则加载数据
	if _, ok := data[year]; !ok {
		if err := loadJSONData(year); err != nil {
			http.Error(w, "未找到数据", http.StatusNotFound)
			return
		}
	}

	// 返回节假日数据
	json.NewEncoder(w).Encode(data[year].Holidays)
}

// GetHolidayByNameAndYear 获取指定年份中某个节日的信息
func GetHolidayByNameAndYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]
	festival := vars["festival"]

	// 如果数据尚未加载，则加载数据
	if _, ok := data[year]; !ok {
		if err := loadJSONData(year); err != nil {
			http.Error(w, "未找到数据", http.StatusNotFound)
			return
		}
	}

	// 查找并返回指定节日的信息
	for _, holiday := range data[year].Holidays {
		if strings.EqualFold(holiday.Name, festival) {
			json.NewEncoder(w).Encode(holiday)
			return
		}
	}

	// 如果未找到节日信息，返回 404 错误
	http.Error(w, "未找到节日", http.StatusNotFound)
}

// GetWorkdaysByYear 获取指定年份的所有工作日
func GetWorkdaysByYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]

	// 如果数据尚未加载，则加载数据
	if _, ok := data[year]; !ok {
		if err := loadJSONData(year); err != nil {
			http.Error(w, "未找到数据", http.StatusNotFound)
			return
		}
	}

	// 计算工作日并返回
	workdays, err := calculateWorkdays(year)
	if err != nil {
		http.Error(w, "计算工作日时出错", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(workdays)
}

// calculateWorkdays 计算指定年份的所有工作日
func calculateWorkdays(year string) ([]string, error) {
	// 解析年份的起始和结束日期
	startDate, err := time.Parse("2006-01-02", year+"-01-01")
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", year+"-12-31")
	if err != nil {
		return nil, err
	}

	// 创建一个映射来存储节假日
	holidayMap := make(map[string]struct{})
	for _, holiday := range data[year].Holidays {
		start, _ := time.Parse("2006-01-02", holiday.Start)
		end, _ := time.Parse("2006-01-02", holiday.End)
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			holidayMap[d.Format("2006-01-02")] = struct{}{}
		}
	}

	// 创建一个映射来存储调休工作日
	workdayMap := make(map[string]struct{})
	for _, workday := range data[year].Work {
		workdayMap[workday.Date] = struct{}{}
	}

	// 初始化工作日数组
	workdays := []string{}
	// 遍历年份中的每一天
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		day := d.Weekday()
		// 检查是否为周末
		if day == time.Saturday || day == time.Sunday {
			// 如果是周末但不在调休工作日映射中，则跳过
			if _, ok := workdayMap[d.Format("2006-01-02")]; !ok {
				continue
			}
		} else {
			// 如果是工作日但在节假日映射中，且不在调休工作日映射中，则跳过
			if _, ok := holidayMap[d.Format("2006-01-02")]; ok {
				if _, ok := workdayMap[d.Format("2006-01-02")]; !ok {
					continue
				}
			}
		}
		// 将有效工作日添加到工作日数组中
		workdays = append(workdays, d.Format("2006-01-02"))
	}

	return workdays, nil
}
