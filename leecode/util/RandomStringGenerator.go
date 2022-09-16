package util

import (
	"math/rand"
	"time"
)

type RandomStringGenerator struct {
	//生成随机字符串长度
	Length int
	//是否包含数字
	ContainNum bool
	//大小写
	LetterCase   byte
	LetterSource string
}

const (
	SmallCaseLetters   = "abcdefghijklmnopqrstuvwxyz"
	CapitalCaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers            = "0123456789"

	CaseSmall   = 0
	CaseCapital = 1
	CaseAll     = 2
)

func DefaultGenerator() *RandomStringGenerator {
	return &RandomStringGenerator{Length: 20, ContainNum: false, LetterCase: CaseSmall}
}

func NewGenerator(length int, containNum bool, letterCase byte) *RandomStringGenerator {
	return &RandomStringGenerator{Length: length, ContainNum: containNum, LetterCase: letterCase}
}

func (g *RandomStringGenerator) RandomString() string {
	if g.LetterSource == "" {
		switch g.LetterCase {
		case CaseSmall:
			g.LetterSource = SmallCaseLetters
		case CaseCapital:
			g.LetterSource = CapitalCaseLetters
		case CaseAll:
			g.LetterSource = SmallCaseLetters + CapitalCaseLetters
		default:
			g.LetterSource = SmallCaseLetters
		}

		if g.ContainNum {
			g.LetterSource += Numbers
		}
	}

	charArr := []byte(g.LetterSource)
	b := make([]byte, g.Length)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = charArr[rand.Intn(len(charArr))]
	}

	return string(b)
}
