package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// 定义营养数据结构
type NutritionData struct {
	ID              string
	CateID          string
	FoodName        string
	AliasName       string
	EnglishName     string
	EdiblePart      string
	Water           string
	Energy          string
	Protein         string
	Fat             string
	Cholesterol     string
	Ash             string
	Carbohydrate    string
	DietaryFiber    string
	Carotene        string
	Retinol         string
	VitaminA        string
	VitaminE        string
	Thiamin         string
	Riboflavin      string
	Niacin          string
	VitaminC        string
	Calcium         string
	Phosphorus      string
	Potassium       string
	Sodium          string
	Magnesium       string
	Iron            string
	Zinc            string
	Selenium        string
	Copper          string
	Manganese       string
	Iodine          string
	SFA             string
	MUFA            string
	PUFA            string
	FattyAcidsTotal string
}

// 从CSV行创建NutritionData实例
func createNutritionData(row []string, headers map[int]string) NutritionData {
	data := NutritionData{}
	for i, value := range row {
		// 移除双引号
		value = strings.Trim(value, `"`)
		// 移除嵌套的引号
		value = strings.ReplaceAll(value, `""`, `"`)
		switch headers[i] {
		case "id":
			data.ID = value
		case "cate_id":
			data.CateID = value
		case "food_name":
			data.FoodName = value
		case "alias_name":
			data.AliasName = value
		case "english_name":
			data.EnglishName = value
		case "edible_part":
			data.EdiblePart = value
		case "water":
			data.Water = value
		case "energy":
			data.Energy = value
		case "protein":
			data.Protein = value
		case "fat":
			data.Fat = value
		case "cholesterol":
			data.Cholesterol = value
		case "ash":
			data.Ash = value
		case "carbohydrate":
			data.Carbohydrate = value
		case "dietary_fiber":
			data.DietaryFiber = value
		case "carotene":
			data.Carotene = value
		case "shihuangchun": // 处理nutrition_all.csv中的视黄醇字段名
			data.Retinol = value
		case "retinol":
			data.Retinol = value
		case "vitamin_a":
			data.VitaminA = value
		case "vitamin_e":
			data.VitaminE = value
		case "thiamin":
			data.Thiamin = value
		case "riboflavin":
			data.Riboflavin = value
		case "niacin":
			data.Niacin = value
		case "vitamin_c":
			data.VitaminC = value
		case "calcium":
			data.Calcium = value
		case "phosphorus":
			data.Phosphorus = value
		case "potassium":
			data.Potassium = value
		case "sodium":
			data.Sodium = value
		case "magnesium":
			data.Magnesium = value
		case "iron":
			data.Iron = value
		case "zinc":
			data.Zinc = value
		case "selenium":
			data.Selenium = value
		case "copper":
			data.Copper = value
		case "manganese":
			data.Manganese = value
		case "iodine":
			data.Iodine = value
		case "sfa":
			data.SFA = value
		case "mufa":
			data.MUFA = value
		case "pufa":
			data.PUFA = value
		case "fatty_acids_total":
			data.FattyAcidsTotal = value
		}
	}
	return data
}

// 提取数值
func extractNumericValue(value string) float64 {
	if value == "" || value == "—" || value == "—" {
		return 0
	}
	// 使用正则表达式提取数值部分
	re := regexp.MustCompile(`[\d\.]+`)
	matches := re.FindString(value)
	if matches == "" {
		return 0
	}
	f, err := strconv.ParseFloat(matches, 64)
	if err != nil {
		return 0
	}
	return f
}

// 比较两个营养值，返回较高的那个
func maxNutritionValue(val1, val2 string) string {
	if val1 == "" {
		return val2
	}
	if val2 == "" {
		return val1
	}
	// 处理特殊值
	if val1 == "—" && val2 != "—" {
		return val2
	}
	if val2 == "—" && val1 != "—" {
		return val1
	}
	if val1 == "—" && val2 == "—" {
		return val1
	}

	// 提取数值进行比较
	num1 := extractNumericValue(val1)
	num2 := extractNumericValue(val2)

	if num1 >= num2 {
		return val1
	}
	return val2
}

// 从CSV文件读取营养数据
func readNutritionData(filePath string) ([]NutritionData, map[string]NutritionData, error) {
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
		return []NutritionData{}, make(map[string]NutritionData), nil
	}

	// 解析表头
	headers := make(map[int]string)
	for i, header := range lines[0] {
		headers[i] = strings.Trim(strings.ToLower(header), `"`)
	}

	// 读取数据行
	data := make([]NutritionData, 0)
	foodMap := make(map[string]NutritionData)
	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		rowData := createNutritionData(line, headers)
		if rowData.FoodName != "" {
			data = append(data, rowData)
			foodMap[strings.TrimSpace(rowData.FoodName)] = rowData
		}
	}

	return data, foodMap, nil
}

// 将营养数据写入CSV文件
func writeNutritionData(data []NutritionData, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{
		"id", "cate_id", "food_name", "alias_name", "english_name", "edible_part",
		"water", "energy", "protein", "fat", "cholesterol", "ash",
		"carbohydrate", "dietary_fiber", "carotene", "retinol", "vitamin_a", "vitamin_e",
		"thiamin", "riboflavin", "niacin", "vitamin_c", "calcium",
		"phosphorus", "potassium", "sodium", "magnesium", "iron",
		"zinc", "selenium", "copper", "manganese", "iodine",
		"sfa", "mufa", "pufa", "fatty_acids_total",
	}
	writer.Write(headers)

	// 写入数据行
	for _, row := range data {
		record := []string{
			row.ID, row.CateID, row.FoodName, row.AliasName, row.EnglishName, row.EdiblePart,
			row.Water, row.Energy, row.Protein, row.Fat, row.Cholesterol, row.Ash,
			row.Carbohydrate, row.DietaryFiber, row.Carotene, row.Retinol, row.VitaminA, row.VitaminE,
			row.Thiamin, row.Riboflavin, row.Niacin, row.VitaminC, row.Calcium,
			row.Phosphorus, row.Potassium, row.Sodium, row.Magnesium, row.Iron,
			row.Zinc, row.Selenium, row.Copper, row.Manganese, row.Iodine,
			row.SFA, row.MUFA, row.PUFA, row.FattyAcidsTotal,
		}
		// 添加双引号
		for i, val := range record {
			record[i] = fmt.Sprintf(`"%s"`, val)
		}
		writer.Write(record)
	}

	return nil
}

func main() {
	// 定义文件路径
	dataDir := "../data"
	jFoodNutritionFile := filepath.Join(dataDir, "j_food_nutrition.csv")
	nutritionAllFile := filepath.Join(dataDir, "nutrition_all.csv")
	outputFile := filepath.Join(dataDir, "j_food_nutrition_updated.csv")

	fmt.Println("开始更新营养数据...")
	fmt.Printf("基础文件: %s\n", jFoodNutritionFile)
	fmt.Printf("更新源文件: %s\n", nutritionAllFile)

	// 读取j_food_nutrition.csv数据
	jFoodData, _, err := readNutritionData(jFoodNutritionFile)
	if err != nil {
		fmt.Printf("读取基础文件失败: %v\n", err)
		return
	}

	// 读取nutrition_all.csv数据
	_, nutritionAllMap, err := readNutritionData(nutritionAllFile)
	if err != nil {
		fmt.Printf("读取更新源文件失败: %v\n", err)
		return
	}

	fmt.Printf("基础文件记录数: %d\n", len(jFoodData))
	fmt.Printf("更新源文件记录数: %d\n", len(nutritionAllMap))

	// 统计更新的记录数
	updatedCount := 0

	// 更新j_food_nutrition.csv中的数据
	for i, data := range jFoodData {
		foodName := strings.TrimSpace(data.FoodName)
		if updateData, exists := nutritionAllMap[foodName]; exists {
			// 找到匹配的食物，更新每个字段，保留较大的值
			updatedData := NutritionData{
				ID:              data.ID,
				CateID:          data.CateID,
				FoodName:        data.FoodName,
				AliasName:       data.AliasName,
				EnglishName:     data.EnglishName,
				EdiblePart:      data.EdiblePart,
				Water:           maxNutritionValue(data.Water, updateData.Water),
				Energy:          maxNutritionValue(data.Energy, updateData.Energy),
				Protein:         maxNutritionValue(data.Protein, updateData.Protein),
				Fat:             maxNutritionValue(data.Fat, updateData.Fat),
				Cholesterol:     maxNutritionValue(data.Cholesterol, updateData.Cholesterol),
				Ash:             maxNutritionValue(data.Ash, updateData.Ash),
				Carbohydrate:    maxNutritionValue(data.Carbohydrate, updateData.Carbohydrate),
				DietaryFiber:    maxNutritionValue(data.DietaryFiber, updateData.DietaryFiber),
				Carotene:        maxNutritionValue(data.Carotene, updateData.Carotene),
				Retinol:         maxNutritionValue(data.Retinol, updateData.Retinol),
				VitaminA:        maxNutritionValue(data.VitaminA, updateData.VitaminA),
				VitaminE:        maxNutritionValue(data.VitaminE, updateData.VitaminE),
				Thiamin:         maxNutritionValue(data.Thiamin, updateData.Thiamin),
				Riboflavin:      maxNutritionValue(data.Riboflavin, updateData.Riboflavin),
				Niacin:          maxNutritionValue(data.Niacin, updateData.Niacin),
				VitaminC:        maxNutritionValue(data.VitaminC, updateData.VitaminC),
				Calcium:         maxNutritionValue(data.Calcium, updateData.Calcium),
				Phosphorus:      maxNutritionValue(data.Phosphorus, updateData.Phosphorus),
				Potassium:       maxNutritionValue(data.Potassium, updateData.Potassium),
				Sodium:          maxNutritionValue(data.Sodium, updateData.Sodium),
				Magnesium:       maxNutritionValue(data.Magnesium, updateData.Magnesium),
				Iron:            maxNutritionValue(data.Iron, updateData.Iron),
				Zinc:            maxNutritionValue(data.Zinc, updateData.Zinc),
				Selenium:        maxNutritionValue(data.Selenium, updateData.Selenium),
				Copper:          maxNutritionValue(data.Copper, updateData.Copper),
				Manganese:       maxNutritionValue(data.Manganese, updateData.Manganese),
				Iodine:          maxNutritionValue(data.Iodine, updateData.Iodine),
				SFA:             maxNutritionValue(data.SFA, updateData.SFA),
				MUFA:            maxNutritionValue(data.MUFA, updateData.MUFA),
				PUFA:            maxNutritionValue(data.PUFA, updateData.PUFA),
				FattyAcidsTotal: maxNutritionValue(data.FattyAcidsTotal, updateData.FattyAcidsTotal),
			}

			// 更新数据
			jFoodData[i] = updatedData
			updatedCount++
		}
		// 不匹配的食物保持原样
	}

	fmt.Printf("更新的记录数: %d\n", updatedCount)

	// 写入更新后的数据
	err = writeNutritionData(jFoodData, outputFile)
	if err != nil {
		fmt.Printf("写入更新后的数据失败: %v\n", err)
		return
	}

	fmt.Printf("更新完成！结果保存到: %s\n", outputFile)
}
