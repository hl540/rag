package textsplitter

import (
	"errors"
	"strings"
	"unicode"
)

// CharacterTextSplitter 是一个基于字符数量的简单文本分割器
type CharacterTextSplitter struct {
	ChunkSize    int    // 每段最大长度（字符）
	ChunkOverlap int    // 每段之间的重叠部分长度
	Separator    string // 分隔符，用于在合适的位置分割
}

// NewCharacterTextSplitter 创建一个新的字符分割器
func NewCharacterTextSplitter(size, overlap int, separator string) (TextSplitter, error) {
	if size <= 0 {
		return nil, errors.New("chunk size must be positive")
	}
	if overlap < 0 {
		return nil, errors.New("chunk overlap must be non-negative")
	}
	if overlap >= size {
		return nil, errors.New("chunk overlap must be less than chunk size")
	}

	return &CharacterTextSplitter{
		ChunkSize:    size,
		ChunkOverlap: overlap,
		Separator:    separator,
	}, nil
}

// SplitText 将文本分割成多个块
func (c *CharacterTextSplitter) SplitText(text string) []string {
	if text == "" {
		return nil
	}

	// 预处理：去除多余空白
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	var chunks []string
	textRunes := []rune(text)
	textLength := len(textRunes)

	// 如果文本长度小于等于块大小，直接返回
	if textLength <= c.ChunkSize {
		return []string{text}
	}

	// 计算步长
	step := c.ChunkSize - c.ChunkOverlap
	if step <= 0 {
		step = c.ChunkSize / 2 // 如果重叠太大，使用一半大小作为步长
	}

	// 分割文本
	for i := 0; i < textLength; i += step {
		end := i + c.ChunkSize
		if end > textLength {
			end = textLength
		}

		// 如果这不是最后一块，尝试在分隔符处分割
		if end < textLength && c.Separator != "" {
			// 在指定范围内查找最后一个分隔符
			lastSeparator := -1
			for j := end; j > i; j-- {
				if j >= len(textRunes) {
					continue
				}
				// 检查是否是分隔符
				if c.isSeparator(textRunes[j-1]) {
					lastSeparator = j
					break
				}
			}

			// 如果找到了分隔符，在分隔符处分割
			if lastSeparator > i {
				end = lastSeparator
			}
		}

		chunk := textRunes[i:end]
		if len(chunk) > 0 {
			chunks = append(chunks, string(chunk))
		}
	}

	return chunks
}

// isSeparator 判断是否是分隔符
func (c *CharacterTextSplitter) isSeparator(r rune) bool {
	if c.Separator == "" {
		return false
	}

	// 检查是否是空白字符
	if unicode.IsSpace(r) {
		return true
	}

	// 检查是否是标点符号
	if unicode.IsPunct(r) {
		return true
	}

	// 检查是否是换行符
	if r == '\n' || r == '\r' {
		return true
	}

	return false
}
