package nash

import "testing"

type fmtTestTable struct {
	input, expected string
}

func testFmt(input string, expected string, t *testing.T) {
	p := NewParser("fmt test", input)

	tree, err := p.Parse()

	if err != nil {
		t.Error(err)
		return
	}

	fmtval := tree.String()

	if fmtval != expected {
		t.Errorf("Fmt differ: '%s' != '%s'", fmtval, expected)
		return
	}
}

func testFmtTable(testTable []fmtTestTable, t *testing.T) {
	for _, test := range testTable {
		testFmt(test.input, test.expected, t)
	}
}

func TestFmtVariables(t *testing.T) {
	testTable := []fmtTestTable{

		// correct adjust of spaces
		fmtTestTable{`test = "a"`, `test = "a"`},
		fmtTestTable{`test="a"`, `test = "a"`},
		fmtTestTable{`test= "a"`, `test = "a"`},
		fmtTestTable{`test  ="a"`, `test = "a"`},
		fmtTestTable{`test =    "a"`, `test = "a"`},
		fmtTestTable{`test	="a"`, `test = "a"`},
		fmtTestTable{`test		="a"`, `test = "a"`},
		fmtTestTable{`test =	"a"`, `test = "a"`},
		fmtTestTable{`test =		"a"`, `test = "a"`},
		fmtTestTable{`test = ()`, `test = ()`},
		fmtTestTable{`test=()`, `test = ()`},
		fmtTestTable{`test =()`, `test = ()`},
		fmtTestTable{`test	=()`, `test = ()`},
		fmtTestTable{`test=	()`, `test = ()`},
		fmtTestTable{`test = (plan9)`, `test = (plan9)`},
		fmtTestTable{`test=(plan9)`, `test = (plan9)`},
		fmtTestTable{`test      = (plan9)`, `test = (plan9)`},
		fmtTestTable{`test	= (plan9)`, `test = (plan9)`},
		fmtTestTable{`test	=	(plan9)`, `test = (plan9)`},
		fmtTestTable{`test = (	plan9)`, `test = (plan9)`},
		fmtTestTable{`test = (     plan9)`, `test = (plan9)`},
		fmtTestTable{`test = (plan9     )`, `test = (plan9)`},
		fmtTestTable{`test = (plan9 from bell labs)`, `test = (plan9 from bell labs)`},
		fmtTestTable{`test = (plan9         from bell labs)`, `test = (plan9 from bell labs)`},
		fmtTestTable{`test = (plan9         from         bell         labs)`, `test = (plan9 from bell labs)`},
		fmtTestTable{`test = (plan9	from	bell	labs)`, `test = (plan9 from bell labs)`},
		fmtTestTable{`test = (
	plan9
	from
	bell
	labs
)`, `test = (plan9 from bell labs)`},
		fmtTestTable{`test = (plan9 from bell labs windows linux freebsd netbsd openbsd)`, `test = (
	plan9
	from
	bell
	labs
	windows
	linux
	freebsd
	netbsd
	openbsd
)`},

		fmtTestTable{`IFS = ("\n")`, `IFS = ("\n")`},

		// multiple variables
		fmtTestTable{`test = "a"
testb = "b"`, `test = "a"
testb = "b"`},
	}

	testFmtTable(testTable, t)
}
