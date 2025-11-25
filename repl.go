package main

import (
	"strings"
)

func cleanInput(text string) []string {
    // Fields сама разобьет по пробелам и уберет дублирующиеся пробелы
    fields := strings.Fields(text)
    
    // Если слов много, лучше сразу выделить память (оптимизация)
    output := make([]string, 0, len(fields))
    
    for _, word := range fields {
        // Просто лоуеркейсим и добавляем
        output = append(output, strings.ToLower(word))
    }
    
    return output
}
