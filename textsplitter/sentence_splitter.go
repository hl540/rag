package textsplitter

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// SentenceSplitter 是一个基于句子的文本分割器
type SentenceSplitter struct {
	ChunkSize    int  // 每段最大长度（字符）
	ChunkOverlap int  // 每段之间的重叠部分长度
	RespectLine  bool // 是否优先按换行分段
	MinChunkSize int  // 最小块大小，避免过小的块
}

// NewSentenceSplitter 创建一个新的句子分割器
func NewSentenceSplitter(size, overlap int, line bool) (TextSplitter, error) {
	if size <= 0 {
		return nil, errors.New("chunk size must be positive")
	}
	if overlap < 0 {
		return nil, errors.New("chunk overlap must be non-negative")
	}
	if overlap >= size {
		return nil, errors.New("chunk overlap must be less than chunk size")
	}

	return &SentenceSplitter{
		ChunkSize:    size,
		ChunkOverlap: overlap,
		RespectLine:  line,
		MinChunkSize: size / 4, // 默认最小块大小为最大块大小的 1/4
	}, nil
}

// SplitText 将文本分割成多个块
func (s *SentenceSplitter) SplitText(text string) []string {
	if text == "" {
		return nil
	}

	// 预处理：统一换行符，去除多余空白
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.TrimSpace(text)

	var sentences []string
	if s.RespectLine {
		// 按行分割，但保持段落的完整性
		paragraphs := strings.Split(text, "\n\n")
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para == "" {
				continue
			}
			// 如果段落长度超过块大小，需要进一步分割
			if len([]rune(para)) > s.ChunkSize {
				// 使用更细粒度的分割
				subSentences := s.splitIntoSentences(para)
				sentences = append(sentences, subSentences...)
			} else {
				sentences = append(sentences, para)
			}
		}
	} else {
		sentences = s.splitIntoSentences(text)
	}

	var chunks []string
	var currentChunk []rune
	currentLength := 0

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}

		sentRunes := []rune(sent)
		sentLength := len(sentRunes)

		// 如果当前句子加上已有内容超过块大小
		if currentLength+sentLength > s.ChunkSize {
			// 如果当前块不为空，保存它
			if currentLength > 0 {
				chunk := string(currentChunk)
				if len([]rune(chunk)) >= s.MinChunkSize {
					chunks = append(chunks, chunk)
				}

				// 处理重叠
				if s.ChunkOverlap > 0 {
					// 从当前块的末尾开始，向前找到最后一个完整句子
					overlapStart := len(currentChunk) - s.ChunkOverlap
					if overlapStart < 0 {
						overlapStart = 0
					}
					// 找到最后一个句子的开始位置
					for i := overlapStart; i < len(currentChunk); i++ {
						if i == 0 || isSentenceEnd(currentChunk[i-1]) {
							overlapStart = i
							break
						}
					}
					currentChunk = currentChunk[overlapStart:]
					currentLength = len(currentChunk)
				} else {
					currentChunk = nil
					currentLength = 0
				}
			}

			// 如果单个句子超过块大小，需要进一步分割
			if sentLength > s.ChunkSize {
				// 将长句子分割成更小的块，确保在字符边界处分割
				subChunks := s.splitLongSentence(sentRunes)
				chunks = append(chunks, subChunks...)
				continue
			}
		}

		// 添加当前句子到块中
		if currentLength > 0 {
			currentChunk = append(currentChunk, []rune(" ")...)
			currentLength++
		}
		currentChunk = append(currentChunk, sentRunes...)
		currentLength += sentLength
	}

	// 处理最后一个块
	if currentLength > 0 {
		chunk := string(currentChunk)
		if len([]rune(chunk)) >= s.MinChunkSize {
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}

// splitLongSentence 安全地分割长句子，确保在字符边界处分割
func (s *SentenceSplitter) splitLongSentence(sentRunes []rune) []string {
	var chunks []string
	sentLength := len(sentRunes)

	// 计算步长，考虑重叠
	step := s.ChunkSize - s.ChunkOverlap
	if step <= 0 {
		step = s.ChunkSize / 2 // 如果重叠太大，使用一半大小作为步长
	}

	for i := 0; i < sentLength; i += step {
		end := i + s.ChunkSize
		if end > sentLength {
			end = sentLength
		}

		// 如果这不是最后一块，尝试在句子边界处分割
		if end < sentLength {
			// 向前查找最近的句子结束符
			for j := end; j > i; j-- {
				if isSentenceEnd(sentRunes[j-1]) {
					end = j
					break
				}
			}

			// 如果没找到句子结束符，尝试在标点符号处分割
			if end == i+s.ChunkSize {
				for j := end; j > i; j-- {
					if isPunctuation(sentRunes[j-1]) {
						end = j
						break
					}
				}
			}

			// 如果还是没找到合适的分割点，确保至少有一些内容
			if end <= i {
				end = i + s.ChunkSize
				if end > sentLength {
					end = sentLength
				}
			}
		}

		chunk := sentRunes[i:end]
		if len(chunk) >= s.MinChunkSize {
			chunks = append(chunks, string(chunk))
		}
	}

	return chunks
}

// splitIntoSentences 将文本分割成句子
func (s *SentenceSplitter) splitIntoSentences(text string) []string {
	// 使用更完善的中文分句正则表达式
	re := regexp.MustCompile(`([^。！？\.\!\?]+[。！？\.\!\?])`)
	sentences := re.FindAllString(text, -1)

	// 处理没有标点符号的长句
	var result []string
	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}
		// 如果句子太长，在适当的位置分割
		if len([]rune(sent)) > s.ChunkSize {
			// 尝试在标点符号处分割
			parts := regexp.MustCompile(`[，；：、]`).Split(sent, -1)
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					result = append(result, part)
				}
			}
		} else {
			result = append(result, sent)
		}
	}
	return result
}

// isSentenceEnd 判断是否是句子结束符
func isSentenceEnd(r rune) bool {
	return r == '。' || r == '！' || r == '？' || r == '.' || r == '!' || r == '?'
}

// isPunctuation 判断是否是标点符号
func isPunctuation(r rune) bool {
	return unicode.IsPunct(r) || r == '，' || r == '；' || r == '：' || r == '、'
}
