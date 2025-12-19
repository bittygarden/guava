package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 高膳食纤维食物数据（从搜索结果获取）
var highFiberFoods = map[string]float64{
	"奇亚籽":   34.4,
	"火麻仁":   27.6,
	"干香菇":   26.5,
	"大杏仁":   18.5,
	"鹰嘴豆":   17.4,
	"黑豆":    15.5,
	"青豆(干)": 12.6,
	"松子":    12.4,
	"腰果":    10.4,
	"核桃":    9.5,
	"黄花菜":   7.7,
	"红豆":    7.7,
	"库尔勒香梨": 6.7,
	"荞麦":    6.5,
	"藜麦":    6.5,
	"燕麦":    10.0,
	"绿豆":    6.4,
	"糙米":    3.5,
}

// 从CSV行创建NutritionData实例，保留所有字段
func createNutritionData(row []string, headers map[int]string) map[string]string {
	data := make(map[string]string)
	for i, value := range row {
		// 移除双引号
		value = strings.Trim(value, `"`)
		// 处理嵌套的双引号
		value = strings.ReplaceAll(value, `""`, `"`)
		header := strings.Trim(strings.ToLower(row[i]), `"`)
		data[header] = value
	}
	return data
}

// 从CSV文件读取营养数据
func readNutritionData(filePath string) ([]map[string]string, map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(lines) < 2 {
		return []map[string]string{}, nil, nil
	}

	// 解析表头
	headers := make(map[string]string)
	for _, header := range lines[0] {
		cleanHeader := strings.Trim(strings.ToLower(header), `"`)
		headers[cleanHeader] = header
	}

	// 读取数据行
	data := make([]map[string]string, 0)
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		rowData := make(map[string]string)
		for i, value := range line {
			// 移除双引号
			value = strings.Trim(value, `"`)
			// 处理嵌套的双引号
			value = strings.ReplaceAll(value, `""`, `"`)
			header := strings.Trim(strings.ToLower(lines[0][i]), `"`)
			rowData[header] = value
		}
		if rowData["food_name"] != "" {
			data = append(data, rowData)
		}
	}

	return data, headers, nil
}

// 将营养数据写入CSV文件
func writeNutritionData(data []map[string]string, headers map[string]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	originalHeaders := []string{}
	for _, header := range headers {
		originalHeaders = append(originalHeaders, header)
	}
	writer.Write(originalHeaders)

	// 写入数据行
	for _, row := range data {
		record := []string{}
		for _, header := range originalHeaders {
			cleanHeader := strings.Trim(strings.ToLower(header), `"`)
			value := row[cleanHeader]
			// 添加双引号
			record = append(record, fmt.Sprintf(`"%s"`, value))
		}
		writer.Write(record)
	}

	return nil
}

// 提取数值
func extractNumericValue(value string) float64 {
	if value == "" || value == "—" || value == "—" {
		return 0
	}
	// 移除单位
	value = strings.TrimRight(value, "g")
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return f
}

func main() {
	// 定义文件路径
	dataDir := "../data"
	inputFile := filepath.Join(dataDir, "j_food_nutrition_updated.csv")
	outputFile := filepath.Join(dataDir, "j_food_nutrition_updated_fiber.csv")

	fmt.Println("开始更新现有食物的膳食纤维含量...")

	// 读取营养数据
	nutritionData, headers, err := readNutritionData(inputFile)
	if err != nil {
		fmt.Printf("读取营养数据文件失败: %v\n", err)
		return
	}

	fmt.Printf("营养数据记录数: %d\n", len(nutritionData))

	// 统计更新的记录数
	updatedCount := 0

	// 更新现有食物的膳食纤维含量
	for i, data := range nutritionData {
		foodName := data["food_name"]
		if fiberValue, exists := highFiberFoods[foodName]; exists {
			// 获取当前膳食纤维含量
			currentFiberValue := extractNumericValue(data["dietary_fiber"])
			if fiberValue > currentFiberValue {
				// 更新膳食纤维含量
				nutritionData[i]["dietary_fiber"] = fmt.Sprintf("%.1fg", fiberValue)
				updatedCount++
				fmt.Printf("更新 %s 的膳食纤维含量: %.1fg -> %.1fg\n", foodName, currentFiberValue, fiberValue)
			}
		}
	}

	fmt.Printf("共更新了 %d 条记录\n", updatedCount)

	// 写入更新后的数据
	err = writeNutritionData(nutritionData, headers, outputFile)
	if err != nil {
		fmt.Printf("写入更新后的数据失败: %v\n", err)
		return
	}

	fmt.Printf("更新完成！结果保存到: %s\n", outputFile)
}
