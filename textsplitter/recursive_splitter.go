package textsplitter

import (
	"errors"
	"strings"
)

// RecursiveCharacterTextSplitter 是一个递归字符文本分割器
// 它按照优先级顺序尝试不同的分隔符来分割文本
type RecursiveCharacterTextSplitter struct {
	ChunkSize    int      // 每段最大长度（字符）
	ChunkOverlap int      // 每段之间的重叠部分长度
	Separators   []string // 分隔符列表，按优先级排序
}

// NewRecursiveCharacterTextSplitter 创建一个新的递归字符分割器
func NewRecursiveCharacterTextSplitter(size, overlap int, separators []string) (TextSplitter, error) {
	if size <= 0 {
		return nil, errors.New("chunk size must be positive")
	}
	if overlap < 0 {
		return nil, errors.New("chunk overlap must be non-negative")
	}
	if overlap >= size {
		return nil, errors.New("chunk overlap must be less than chunk size")
	}
	if len(separators) == 0 {
		return nil, errors.New("separators list cannot be empty")
	}

	return &RecursiveCharacterTextSplitter{
		ChunkSize:    size,
		ChunkOverlap: overlap,
		Separators:   separators,
	}, nil
}

// NewRecursiveCharacterTextSplitterWithDefaults 使用默认分隔符创建递归字符分割器
func NewRecursiveCharacterTextSplitterWithDefaults(size, overlap int) (TextSplitter, error) {
	// 默认分隔符，按优先级排序
	defaultSeparators := []string{
		"\n\n", // 段落分隔
		"\n",   // 换行
		"。",    // 中文句号
		"！",    // 中文感叹号
		"？",    // 中文问号
		".",    // 英文句号
		"!",    // 英文感叹号
		"?",    // 英文问号
		"；",    // 中文分号
		";",    // 英文分号
		"：",    // 中文冒号
		":",    // 英文冒号
		"，",    // 中文逗号
		",",    // 英文逗号
		" ",    // 空格
		"",     // 无分隔符（最后手段）
	}

	return NewRecursiveCharacterTextSplitter(size, overlap, defaultSeparators)
}

// SplitText 将文本分割成多个块
func (r *RecursiveCharacterTextSplitter) SplitText(text string) []string {
	if text == "" {
		return nil
	}

	// 预处理：去除多余空白
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// 如果文本长度小于等于块大小，直接返回
	if len([]rune(text)) <= r.ChunkSize {
		return []string{text}
	}

	// 递归分割
	return r.recursiveSplit(text)
}

// recursiveSplit 递归地分割文本
func (r *RecursiveCharacterTextSplitter) recursiveSplit(text string) []string {
	// 如果文本长度小于等于块大小，直接返回
	if len([]rune(text)) <= r.ChunkSize {
		return []string{text}
	}

	// 尝试使用每个分隔符
	for _, separator := range r.Separators {
		if separator == "" {
			// 如果没有分隔符，使用字符分割
			return r.characterSplit(text)
		}

		// 使用当前分隔符分割文本
		parts := strings.Split(text, separator)

		// 如果分割后只有一个部分，尝试下一个分隔符
		if len(parts) == 1 {
			continue
		}

		// 处理分割后的部分
		var chunks []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// 如果部分仍然太大，递归分割
			if len([]rune(part)) > r.ChunkSize {
				subChunks := r.recursiveSplit(part)
				chunks = append(chunks, subChunks...)
			} else {
				chunks = append(chunks, part)
			}
		}

		// 如果成功分割，返回结果
		if len(chunks) > 0 {
			return r.mergeChunks(chunks)
		}
	}

	// 如果所有分隔符都失败了，使用字符分割
	return r.characterSplit(text)
}

// characterSplit 使用字符数量分割文本
func (r *RecursiveCharacterTextSplitter) characterSplit(text string) []string {
	var chunks []string
	textRunes := []rune(text)
	textLength := len(textRunes)

	// 计算步长
	step := r.ChunkSize - r.ChunkOverlap
	if step <= 0 {
		step = r.ChunkSize / 2
	}

	// 分割文本
	for i := 0; i < textLength; i += step {
		end := i + r.ChunkSize
		if end > textLength {
			end = textLength
		}

		chunk := textRunes[i:end]
		if len(chunk) > 0 {
			chunks = append(chunks, string(chunk))
		}
	}

	return chunks
}

// mergeChunks 合并小块并处理重叠
func (r *RecursiveCharacterTextSplitter) mergeChunks(chunks []string) []string {
	if len(chunks) == 0 {
		return nil
	}

	var result []string
	var currentChunk []rune
	currentLength := 0

	for _, chunk := range chunks {
		chunkRunes := []rune(chunk)
		chunkLength := len(chunkRunes)

		// 如果当前块加上新块超过大小限制
		if currentLength+chunkLength > r.ChunkSize {
			// 保存当前块
			if currentLength > 0 {
				result = append(result, string(currentChunk))
			}

			// 处理重叠
			if r.ChunkOverlap > 0 && currentLength > r.ChunkOverlap {
				// 保留重叠部分
				overlapStart := currentLength - r.ChunkOverlap
				if overlapStart < 0 {
					overlapStart = 0
				}
				currentChunk = currentChunk[overlapStart:]
				currentLength = len(currentChunk)
			} else {
				currentChunk = nil
				currentLength = 0
			}
		}

		// 添加新块
		if currentLength > 0 {
			currentChunk = append(currentChunk, []rune(" ")...)
			currentLength++
		}
		currentChunk = append(currentChunk, chunkRunes...)
		currentLength += chunkLength
	}

	// 处理最后一个块
	if currentLength > 0 {
		result = append(result, string(currentChunk))
	}

	return result
}
