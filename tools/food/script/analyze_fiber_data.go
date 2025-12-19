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
	
	fmt.Printf("开始分析 %s 中的膳食纤维含量数据...\n\n", filePath)
	
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
	
	// 统计膳食纤维数据情况
	var totalFoods, foodsWithFiber, foodsWithoutFiber int
	var foodsWithFiberList, foodsWithoutFiberList []string
	
	for _, line := range lines[1:] {
		if len(line) <= foodNameIndex {
			continue
		}
		
		totalFoods++
		foodName := strings.Trim(line[foodNameIndex], `"`)
		
		if len(line) <= dietaryFiberIndex {
			foodsWithoutFiber++
			foodsWithoutFiberList = append(foodsWithoutFiberList, foodName)
			continue
		}
		
		fiberContent := strings.Trim(line[dietaryFiberIndex], `"`)
		if fiberContent != "" && fiberContent != "—" {
			foodsWithFiber++
			foodsWithFiberList = append(foodsWithFiberList, fmt.Sprintf("%s: %s", foodName, fiberContent))
		} else {
			foodsWithoutFiber++
			foodsWithoutFiberList = append(foodsWithoutFiberList, foodName)
		}
	}
	
	// 输出统计结果
	fmt.Printf("统计结果:\n")
	fmt.Printf("总食物数量: %d\n", totalFoods)
	fmt.Printf("已有膳食纤维数据的食物数量: %d\n", foodsWithFiber)
	fmt.Printf("缺少膳食纤维数据的食物数量: %d\n\n", foodsWithoutFiber)
	
	fmt.Printf("已有膳食纤维数据的食物 (前20个):\n")
	for i, food := range foodsWithFiberList {
		if i >= 20 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, food)
	}
	
	if len(foodsWithFiberList) > 20 {
		fmt.Printf("... 还有 %d 个食物未显示\n\n", len(foodsWithFiberList)-20)
	} else {
		fmt.Printf("\n")
	}
	
	fmt.Printf("缺少膳食纤维数据的食物 (前20个):\n")
	for i, food := range foodsWithoutFiberList {
		if i >= 20 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, food)
	}
	
	if len(foodsWithoutFiberList) > 20 {
		fmt.Printf("... 还有 %d 个食物未显示\n\n", len(foodsWithoutFiberList)-20)
	} else {
		fmt.Printf("\n")
	}
	
	fmt.Printf("分析完成！\n")
}