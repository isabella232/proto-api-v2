// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal_gengo is internal to the protobuf module.
package internal_gengo

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/gogo/protobuf/internal/encoding/tag"
	"github.com/golang/protobuf/v2/protogen"
	"github.com/golang/protobuf/v2/reflect/protoreflect"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// proto package is introduced; the generated code references
// a constant, proto.ProtoPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 3

const (
	fmtPackage      = protogen.GoImportPath("fmt")
	mathPackage     = protogen.GoImportPath("math")
	protoPackage    = protogen.GoImportPath("github.com/golang/protobuf/proto")
	protoapiPackage = protogen.GoImportPath("github.com/golang/protobuf/protoapi")
)

type fileInfo struct {
	*protogen.File
	descriptorVar string // var containing the gzipped FileDescriptorProto

	allEnums         []*protogen.Enum
	allEnumsByPtr    map[*protogen.Enum]int // value is index into allEnums
	allMessages      []*protogen.Message
	allMessagesByPtr map[*protogen.Message]int // value is index into allMessages
	allExtensions    []*protogen.Extension
}

// protoPackage returns the package to import, which is either the protoPackage
// or the protoapiPackage constant.
//
// This special casing exists because we are unable to move InternalMessageInfo
// to protoapi since the implementation behind that logic is heavy and
// too intricately connected to other parts of the proto package.
// The descriptor proto is special in that it avoids using InternalMessageInfo
// so that it is able to depend solely on protoapi and break its dependency
// on the proto package. It is still semantically correct for descriptor to
// avoid using InternalMessageInfo, but it does incur some performance penalty.
// This is acceptable for descriptor, which is a single proto file and is not
// known to be in the hot path for any code.
//
// TODO: Remove this special-casing when the table-driven implementation has
// been ported over to v2.
func (f *fileInfo) protoPackage() protogen.GoImportPath {
	if isDescriptor(f.File) {
		return protoapiPackage
	}
	return protoPackage
}

// GenerateFile generates the contents of a .pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	f := &fileInfo{
		File: file,
	}

	// Collect all enums, messages, and extensions in a breadth-first order.
	f.allEnums = append(f.allEnums, f.Enums...)
	f.allMessages = append(f.allMessages, f.Messages...)
	f.allExtensions = append(f.allExtensions, f.Extensions...)
	walkMessages(f.Messages, func(m *protogen.Message) {
		f.allEnums = append(f.allEnums, m.Enums...)
		f.allMessages = append(f.allMessages, m.Messages...)
		f.allExtensions = append(f.allExtensions, m.Extensions...)
	})

	// Derive a reverse mapping of enum and message pointers to their index
	// in allEnums and allMessages.
	if len(f.allEnums) > 0 {
		f.allEnumsByPtr = make(map[*protogen.Enum]int)
		for i, e := range f.allEnums {
			f.allEnumsByPtr[e] = i
		}
	}
	if len(f.allMessages) > 0 {
		f.allMessagesByPtr = make(map[*protogen.Message]int)
		for i, m := range f.allMessages {
			f.allMessagesByPtr[m] = i
		}
	}

	// Determine the name of the var holding the file descriptor:
	//
	//     fileDescriptor_<hash of filename>
	filenameHash := sha256.Sum256([]byte(f.Desc.Path()))
	f.descriptorVar = fmt.Sprintf("fileDescriptor_%s", hex.EncodeToString(filenameHash[:8]))

	g.P("// Code generated by protoc-gen-go. DO NOT EDIT.")
	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", f.Desc.Path())
	}
	g.P()
	const filePackageField = 2 // FileDescriptorProto.package
	g.PrintLeadingComments(protogen.Location{
		SourceFile: f.Proto.GetName(),
		Path:       []int32{filePackageField},
	})
	g.P()
	g.P("package ", f.GoPackageName)
	g.P()

	if !isDescriptor(file) {
		g.P("// This is a compile-time assertion to ensure that this generated file")
		g.P("// is compatible with the proto package it is being compiled against.")
		g.P("// A compilation error at this line likely means your copy of the")
		g.P("// proto package needs to be updated.")
		g.P("const _ = ", protoPackage.Ident(fmt.Sprintf("ProtoPackageIsVersion%d", generatedCodeVersion)),
			"// please upgrade the proto package")
		g.P()
	}

	for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
		genImport(gen, g, f, imps.Get(i))
	}
	for _, enum := range f.allEnums {
		genEnum(gen, g, f, enum)
	}
	for _, message := range f.allMessages {
		genMessage(gen, g, f, message)
	}
	for _, extension := range f.allExtensions {
		genExtension(gen, g, f, extension)
	}

	genInitFunction(gen, g, f)
	genFileDescriptor(gen, g, f)
	genReflectInitFunction(gen, g, f)
	genReflectFileDescriptor(gen, g, f)
}

// walkMessages calls f on each message and all of its descendants.
func walkMessages(messages []*protogen.Message, f func(*protogen.Message)) {
	for _, m := range messages {
		f(m)
		walkMessages(m.Messages, f)
	}
}

func genImport(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, imp protoreflect.FileImport) {
	impFile, ok := gen.FileByName(imp.Path())
	if !ok {
		return
	}
	if impFile.GoImportPath == f.GoImportPath {
		// Don't generate imports or aliases for types in the same Go package.
		return
	}
	// Generate imports for all non-weak dependencies, even if they are not
	// referenced, because other code and tools depend on having the
	// full transitive closure of protocol buffer types in the binary.
	if !imp.IsWeak {
		g.Import(impFile.GoImportPath)
	}
	if !imp.IsPublic {
		return
	}

	// Generate public imports by generating the imported file, parsing it,
	// and extracting every symbol that should receive a forwarding declaration.
	impGen := gen.NewGeneratedFile("temp.go", impFile.GoImportPath)
	impGen.Skip()
	GenerateFile(gen, impFile, impGen)
	b, err := impGen.Content()
	if err != nil {
		gen.Error(err)
		return
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", b, parser.ParseComments)
	if err != nil {
		gen.Error(err)
		return
	}
	genForward := func(tok token.Token, name string, expr ast.Expr) {
		// Don't import unexported symbols.
		r, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsUpper(r) {
			return
		}
		// Don't import the FileDescriptor.
		if name == impFile.GoDescriptorIdent.GoName {
			return
		}
		// Don't import decls referencing a symbol defined in another package.
		// i.e., don't import decls which are themselves public imports:
		//
		//	type T = somepackage.T
		if _, ok := expr.(*ast.SelectorExpr); ok {
			return
		}
		g.P(tok, " ", name, " = ", impFile.GoImportPath.Ident(name))
	}
	g.P("// Symbols defined in public import of ", imp.Path())
	g.P()
	for _, decl := range astFile.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					genForward(decl.Tok, spec.Name.Name, spec.Type)
				case *ast.ValueSpec:
					for i, name := range spec.Names {
						var expr ast.Expr
						if i < len(spec.Values) {
							expr = spec.Values[i]
						}
						genForward(decl.Tok, name.Name, expr)
					}
				case *ast.ImportSpec:
				default:
					panic(fmt.Sprintf("can't generate forward for spec type %T", spec))
				}
			}
		}
	}
	g.P()
}

func genFileDescriptor(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo) {
	// Trim the source_code_info from the descriptor.
	// Marshal and gzip it.
	descProto := proto.Clone(f.Proto).(*descriptorpb.FileDescriptorProto)
	descProto.SourceCodeInfo = nil
	b, err := proto.Marshal(descProto)
	if err != nil {
		gen.Error(err)
		return
	}
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	b = buf.Bytes()

	g.P("var ", f.descriptorVar, " = []byte{")
	g.P("// ", len(b), " bytes of a gzipped FileDescriptorProto")
	for len(b) > 0 {
		n := 16
		if n > len(b) {
			n = len(b)
		}

		s := ""
		for _, c := range b[:n] {
			s += fmt.Sprintf("0x%02x,", c)
		}
		g.P(s)

		b = b[n:]
	}
	g.P("}")
	g.P()
}

func genEnum(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, enum *protogen.Enum) {
	g.PrintLeadingComments(enum.Location)
	g.Annotate(enum.GoIdent.GoName, enum.Location)
	g.P("type ", enum.GoIdent, " int32",
		deprecationComment(enum.Desc.Options().(*descriptorpb.EnumOptions).GetDeprecated()))
	g.P("const (")
	for _, value := range enum.Values {
		g.PrintLeadingComments(value.Location)
		g.Annotate(value.GoIdent.GoName, value.Location)
		g.P(value.GoIdent, " ", enum.GoIdent, " = ", value.Desc.Number(),
			deprecationComment(value.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated()))
	}
	g.P(")")
	g.P()

	// Generate support for protobuf reflection.
	genReflectEnum(gen, g, f, enum)

	nameMap := enum.GoIdent.GoName + "_name"
	g.P("var ", nameMap, " = map[int32]string{")
	generated := make(map[protoreflect.EnumNumber]bool)
	for _, value := range enum.Values {
		duplicate := ""
		if _, present := generated[value.Desc.Number()]; present {
			duplicate = "// Duplicate value: "
		}
		g.P(duplicate, value.Desc.Number(), ": ", strconv.Quote(string(value.Desc.Name())), ",")
		generated[value.Desc.Number()] = true
	}
	g.P("}")
	g.P()
	valueMap := enum.GoIdent.GoName + "_value"
	g.P("var ", valueMap, " = map[string]int32{")
	for _, value := range enum.Values {
		g.P(strconv.Quote(string(value.Desc.Name())), ": ", value.Desc.Number(), ",")
	}
	g.P("}")
	g.P()
	if enum.Desc.Syntax() != protoreflect.Proto3 {
		g.P("func (x ", enum.GoIdent, ") Enum() *", enum.GoIdent, " {")
		g.P("p := new(", enum.GoIdent, ")")
		g.P("*p = x")
		g.P("return p")
		g.P("}")
		g.P()
	}
	g.P("func (x ", enum.GoIdent, ") String() string {")
	g.P("return ", f.protoPackage().Ident("EnumName"), "(", enum.GoIdent, "_name, int32(x))")
	g.P("}")
	g.P()

	if enum.Desc.Syntax() == protoreflect.Proto2 {
		g.P("func (x *", enum.GoIdent, ") UnmarshalJSON(data []byte) error {")
		g.P("value, err := ", f.protoPackage().Ident("UnmarshalJSONEnum"), "(", enum.GoIdent, `_value, data, "`, enum.GoIdent, `")`)
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("*x = ", enum.GoIdent, "(value)")
		g.P("return nil")
		g.P("}")
		g.P()
	}

	var indexes []string
	for i := 1; i < len(enum.Location.Path); i += 2 {
		indexes = append(indexes, strconv.Itoa(int(enum.Location.Path[i])))
	}
	g.P("func (", enum.GoIdent, ") EnumDescriptor() ([]byte, []int) {")
	g.P("return ", f.descriptorVar, ", []int{", strings.Join(indexes, ","), "}")
	g.P("}")
	g.P()

	genWellKnownType(g, "", enum.GoIdent, enum.Desc)
}

// enumRegistryName returns the name used to register an enum with the proto
// package registry.
//
// Confusingly, this is <proto_package>.<go_ident>. This probably should have
// been the full name of the proto enum type instead, but changing it at this
// point would require thought.
func enumRegistryName(enum *protogen.Enum) string {
	// Find the FileDescriptor for this enum.
	var desc protoreflect.Descriptor = enum.Desc
	for {
		p, ok := desc.Parent()
		if !ok {
			break
		}
		desc = p
	}
	fdesc := desc.(protoreflect.FileDescriptor)
	if fdesc.Package() == "" {
		return enum.GoIdent.GoName
	}
	return string(fdesc.Package()) + "." + enum.GoIdent.GoName
}

func genMessage(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, message *protogen.Message) {
	if message.Desc.IsMapEntry() {
		return
	}

	hasComment := g.PrintLeadingComments(message.Location)
	if message.Desc.Options().(*descriptorpb.MessageOptions).GetDeprecated() {
		if hasComment {
			g.P("//")
		}
		g.P(deprecationComment(true))
	}
	g.Annotate(message.GoIdent.GoName, message.Location)
	g.P("type ", message.GoIdent, " struct {")
	for _, field := range message.Fields {
		if field.OneofType != nil {
			// It would be a bit simpler to iterate over the oneofs below,
			// but generating the field here keeps the contents of the Go
			// struct in the same order as the contents of the source
			// .proto file.
			if field == field.OneofType.Fields[0] {
				genOneofField(gen, g, f, message, field.OneofType)
			}
			continue
		}
		g.PrintLeadingComments(field.Location)
		goType, pointer := fieldGoType(g, field)
		if pointer {
			goType = "*" + goType
		}
		tags := []string{
			fmt.Sprintf("protobuf:%q", fieldProtobufTag(field)),
			fmt.Sprintf("json:%q", fieldJSONTag(field)),
		}
		if field.Desc.IsMap() {
			key := field.MessageType.Fields[0]
			val := field.MessageType.Fields[1]
			tags = append(tags,
				fmt.Sprintf("protobuf_key:%q", fieldProtobufTag(key)),
				fmt.Sprintf("protobuf_val:%q", fieldProtobufTag(val)),
			)
		}
		g.Annotate(message.GoIdent.GoName+"."+field.GoName, field.Location)
		g.P(field.GoName, " ", goType, " `", strings.Join(tags, " "), "`",
			deprecationComment(field.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated()))
	}
	g.P("XXX_NoUnkeyedLiteral struct{} `json:\"-\"`")

	if message.Desc.ExtensionRanges().Len() > 0 {
		var tags []string
		if message.Desc.Options().(*descriptorpb.MessageOptions).GetMessageSetWireFormat() {
			tags = append(tags, `protobuf_messageset:"1"`)
		}
		tags = append(tags, `json:"-"`)
		g.P(f.protoPackage().Ident("XXX_InternalExtensions"), " `", strings.Join(tags, " "), "`")
	}
	g.P("XXX_unrecognized []byte `json:\"-\"`")
	g.P("XXX_sizecache int32 `json:\"-\"`")
	g.P("}")
	g.P()

	// Generate support for protobuf reflection.
	genReflectMessage(gen, g, f, message)

	// Reset
	g.P("func (m *", message.GoIdent, ") Reset() { *m = ", message.GoIdent, "{} }")
	// String
	g.P("func (m *", message.GoIdent, ") String() string { return ", f.protoPackage().Ident("CompactTextString"), "(m) }")
	// ProtoMessage
	g.P("func (*", message.GoIdent, ") ProtoMessage() {}")
	// Descriptor
	var indexes []string
	for i := 1; i < len(message.Location.Path); i += 2 {
		indexes = append(indexes, strconv.Itoa(int(message.Location.Path[i])))
	}
	g.P("func (*", message.GoIdent, ") Descriptor() ([]byte, []int) {")
	g.P("return ", f.descriptorVar, ", []int{", strings.Join(indexes, ","), "}")
	g.P("}")
	g.P()

	// ExtensionRangeArray
	if extranges := message.Desc.ExtensionRanges(); extranges.Len() > 0 {
		protoExtRange := f.protoPackage().Ident("ExtensionRange")
		extRangeVar := "extRange_" + message.GoIdent.GoName
		g.P("var ", extRangeVar, " = []", protoExtRange, " {")
		for i := 0; i < extranges.Len(); i++ {
			r := extranges.Get(i)
			g.P("{Start:", r[0], ", End:", r[1]-1 /* inclusive */, "},")
		}
		g.P("}")
		g.P()
		g.P("func (*", message.GoIdent, ") ExtensionRangeArray() []", protoExtRange, " {")
		g.P("return ", extRangeVar)
		g.P("}")
		g.P()
	}

	genWellKnownType(g, "*", message.GoIdent, message.Desc)

	// Table-driven proto support.
	//
	// TODO: It does not scale to keep adding another method for every
	// operation on protos that we want to switch over to using the
	// table-driven approach. Instead, we should only add a single method
	// that allows getting access to the *InternalMessageInfo struct and then
	// calling Unmarshal, Marshal, Merge, Size, and Discard directly on that.
	if !isDescriptor(f.File) {
		// NOTE: We avoid adding table-driven support for descriptor proto
		// since this depends on the v1 proto package, which would eventually
		// need to depend on the descriptor itself.
		messageInfoVar := "xxx_messageInfo_" + message.GoIdent.GoName
		// XXX_Unmarshal
		g.P("func (m *", message.GoIdent, ") XXX_Unmarshal(b []byte) error {")
		g.P("return ", messageInfoVar, ".Unmarshal(m, b)")
		g.P("}")
		// XXX_Marshal
		g.P("func (m *", message.GoIdent, ") XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {")
		g.P("return ", messageInfoVar, ".Marshal(b, m, deterministic)")
		g.P("}")
		// XXX_Merge
		g.P("func (m *", message.GoIdent, ") XXX_Merge(src proto.Message) {")
		g.P(messageInfoVar, ".Merge(m, src)")
		g.P("}")
		// XXX_Size
		g.P("func (m *", message.GoIdent, ") XXX_Size() int {")
		g.P("return ", messageInfoVar, ".Size(m)")
		g.P("}")
		// XXX_DiscardUnknown
		g.P("func (m *", message.GoIdent, ") XXX_DiscardUnknown() {")
		g.P(messageInfoVar, ".DiscardUnknown(m)")
		g.P("}")
		g.P()
		g.P("var ", messageInfoVar, " ", protoPackage.Ident("InternalMessageInfo"))
		g.P()
	}

	// Constants and vars holding the default values of fields.
	for _, field := range message.Fields {
		if !field.Desc.HasDefault() {
			continue
		}
		defVarName := "Default_" + message.GoIdent.GoName + "_" + field.GoName
		def := field.Desc.Default()
		switch field.Desc.Kind() {
		case protoreflect.StringKind:
			g.P("const ", defVarName, " string = ", strconv.Quote(def.String()))
		case protoreflect.BytesKind:
			g.P("var ", defVarName, " []byte = []byte(", strconv.Quote(string(def.Bytes())), ")")
		case protoreflect.EnumKind:
			evalueDesc := field.Desc.DefaultEnumValue()
			enum := field.EnumType
			evalue := enum.Values[evalueDesc.Index()]
			g.P("const ", defVarName, " ", field.EnumType.GoIdent, " = ", evalue.GoIdent)
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			// Floating point numbers need extra handling for -Inf/Inf/NaN.
			f := field.Desc.Default().Float()
			goType := "float64"
			if field.Desc.Kind() == protoreflect.FloatKind {
				goType = "float32"
			}
			// funcCall returns a call to a function in the math package,
			// possibly converting the result to float32.
			funcCall := func(fn, param string) string {
				s := g.QualifiedGoIdent(mathPackage.Ident(fn)) + param
				if goType != "float64" {
					s = goType + "(" + s + ")"
				}
				return s
			}
			switch {
			case math.IsInf(f, -1):
				g.P("var ", defVarName, " ", goType, " = ", funcCall("Inf", "(-1)"))
			case math.IsInf(f, 1):
				g.P("var ", defVarName, " ", goType, " = ", funcCall("Inf", "(1)"))
			case math.IsNaN(f):
				g.P("var ", defVarName, " ", goType, " = ", funcCall("NaN", "()"))
			default:
				g.P("const ", defVarName, " ", goType, " = ", field.Desc.Default().Interface())
			}
		default:
			goType, _ := fieldGoType(g, field)
			g.P("const ", defVarName, " ", goType, " = ", def.Interface())
		}
	}
	g.P()

	// Getters.
	for _, field := range message.Fields {
		if field.OneofType != nil {
			if field == field.OneofType.Fields[0] {
				genOneofTypes(gen, g, f, message, field.OneofType)
			}
		}
		goType, pointer := fieldGoType(g, field)
		defaultValue := fieldDefaultValue(g, message, field)
		if field.Desc.Options().(*descriptorpb.FieldOptions).GetDeprecated() {
			g.P(deprecationComment(true))
		}
		g.Annotate(message.GoIdent.GoName+".Get"+field.GoName, field.Location)
		g.P("func (m *", message.GoIdent, ") Get", field.GoName, "() ", goType, " {")
		if field.OneofType != nil {
			g.P("if x, ok := m.Get", field.OneofType.GoName, "().(*", fieldOneofType(field), "); ok {")
			g.P("return x.", field.GoName)
			g.P("}")
		} else {
			if field.Desc.Syntax() == protoreflect.Proto3 || defaultValue == "nil" {
				g.P("if m != nil {")
			} else {
				g.P("if m != nil && m.", field.GoName, " != nil {")
			}
			star := ""
			if pointer {
				star = "*"
			}
			g.P("return ", star, " m.", field.GoName)
			g.P("}")
		}
		g.P("return ", defaultValue)
		g.P("}")
		g.P()
	}

	if len(message.Oneofs) > 0 {
		genOneofWrappers(gen, g, f, message)
	}
}

// fieldGoType returns the Go type used for a field.
//
// If it returns pointer=true, the struct field is a pointer to the type.
func fieldGoType(g *protogen.GeneratedFile, field *protogen.Field) (goType string, pointer bool) {
	pointer = true
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.EnumType.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if field.Desc.IsMap() {
			keyType, _ := fieldGoType(g, field.MessageType.Fields[0])
			valType, _ := fieldGoType(g, field.MessageType.Fields[1])
			return fmt.Sprintf("map[%v]%v", keyType, valType), false
		}
		goType = "*" + g.QualifiedGoIdent(field.MessageType.GoIdent)
		pointer = false
	}
	if field.Desc.Cardinality() == protoreflect.Repeated {
		goType = "[]" + goType
		pointer = false
	}
	// Extension fields always have pointer type, even when defined in a proto3 file.
	if field.Desc.Syntax() == protoreflect.Proto3 && field.Desc.ExtendedType() == nil {
		pointer = false
	}
	return goType, pointer
}

func fieldProtobufTag(field *protogen.Field) string {
	var enumName string
	if field.Desc.Kind() == protoreflect.EnumKind {
		enumName = enumRegistryName(field.EnumType)
	}
	return tag.Marshal(field.Desc, enumName)
}

func fieldDefaultValue(g *protogen.GeneratedFile, message *protogen.Message, field *protogen.Field) string {
	if field.Desc.Cardinality() == protoreflect.Repeated {
		return "nil"
	}
	if field.Desc.HasDefault() {
		defVarName := "Default_" + message.GoIdent.GoName + "_" + field.GoName
		if field.Desc.Kind() == protoreflect.BytesKind {
			return "append([]byte(nil), " + defVarName + "...)"
		}
		return defVarName
	}
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return `""`
	case protoreflect.MessageKind, protoreflect.GroupKind, protoreflect.BytesKind:
		return "nil"
	case protoreflect.EnumKind:
		return g.QualifiedGoIdent(field.EnumType.Values[0].GoIdent)
	default:
		return "0"
	}
}

func fieldJSONTag(field *protogen.Field) string {
	return string(field.Desc.Name()) + ",omitempty"
}

func genExtension(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, extension *protogen.Extension) {
	// Special case for proto2 message sets: If this extension is extending
	// proto2.bridge.MessageSet, and its final name component is "message_set_extension",
	// then drop that last component.
	//
	// TODO: This should be implemented in the text formatter rather than the generator.
	// In addition, the situation for when to apply this special case is implemented
	// differently in other languages:
	// https://github.com/google/protobuf/blob/aff10976/src/google/protobuf/text_format.cc#L1560
	name := extension.Desc.FullName()
	if n, ok := isExtensionMessageSetElement(extension); ok {
		name = n
	}

	g.P("var ", extensionVar(f.File, extension), " = &", f.protoPackage().Ident("ExtensionDesc"), "{")
	g.P("ExtendedType: (*", extension.ExtendedType.GoIdent, ")(nil),")
	goType, pointer := fieldGoType(g, extension)
	if pointer {
		goType = "*" + goType
	}
	g.P("ExtensionType: (", goType, ")(nil),")
	g.P("Field: ", extension.Desc.Number(), ",")
	g.P("Name: ", strconv.Quote(string(name)), ",")
	g.P("Tag: ", strconv.Quote(fieldProtobufTag(extension)), ",")
	g.P("Filename: ", strconv.Quote(f.Desc.Path()), ",")
	g.P("}")
	g.P()
}

// isExtensionMessageSetELement returns the adjusted name of an extension
// which extends proto2.bridge.MessageSet.
func isExtensionMessageSetElement(extension *protogen.Extension) (name protoreflect.FullName, ok bool) {
	opts := extension.ExtendedType.Desc.Options().(*descriptorpb.MessageOptions)
	if !opts.GetMessageSetWireFormat() || extension.Desc.Name() != "message_set_extension" {
		return "", false
	}
	if extension.ParentMessage == nil {
		// This case shouldn't be given special handling at all--we're
		// only supposed to drop the ".message_set_extension" for
		// extensions defined within a message (i.e., the extension
		// takes the message's name).
		//
		// This matches the behavior of the v1 generator, however.
		//
		// TODO: See if we can drop this case.
		name = extension.Desc.FullName()
		name = name[:len(name)-len("message_set_extension")]
		return name, true
	}
	return extension.Desc.FullName().Parent(), true
}

// extensionVar returns the var holding the ExtensionDesc for an extension.
func extensionVar(f *protogen.File, extension *protogen.Extension) protogen.GoIdent {
	name := "E_"
	if extension.ParentMessage != nil {
		name += extension.ParentMessage.GoIdent.GoName + "_"
	}
	name += extension.GoName
	return f.GoImportPath.Ident(name)
}

// genInitFunction generates an init function that registers the types in the
// generated file with the proto package.
func genInitFunction(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo) {
	g.P("func init() {")
	g.P(f.protoPackage().Ident("RegisterFile"), "(", strconv.Quote(f.Desc.Path()), ", ", f.descriptorVar, ")")
	for _, enum := range f.allEnums {
		name := enum.GoIdent.GoName
		g.P(f.protoPackage().Ident("RegisterEnum"), fmt.Sprintf("(%q, %s_name, %s_value)", enumRegistryName(enum), name, name))
	}
	for _, message := range f.allMessages {
		if message.Desc.IsMapEntry() {
			continue
		}

		name := message.GoIdent.GoName
		g.P(f.protoPackage().Ident("RegisterType"), fmt.Sprintf("((*%s)(nil), %q)", name, message.Desc.FullName()))

		// Types of map fields, sorted by the name of the field message type.
		var mapFields []*protogen.Field
		for _, field := range message.Fields {
			if field.Desc.IsMap() {
				mapFields = append(mapFields, field)
			}
		}
		sort.Slice(mapFields, func(i, j int) bool {
			ni := mapFields[i].MessageType.Desc.FullName()
			nj := mapFields[j].MessageType.Desc.FullName()
			return ni < nj
		})
		for _, field := range mapFields {
			typeName := string(field.MessageType.Desc.FullName())
			goType, _ := fieldGoType(g, field)
			g.P(f.protoPackage().Ident("RegisterMapType"), fmt.Sprintf("((%v)(nil), %q)", goType, typeName))
		}
	}
	for _, extension := range f.allExtensions {
		g.P(f.protoPackage().Ident("RegisterExtension"), "(", extensionVar(f.File, extension), ")")
	}
	g.P("}")
	g.P()
}

// deprecationComment returns a standard deprecation comment if deprecated is true.
func deprecationComment(deprecated bool) string {
	if !deprecated {
		return ""
	}
	return "// Deprecated: Do not use."
}

func genWellKnownType(g *protogen.GeneratedFile, ptr string, ident protogen.GoIdent, desc protoreflect.Descriptor) {
	if wellKnownTypes[desc.FullName()] {
		g.P("func (", ptr, ident, `) XXX_WellKnownType() string { return "`, desc.Name(), `" }`)
		g.P()
	}
}

// Names of messages and enums for which we will generate XXX_WellKnownType methods.
var wellKnownTypes = map[protoreflect.FullName]bool{
	"google.protobuf.Any":       true,
	"google.protobuf.Duration":  true,
	"google.protobuf.Empty":     true,
	"google.protobuf.Struct":    true,
	"google.protobuf.Timestamp": true,

	"google.protobuf.BoolValue":   true,
	"google.protobuf.BytesValue":  true,
	"google.protobuf.DoubleValue": true,
	"google.protobuf.FloatValue":  true,
	"google.protobuf.Int32Value":  true,
	"google.protobuf.Int64Value":  true,
	"google.protobuf.ListValue":   true,
	"google.protobuf.NullValue":   true,
	"google.protobuf.StringValue": true,
	"google.protobuf.UInt32Value": true,
	"google.protobuf.UInt64Value": true,
	"google.protobuf.Value":       true,
}

// genOneofField generates the struct field for a oneof.
func genOneofField(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, message *protogen.Message, oneof *protogen.Oneof) {
	if g.PrintLeadingComments(oneof.Location) {
		g.P("//")
	}
	g.P("// Types that are valid to be assigned to ", oneofFieldName(oneof), ":")
	for _, field := range oneof.Fields {
		g.PrintLeadingComments(field.Location)
		g.P("//\t*", fieldOneofType(field))
	}
	g.Annotate(message.GoIdent.GoName+"."+oneofFieldName(oneof), oneof.Location)
	g.P(oneofFieldName(oneof), " ", oneofInterfaceName(oneof), " `protobuf_oneof:\"", oneof.Desc.Name(), "\"`")
}

// genOneofTypes generates the interface type used for a oneof field,
// and the wrapper types that satisfy that interface.
//
// It also generates the getter method for the parent oneof field
// (but not the member fields).
func genOneofTypes(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, message *protogen.Message, oneof *protogen.Oneof) {
	ifName := oneofInterfaceName(oneof)
	g.P("type ", ifName, " interface {")
	g.P(ifName, "()")
	g.P("}")
	g.P()
	for _, field := range oneof.Fields {
		name := fieldOneofType(field)
		g.Annotate(name.GoName, field.Location)
		g.Annotate(name.GoName+"."+field.GoName, field.Location)
		g.P("type ", name, " struct {")
		goType, _ := fieldGoType(g, field)
		tags := []string{
			fmt.Sprintf("protobuf:%q", fieldProtobufTag(field)),
		}
		g.P(field.GoName, " ", goType, " `", strings.Join(tags, " "), "`")
		g.P("}")
		g.P()
	}
	for _, field := range oneof.Fields {
		g.P("func (*", fieldOneofType(field), ") ", ifName, "() {}")
		g.P()
	}
	g.Annotate(message.GoIdent.GoName+".Get"+oneof.GoName, oneof.Location)
	g.P("func (m *", message.GoIdent.GoName, ") Get", oneof.GoName, "() ", ifName, " {")
	g.P("if m != nil {")
	g.P("return m.", oneofFieldName(oneof))
	g.P("}")
	g.P("return nil")
	g.P("}")
	g.P()
}

// oneofFieldName returns the name of the struct field holding the oneof value.
//
// This function is trivial, but pulling out the name like this makes it easier
// to experiment with alternative oneof implementations.
func oneofFieldName(oneof *protogen.Oneof) string {
	return oneof.GoName
}

// oneofInterfaceName returns the name of the interface type implemented by
// the oneof field value types.
func oneofInterfaceName(oneof *protogen.Oneof) string {
	return fmt.Sprintf("is%s_%s", oneof.ParentMessage.GoIdent.GoName, oneof.GoName)
}

// genOneofWrappers generates the XXX_OneofWrappers method for a message.
func genOneofWrappers(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, message *protogen.Message) {
	g.P("// XXX_OneofWrappers is for the internal use of the proto package.")
	g.P("func (*", message.GoIdent.GoName, ") XXX_OneofWrappers() []interface{} {")
	g.P("return []interface{}{")
	for _, oneof := range message.Oneofs {
		for _, field := range oneof.Fields {
			g.P("(*", fieldOneofType(field), ")(nil),")
		}
	}
	g.P("}")
	g.P("}")
	g.P()
}

// fieldOneofType returns the wrapper type used to represent a field in a oneof.
func fieldOneofType(field *protogen.Field) protogen.GoIdent {
	ident := protogen.GoIdent{
		GoImportPath: field.ParentMessage.GoIdent.GoImportPath,
		GoName:       field.ParentMessage.GoIdent.GoName + "_" + field.GoName,
	}
	// Check for collisions with nested messages or enums.
	//
	// This conflict resolution is incomplete: Among other things, it
	// does not consider collisions with other oneof field types.
	//
	// TODO: Consider dropping this entirely. Detecting conflicts and
	// producing an error is almost certainly better than permuting
	// field and type names in mostly unpredictable ways.
Loop:
	for {
		for _, message := range field.ParentMessage.Messages {
			if message.GoIdent == ident {
				ident.GoName += "_"
				continue Loop
			}
		}
		for _, enum := range field.ParentMessage.Enums {
			if enum.GoIdent == ident {
				ident.GoName += "_"
				continue Loop
			}
		}
		return ident
	}
}
