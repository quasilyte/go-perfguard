// Code generated by "precompile.go". DO NOT EDIT.

package rulesdata

import "github.com/quasilyte/go-ruleguard/ruleguard/ir"

var Opt = &ir.File{
	PkgPath:       "gorules",
	CustomDecls:   []string{},
	BundleImports: []ir.BundleImport{},
	RuleGroups: []ir.RuleGroup{
		{
			Line:        9,
			Name:        "stringConcatAssign",
			MatcherName: "m",
			DocTags:     []string{"o2", "score5"},
			DocSummary:  "Detects string concat in hot paths",
			Rules: []ir.Rule{{
				Line:           10,
				SyntaxPatterns: []ir.PatternString{{Line: 10, Value: "$s += $_"}},
				ReportTemplate: "string concat on the hot path",
				WhereExpr: ir.FilterExpr{
					Line:  11,
					Op:    ir.FilterVarTypeIsOp,
					Src:   "m[\"s\"].Type.Is(`string`)",
					Value: "s",
					Args:  []ir.FilterExpr{{Line: 11, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
				},
			}},
		},
		{
			Line:        17,
			Name:        "regexpCompile",
			MatcherName: "m",
			DocTags:     []string{"o1", "score4"},
			DocSummary:  "Detects regexp compilation on hot execution paths",
			Rules: []ir.Rule{{
				Line: 22,
				SyntaxPatterns: []ir.PatternString{
					{Line: 24, Value: "regexp.Compile($*_)"},
					{Line: 25, Value: "regexp.MustCompile($*_)"},
					{Line: 26, Value: "regexp.CompilePOSIX($*_)"},
					{Line: 27, Value: "regexp.MustCompilePOSIX($*_)"},
					{Line: 29, Value: "regexp.Match($*_)"},
					{Line: 30, Value: "regexp.MatchString($*_)"},
					{Line: 31, Value: "regexp.MatchReader($*_)"},
				},
				ReportTemplate: "regexp compilation should be avoided on the hot paths",
			}},
		},
		{
			Line:        37,
			Name:        "sprintfConcat2",
			MatcherName: "m",
			DocTags:     []string{"o2", "score2"},
			DocSummary:  "Detects sprint calls that can be rewritten as a string concat",
			Rules: []ir.Rule{
				{
					Line:            42,
					SyntaxPatterns:  []ir.PatternString{{Line: 42, Value: "fmt.Sprintf(\"%s=%s\", $x, $y)"}},
					ReportTemplate:  "$$ => $x + \"=\" + $y",
					SuggestTemplate: "$x + \"=\" + $y",
					WhereExpr: ir.FilterExpr{
						Line: 43,
						Op:   ir.FilterAndOp,
						Src:  "m[\"x\"].Type.Is(`string`) && m[\"y\"].Type.Is(`string`)",
						Args: []ir.FilterExpr{
							{
								Line:  43,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"x\"].Type.Is(`string`)",
								Value: "x",
								Args:  []ir.FilterExpr{{Line: 43, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
							{
								Line:  43,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"y\"].Type.Is(`string`)",
								Value: "y",
								Args:  []ir.FilterExpr{{Line: 43, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
						},
					},
				},
				{
					Line:            46,
					SyntaxPatterns:  []ir.PatternString{{Line: 46, Value: "fmt.Sprintf(\"%s.%s\", $x, $y)"}},
					ReportTemplate:  "$$ => $x + \".\" + $y",
					SuggestTemplate: "$x + \".\" + $y",
					WhereExpr: ir.FilterExpr{
						Line: 47,
						Op:   ir.FilterAndOp,
						Src:  "m[\"x\"].Type.Is(`string`) && m[\"y\"].Type.Is(`string`)",
						Args: []ir.FilterExpr{
							{
								Line:  47,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"x\"].Type.Is(`string`)",
								Value: "x",
								Args:  []ir.FilterExpr{{Line: 47, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
							{
								Line:  47,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"y\"].Type.Is(`string`)",
								Value: "y",
								Args:  []ir.FilterExpr{{Line: 47, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
						},
					},
				},
				{
					Line:            50,
					SyntaxPatterns:  []ir.PatternString{{Line: 50, Value: "fmt.Sprintf(\"%s/%s\", $x, $y)"}},
					ReportTemplate:  "$$ => $x + \"/\" + $y",
					SuggestTemplate: "$x + \"/\" + $y",
					WhereExpr: ir.FilterExpr{
						Line: 51,
						Op:   ir.FilterAndOp,
						Src:  "m[\"x\"].Type.Is(`string`) && m[\"y\"].Type.Is(`string`)",
						Args: []ir.FilterExpr{
							{
								Line:  51,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"x\"].Type.Is(`string`)",
								Value: "x",
								Args:  []ir.FilterExpr{{Line: 51, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
							{
								Line:  51,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"y\"].Type.Is(`string`)",
								Value: "y",
								Args:  []ir.FilterExpr{{Line: 51, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
						},
					},
				},
				{
					Line:            54,
					SyntaxPatterns:  []ir.PatternString{{Line: 54, Value: "fmt.Sprintf(\"%s:%s\", $x, $y)"}},
					ReportTemplate:  "$$ => $x + \":\" + $y",
					SuggestTemplate: "$x + \":\" + $y",
					WhereExpr: ir.FilterExpr{
						Line: 55,
						Op:   ir.FilterAndOp,
						Src:  "m[\"x\"].Type.Is(`string`) && m[\"y\"].Type.Is(`string`)",
						Args: []ir.FilterExpr{
							{
								Line:  55,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"x\"].Type.Is(`string`)",
								Value: "x",
								Args:  []ir.FilterExpr{{Line: 55, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
							{
								Line:  55,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"y\"].Type.Is(`string`)",
								Value: "y",
								Args:  []ir.FilterExpr{{Line: 55, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
						},
					},
				},
				{
					Line:            58,
					SyntaxPatterns:  []ir.PatternString{{Line: 58, Value: "fmt.Sprintf(\"%s: %s\", $x, $y)"}},
					ReportTemplate:  "$$ => $x + \": \" + $y",
					SuggestTemplate: "$x + \": \" + $y",
					WhereExpr: ir.FilterExpr{
						Line: 59,
						Op:   ir.FilterAndOp,
						Src:  "m[\"x\"].Type.Is(`string`) && m[\"y\"].Type.Is(`string`)",
						Args: []ir.FilterExpr{
							{
								Line:  59,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"x\"].Type.Is(`string`)",
								Value: "x",
								Args:  []ir.FilterExpr{{Line: 59, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
							{
								Line:  59,
								Op:    ir.FilterVarTypeIsOp,
								Src:   "m[\"y\"].Type.Is(`string`)",
								Value: "y",
								Args:  []ir.FilterExpr{{Line: 59, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
							},
						},
					},
				},
			},
		},
		{
			Line:        67,
			Name:        "writeString2",
			MatcherName: "m",
			DocTags:     []string{"o2", "score3"},
			DocSummary:  "Detects Write calls that should be rewritten as io.WriteString",
			DocBefore:   "w.Write([]byte(s))",
			DocAfter:    "io.WriteString(w, s)",
			Rules: []ir.Rule{{
				Line:            68,
				SyntaxPatterns:  []ir.PatternString{{Line: 68, Value: "$w.Write([]byte($s))"}},
				ReportTemplate:  "$$ => io.WriteString($w, $s)",
				SuggestTemplate: "io.WriteString($w, $s)",
				WhereExpr: ir.FilterExpr{
					Line: 69,
					Op:   ir.FilterAndOp,
					Src:  "m[\"w\"].Type.Is(\"io.Writer\") && m[\"s\"].Type.Is(`string`) && m[\"s\"].Const",
					Args: []ir.FilterExpr{
						{
							Line: 69,
							Op:   ir.FilterAndOp,
							Src:  "m[\"w\"].Type.Is(\"io.Writer\") && m[\"s\"].Type.Is(`string`)",
							Args: []ir.FilterExpr{
								{
									Line:  69,
									Op:    ir.FilterVarTypeIsOp,
									Src:   "m[\"w\"].Type.Is(\"io.Writer\")",
									Value: "w",
									Args:  []ir.FilterExpr{{Line: 69, Op: ir.FilterStringOp, Src: "\"io.Writer\"", Value: "io.Writer"}},
								},
								{
									Line:  69,
									Op:    ir.FilterVarTypeIsOp,
									Src:   "m[\"s\"].Type.Is(`string`)",
									Value: "s",
									Args:  []ir.FilterExpr{{Line: 69, Op: ir.FilterStringOp, Src: "`string`", Value: "string"}},
								},
							},
						},
						{
							Line:  69,
							Op:    ir.FilterVarConstOp,
							Src:   "m[\"s\"].Const",
							Value: "s",
						},
					},
				},
			}},
		},
		{
			Line:        75,
			Name:        "rangeValueCopy",
			MatcherName: "m",
			DocTags:     []string{"o1", "score2"},
			DocSummary:  "Detects range loops that copy large value on every iteration",
			Rules: []ir.Rule{{
				Line: 77,
				SyntaxPatterns: []ir.PatternString{
					{Line: 77, Value: "for $_, $v := range $_"},
					{Line: 77, Value: "for $_, $v = range $_"},
				},
				ReportTemplate: "every iteration copies a large object into $v",
				WhereExpr: ir.FilterExpr{
					Line: 78,
					Op:   ir.FilterGtOp,
					Src:  "m[\"v\"].Type.Size > 128",
					Args: []ir.FilterExpr{
						{
							Line:  78,
							Op:    ir.FilterVarTypeSizeOp,
							Src:   "m[\"v\"].Type.Size",
							Value: "v",
						},
						{
							Line:  78,
							Op:    ir.FilterIntOp,
							Src:   "128",
							Value: int64(128),
						},
					},
				},
			}},
		},
		{
			Line:        84,
			Name:        "constErrorNew",
			MatcherName: "m",
			DocTags:     []string{"o1", "score3"},
			DocSummary:  "Detects errors.New that can be allocated exactly once",
			Rules: []ir.Rule{{
				Line:           85,
				SyntaxPatterns: []ir.PatternString{{Line: 85, Value: "errors.New($x)"}},
				ReportTemplate: "errors with const message can be a global var, allocated only once",
				WhereExpr: ir.FilterExpr{
					Line:  86,
					Op:    ir.FilterVarConstOp,
					Src:   "m[\"x\"].Const",
					Value: "x",
				},
			}},
		},
	},
}

