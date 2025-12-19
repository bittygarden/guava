package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 定义文件路径
	inputFilePath := "../data/j_food_nutrition.csv"
	outputFilePath := "../data/existing_fiber_data.csv"
	
	fmt.Printf("开始提取 %s 中的现有膳食纤维数据...\n\n", inputFilePath)
	
	// 打开输入CSV文件
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("打开输入文件失败: %v\n", err)
		return
	}
	defer inputFile.Close()
	
	// 创建CSV读取器
	reader := csv.NewReader(inputFile)
	
	// 读取所有行
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("读取输入文件失败: %v\n", err)
		return
	}
	
	if len(lines) < 2 {
		fmt.Println("输入文件中没有食物数据\n")
		return
	}
	
	// 解析表头，查找food_name和dietary_fiber列
	headerRow := lines[0]
	foodNameIndex := -1
	dietaryFiberIndex := -1
	
	for i, header := range headerRow {
		cleanHeader := strings.Trim(strings.ToLower(header), `"`)
		if cleanHeader == "food_name" {
			foodNameIndex = i
		} else if cleanHeader == "dietary_fiber" {
			dietaryFiberIndex = i
		}
	}
	
	if foodNameIndex == -1 {
		fmt.Println("未找到food_name列\n")
		return
	}
	
	if dietaryFiberIndex == -1 {
		fmt.Println("未找到dietary_fiber列\n")
		return
	}
	
	// 创建输出文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("创建输出文件失败: %v\n", err)
		return
	}
	defer outputFile.Close()
	
	// 创建CSV写入器
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	
	// 写入输出文件表头
	if err := writer.Write([]string{"food_name", "dietary_fiber"}); err != nil {
		fmt.Printf("写入表头失败: %v\n", err)
		return
	}
	
	// 提取并写入现有膳食纤维数据
	var extractedCount int
	
	for _, line := range lines[1:] {
		if len(line) <= foodNameIndex || len(line) <= dietaryFiberIndex {
			continue
		}
		
		foodName := strings.Trim(line[foodNameIndex], `"`)
		fiberContent := strings.Trim(line[dietaryFiberIndex], `"`)
		
		if foodName != "" && fiberContent != "" && fiberContent != "—" {
			if err := writer.Write([]string{foodName, fiberContent}); err != nil {
				fmt.Printf("写入数据失败: %v\n", err)
				return
			}
			extractedCount++
		}
	}
	
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Printf("刷新写入器失败: %v\n", err)
		return
	}
	
	fmt.Printf("已成功提取 %d 条膳食纤维数据到 %s\n\n", extractedCount, outputFilePath)
	fmt.Printf("提取完成！\n")
}