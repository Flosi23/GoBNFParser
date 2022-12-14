package BNF

import "errors"

type Rule struct {
	name    string
	options [][]Def
}

type Def struct {
	isVariable bool
	isSelf     bool
	rule       Rule
	symbol     string
}

type ParseTreeNode struct {
	ruleName string
	value    string
	children []ParseTreeNode
}

func optionParseString(parentRule Rule, option []Def, ruleTree ParseTreeNode, input string) (ParseTreeNode, error) {
	currentString := input

	for optionNum, def := range option {
		if def.isVariable {
			var varErr error

			rule := def.rule
			if def.isSelf {
				rule = parentRule
			}

			for pos := range currentString {
				tree, err := rule.ParseString(currentString[0 : pos+1])
				varErr = err

				if err == nil && !(len(currentString[pos+1:]) > 0 && optionNum+1 == len(option)) {
					currentString = currentString[pos+1:]
					ruleTree.children = append(ruleTree.children, tree)
					break
				}
			}

			if varErr != nil {
				return ParseTreeNode{}, varErr
			}
		} else {
			if len(currentString) < len(def.symbol) {
				return ParseTreeNode{}, errors.New("no_match")
			}

			if currentString[0:len(def.symbol)] == def.symbol {
				currentString = currentString[len(def.symbol):]
				ruleTree.children = append(ruleTree.children, ParseTreeNode{
					value: def.symbol,
				})
			} else {
				return ParseTreeNode{}, errors.New("no_match")
			}
		}
	}

	if len(currentString) > 0 {
		return ParseTreeNode{}, errors.New("no_match")
	}

	return ruleTree, nil
}

func (rule Rule) ParseString(input string) (ParseTreeNode, error) {
	for _, option := range rule.options {
		tree := ParseTreeNode{
			ruleName: rule.name,
		}

		tree, err := optionParseString(rule, option, tree, input)

		if err == nil {
			return tree, nil
		}
	}

	return ParseTreeNode{}, errors.New("invalid input")
}

func (treeOne ParseTreeNode) equals(treeTwo ParseTreeNode) bool {
	isValueEqual := treeOne.value == treeTwo.value
	isRuleNameEqual := treeOne.ruleName == treeTwo.ruleName
	areChildrenEqual := true

	if len(treeOne.children) == len(treeTwo.children) {
		for i := range treeOne.children {
			if !treeOne.children[i].equals(treeOne.children[i]) {
				areChildrenEqual = false
				break
			}
		}
	} else {
		areChildrenEqual = false
	}

	return isValueEqual && isRuleNameEqual && areChildrenEqual
}
