package java

import (
	"os"
	"testing"

	"github.com/vifraa/gopom"

	"github.com/anchore/syft/syft/pkg"
	"github.com/stretchr/testify/assert"
)

func Test_parserPomXML(t *testing.T) {
	tests := []struct {
		input    string
		expected []*pkg.Package
	}{
		{
			input: "test-fixtures/pom/pom.xml",
			expected: []*pkg.Package{
				{
					Name:         "joda-time",
					Version:      "2.9.2",
					FoundBy:      javaPomCataloger,
					Language:     pkg.Java,
					Type:         pkg.JavaPkg,
					MetadataType: pkg.JavaMetadataType,
					Metadata: pkg.JavaMetadata{
						PURL: "pkg:maven/joda-time/joda-time@2.9.2",
					},
				},
				{
					Name:         "junit",
					Version:      "4.12",
					FoundBy:      "java-pom-cataloger",
					Language:     pkg.Java,
					Type:         pkg.JavaPkg,
					MetadataType: pkg.JavaMetadataType,
					Metadata: pkg.JavaMetadata{
						PURL: "pkg:maven/junit/junit@4.12",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			fixture, err := os.Open(test.input)
			assert.NoError(t, err)

			actual, relationships, err := parserPomXML(fixture.Name(), fixture)
			assert.NoError(t, err)
			assert.Nil(t, relationships)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_parsePomXMLProject(t *testing.T) {
	tests := []struct {
		expected pkg.PomProject
	}{
		{
			expected: pkg.PomProject{
				Path: "test-fixtures/pom/commons-codec.pom.xml",
				Parent: &pkg.PomParent{
					GroupID:    "org.apache.commons",
					ArtifactID: "commons-parent",
					Version:    "42",
				},
				GroupID:     "commons-codec",
				ArtifactID:  "commons-codec",
				Version:     "1.11",
				Name:        "Apache Commons Codec",
				Description: "The Apache Commons Codec package contains simple encoder and decoders for various formats such as Base64 and Hexadecimal.  In addition to these widely used encoders and decoders, the codec package also maintains a collection of phonetic encoding utilities.",
				URL:         "http://commons.apache.org/proper/commons-codec/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.expected.Path, func(t *testing.T) {
			fixture, err := os.Open(test.expected.Path)
			assert.NoError(t, err)

			actual, err := parsePomXMLProject(fixture.Name(), fixture)
			assert.NoError(t, err)

			assert.Equal(t, &test.expected, actual)
		})
	}
}

func Test_pomParent(t *testing.T) {
	tests := []struct {
		name     string
		input    gopom.Parent
		expected *pkg.PomParent
	}{
		{
			name: "only group ID",
			input: gopom.Parent{
				GroupID: "org.something",
			},
			expected: &pkg.PomParent{
				GroupID: "org.something",
			},
		},
		{
			name: "only artifact ID",
			input: gopom.Parent{
				ArtifactID: "something",
			},
			expected: &pkg.PomParent{
				ArtifactID: "something",
			},
		},
		{
			name: "only Version",
			input: gopom.Parent{
				Version: "something",
			},
			expected: &pkg.PomParent{
				Version: "something",
			},
		},
		{
			name:     "empty",
			input:    gopom.Parent{},
			expected: nil,
		},
		{
			name: "unused field",
			input: gopom.Parent{
				RelativePath: "something",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, pomParent(test.input))
		})
	}
}

func Test_cleanDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "indent + multiline",
			input: `        The Apache Commons Codec package contains simple encoder and decoders for
        various formats such as Base64 and Hexadecimal.  In addition to these
        widely used encoders and decoders, the codec package also maintains a
        collection of phonetic encoding utilities.`,
			expected: "The Apache Commons Codec package contains simple encoder and decoders for various formats such as Base64 and Hexadecimal.  In addition to these widely used encoders and decoders, the codec package also maintains a collection of phonetic encoding utilities.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, cleanDescription(test.input))
		})
	}
}
