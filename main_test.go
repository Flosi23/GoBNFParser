package BNF

import "testing"

func TestParseString_OneOption_NoVars(t *testing.T) {
	rule := Rule{
		name:    "Zero",
		options: [][]Def{{{symbol: "0"}}},
	}
	want := ParseTreeNode{
		ruleName: "Zero",
		children: []ParseTreeNode{
			{value: "0"},
		},
	}
	result, _ := rule.ParseString("0")

	if !want.equals(result) {
		t.Errorf("rule.parseString('0') wanted: %+v\n got: %+v\n", want, result)
	}
}

func TestParseString_MultipleOptions_NoVars(t *testing.T) {
	rule := Rule{
		name:    "ZeroOrOne",
		options: [][]Def{{{symbol: "0", isVariable: false}}, {{symbol: "1", isVariable: false}}},
	}
	want := ParseTreeNode{
		ruleName: "ZeroOrOne",
		children: []ParseTreeNode{
			{value: "1"},
		},
	}
	result, _ := rule.ParseString("1")

	if !want.equals(result) {
		t.Errorf("rule.parseString('1') wanted: %+v\n got: %+v\n", want, result)
	}
}

func TestParseString_OneOption_OneVar(t *testing.T) {
	abcRule := Rule{
		name: "abc", options: [][]Def{{{symbol: "abc", isVariable: false}}},
	}
	abcWrapper := Rule{
		name:    "abcWrapper",
		options: [][]Def{{{rule: abcRule, isVariable: true}}},
	}
	want := ParseTreeNode{
		ruleName: "abcWrapper",
		children: []ParseTreeNode{
			{
				ruleName: "abc",
				children: []ParseTreeNode{
					{value: "abc"},
				},
			},
		},
	}
	result, _ := abcWrapper.ParseString("abc")

	if !want.equals(result) {
		t.Errorf("rule.parseString('abc') wanted: %+v\n got: %+v\n", want, result)
	}
}
