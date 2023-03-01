package gstrconv

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

var (
	unitCn2Arabic = map[string]int64{
		"十": 10,
		"百": 100,
		"千": 1000,
		"万": 10_000,
		"亿": 100_000_000,

		"十万": 100_000,
		"百万": 1_000_000,
		"千万": 10_000_000,

		"十亿": 1_000_000_000,
		"百亿": 10_000_000_000,
		"千亿": 100_000_000_000,
		"万亿": 1_000_000_000_000,
	}

	digitCn2Arabic = map[string]int64{
		"零": 0,
		"一": 1,
		"二": 2,
		"三": 3,
		"四": 4,
		"五": 5,
		"六": 6,
		"七": 7,
		"八": 8,
		"九": 9,
	}
)

func isDigitCn(c rune) bool {
	_, ok := digitCn2Arabic[string(c)]
	return ok
}

func isUnitCn(s string) bool {
	_, ok := unitCn2Arabic[s]
	return ok
}

// 类似 strconv.Atoi，但增加支持中文数字。
// 入参可以是中文数字或阿拉伯数字，如 `一百二十三`，`123`，但不支持混合，如 `一百零1` 则返回错误。
// 「兆」在不同上下文有歧义，因此只支持解释最大数字为 9_9999_9999_9999 / 9_999_999_999_999 / 九万九千九百九十九亿 九千九百九十九万 九千九百九十九 。
// https://en.wikipedia.org/wiki/Chinese_numerals
func AtoiCn(s string) (rs int64, err error) {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "_", "")
	str := s

	if str == "" {
		return 0, ErrInvalidData
	}

	stackValue := []int64{} // 数值 0-9
	stackUnit := []int64{}  // 单位 十百千万等 10 100 1000 10000
	part := []rune{}
	for len(str) > 0 {
		codePoint, size := utf8.DecodeRuneInString(str)
		if unicode.IsDigit(codePoint) {
			rs, _ = strconv.ParseInt(s, 10, 64)
			return rs, nil
		}

		if size != 3 {
			// 中文数字占3个字节，中文阿拉伯混合或其他异常情况时不是3
			return 0, ErrInvalidData
		}

		char := string(codePoint)
		digitCn, okDigitCn := digitCn2Arabic[char]
		unitCn, okUnitCn := unitCn2Arabic[char]

		partN := len(part)

		if okDigitCn && digitCn == 0 {
			// 解释 "一百零二" 中的 "零" 可跳过
			str = str[size:]
			continue
		}

		if okDigitCn {
			if partN > 0 {
				chars := string(part)

				if isUnitCn(chars) {
					// 对 零分割值特殊处理，如 "三百万零一百五十一"
					// stackValue=[301, 1] stackUnit=[10000, 10000]
					// 需要判断确保单位栈左到右越来越小
					stackUnitN := len(stackUnit)
					if stackUnitN > 0 {
						lastUnit := stackUnit[stackUnitN-1]
						thisUnit := unitCn2Arabic[chars]
						if thisUnit < lastUnit {
							stackUnit = append(stackUnit, thisUnit)
						}
					} else {
						stackUnit = append(stackUnit, unitCn2Arabic[chars])
					}
					part = []rune{}
				}
			}

			part = append(part, codePoint)
		} else if okUnitCn {
			if partN == 0 {
				// 如 s="十"
				stackValue = append(stackValue, 1)
			} else if partN > 0 {
				lastCodePoint := part[partN-1]
				if isDigitCn(lastCodePoint) {
					stackValueN := len(stackValue)
					stackUnitN := len(stackUnit)

					if stackValueN > 0 && stackValueN == stackUnitN && unitCn > stackUnit[stackUnitN-1] {
						// 单位栈中最后一个值比目前还小，说明需要退回累加
						// 比如 "三百零一万"，栈会类似 stackValue=[3,1], stackUnit=[100,10000]
						// 左到右应该越来越小，应该改为 stackValue=[301] stackUnit=[10000]
						lastValue := stackValue[stackValueN-1]
						lastUnit := stackUnit[stackUnitN-1]

						stackValue = stackValue[:stackValueN-1]
						stackUnit = stackUnit[:stackUnitN-1]

						value := lastValue*lastUnit + digitCn2Arabic[string(lastCodePoint)]
						unit := unitCn

						stackValue = append(stackValue, value)
						stackUnit = append(stackUnit, unit)
					} else {
						stackValue = append(stackValue, digitCn2Arabic[string(lastCodePoint)])
					}
					part = []rune{}
				}
			}

			part = append(part, codePoint)
		} else {
			// 非有效数字字符集
			return 0, ErrInvalidData
		}

		str = str[size:]
	}

	n := len(part)
	if n > 2 {
		return 0, ErrInvalidData
	}
	// 遍历后，part 缓存可能残留字符串中最后一个或两个字符，如 "一" "十" "一百零一" "三千万"
	if n > 0 {
		codePoint := part[0]
		chars := string(part)
		if isDigitCn(codePoint) {
			stackValue = append(stackValue, digitCn2Arabic[chars])
		} else if isUnitCn(chars) {
			// 如 s="十"
			stackUnit = append(stackUnit, unitCn2Arabic[chars])
		}
	}

	// 补全让两个栈长度一样
	stackValueN := len(stackValue)
	stackUnitN := len(stackUnit)
	if stackValueN > stackUnitN {
		stackUnit = append(stackUnit, 1)
	} else if stackUnitN > stackValueN {
		stackValue = append(stackValue, 1)
	}

	// 累加各个位
	for len(stackValue) > 0 {
		n := len(stackValue) - 1
		rs += stackValue[n] * stackUnit[n]
		stackValue = stackValue[:n]
		stackUnit = stackUnit[:n]
	}

	return rs, nil
}
