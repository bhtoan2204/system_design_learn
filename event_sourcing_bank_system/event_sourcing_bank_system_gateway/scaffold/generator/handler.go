package generator

import (
	"event_sourcing_bank_system_gateway/scaffold/model"
	"event_sourcing_bank_system_gateway/scaffold/parser"
	"event_sourcing_bank_system_gateway/scaffold/util"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const HANDLER_PARENT_PATH = "application/routing/delivery/handler"
const ROUTING_REGISTRY_PATH = "application/routing/delivery/registry.go"

func GenerateHandlers(protos []*model.Proto) ([]*model.Handler, []string, error) {
	// because messages can be defined in a separate proto file
	// we need to create a map of all messages and use it for generating handler for each service
	messages := parser.MergeProtoMessages(protos...)
	var allHandlers []*model.Handler
	handlerFolderPaths := make(map[string]bool) // use a map to collect non-duplicated handler folder paths
	for _, proto := range protos {
		handlers, err := generateHandlers(proto, messages)
		if err != nil {
			fmt.Println("failed to generate handler for:", proto.FilePath, err)
		} else {
			allHandlers = append(allHandlers, handlers...)
			for _, handler := range handlers {
				handlerFolderPaths[handler.FolderPath] = true
			}
		}
	}

	folderPaths := make([]string, 0, len(handlerFolderPaths))
	for path := range handlerFolderPaths {
		folderPaths = append(folderPaths, path)
	}

	return allHandlers, folderPaths, nil
}

func generateHandlers(proto *model.Proto, messages map[string]map[string]*model.Message) ([]*model.Handler, error) {
	// proto.FilePath must follow this format: proto/something/something.proto
	var segments []string
	if runtime.GOOS == "windows" {
		segments = strings.Split(proto.FilePath, "\\")
	} else {
		segments = strings.Split(proto.FilePath, "/")
	}

	if len(segments) < 3 {
		return nil, fmt.Errorf("unsupported proto file path: %s", proto.FilePath)
	}
	parentFolderPath := HANDLER_PARENT_PATH

	//skip: proto(root) and .proto(file)
	for i := 1; i < len(segments)-1; i++ {
		folderPath := parentFolderPath + "/" + segments[i]
		if proto.Package == "validate" {
			break
		}
		err := util.CreateFolderIfNotExist(folderPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new folder at: %s, skip generating handler, error: %e", folderPath, err)
		}
		parentFolderPath = folderPath
	}

	if segments[len(segments)-2] != proto.Package {
		parentFolderPath += "/" + proto.Package
		err := util.CreateFolderIfNotExist(parentFolderPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new folder at: %s, skip generating handler, error: %e", parentFolderPath, err)
		}
		fmt.Printf("create new directory if not exist = %s \n", parentFolderPath)
	}

	handlers := []*model.Handler{}
	generationErrors := []error{}
	for _, service := range proto.Services {
		for _, api := range service.Apis {

			// TODO: optimize this paths
			paths := strings.Split(service.Path, ",")
			handler := model.Handler{
				Name:   api.Name,
				Method: api.Method,
				Accept: api.Accept,
				// Paths:         paths[0] + api.Path,
				Permission:   api.Permission,
				ServiceName:  service.Name,
				ProtoFolder:  strings.Join(segments[1:len(segments)-1], "/"),
				ProtoPackage: service.Package,
				RequestType:  api.RequestType,
				ResponseType: api.ResponseType,
				FolderPath:   parentFolderPath,
			}
			for _, path := range paths {
				handler.Paths = append(handler.Paths, path+api.Path)
			}
			var err error
			switch api.Method {
			case "GET":
				msgReqType := findMessage(api.RequestType, proto, messages)
				err = generateGetHandler(&handler, msgReqType)
			case "POST":
				msgReqType := findMessage(api.RequestType, proto, messages)
				err = generatePostHandler(&handler, msgReqType)
			case "PUT":
				msgReqType := findMessage(api.RequestType, proto, messages)
				err = generatePutHandler(&handler, msgReqType)
			default:
				err = fmt.Errorf("failed to generate handler for: %s, unsupported method: %s", api.Name, api.Method)
			}
			if err != nil {
				generationErrors = append(generationErrors, err)
			} else {
				handlers = append(handlers, &handler)
			}
		}
	}

	var err error
	if len(generationErrors) > 0 {
		err = generationErrors[0]
		for _, e := range generationErrors[1:] {
			err = fmt.Errorf("%v; %v", err, e)
		}

		return nil, err
	}

	return handlers, nil
}

func findMessage(reqType string, proto *model.Proto, messages map[string]map[string]*model.Message) *model.Message {
	if val, ok := messages[proto.FilePath][reqType]; ok {
		return val
	}

	for _, msgMap := range messages {
		if val, ok := msgMap[reqType]; ok {
			return val
		}
	}

	return nil
}

func generatePostHandler(handler *model.Handler, requestType *model.Message) error {
	var templatePath string
	if handler.Accept == "multipart/form-data" {
		templatePath = "scaffold/template/handler_post_file.tmpl"
	} else if handler.Type == "webhook" {
		templatePath = "scaffold/template/handler_post_webhook.tmpl"
	} else {
		templatePath = "scaffold/template/handler_post.tmpl"
	}

	template, err := util.ReadFileAsString(templatePath)
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<request_type_lower_1st>", util.LowerFirstChar(handler.RequestType))
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	if strings.EqualFold(handler.ResponseType, "google.protobuf.Empty") {
		template = strings.ReplaceAll(template, "<ext_import>", `_ "github.com/golang/protobuf/ptypes/empty"`)
		template = strings.ReplaceAll(template, "<proto_custom_package>", "empty")
		template = strings.ReplaceAll(template, "<response_type>", "Empty")
	} else {
		template = strings.ReplaceAll(template, "<ext_import>", "")
		template = strings.ReplaceAll(template, "<proto_custom_package>", handler.ProtoPackage)
		template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	}
	if len(handler.Paths) == 2 {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[1])
	} else {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[0])
	}
	if model.APIPublic[handler.Paths[0]] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}
	if len(handler.Accept) > 0 {
		template = strings.ReplaceAll(template, "<accept>", handler.Accept)
		if handler.Accept == "multipart/form-data" {
			template = strings.ReplaceAll(template, "bind", "1")
		}
	} else {
		template = strings.ReplaceAll(template, "<accept>", "json")
	}
	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)

	// handle params
	paramDecTemplate, err := util.ReadFileAsString("scaffold/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	paramParsingIdTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_id.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingFileTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_file.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingEnumTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_enum.tmpl")
	if err != nil {
		fmt.Println("failed to read template enum:", err)
		return err
	}

	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	for _, field := range requestType.Fields {
		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else if handler.Accept == "multipart/form-data" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "formData")
			if field.Location == "file" {
				declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
			} else {
				declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
			}
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "body")
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}

		field_type := field.Type
		if strings.Contains(field_type, "Request") {
			field_type = handler.ProtoPackage + "." + field_type
		} else {
			switch field_type {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "byte", "rune", "string", "bool", "enum":
			default:
				field_type = handler.ProtoPackage + "." + field_type
			}
		}
		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "[]"+field_type)
		} else if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
		} else if field.Location == "file" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "file")
			// } else if field.Location == "sub-message" {
			// 	ptype := handler.ProtoPackage + "." + strings.Trim(field_type, "*")
			// 	declaration = strings.ReplaceAll(declaration, "<param_type>", ptype)
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field_type)
		}

		paramDeclarations = append(paramDeclarations, declaration)

		if field.Location == "path" {
			parsing := paramParsingIdTemplate
			hasIntField = true

			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			parsing = strings.ReplaceAll(parsing, "<param_name_cammel_case>", util.SnakeToCamel(field.Name))
			paramContext := "Query"
			if field.Location == "path" {
				paramContext = "Param"
			}
			parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
			parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
			paramParsings = append(paramParsings, parsing)
		} else if field.Location == "file" {
			parsing := paramParsingFileTemplate
			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			parsing = strings.ReplaceAll(parsing, "<param_name_cammel_case>", util.SnakeToCamel(field.Name))
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))

			paramParsings = append(paramParsings, parsing)
		} else if field.Location == "enum" {
			parsing := paramParsingEnumTemplate

			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			parsing = strings.ReplaceAll(parsing, "<param_enum_object>", fmt.Sprintf("%s.%s", handler.ProtoPackage, field.Type))
			parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
			parsing = strings.ReplaceAll(parsing, "<param_context>", "Query")
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
			paramParsings = append(paramParsings, parsing)
		}
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}
	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		template = strings.ReplaceAll(template, "\"strconv\"", "")
	}

	fpath := handler.FolderPath + "/" + util.ToSnakeCase(handlerNameLower1st) + ".go"

	err = os.WriteFile(fpath, []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func generatePutHandler(handler *model.Handler, requestType *model.Message) error {
	template, err := util.ReadFileAsString("scaffold/template/handler_put.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<request_type_lower_1st>", util.LowerFirstChar(handler.RequestType))
	template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	if len(handler.Paths) == 2 {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[1])
	} else {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[0])
	}
	if model.APIPublic[handler.Paths[0]] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}

	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)
	// handle params
	paramDecTemplate, err := util.ReadFileAsString("scaffold/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingIdTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_id.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	for _, field := range requestType.Fields {
		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "body")
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}
		field_type := field.Type
		if strings.Contains(field_type, "Request") {
			field_type = handler.ProtoPackage + "." + field_type
		} else {
			switch field_type {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "byte", "rune", "string", "bool", "enum":
			default:
				field_type = handler.ProtoPackage + "." + field_type
			}
		}
		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "[]"+field_type)
		} else if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
		} else if field.Location == "file" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "file")
			// } else if field.Location == "sub-message" {
			// 	ptype := handler.ProtoPackage + "." + strings.Trim(field_type, "*")
			// 	declaration = strings.ReplaceAll(declaration, "<param_type>", ptype)
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field_type)
		}

		paramDeclarations = append(paramDeclarations, declaration)

		if field.Location == "path" {
			parsing := paramParsingIdTemplate
			hasIntField = true

			parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
			parsing = strings.ReplaceAll(parsing, "<param_name_cammel_case>", util.SnakeToCamel(field.Name))
			paramContext := "Query"
			if field.Location == "path" {
				paramContext = "Param"
			}
			parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
			parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
			parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
			paramParsings = append(paramParsings, parsing)
		}
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}
	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		template = strings.ReplaceAll(template, "\"strconv\"", "")
	}

	err = os.WriteFile(handler.FolderPath+"/"+util.ToSnakeCase(handlerNameLower1st)+".go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func generateGetHandler(handler *model.Handler, requestType *model.Message) error {
	template, err := util.ReadFileAsString("scaffold/template/handler_get.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	template = strings.ReplaceAll(template, "<proto_folder>", handler.ProtoFolder)
	template = strings.ReplaceAll(template, "<proto_package>", handler.ProtoPackage)
	template = strings.ReplaceAll(template, "<handler_name>", handler.Name)
	handlerNameLower1st := util.LowerFirstChar(handler.Name)
	template = strings.ReplaceAll(template, "<handler_name_lower_1st>", handlerNameLower1st)
	template = strings.ReplaceAll(template, "<request_type>", handler.RequestType)
	template = strings.ReplaceAll(template, "<response_type>", handler.ResponseType)
	if len(handler.Paths) == 2 {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[1])
	} else {
		template = strings.ReplaceAll(template, "<path>", handler.Paths[0])
	}

	if model.APIPublic[handler.Paths[0]] {
		template = strings.ReplaceAll(template, "<summary>", "api public")
	} else {
		template = strings.ReplaceAll(template, "<summary>", "permission: "+handler.Permission)
	}
	template = strings.ReplaceAll(template, "<tags>", handler.ServiceName)

	// handle params
	paramDecTemplate, err := util.ReadFileAsString("scaffold/template/param_declaration.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingIntTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_int.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingArrIntTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_arr_int.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingStringTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_string.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingEnumTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_enum.tmpl")
	if err != nil {
		fmt.Println("failed to read template enum:", err)
		return err
	}
	paramParsingArrStringTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_arr_string.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingBoolTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_bool.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramParsingGGTimestampTemplate, err := util.ReadFileAsString("scaffold/template/param_parsing_gg_timestamp.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	paramDeclarations := []string{}
	paramParsings := []string{}
	hasIntField := false
	hasRepeatField := false
	hasBoolField := false
	for _, field := range requestType.Fields {
		declaration := paramDecTemplate
		declaration = strings.ReplaceAll(declaration, "<param_name>", field.Name)
		if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_location>", "query")
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_location>", field.Location)
		}
		if field.Repeated {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
			declaration = strings.ReplaceAll(declaration, "<param_description>", "value1,value2,value3,...")
		} else if field.Location == "enum" {
			declaration = strings.ReplaceAll(declaration, "<param_type>", "string")
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_type>", field.Type)
		}
		if field.Location == "path" {
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(true))
		} else {
			declaration = strings.ReplaceAll(declaration, "<param_required>", strconv.FormatBool(field.Required))
		}
		declaration = strings.ReplaceAll(declaration, "<param_description>", " ")
		paramDeclarations = append(paramDeclarations, declaration)

		parsing := paramParsingStringTemplate
		if strings.HasPrefix(field.Type, "int") {
			hasIntField = true
			parsing = paramParsingIntTemplate
			if field.Repeated {
				hasRepeatField = true
				parsing = paramParsingArrIntTemplate
			}
		} else if field.Type == "bool" {
			hasBoolField = true
			parsing = paramParsingBoolTemplate
		} else if field.Type == "google.protobuf.Timestamp" {
			parsing = paramParsingGGTimestampTemplate
		} else if field.Type == "string" {
			if field.Repeated {
				hasRepeatField = true
				parsing = paramParsingArrStringTemplate
			}
		} else if field.Location == "enum" {
			parsing = paramParsingEnumTemplate
		} else {
			fmt.Println("unsupported param type:", field.Type, ", fallback to string type")
		}

		parsing = strings.ReplaceAll(parsing, "<param_name>", field.Name)
		parsing = strings.ReplaceAll(parsing, "<param_name_cammel_case>", util.SnakeToCamel(field.Name))
		paramContext := "Query"
		if field.Location == "path" {
			paramContext = "Param"
		}
		parsing = strings.ReplaceAll(parsing, "<param_enum_object>", fmt.Sprintf("%s.%s", handler.ProtoPackage, field.Type))
		parsing = strings.ReplaceAll(parsing, "<param_context>", paramContext)
		parsing = strings.ReplaceAll(parsing, "<param_bit>", strconv.FormatInt(field.Bit, 10))
		parsing = strings.ReplaceAll(parsing, "<param_name_upper_case_word>", util.ConvertSliceStrToUCWord(strings.Split(field.Name, "_")))
		paramParsings = append(paramParsings, parsing)
	}
	if len(paramDeclarations) > 0 {
		template = strings.ReplaceAll(template, "<param_declarations>", strings.Join(paramDeclarations, "\n"))
	} else {
		template = strings.ReplaceAll(template, "\n<param_declarations>", strings.Join(paramDeclarations, "\n"))
	}

	template = strings.ReplaceAll(template, "<param_parsings>", strings.Join(paramParsings, "\n"))
	if !hasIntField {
		template = strings.ReplaceAll(template, "\"strconv\"", "")
	}
	if !hasRepeatField {
		template = strings.ReplaceAll(template, "\"strings\"", "")
	}
	if !hasBoolField {
		template = strings.ReplaceAll(template, "\"reflect\"", "")
	}

	err = os.WriteFile(handler.FolderPath+"/"+util.ToSnakeCase(handlerNameLower1st)+".go", []byte(template), 0644)
	if err != nil {
		fmt.Println("failed to write handler file:", err)
		return err
	}

	return nil
}

func GenerateRoutingRegistry(handlers []*model.Handler, handlerFolderPaths []string) error {
	fmt.Println("Generating routing registry")
	template, err := util.ReadFileAsString("scaffold/template/routing_registry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	handlerImportTemplate, err := util.ReadFileAsString("scaffold/template/handler_import.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	handlerImportPaths := []string{}
	for _, path := range handlerFolderPaths {
		handlerImportPath := handlerImportTemplate
		handlerImportPath = strings.ReplaceAll(handlerImportPath, "<handler_folder>", path)
		handlerImportPaths = append(handlerImportPaths, handlerImportPath)
	}

	routingConfigTemplate, err := util.ReadFileAsString("scaffold/template/routing_config.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}
	routingConfigs := []string{}
	for _, handler := range handlers {
		for _, path := range handler.Paths {
			routingConfig := routingConfigTemplate
			routingConfig = strings.ReplaceAll(routingConfig, "<method>", handler.Method)
			routingConfig = strings.ReplaceAll(routingConfig, "<path>", path)
			routingConfig = strings.ReplaceAll(routingConfig, "<handler_package>", handler.ProtoPackage)
			routingConfig = strings.ReplaceAll(routingConfig, "<handler_name>", handler.Name)
			routingConfig = strings.ReplaceAll(routingConfig, "<service_name>", handler.ServiceName)
			routingConfig = strings.ReplaceAll(routingConfig, "<action>", handler.Name)
			routingConfig = strings.ReplaceAll(routingConfig, "<permission>", handler.Permission)
			routingConfigs = append(routingConfigs, routingConfig)
		}

	}
	template = strings.ReplaceAll(template, "<handler_paths>", strings.Join(handlerImportPaths, "\n"))
	template = strings.ReplaceAll(template, "<routing_configs>", strings.Join(routingConfigs, "\n"))

	err = util.WriteFileWithFormat([]byte(template), ROUTING_REGISTRY_PATH)
	if err != nil {
		fmt.Println("failed to write routing registry file:", err)
		return err
	}

	return nil
}
