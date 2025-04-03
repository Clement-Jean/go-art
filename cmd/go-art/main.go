package main

import (
	"embed"
	"log"
	"os"
	"text/template"
)

//go:embed tree.tmpl
var codeTmpl embed.FS

type Tree struct {
	KeysConstraint              string
	Name                        string
	NodeName                    string
	KeyName                     string
	AddNullByte                 bool
	ComparableKeys, CompoundKey bool

	HasPrefix bool
}

func main() {
	trees := []Tree{
		{
			KeysConstraint: "chars",
			Name:           "alphaSortedTree",
			NodeName:       "alphaLeafNode",
			KeyName:        "AlphabeticalOrderKey",

			AddNullByte:    true,
			ComparableKeys: false,
			HasPrefix:      true,
			CompoundKey:    false,
		},
		{
			KeysConstraint: "uints",
			Name:           "unsignedSortedTree",
			NodeName:       "unsignedLeafNode",
			KeyName:        "UnsignedBinaryKey",

			AddNullByte:    false,
			ComparableKeys: true,
			CompoundKey:    false,
		},
		{
			KeysConstraint: "ints",
			Name:           "signedSortedTree",
			NodeName:       "signedLeafNode",
			KeyName:        "SignedBinaryKey",

			AddNullByte:    false,
			ComparableKeys: true,
			CompoundKey:    false,
		},
		{
			KeysConstraint: "floats",
			Name:           "floatSortedTree",
			NodeName:       "floatLeafNode",
			KeyName:        "FloatBinaryKey",

			AddNullByte:    false,
			ComparableKeys: true,
			CompoundKey:    false,
		},
		{
			KeysConstraint: "any",
			Name:           "compoundSortedTree",
			NodeName:       "compoundLeafNode",
			KeyName:        "BinaryComparableKey",

			AddNullByte:    false,
			ComparableKeys: false,
			CompoundKey:    true,
		},
	}

	tmpl, err := template.ParseFS(codeTmpl, "tree.tmpl")
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("trees.go", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, trees)
	if err != nil {
		panic(err)
	}
}
