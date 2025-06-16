package textsplitter

type TextSplitter interface {
	SplitText(text string) []string
}
