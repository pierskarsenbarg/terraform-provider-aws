//go:generate go run ../../generate/listpages/main.go -ListOps=GetApis,GetApiMappings,GetDomainNames,GetVpcLinks
//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=GetTags -ServiceTagsMap -UpdateTags
//go:generate go run ../../generate/servicepackage/main.go
// ONLY generate directives and package declaration! Do not add anything else to this file.

package apigatewayv2
