package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 定义中文到英文的营养成分映射
var nutrientMapCNtoEN = map[string]string{
	"热量":      "energy",
	"蛋白质":     "protein",
	"脂肪":      "fat",
	"碳水化合物":   "carbohydrate",
	"膳食纤维":    "dietary_fiber",
	"维生素A":    "vitamin_a",
	"视黄醇当量":   "shihuangchun",
	"硫胺素":     "thiamin",
	"核黄素":     "riboflavin",
	"烟酸":      "niacin",
	"维生素C":    "vitamin_c",
	"维生素E":    "vitamin_e",
	"钙":       "calcium",
	"镁":       "magnesium",
	"铁":       "iron",
	"锰":       "manganese",
	"锌":       "zinc",
	"铜":       "copper",
	"磷":       "phosphorus",
	"钾":       "potassium",
	"钠":       "sodium",
	"硒":       "selenium",
	"胆固醇":     "cholesterol",
	"水分":      "water",
	"灰分":      "ash",
	"胡罗卜素":    "carotene",
	"碘":       "iodine",
	"饱和脂肪酸":   "sfa",
	"单不饱和脂肪酸": "mufa",
	"多不饱和脂肪酸": "pufa",
	"总脂肪酸":    "fatty_acids_total",
}

// NutritionData 定义营养数据结构
type NutritionData struct {
	ID              string `json:"id"`
	CateID          string `json:"cate_id"`
	FoodName        string `json:"food_name"`
	AliasName       string `json:"alias_name"`
	EnglishName     string `json:"english_name"`
	EdiblePart      string `json:"edible_part"`
	Water           string `json:"water"`
	Energy          string `json:"energy"`
	Protein         string `json:"protein"`
	Fat             string `json:"fat"`
	Cholesterol     string `json:"cholesterol"`
	Ash             string `json:"ash"`
	Carbohydrate    string `json:"carbohydrate"`
	DietaryFiber    string `json:"dietary_fiber"`
	Carotene        string `json:"carotene"`
	VitaminA        string `json:"vitamin_a"`
	Shihuangchun    string `json:"shihuangchun"`
	VitaminE        string `json:"vitamin_e"`
	Thiamin         string `json:"thiamin"`
	Riboflavin      string `json:"riboflavin"`
	Niacin          string `json:"niacin"`
	VitaminC        string `json:"vitamin_c"`
	Calcium         string `json:"calcium"`
	Phosphorus      string `json:"phosphorus"`
	Potassium       string `json:"potassium"`
	Sodium          string `json:"sodium"`
	Magnesium       string `json:"magnesium"`
	Iron            string `json:"iron"`
	Zinc            string `json:"zinc"`
	Selenium        string `json:"selenium"`
	Copper          string `json:"copper"`
	Manganese       string `json:"manganese"`
	Iodine          string `json:"iodine"`
	SFA             string `json:"sfa"`
	MUFA            string `json:"mufa"`
	PUFA            string `json:"pufa"`
	FattyAcidsTotal string `json:"fatty_acids_total"`
}

// CrawlNutritionData 从指定URL爬取营养数据
func CrawlNutritionData(id string) (*NutritionData, error) {
	url := fmt.Sprintf("https://yingyang.supfree.net/wochuo.asp?id=%s", id)

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加必要的HTTP头部
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 添加Cookie
	req.AddCookie(&http.Cookie{Name: "ASPSESSIONIDQABSSAQD", Value: "HMEJDNDCOEIJJEFELKJDEDNO"})
	req.AddCookie(&http.Cookie{Name: "server_name_session", Value: "09ac9d2a18cf56bf011e3f1496d86d8e"})

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d %s", resp.StatusCode, resp.Status)
	}

	// 使用GBK编码读取响应体
	reader := bufio.NewReader(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder()))
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 初始化数据结构
	data := &NutritionData{
		ID: id,
	}

	// 提取食物名称：从面包屑导航中获取，如"当前位置：在线查询网 > 食物营养成分查询 > 小麦"中的"小麦"
	breadnav := doc.Find("#breadnav").Text()
	// 解析面包屑，提取最后一个元素
	parts := strings.Split(breadnav, " > ")
	if len(parts) > 1 {
		// 获取最后一个部分并去除空格和可能的HTML标签
		foodName := strings.TrimSpace(parts[len(parts)-1])
		// 去除可能的HTML标签
		foodName = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(foodName, "")
		data.FoodName = foodName
	} else {
		// 如果面包屑解析失败，尝试从h1标签获取
		data.FoodName = strings.TrimSpace(doc.Find("h1").Text())
	}
	fmt.Printf("食物名称: %s\n", data.FoodName)

	// 提取营养数据表格（根据参考HTML结构，营养数据在第一个表格中）
	table := doc.Find("table").First()
	if table.Length() == 0 {
		return nil, fmt.Errorf("未找到营养数据表格")
	}

	// 解析表格数据
	var nutrientMap = make(map[string]string)
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		// 页面返回的是3列布局的表格，每行有3组营养成分
		for j := 0; j < cells.Length(); j += 2 {
			// 每组包含营养成分名称和值（2个单元格）
			nutrient := strings.TrimSpace(cells.Eq(j).Text())
			value := strings.TrimSpace(cells.Eq(j + 1).Text())
			if nutrient != "" && value != "" {
				nutrientMap[nutrient] = value
				fmt.Printf("找到营养成分: %s = %s\n", nutrient, value)
			}
		}
	})

	// 创建一个临时映射来存储营养数据
	nutrientDataMap := make(map[string]string)

	// 使用中文到英文的映射表处理营养数据
	for cnNutrient, value := range nutrientMap {
		if enNutrient, exists := nutrientMapCNtoEN[cnNutrient]; exists {
			if enNutrient == "energy" {
				// 转换能量单位
				nutrientDataMap[enNutrient] = convertEnergy(value)
			} else {
				nutrientDataMap[enNutrient] = extractUnitValue(value)
			}
		}
	}

	// 将映射数据赋值给结构体
	data.Water = nutrientDataMap["water"]
	data.Energy = nutrientDataMap["energy"]
	data.Protein = nutrientDataMap["protein"]
	data.Fat = nutrientDataMap["fat"]
	data.Carbohydrate = nutrientDataMap["carbohydrate"]
	data.DietaryFiber = nutrientDataMap["dietary_fiber"]
	data.Shihuangchun = nutrientDataMap["shihuangchun"]
	data.VitaminA = nutrientDataMap["vitamin_a"]
	data.VitaminC = nutrientDataMap["vitamin_c"]
	data.VitaminE = nutrientDataMap["vitamin_e"]
	data.Thiamin = nutrientDataMap["thiamin"]
	data.Riboflavin = nutrientDataMap["riboflavin"]
	data.Niacin = nutrientDataMap["niacin"]
	data.Calcium = nutrientDataMap["calcium"]
	data.Magnesium = nutrientDataMap["magnesium"]
	data.Iron = nutrientDataMap["iron"]
	data.Manganese = nutrientDataMap["manganese"]
	data.Zinc = nutrientDataMap["zinc"]
	data.Copper = nutrientDataMap["copper"]
	data.Phosphorus = nutrientDataMap["phosphorus"]
	data.Potassium = nutrientDataMap["potassium"]
	data.Sodium = nutrientDataMap["sodium"]
	data.Selenium = nutrientDataMap["selenium"]
	data.Cholesterol = nutrientDataMap["cholesterol"]
	data.Ash = nutrientDataMap["ash"]
	data.Carotene = nutrientDataMap["carotene"]
	data.Iodine = nutrientDataMap["iodine"]
	data.SFA = nutrientDataMap["sfa"]
	data.MUFA = nutrientDataMap["mufa"]
	data.PUFA = nutrientDataMap["pufa"]
	data.FattyAcidsTotal = nutrientDataMap["fatty_acids_total"]

	return data, nil
}

// extractUnitValue 提取带单位的值，并将中文单位转换为英文单位
func extractUnitValue(value string) string {
	if value == "" {
		return ""
	}
	// 移除多余的空格
	value = strings.TrimSpace(value)

	// 将中文单位转换为英文单位
	value = strings.ReplaceAll(value, "(毫克)", "mg")
	value = strings.ReplaceAll(value, "(克)", "g")
	value = strings.ReplaceAll(value, "(微克)", "μg")
	value = strings.ReplaceAll(value, "毫克", "mg")
	value = strings.ReplaceAll(value, "克", "g")
	value = strings.ReplaceAll(value, "微克", "μg")
	value = strings.ReplaceAll(value, "微g", "μg")

	return value
}

// convertEnergy 将千卡转换为千焦
func convertEnergy(calories string) string {
	if calories == "" {
		return ""
	}
	// 提取数值
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(calories)
	if len(matches) < 2 {
		return calories
	}
	// 转换为千焦 (1千卡 = 4.184千焦)
	calorieVal := parseFloat(matches[1])
	kJ := calorieVal * 4.184
	return fmt.Sprintf("%.0fkJ", kJ)
}

// parseFloat 解析字符串为浮点数
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func main() {
	// 定义CSV文件名
	filename := "/Users/sunfeilong/workspace/github/guava/tools/food/data/nutrition_all.csv"

	// 先写入CSV表头
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{"id", "cate_id", "food_name", "alias_name", "english_name", "edible_part",
		"water", "energy", "protein", "fat", "cholesterol", "ash", "carbohydrate", "dietary_fiber",
		"carotene", "shihuangchun", "vitamin_a", "vitamin_e", "thiamin", "riboflavin", "niacin", "vitamin_c",
		"calcium", "phosphorus", "potassium", "sodium", "magnesium", "iron", "zinc", "selenium",
		"copper", "manganese", "iodine", "sfa", "mufa", "pufa", "fatty_acids_total"}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("写入表头失败: %v\n", err)
		return
	}

	// 循环抓取从1到1450的数据
	successCount := 0
	failCount := 0

	for i := 1; i <= 1450; i++ {
		id := fmt.Sprintf("%d", i)
		fmt.Printf("尝试爬取ID为%s的数据...\n", id)

		data, err := CrawlNutritionData(id)
		if err != nil {
			fmt.Printf("爬取ID为%s的数据失败: %v\n", id, err)
			failCount++
			continue
		}

		// 写入数据行
		row := []string{
			data.ID,
			data.CateID,
			data.FoodName,
			data.AliasName,
			data.EnglishName,
			data.EdiblePart,
			data.Water,
			data.Energy,
			data.Protein,
			data.Fat,
			data.Cholesterol,
			data.Ash,
			data.Carbohydrate,
			data.DietaryFiber,
			data.Carotene,
			data.Shihuangchun,
			data.VitaminA,
			data.VitaminE,
			data.Thiamin,
			data.Riboflavin,
			data.Niacin,
			data.VitaminC,
			data.Calcium,
			data.Phosphorus,
			data.Potassium,
			data.Sodium,
			data.Magnesium,
			data.Iron,
			data.Zinc,
			data.Selenium,
			data.Copper,
			data.Manganese,
			data.Iodine,
			data.SFA,
			data.MUFA,
			data.PUFA,
			data.FattyAcidsTotal,
		}

		if err := writer.Write(row); err != nil {
			fmt.Printf("写入ID为%s的数据失败: %v\n", id, err)
			failCount++
			continue
		}

		writer.Flush() // 立即刷新缓冲区
		successCount++
		fmt.Printf("成功爬取并写入ID为%s的数据\n", id)
		fmt.Printf("食物名称: %s\n", data.FoodName)
		fmt.Printf("进度: %d/1450, 成功: %d, 失败: %d\n", i, successCount, failCount)
		time.Sleep(time.Second * 1)
	}

	fmt.Printf("爬取完成！总共尝试了1450个ID，成功%d个，失败%d个\n", successCount, failCount)
	fmt.Printf("数据已写入文件: %s\n", filename)
}
