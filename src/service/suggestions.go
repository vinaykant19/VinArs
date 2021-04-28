package service

func GetSearchHelp() []interface{} {
	var searchCategories []interface{}
	searchCategories = append(searchCategories, map[string]string{"text": "blood", "image":"blood.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "pregnancy", "image":"pregnancy.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "sugar", "image":"sugar.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "kidny", "image":"kidny.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "dna", "image":"dna.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "x-ray", "image":"x-ray.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "urine", "image":"urine.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "c-t scan", "image":"c-t-scan.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "blood", "image":"blood.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "lung", "image":"lung.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "sonography", "image":"sonography.jpg"})
	searchCategories = append(searchCategories, map[string]string{"text": "stool", "image":"stool.jpg"})
	return searchCategories
}