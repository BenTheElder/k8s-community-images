package community_images

import (
	"testing"

	semver "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SplitStringOnLen(t *testing.T) {
	tests := []struct {
		name              string
		allSemverTags     []string
		expectedOutliers  []string
		expectedRemaining []string
	}{
		{
			name: "elasticsearch",
			allSemverTags: []string{
				"0.0.1",
				"0.0.2",
				"0.8-alpha4",
				"1.0.0",
				"1.0.1",
				"43.0.0",
			},
			expectedOutliers: []string{
				"43.0.0",
			},
			expectedRemaining: []string{
				"0.0.1",
				"0.0.2",
				"0.8-alpha4",
				"1.0.0",
				"1.0.1",
			},
		},
		{
			name: "no outliers",
			allSemverTags: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"4.0.1",
			},
			expectedOutliers: []string{},
			expectedRemaining: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"4.0.1",
			},
		},
		{

			name: "just barely an outlier",
			allSemverTags: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"3.0.1",
				"5.0.0-beta",
				"8.0.0",
			},
			expectedOutliers: []string{
				"8.0.0",
			},
			expectedRemaining: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"3.0.1",
				"5.0.0-beta",
			},
		},
		{
			name: "multiple spread outliers",
			allSemverTags: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"4.0.1",
				"7.0.0",
				"11.0.0",
			},
			expectedOutliers: []string{
				"7.0.0",
				"11.0.0",
			},
			expectedRemaining: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"4.0.1",
			},
		},
		{
			name: "multiple ranges",
			allSemverTags: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
				"7.0.1",
				"7.0.2",
				"7.0.3",
				"11.0.0",
				"12.0-beta",
			},
			expectedOutliers: []string{
				"7.0.1",
				"7.0.2",
				"7.0.3",
				"11.0.0",
				"12.0-beta",
			},
			expectedRemaining: []string{
				"0.0.1",
				"1.0.0",
				"2.0.0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			allSemvers := []*semver.Version{}
			for _, s := range test.allSemverTags {
				v := semver.Must(semver.NewVersion(s))
				allSemvers = append(allSemvers, v)
			}

			actualOutliers, actualRemaining, err := splitOutlierSemvers(allSemvers)
			req.NoError(err)

			actualOutliersStrings := []string{}
			for _, a := range actualOutliers {
				actualOutliersStrings = append(actualOutliersStrings, a.Original())
			}

			actualRemainingStrings := []string{}
			for _, a := range actualRemaining {
				actualRemainingStrings = append(actualRemainingStrings, a.Original())
			}

			assert.Equal(t, test.expectedOutliers, actualOutliersStrings)
			assert.Equal(t, test.expectedRemaining, actualRemainingStrings)
		})
	}
}
