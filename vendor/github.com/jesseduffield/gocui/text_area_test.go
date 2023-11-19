package gocui

import (
	"reflect"
	"testing"
)

func Test_AutoWrapContent(t *testing.T) {
	tests := []struct {
		name                   string
		content                string
		autoWrapWidth          int
		expectedWrappedContent string
		expectedCursorMapping  []CursorMapping
	}{
		{
			name:                   "empty content",
			content:                "",
			autoWrapWidth:          7,
			expectedWrappedContent: "",
			expectedCursorMapping:  []CursorMapping{},
		},
		{
			name:                   "no wrapping necessary",
			content:                "abcde",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde",
			expectedCursorMapping:  []CursorMapping{},
		},
		{
			name:                   "wrap at whitespace",
			content:                "abcde xyz",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde \nxyz",
			expectedCursorMapping:  []CursorMapping{{6, 7}},
		},
		{
			name:                   "lots of whitespace is preserved at end of line",
			content:                "abcde      xyz",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde      \nxyz",
			expectedCursorMapping:  []CursorMapping{{11, 12}},
		},
		{
			name:                   "don't wrap inside long word when there's no whitespace",
			content:                "abc defghijklmn opq",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc \ndefghijklmn \nopq",
			expectedCursorMapping:  []CursorMapping{{4, 5}, {16, 18}},
		},
		{
			name:                   "hard line breaks",
			content:                "abc\ndef\n",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc\ndef\n",
			expectedCursorMapping:  []CursorMapping{},
		},
		{
			name:                   "mixture of hard and soft line breaks",
			content:                "abc def ghi jkl mno\npqr stu vwx yz\n",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc def \nghi jkl \nmno\npqr stu \nvwx yz\n",
			expectedCursorMapping:  []CursorMapping{{8, 9}, {16, 18}, {28, 31}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedContent, cursorMapping := AutoWrapContent([]rune(tt.content), tt.autoWrapWidth)
			if !reflect.DeepEqual(wrappedContent, []rune(tt.expectedWrappedContent)) {
				t.Errorf("autoWrapContentImpl() wrappedContent = %v, expected %v", string(wrappedContent), tt.expectedWrappedContent)
			}
			if !reflect.DeepEqual(cursorMapping, tt.expectedCursorMapping) {
				t.Errorf("autoWrapContentImpl() cursorMapping = %v, expected %v", cursorMapping, tt.expectedCursorMapping)
			}

			// As a sanity check, run through all runes of the original content,
			// convert the cursor to the wrapped cursor, and check that the rune
			// in the wrapped content at that position is the same:
			for i, r := range tt.content {
				wrappedIndex := origCursorToWrappedCursor(i, cursorMapping)
				if r != wrappedContent[wrappedIndex] {
					t.Errorf("Runes in orig content and wrapped content don't match at %d: expected %v, got %v", i, r, wrappedContent[wrappedIndex])
				}

				// Also, check that converting the wrapped position back to the
				// orig position yields the original value again:
				origIndexAgain := wrappedCursorToOrigCursor(wrappedIndex, cursorMapping)
				if i != origIndexAgain {
					t.Errorf("wrappedCursorToOrigCursor doesn't yield original position: expected %d, got %d", i, origIndexAgain)
				}
			}
		})
	}
}
