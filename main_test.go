package BNF

import (
	"testing"
)

func throwWantedXGotYError(testCase string, wanted any, got any, t *testing.T) {
	t.Errorf("%s wanted: %+v\n got %+v'n", testCase, wanted, got)
}

func throwInvalidGrammarNoErrorError(testCase string, t *testing.T) {
	t.Errorf("Expected %s to throw an error", testCase)
}

func TestParseString_String(t *testing.T) {
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
		throwWantedXGotYError("rule.ParseString('0')", want, result, t)
	}

	_, err := rule.ParseString("1")

	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString('1')", t)
	}
}

func TestParseString_MultipleOptions(t *testing.T) {
	rule := Rule{
		name:    "ZeroOrOne",
		options: [][]Def{{{symbol: "0", isVariable: false}}, {{symbol: "1", isVariable: false}}},
	}

	wantOne := ParseTreeNode{
		ruleName: "ZeroOrOne",
		children: []ParseTreeNode{
			{value: "1"},
		},
	}
	result, _ := rule.ParseString("1")
	if !wantOne.equals(result) {
		throwWantedXGotYError("rule.parseString('1')", wantOne, result, t)
	}

	wantTwo := ParseTreeNode{
		ruleName: "ZeroOrOne",
		children: []ParseTreeNode{
			{value: "1"},
		},
	}
	result, _ = rule.ParseString("0")
	if !wantTwo.equals(result) {
		throwWantedXGotYError("rule.parseString('0')", wantTwo, result, t)
	}

	_, err := rule.ParseString("3")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.parseString('3')", t)
	}
}

func TestParseString_Var(t *testing.T) {
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
		throwWantedXGotYError("rule.parseString('abc')", want, result, t)
	}

	_, err := abcWrapper.ParseString("bac")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.parseString('bac')", t)
	}
}

func TestParseString_StringVar(t *testing.T) {
	oneRule := Rule{
		name: "One",
		options: [][]Def{
			{{symbol: "1"}},
		},
	}
	minusOneRule := Rule{
		name: "MinusOne",
		options: [][]Def{
			{{symbol: "-"}, {isVariable: true, rule: oneRule}},
		},
	}

	want := ParseTreeNode{
		ruleName: "MinusOne",
		children: []ParseTreeNode{
			{value: "-"},
			{
				ruleName: "One",
				children: []ParseTreeNode{
					{value: "1"},
				},
			},
		},
	}

	result, _ := minusOneRule.ParseString("-1")
	if !want.equals(result) {
		throwWantedXGotYError("rule.parseString('-1')", want, result, t)
	}

	_, err := minusOneRule.ParseString("+")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.parseString('+')", t)
	}
}

func TestParseString_VarStringVar(t *testing.T) {
	digitRule := Rule{
		name: "Digit",
		options: [][]Def{
			{{symbol: "1"}},
			{{symbol: "2"}},
		},
	}
	additionRule := Rule{
		name: "Addition",
		options: [][]Def{
			{
				{isVariable: true, rule: digitRule},
				{symbol: "+"},
				{isVariable: true, rule: digitRule},
			},
		},
	}

	want := ParseTreeNode{
		ruleName: "Addition",
		children: []ParseTreeNode{
			{
				ruleName: "Digit",
				children: []ParseTreeNode{{value: "2"}},
			},
			{value: "+"},
			{
				ruleName: "Digit",
				children: []ParseTreeNode{{value: "1"}},
			},
		},
	}

	result, _ := additionRule.ParseString("2+1")
	if !want.equals(result) {
		throwWantedXGotYError("rule.ParseString('2+1')", want, result, t)
	}

	_, err := additionRule.ParseString("23-45")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString(23-45)", t)
	}
}

func TestParseString_VarVar(t *testing.T) {
	firstnameRule := Rule{
		name: "First name",
		options: [][]Def{
			{{symbol: "Simon"}},
		},
	}
	surnameRule := Rule{
		name: "Surname",
		options: [][]Def{
			{{symbol: "Weckler"}},
		},
	}
	fullnameRule := Rule{
		name: "Full name",
		options: [][]Def{
			{{rule: firstnameRule, isVariable: true}, {rule: surnameRule, isVariable: true}},
		},
	}

	want := ParseTreeNode{
		ruleName: "Full name",
		children: []ParseTreeNode{
			{
				ruleName: "First name",
				children: []ParseTreeNode{
					{value: "Simon"},
				},
			},
			{
				ruleName: "Surname",
				children: []ParseTreeNode{
					{value: "Weckler"},
				},
			},
		},
	}

	result, _ := fullnameRule.ParseString("SimonWeckler")
	if !want.equals(result) {
		throwWantedXGotYError("rule.ParseString('SimonWeckler')", want, result, t)
	}

	_, err := fullnameRule.ParseString("JakobWeckler")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString('JakobWeckler')", t)
	}

	_, err = fullnameRule.ParseString("SimonScholz")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString('SimonScholz')", t)
	}

	_, err = fullnameRule.ParseString("JakobScholz")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString('JakobScholz')", t)
	}
}

func TestParseString_Recursion(t *testing.T) {
	zeroOrOneRule := Rule{
		name: "ZeroOrOne",
		options: [][]Def{
			{{symbol: "0"}},
			{{symbol: "1"}},
		},
	}
	machineLanguageRule := Rule{
		name: "Machine Language",
		options: [][]Def{
			{{isVariable: true, rule: zeroOrOneRule}},
			{{isVariable: true, rule: zeroOrOneRule}, {isVariable: true, isSelf: true}},
		},
	}

	want := ParseTreeNode{
		ruleName: "Machine Language",
		children: []ParseTreeNode{
			{
				ruleName: "ZeroOrOne",
				children: []ParseTreeNode{
					{value: "1"},
				},
			},
			{
				ruleName: "Machine Language",
				children: []ParseTreeNode{
					{
						ruleName: "ZeroOrOne",
						children: []ParseTreeNode{
							{value: "1"},
						},
					},
					{
						ruleName: "Machine Language",
						children: []ParseTreeNode{
							{
								ruleName: "ZeroOrOne",
								children: []ParseTreeNode{
									{value: "0"},
								},
							},
						},
					},
				},
			},
		},
	}

	result, _ := machineLanguageRule.ParseString("110")
	if !result.equals(want) {
		throwWantedXGotYError("rule.ParseString('110')", want, result, t)
	}

	_, err := machineLanguageRule.ParseString("12")
	if err == nil {
		throwInvalidGrammarNoErrorError("rule.ParseString('12')", t)
	}
}
