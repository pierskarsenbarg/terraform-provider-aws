//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=ListTagsForResource -ListTagsInIDElem=ResourceArn -ListTagsOutTagsElem=TagList -ServiceTagsSlice -TagOp=TagResource -TagInTagsElem=TagList -TagInIDElem=ResourceArn -UpdateTags -TagType=Tag

//go:generate go run ../../generate/servicepackage/main.go
// ONLY generate directives and package declaration! Do not add anything else to this file.

package fms
