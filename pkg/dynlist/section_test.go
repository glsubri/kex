package dynlist

import (
	"fmt"
	"testing"
)

func TestScrollDown(t *testing.T) {
	tcs := []struct {
		initialSection  section
		by              int
		maxHeight       int
		expectedSection section
	}{
		{initialSection: section{0, 10}, by: 5, maxHeight: 100, expectedSection: section{5, 15}},
		{initialSection: section{0, 10}, by: 21, maxHeight: 100, expectedSection: section{21, 31}},
		{initialSection: section{0, 10}, by: 5, maxHeight: 12, expectedSection: section{2, 12}},
		{initialSection: section{0, 10}, by: 2, maxHeight: 12, expectedSection: section{2, 12}},
		{initialSection: section{0, 10}, by: 3, maxHeight: 12, expectedSection: section{2, 12}},
		{initialSection: section{0, 45}, by: 3, maxHeight: 12, expectedSection: section{0, 45}},
	}

	for _, tc := range tcs {
		name := fmt.Sprintf(
			"[%d:%d] + %d with max %d",
			tc.initialSection.start,
			tc.initialSection.end,
			tc.by,
			tc.maxHeight,
		)

		t.Run(name, func(t *testing.T) {
			res := tc.initialSection.scrollDown(tc.by, tc.maxHeight)
			if res != tc.expectedSection {
				t.Errorf("has %+v, expected %+v", res, tc.expectedSection)
			}
		})
	}
}

func TestScrollUp(t *testing.T) {
	tcs := []struct {
		initialSection  section
		by              int
		expectedSection section
	}{
		{initialSection: section{0, 10}, by: 5, expectedSection: section{0, 10}},
		{initialSection: section{10, 20}, by: 5, expectedSection: section{5, 15}},
	}

	for _, tc := range tcs {
		name := fmt.Sprintf(
			"[%d:%d] - %d",
			tc.initialSection.start,
			tc.initialSection.end,
			tc.by,
		)

		t.Run(name, func(t *testing.T) {
			res := tc.initialSection.scrollUp(tc.by)
			if res != tc.expectedSection {
				t.Errorf("has %+v, expected %+v", res, tc.expectedSection)
			}
		})
	}
}

func TestHeight(t *testing.T) {
	tcs := []struct {
		section section
		height  int
	}{
		{section: section{0, 1}, height: 1},
		{section: section{1, 2}, height: 1},
		{section: section{0, 10}, height: 10},
		{section: section{10, 20}, height: 10},
	}

	for _, tc := range tcs {
		name := fmt.Sprintf("height of [%d:%d] is %d", tc.section.start, tc.section.end, tc.height)

		t.Run(name, func(t *testing.T) {
			res := tc.section.height()
			if res != tc.height {
				t.Errorf("has %d, expected %d", res, tc.height)
			}
		})
	}
}
