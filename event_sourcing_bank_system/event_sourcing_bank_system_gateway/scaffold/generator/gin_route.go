package generator

import (
	"event_sourcing_bank_system_gateway/scaffold/model"
	"event_sourcing_bank_system_gateway/scaffold/util"
	"fmt"
	"strings"
)

const GIN_ROUTES_PATH = "application/routing/delivery/routes.go"

func GenerateGinRoutes(protos []*model.Proto) error {
	fmt.Println("Generating gin routes")
	template, err := util.ReadFileAsString("scaffold/template/gin_routes.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	routingTemplate, err := util.ReadFileAsString("scaffold/template/gin_route_entry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	pubRoutingTemplate, err := util.ReadFileAsString("scaffold/template/gin_route_entry_public.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	entryTemplate, err := util.ReadFileAsString("scaffold/template/route_entry.tmpl")
	if err != nil {
		fmt.Println("failed to read template file:", err)
		return err
	}

	funcs := []string{}
	for _, proto := range protos {
		for _, service := range proto.Services {
			path := strings.Split(service.Path, ",")
			routes, proutes := []string{}, []string{}
			routing, prouting := routingTemplate, pubRoutingTemplate
			if service.Name == "LoginService" {
				continue // skip login service routes, write manually in controller.go
			}
			name := model.ServicePath[path[0]]

			if _, ok := model.PathPublic[path[0]]; ok {
				prouting = strings.ReplaceAll(prouting, "<name>", name)
			}
			routing = strings.ReplaceAll(routing, "<name>", name)
			// routes = append(routes, group)
			for _, api := range service.Apis {
				entry := entryTemplate
				entry = strings.ReplaceAll(entry, "<method>", api.Method)
				entry = strings.ReplaceAll(entry, "<path>", api.Path)
				if v := model.APIPublic[fmt.Sprintf("%s%s", path[0], api.Path)]; v {
					proutes = append(proutes, entry)
				} else {
					routes = append(routes, entry)
				}
			}
			routing = strings.ReplaceAll(routing, "<routes>", strings.Join(routes, "\n"))
			if len(proutes) > 0 {
				prouting = strings.ReplaceAll(prouting, "<routes>", strings.Join(proutes, "\n"))
				funcs = append(funcs, prouting)
			}
			funcs = append(funcs, routing)
		}
	}
	template = strings.ReplaceAll(template, "<funcs>", strings.Join(funcs, "\n\n"))

	err = util.WriteFileWithFormat([]byte(template), GIN_ROUTES_PATH)
	if err != nil {
		fmt.Println("failed to write gin routes file:", err)
		return err
	}

	return nil
}
