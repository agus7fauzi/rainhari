/*
Copyright 2022 Agus Imam Fauzi <agus7fauzi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPackage     = protogen.GoImportPath("context")
	serverPackage      = protogen.GoImportPath("github.com/agus7fauzi/rainhari/core/server")
	clientPackage      = protogen.GoImportPath("github.com/agus7fauzi/rainhari/core/client")
	deprecationComment = "// Deprecated: Do not use."
)

var reversedClientName = map[string]bool{}

// GenerateFile generates *.rainhari.pb.go files containing gRPC service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".rainhari.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-rainhari. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	generateCoreCode(gen, file, g)

	return g
}

func generateCoreCode(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	// Iterate through all the `service`s in the passed-in file
	for i, service := range file.Services {
		generateService(gen, file, g, service, i)
	}
}

func generateService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, index int) {
	path := "6," + string(index)
	originSvcName := service.GoName
	serviceName := strings.ToLower(originSvcName)

	if pkg := file.GoPackageName; pkg != "" {
		serviceName = string(pkg)
	}

	cCaseSvcName := camelCase(originSvcName)
	serviceAlias := cCaseSvcName + "Service"

	// Strip suffix
	if strings.HasSuffix(serviceAlias, "ServiceService") {
		serviceAlias = strings.TrimSuffix(serviceAlias, "Service")
	}

	// Client interface
	g.P("type ", serviceAlias, " interface {")
	for i, method := range service.Methods {
		g.P(method.Comments.Leading, path+",2,"+string(i))
		g.Annotate(serviceName+"."+method.GoName, method.Location)
		if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		g.P(method.Comments.Leading, generateClientSignature(serviceName, method, g))
	}
	g.P("}")
	g.P()

	// Client struct
	g.P("type ", unexport(serviceAlias), " struct {")
	g.P("c " + g.QualifiedGoIdent(clientPackage.Ident("Client")))
	g.P("name string")
	g.P("}")
	g.P()

	// NewClient factory
	g.P("func New", serviceAlias, " (name string, c "+g.QualifiedGoIdent(clientPackage.Ident("Client"))+") ", serviceAlias, " {")
	g.P("return &", unexport(serviceAlias), "{")
	g.P("c: c,")
	g.P("name: name,")
	g.P("}")
	g.P("}")
	g.P()

	var methodIndex, streamIndex int
	serviceDesc := "_" + serviceName + "_serviceDesc"
	// Client method implementations
	for _, method := range service.Methods {
		var descExpr string
		if !method.Desc.IsStreamingServer() {
			// Unary RPC method
			descExpr = "&" + serviceDesc + ".Methods[" + string(rune(methodIndex)) + "]"
			methodIndex++
		} else {
			// Streaming RPC method
			descExpr = "&" + serviceDesc + ".Methods[" + string(rune(streamIndex)) + "]"
			streamIndex++
		}
		generateClientMethod(gen, file, g, method, serviceName, cCaseSvcName, serviceDesc, descExpr)
	}
}

func generateClientMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, method *protogen.Method, reqSvc, cCaseSvcName, serviceDesc string, descExpr string) {
	reqMethod := cCaseSvcName + "." + method.GoName
	methName := camelCase(method.GoName)
	inType := method.Input.GoIdent
	outType := method.Output.GoIdent

	serviceAlias := cCaseSvcName + "Service"

	// strip suffix
	if strings.HasSuffix(serviceAlias, "ServiceService") {
		serviceAlias = strings.TrimSuffix(serviceAlias, "Service")
	}

	g.P("func (c *", unexport(serviceAlias), ") ", generateClientSignature(cCaseSvcName, method, g), "{")
	if !method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
		g.P(`req := c.c.NewRequest(c.name, `, reqMethod, `", in)`)
		g.P("out := new(", outType, ")")
		g.P("err := ", `c.c.Call(ctx, req, out, opts...)`)
		g.P("if err != nil { return nil, err }")
		g.P("return out, nil")
		g.P("}")
		g.P()
		return
	}
	streamType := unexport(cCaseSvcName) + methName
	g.P(`req := c.c.NewRequest(c.name, "`, reqMethod, `", &`, inType, `{})`)
	g.P("stream, err := c.c.Stream(ctx, req, opts...)")
	g.P("if err != nil { return nil, err }")

	if !method.Desc.IsStreamingClient() {
		g.P("if err := stream.Send(in); err != nil { return nil, err }")
	}

	g.P("return &", streamType, "{stream}, nil")
	g.P("}")
	g.P()

	genSend := method.Desc.IsStreamingClient()
	genRecv := method.Desc.IsStreamingServer()

	// Stream auxiliary types and methods
	g.P("type ", cCaseSvcName, "_", methName, "Service interface {")
	g.P("Context() context.Context")
	g.P("SendMsg(interface{}) error")
	g.P("RecvMsg(interface{}) error")
	g.P("Close() error")

	if genSend {
		g.P("Send(*", inType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", outType, ", error)")
	}
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P("stream", g.QualifiedGoIdent(clientPackage.Ident("Stream")))
	g.P("}")
	g.P()
}

func generateClientSignature(serviceName string, method *protogen.Method, g *protogen.GeneratedFile) string {
	originMethodName := method.GoName
	methodName := camelCase(originMethodName)

	if reversedClientName[methodName] {
		methodName += "_"
	}
	reqArg := ", in *" + g.QualifiedGoIdent(method.Input.GoIdent)
	if !method.Desc.IsStreamingClient() {
		reqArg = ""
	}
	rspName := "*" + g.QualifiedGoIdent(method.Output.GoIdent)
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		rspName = serviceName + "_" + camelCase(originMethodName) + "Service"
	}

	return methodName + "(ctx " + g.QualifiedGoIdent(contextPackage.Ident("Context")) + reqArg + ", opts ..." + g.QualifiedGoIdent(clientPackage.Ident("CallOption")) + ") (" + rspName + ", error)"
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func unexport(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToLower(s[:1] + s[1:])
}
