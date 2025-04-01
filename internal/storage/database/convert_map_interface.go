package database

import "fmt"

// convertMapInterfaceToStringInt преобразует карту с ключами и значениями 
// типа interface{} в карту с ключами типа string и значениями типа int32.
// Если тип ключа или значения не соответствует ожидаемому, возвращается ошибка.
func convertMapInterfaceToStringInt(input map[interface{}]interface{}) (map[string]int32, error) {
	result := make(map[string]int32)

	for key, value := range input {
		// Приведение ключа к строке
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected key type: %v", key)
		}

		// Приведение значения к int
		valueInt, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("unexpected value type: %v", value)
		}

		// Добавление в результирующую карту
		result[keyStr] = valueInt
	}

	return result, nil
}

// convertMapInterfaceToStringBool преобразует карту с ключами и значениями 
// типа interface{} в карту с ключами типа string и значениями типа bool.
// Если тип ключа или значения не соответствует ожидаемому, возвращается ошибка.
func convertMapInterfaceToStringBool(input map[interface{}]interface{}) (map[string]bool, error) {
	result := make(map[string]bool)

	for key, value := range input {
		// Приведение ключа к строке
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected key type: %v", key)
		}

		// Приведение значения к int
		valueBool, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("unexpected value type: %v", valueBool)
		}

		// Добавление в результирующую карту
		result[keyStr] = valueBool
	}

	return result, nil
}
