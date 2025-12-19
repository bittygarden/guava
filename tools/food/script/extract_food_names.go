package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 定义文件路径
	filePath := "../data/j_food_nutrition.csv"
	
	fmt.Printf("开始提取 %s 中的所有食物名称...\n\n", filePath)
	
	// 打开CSV文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()
	
	// 创建CSV读取器
	reader := csv.NewReader(file)
	
	// 读取所有行
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}
	
	if len(lines) < 2 {
		fmt.Println("文件中没有食物数据\n")
		return
	}
	
	// 解析表头，查找food_name列
	headerRow := lines[0]
	foodNameIndex := -1
	for i, header := range headerRow {
		cleanHeader := strings.Trim(strings.ToLower(header), `"`)
		if cleanHeader == "food_name" {
			foodNameIndex = i
			break
		}
	}
	
	if foodNameIndex == -1 {
		fmt.Println("未找到food_name列\n")
		return
	}
	
	// 提取所有食物名称
	foodNames := make([]string, 0, len(lines)-1)
	for _, line := range lines[1:] {
		if len(line) > foodNameIndex {
			foodName := strings.Trim(line[foodNameIndex], `"`)
			if foodName != "" {
				foodNames = append(foodNames, foodName)
			}
		}
	}
	
	// 输出结果
	fmt.Printf("共提取到 %d 种食物名称:\n\n", len(foodNames))
	for i, foodName := range foodNames {
		fmt.Printf("%d. %s\n", i+1, foodName)
	}
	
	fmt.Printf("\n提取完成！\n")
}
