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
	KeysConstraint string
	Name           string
	NodeName       string
	KeyName        string
	AddNullByte    bool

	HasRange, HasPrefix bool
}

func main() {
	trees := []Tree{
		{
			KeysConstraint: "chars",
			Name:           "alphaSortedTree",
			NodeName:       "alphaLeafNode",
			KeyName:        "AlphabeticalOrderKey",

			AddNullByte: true,
			HasRange:    true,
			HasPrefix:   true,
		},
		{
			KeysConstraint: "uints",
			Name:           "unsignedSortedTree",
			NodeName:       "unsignedLeafNode",
			KeyName:        "UnsignedBinaryKey",

			AddNullByte: false,
		},
		{
			KeysConstraint: "ints",
			Name:           "signedSortedTree",
			NodeName:       "signedLeafNode",
			KeyName:        "SignedBinaryKey",

			AddNullByte: false,
		},
		{
			KeysConstraint: "floats",
			Name:           "floatSortedTree",
			NodeName:       "floatLeafNode",
			KeyName:        "FloatBinaryKey",

			AddNullByte: false,
		},
		{
			KeysConstraint: "any",
			Name:           "compoundSortedTree",
			NodeName:       "compoundLeafNode",
			KeyName:        "BinaryComparableKey",

			AddNullByte: false,
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
