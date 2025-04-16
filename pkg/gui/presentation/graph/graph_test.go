package graph

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func TestRenderCommitGraph(t *testing.T) {
	hashPool := utils.StringPool{}
	pool := func(s string) *string { return hashPool.Add(s) }

	tests := []struct {
		name           string
		commits        []*models.Commit
		expectedOutput string
	}{
		{
			name: "with some merges",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3")}},
				{Hash: pool("3"), Parents: []*string{pool("4")}},
				{Hash: pool("4"), Parents: []*string{pool("5"), pool("7")}},
				{Hash: pool("7"), Parents: []*string{pool("5")}},
				{Hash: pool("5"), Parents: []*string{pool("8")}},
				{Hash: pool("8"), Parents: []*string{pool("9")}},
				{Hash: pool("9"), Parents: []*string{pool("A"), pool("B")}},
				{Hash: pool("B"), Parents: []*string{pool("D")}},
				{Hash: pool("D"), Parents: []*string{pool("D")}},
				{Hash: pool("A"), Parents: []*string{pool("E")}},
				{Hash: pool("E"), Parents: []*string{pool("F")}},
				{Hash: pool("F"), Parents: []*string{pool("D")}},
				{Hash: pool("D"), Parents: []*string{pool("G")}},
			},
			expectedOutput: `
			1 ◯
			2 ◯
			3 ◯
			4 ⏣─╮
			7 │ ◯
			5 ◯─╯
			8 ◯
			9 ⏣─╮
			B │ ◯
			D │ ◯
			A ◯ │
			E ◯ │
			F ◯ │
			D ◯─╯`,
		},
		{
			name: "with a path that has room to move to the left",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("4"), Parents: []*string{pool("3"), pool("5")}},
				{Hash: pool("3"), Parents: []*string{pool("5")}},
				{Hash: pool("5"), Parents: []*string{pool("6")}},
				{Hash: pool("6"), Parents: []*string{pool("7")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			4 │ ⏣─╮
			3 ◯─╯ │
			5 ◯───╯
			6 ◯`,
		},
		{
			name: "with a new commit",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("4"), Parents: []*string{pool("3"), pool("5")}},
				{Hash: pool("Z"), Parents: []*string{pool("Z")}},
				{Hash: pool("3"), Parents: []*string{pool("5")}},
				{Hash: pool("5"), Parents: []*string{pool("6")}},
				{Hash: pool("6"), Parents: []*string{pool("7")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			4 │ ⏣─╮
			Z │ │ │ ◯
			3 ◯─╯ │ │
			5 ◯───╯ │
			6 ◯ ╭───╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("3"), Parents: []*string{pool("5"), pool("4")}},
				{Hash: pool("5"), Parents: []*string{pool("7"), pool("8")}},
				{Hash: pool("4"), Parents: []*string{pool("7")}},
				{Hash: pool("7"), Parents: []*string{pool("11")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			4 │ ◯─╯ │
			7 ◯─╯ ╭─╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("3"), Parents: []*string{pool("5"), pool("4")}},
				{Hash: pool("5"), Parents: []*string{pool("7"), pool("8")}},
				{Hash: pool("7"), Parents: []*string{pool("4"), pool("A")}},
				{Hash: pool("4"), Parents: []*string{pool("B")}},
				{Hash: pool("B"), Parents: []*string{pool("C")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			4 ◯─┴─╯ │ │
			B ◯ ╭───╯ │`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2"), pool("3")}},
				{Hash: pool("3"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("4"), pool("5")}},
				{Hash: pool("4"), Parents: []*string{pool("6"), pool("7")}},
				{Hash: pool("6"), Parents: []*string{pool("8")}},
			},
			expectedOutput: `
			1 ⏣─╮
			3 │ ◯
			2 ⏣─│
			4 ⏣─│─╮
			6 ◯ │ │`,
		},
		{
			name: "new merge path fills gap before continuing path on right",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2"), pool("3"), pool("4"), pool("5")}},
				{Hash: pool("4"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("A")}},
				{Hash: pool("A"), Parents: []*string{pool("6"), pool("B")}},
				{Hash: pool("B"), Parents: []*string{pool("C")}},
			},
			expectedOutput: `
			1 ⏣─┬─┬─╮
			4 │ │ ◯ │
			2 ◯─│─╯ │
			A ⏣─│─╮ │
			B │ │ ◯ │`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("3"), Parents: []*string{pool("5"), pool("4")}},
				{Hash: pool("5"), Parents: []*string{pool("7"), pool("8")}},
				{Hash: pool("7"), Parents: []*string{pool("4"), pool("A")}},
				{Hash: pool("4"), Parents: []*string{pool("B")}},
				{Hash: pool("B"), Parents: []*string{pool("C")}},
				{Hash: pool("C"), Parents: []*string{pool("D")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			4 ◯─┴─╯ │ │
			B ◯ ╭───╯ │
			C ◯ │ ╭───╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Hash: pool("1"), Parents: []*string{pool("2")}},
				{Hash: pool("2"), Parents: []*string{pool("3"), pool("4")}},
				{Hash: pool("3"), Parents: []*string{pool("5"), pool("4")}},
				{Hash: pool("5"), Parents: []*string{pool("7"), pool("G")}},
				{Hash: pool("7"), Parents: []*string{pool("8"), pool("A")}},
				{Hash: pool("8"), Parents: []*string{pool("4"), pool("E")}},
				{Hash: pool("4"), Parents: []*string{pool("B")}},
				{Hash: pool("B"), Parents: []*string{pool("C")}},
				{Hash: pool("C"), Parents: []*string{pool("D")}},
				{Hash: pool("D"), Parents: []*string{pool("F")}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			8 ⏣─│─│─│─│─╮
			4 ◯─┴─╯ │ │ │
			B ◯ ╭───╯ │ │
			C ◯ │ ╭───╯ │
			D ◯ │ │ ╭───╯`,
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
			lines := RenderCommitGraph(&hashPool, test.commits, hashPool.Add("blah"), getStyle)

			trimmedExpectedOutput := ""
			for _, line := range strings.Split(strings.TrimPrefix(test.expectedOutput, "\n"), "\n") {
				trimmedExpectedOutput += strings.TrimSpace(line) + "\n"
			}

			t.Log("\nexpected: \n" + trimmedExpectedOutput)

			output := ""
			for i, line := range lines {
				description := *test.commits[i].Hash
				output += strings.TrimSpace(description+" "+utils.Decolorise(line)) + "\n"
			}
			t.Log("\nactual: \n" + output)

			assert.Equal(t,
				trimmedExpectedOutput,
				output)
		})
	}
}

func TestRenderPipeSet(t *testing.T) {
	cyan := style.FgCyan
	red := style.FgRed
	green := style.FgGreen
	// blue := style.FgBlue
	yellow := style.FgYellow
	magenta := style.FgMagenta
	nothing := style.Nothing

	hashPool := utils.StringPool{}
	pool := func(s string) *string { return hashPool.Add(s) }

	tests := []struct {
		name           string
		pipes          []Pipe
		commit         *models.Commit
		prevCommit     *models.Commit
		expectedStr    string
		expectedStyles []style.TextStyle
	}{
		{
			name: "single cell",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("b"), kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromHash: pool("b"), toHash: pool("c"), kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Hash: pool("a")},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{green},
		},
		{
			name: "single cell, selected",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("selected"), kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("c"), kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Hash: pool("a")},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{highlightStyle},
		},
		{
			name: "terminating hook and starting hook, selected",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("selected"), kind: TERMINATES, style: cyan},
				{fromPos: 1, toPos: 0, fromHash: pool("c"), toHash: pool("selected"), kind: TERMINATES, style: yellow},
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("d"), kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: pool("selected"), toHash: pool("e"), kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Hash: pool("a")},
			expectedStr: "⏣─╮",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
		{
			name: "terminating hook and starting hook, prioritise the terminating one",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("b"), kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 0, fromHash: pool("c"), toHash: pool("b"), kind: TERMINATES, style: magenta},
				{fromPos: 0, toPos: 0, fromHash: pool("b"), toHash: pool("d"), kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: pool("b"), toHash: pool("e"), kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Hash: pool("a")},
			expectedStr: "⏣─│",
			expectedStyles: []style.TextStyle{
				green, green, magenta,
			},
		},
		{
			name: "starting and terminating pipe sharing some space",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: pool("b1"), toHash: pool("b2"), kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromHash: pool("e1"), toHash: pool("a2"), kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromHash: pool("a2"), toHash: pool("c3"), kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "⏣─│─┬─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, yellow, green, green,
			},
		},
		{
			name: "starting and terminating pipe sharing some space, with selection",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("selected"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: pool("b1"), toHash: pool("b2"), kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromHash: pool("e1"), toHash: pool("selected"), kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromHash: pool("selected"), toHash: pool("c3"), kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "⏣───╮ ╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle, highlightStyle, highlightStyle, nothing, green,
			},
		},
		{
			name: "many terminating pipes",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 0, fromHash: pool("b1"), toHash: pool("a2"), kind: TERMINATES, style: magenta},
				{fromPos: 2, toPos: 0, fromHash: pool("c1"), toHash: pool("a2"), kind: TERMINATES, style: green},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "◯─┴─╯",
			expectedStyles: []style.TextStyle{
				yellow, magenta, magenta, green, green,
			},
		},
		{
			name: "starting pipe passing through",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 3, fromHash: pool("a2"), toHash: pool("d3"), kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: pool("b1"), toHash: pool("b3"), kind: CONTINUES, style: magenta},
				{fromPos: 2, toPos: 2, fromHash: pool("c1"), toHash: pool("c3"), kind: CONTINUES, style: green},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "⏣─│─│─╮",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, green, yellow, yellow,
			},
		},
		{
			name: "starting and terminating path crossing continuing path",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromHash: pool("a2"), toHash: pool("b3"), kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: pool("b1"), toHash: pool("a2"), kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 0, fromHash: pool("c1"), toHash: pool("a2"), kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "⏣─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, green, magenta, magenta,
			},
		},
		{
			name: "another clash of starting and terminating paths",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromHash: pool("a2"), toHash: pool("b3"), kind: STARTS, style: yellow},
				{fromPos: 2, toPos: 2, fromHash: pool("c1"), toHash: pool("c3"), kind: CONTINUES, style: green},
				{fromPos: 3, toPos: 0, fromHash: pool("d1"), toHash: pool("a2"), kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Hash: pool("a1")},
			expectedStr: "⏣─┬─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, yellow, magenta, green, magenta, magenta,
			},
		},
		{
			name: "commit whose previous commit is selected",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: pool("selected")},
			expectedStr: "◯",
			expectedStyles: []style.TextStyle{
				yellow,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromHash: pool("selected"), toHash: pool("b3"), kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Hash: pool("selected")},
			expectedStr: "◯ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, highlightStyle,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit, with continuing pipe inbetween",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("selected"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromHash: pool("z1"), toHash: pool("z3"), kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 2, fromHash: pool("selected"), toHash: pool("b3"), kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Hash: pool("selected")},
			expectedStr: "◯ │ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, green, nothing, highlightStyle,
			},
		},
		{
			name: "when previous commit is selected, not a merge commit, and spawns a continuing pipe",
			pipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a1"), toHash: pool("a2"), kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: pool("a2"), toHash: pool("a3"), kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: pool("a2"), toHash: pool("b3"), kind: STARTS, style: green},
				{fromPos: 1, toPos: 0, fromHash: pool("selected"), toHash: pool("a2"), kind: TERMINATES, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: pool("selected")},
			expectedStr: "⏣─╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualStr := renderPipeSet(test.pipes, pool("selected"), test.prevCommit)
			t.Log("actual cells:")
			t.Log(actualStr)
			expectedStr := ""
			if len([]rune(test.expectedStr)) != len(test.expectedStyles) {
				t.Fatalf("Error in test setup: you have %d characters in the expected output (%s) but have specified %d styles", len([]rune(test.expectedStr)), test.expectedStr, len(test.expectedStyles))
			}
			for i, char := range []rune(test.expectedStr) {
				expectedStr += test.expectedStyles[i].Sprint(string(char))
			}
			expectedStr += " "
			t.Log("expected cells:")
			t.Log(expectedStr)

			assert.Equal(t, expectedStr, actualStr)
		})
	}
}

func TestGetNextPipes(t *testing.T) {
	hashPool := utils.StringPool{}
	pool := func(s string) *string { return hashPool.Add(s) }

	tests := []struct {
		prevPipes []Pipe
		commit    *models.Commit
		expected  []Pipe
	}{
		{
			prevPipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("b"), kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    pool("b"),
				Parents: []*string{pool("c")},
			},
			expected: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("b"), kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromHash: pool("b"), toHash: pool("c"), kind: STARTS, style: style.FgDefault},
			},
		},
		{
			prevPipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("b"), kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromHash: pool("b"), toHash: pool("c"), kind: STARTS, style: style.FgDefault},
				{fromPos: 0, toPos: 1, fromHash: pool("b"), toHash: pool("d"), kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    pool("d"),
				Parents: []*string{pool("e")},
			},
			expected: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("b"), toHash: pool("c"), kind: CONTINUES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromHash: pool("b"), toHash: pool("d"), kind: TERMINATES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromHash: pool("d"), toHash: pool("e"), kind: STARTS, style: style.FgDefault},
			},
		},
		{
			prevPipes: []Pipe{
				{fromPos: 0, toPos: 0, fromHash: pool("a"), toHash: pool("root"), kind: TERMINATES, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    pool("root"),
				Parents: []*string{},
			},
			expected: []Pipe{
				{fromPos: 1, toPos: 1, fromHash: pool("root"), toHash: pool(models.EmptyTreeCommitHash), kind: STARTS, style: style.FgDefault},
			},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, test := range tests {
		getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
		pipes := getNextPipes(&hashPool, test.prevPipes, test.commit, getStyle)
		// rendering cells so that it's easier to see what went wrong
		actualStr := renderPipeSet(pipes, pool("selected"), nil)
		expectedStr := renderPipeSet(test.expected, pool("selected"), nil)
		t.Log("expected cells:")
		t.Log(expectedStr)
		t.Log("actual cells:")
		t.Log(actualStr)
		assert.EqualValues(t, test.expected, pipes)
	}
}

func BenchmarkRenderCommitGraph(b *testing.B) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	hashPool := utils.StringPool{}

	commits := generateCommits(&hashPool, 50)
	getStyle := func(commit *models.Commit) style.TextStyle {
		return authors.AuthorStyle(commit.AuthorName)
	}
	b.ResetTimer()
	for b.Loop() {
		RenderCommitGraph(&hashPool, commits, hashPool.Add("selected"), getStyle)
	}
}

func generateCommits(hashPool *utils.StringPool, count int) []*models.Commit {
	rnd := rand.New(rand.NewSource(1234))
	pool := []*models.Commit{{Hash: hashPool.Add("a"), AuthorName: "A"}}
	commits := make([]*models.Commit, 0, count)
	authorPool := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for len(commits) < count {
		currentCommitIdx := rnd.Intn(len(pool))
		currentCommit := pool[currentCommitIdx]
		pool = append(pool[0:currentCommitIdx], pool[currentCommitIdx+1:]...)
		// I need to pick a random number of parents to add
		parentCount := rnd.Intn(2) + 1

		for j := 0; j < parentCount; j++ {
			reuseParent := rnd.Intn(6) != 1 && j <= len(pool)-1 && j != 0
			var newParent *models.Commit
			if reuseParent {
				newParent = pool[j]
			} else {
				newParent = &models.Commit{
					Hash:       hashPool.Add(fmt.Sprintf("%s%d", *currentCommit.Hash, j)),
					AuthorName: authorPool[rnd.Intn(len(authorPool))],
				}
				pool = append(pool, newParent)
			}
			currentCommit.Parents = append(currentCommit.Parents, newParent.Hash)
		}

		commits = append(commits, currentCommit)
	}

	return commits
}
