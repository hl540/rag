package textsplitter

import (
	"regexp"
	"strings"
)

type SentenceSplitter struct {
	ChunkSize    int  // 每段最大长度（字符）
	ChunkOverlap int  // 每段之间的重叠部分长度
	RespectLine  bool // 是否优先按换行分段
}

func NewSentenceSplitter(size, overlap int, line bool) TextSplitter {
	return &SentenceSplitter{
		ChunkSize:    size,
		ChunkOverlap: overlap,
		RespectLine:  line,
	}
}

func (s *SentenceSplitter) SplitText(text string) []string {
	var sentences []string
	if s.RespectLine {
		sentences = strings.Split(text, "\n")
	} else {
		// 使用正则表达式按中文/英文句号、感叹号、问号分句
		re := regexp.MustCompile(`(?m)([^。！？\.\!\?]*[。！？\.\!\?])`)
		sentences = re.FindAllString(text, -1)
	}

	var chunks []string
	var currentRunes []rune

	for _, sent := range sentences {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}
		sentRunes := []rune(sent)

		if len(currentRunes)+len(sentRunes) > s.ChunkSize {
			chunks = append(chunks, string(currentRunes))

			if s.ChunkOverlap > 0 && len(currentRunes) > s.ChunkOverlap {
				// 保留重叠部分
				currentRunes = currentRunes[len(currentRunes)-s.ChunkOverlap:]
			} else {
				currentRunes = nil
			}
		}
		currentRunes = append(currentRunes, sentRunes...)
	}

	if len(currentRunes) > 0 {
		chunks = append(chunks, string(currentRunes))
	}

	return chunks
}
