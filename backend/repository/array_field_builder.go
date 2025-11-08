package repository

import (
	"strconv"
	"strings"
)

type FieldIndexPair struct {
	Field string
	Index int
}

type ArrayFieldPath struct {
	Pairs []FieldIndexPair
}

func (a *ArrayFieldPath) GetPushPath() string {
	return buildArrayField(a.Pairs, true)
}

func (a *ArrayFieldPath) GetPullPath() string {
	return buildArrayField(a.Pairs, false)
}

func (a *ArrayFieldPath) GetUpdatePath() string {
	return buildArrayField(a.Pairs, false)
}

func (a *ArrayFieldPath) GetLastIndex() int {
	return a.Pairs[len(a.Pairs)-1].Index
}

func buildArrayField(pairs []FieldIndexPair, skipLastIndex bool) string {
	var builder strings.Builder

	lastIndex := len(pairs) - 1
	for i, pair := range pairs {
		if i != 0 {
			builder.WriteString(".")
		}
		builder.WriteString(pair.Field)
		if skipLastIndex && i == lastIndex {
			break
		}
		builder.WriteString(".")
		builder.WriteString(strconv.Itoa(pair.Index))
	}

	return builder.String()
}
