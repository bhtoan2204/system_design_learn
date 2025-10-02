package generator

import (
	"event_sourcing_bank_system_gateway/scaffold/model"
	"event_sourcing_bank_system_gateway/scaffold/util"
	"fmt"
	"strconv"
	"strings"
)

const INIT_SERVER_ROUTES_PATH = "application/init.go"

func GenerateInitRoutes(protos []*model.Proto) error {
	fmt.Println("Generating init server routes")
	template, err := util.ReadFileAsString("scaffold/template/routing_init.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	paramsRoutingTemplate, err := util.ReadFileAsString("scaffold/template/routing_init_params.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	funcsRoutingTemplate, err := util.ReadFileAsString("scaffold/template/routing_init_funcs.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	funcs := map[string]string{}
	params := map[string][]string{}
	for _, proto := range protos {
		if proto.Module == "" {
			continue
		}
		for _, service := range proto.Services {
			serviceName := strings.ReplaceAll(util.ToPascalCase(proto.Module)+util.ToPascalCase(service.Name), "Service", "")
			paths := strings.Split(service.Path, ",")
			// name := util.ToPascalCase(proto.Module + model.ServicePath[service.Path])
			if len(paths) > 1 {
				for i, path := range paths {
					idx := strconv.Itoa(i + 1)
					prouting := paramsRoutingTemplate
					prouting = strings.ReplaceAll(prouting, "<name>", model.ServicePath[paths[0]])
					prouting = strings.ReplaceAll(prouting, "<path>", path)
					prouting = strings.ReplaceAll(prouting, "<service_name>", serviceName+idx)
					params[proto.Module] = append(params[proto.Module], prouting)
				}
			} else {
				prouting := paramsRoutingTemplate
				prouting = strings.ReplaceAll(prouting, "<name>", model.ServicePath[paths[0]])
				prouting = strings.ReplaceAll(prouting, "<path>", service.Path)
				prouting = strings.ReplaceAll(prouting, "<service_name>", serviceName)
				params[proto.Module] = append(params[proto.Module], prouting)
			}
		}
		if _, found := funcs[proto.Module]; !found {
			if module := util.Find(model.IGNORE_PROTO_MODULES, func(m string) bool {
				return m == proto.Module
			}); module == "" {
				frouting := funcsRoutingTemplate
				frouting = strings.ReplaceAll(frouting, "<module_name>", util.ToPascalCase(proto.Module))
				funcs[proto.Module] = frouting
			}
		}
	}

	f := []string{}
	for k := range funcs {
		f = append(f, strings.ReplaceAll(funcs[k], "<params_routing_init>", strings.Join(params[k], "\n\n")))
	}

	template = strings.ReplaceAll(template, "<funcs_routing_init>", strings.Join(f, "\n\n"))

	err = util.WriteFileWithFormat([]byte(template), INIT_SERVER_ROUTES_PATH)
	if err != nil {
		fmt.Println("failed to write gin routes file:", err)
		return err
	}

	return nil
}
