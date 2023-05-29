package sqltostruct

import (
	"fmt"
	"strings"
)

var breaks map[rune]bool = map[rune]bool{
	'(': true,
	')': true,
	' ': true,
	';': true,
	',': true,
}

var dbToGo map[string]string = map[string]string{
	"int":     "int",
	"varchar": "string",
	"jsonb":   "string",
}

// Convert will attempt to read sql create table statement and output a golang struct with db struct-tags
func Convert(input string) {
	noNewlines := strings.ReplaceAll(input, "\n", "")
	s := strings.ToLower(noNewlines)
	t := createTable(s)

	fmt.Printf("type %s struct{\n", snakeToCamel(t.name, true))
	for _, s := range t.statements {
		fmt.Printf("\t%s %s `db:\"%s\"`\n", snakeToCamel(s.name, true), dbToGo[s.sType], s.name)
	}
	fmt.Println("}")

}

func snakeToCamel(s string, firstCap bool) string {
	out := ""
	if firstCap {
		out = strings.ToUpper(string(s[0]))
	} else {
		out = strings.ToLower(string(s[0]))
	}
	s = s[1:]
	for i := 0; i < len(s); i++ {
		if s[i] == '_' {
			out = out + strings.ToUpper(string(s[i+1]))
			i = i + 1
		} else {
			out = out + string(s[i])
		}
	}
	return out
}

type table struct {
	name       string
	statements []statement
}
type statement struct {
	name       string
	sType      string
	sTypeValue string
}

func createTable(s string) table {
	ct := "create table"
	rest := strings.TrimPrefix(s, ct)
	if s == rest {
		panic("input don't start with create table")
	}
	rest = strings.TrimSpace(rest)
	t := table{}
	for i := 0; i < len(rest); i++ {
		if breaks[rune(rest[i])] {
			t.name = rest[0:i]
			rest = rest[i:]
			break
		}
	}
	rest = find('(', rest, true)
	for {
		st, theRest, hasMore := getStatement(rest)
		if st != nil {
			t.statements = append(t.statements, *st)
		}
		rest = theRest
		if !hasMore {
			break
		}
	}
	rest = find(')', rest, true)
	rest = find(';', rest, false)

	return t
}
func getStatement(s string) (*statement, string, bool) {
	rest := strings.TrimSpace(s)
	st := statement{}
	for i := 0; i < len(rest); i++ {
		if breaks[rune(rest[i])] {
			st.name = rest[0:i]
			rest = rest[i:]
			break
		}
	}

	if st.name == "constraint" {
		// don't care about this. Find next break and return
		for i := 0; i < len(rest); i++ {

			if breaks[rune(rest[i])] {
				rest = rest[i:]
				return nil, rest, rest[0] == ','
			}
		}
	}
	rest = strings.TrimSpace(rest)
	for i := 0; i < len(rest); i++ {
		if breaks[rune(rest[i])] {
			st.sType = rest[0:i]
			rest = rest[i:]
			break
		}
	}
	rest = strings.TrimSpace(rest)
	if rest[0] == '(' {
		rest = find('(', rest, true)
		for i := 0; i < len(rest); i++ {
			if breaks[rune(rest[i])] {
				st.sTypeValue = rest[0:i]
				rest = rest[i:]
				break
			}
		}
		rest = find(')', rest, true)
	}
	if rest[0] == ',' {
		return &st, rest[1:], true
	}
	return &st, rest, false

}

func find(toFind rune, s string, panicOnFail bool) string {
	for i := 0; i < len(s); i++ {
		if rune(s[i]) == toFind {
			return s[i+1:]
		}
	}
	if panicOnFail {
		panic(fmt.Sprintf("could not find %c in %s", toFind, s))
	}
	return ""
}
