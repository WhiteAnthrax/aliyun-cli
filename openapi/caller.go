/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package openapi

//
//type Caller struct {
//	profile *config.Profile
//	library *meta.Repository
//	helper  *Helper
//
//	force   bool
//	verbose bool // TODO: next version
//}
//
//type CallerProcessor interface {
//	CallerProcessor(caller *Caller)
//}
//
//func NewCaller(profile *config.Profile, library *meta.Repository) *Caller {
//	return &Caller{
//		profile: profile,
//		library: library,
//		helper:  NewHelper(library),
//	}
//}
//
//func (c *Caller) Validate() error {
//	return c.profile.Validate()
//}
//
////
//// entrance of calling from main
//// will call rpc or restful
//func (c *Caller) Run(ctx *cli.Context, productCode string, apiOrMethod string, path string) error {
//	//
//	// get force call information
//	c.force = ForceFlag.IsAssigned()
//
//	//
//	// get product info
//	product, ok := c.library.GetProduct(productCode)
//	if !ok {
//		if !c.force {
//			return &InvalidProductError{Code: productCode, library: c.library}
//		}
//
//		//
//		// Restful Call
//		// aliyun cs GET /clusters
//		// aliyun cs /clusters --roa GET
//		ok, method, path, err := checkRestfulMethod(ctx, apiOrMethod, path)
//		if ok {
//			if err != nil {
//				return err
//			} else {
//				c.InvokeRestful(ctx, &product, method, path)
//			}
//		} else {
//			c.InvokeRpcForce(ctx, &product, apiOrMethod)
//		}
//	} else {
//		//
//		//
//		if strings.ToLower(product.ApiStyle) == "rpc" {
//			//
//			// Rpc call
//			if path != "" {
//				// ctx.Command().PrintFailed(fmt.Errorf("invalid arguments"), "")
//				return fmt.Errorf("invailed argument")
//			}
//			if c.force {
//				c.InvokeRpcForce(ctx, &product, apiOrMethod)
//				return nil
//			} else {
//				return c.InvokeRpc(ctx, &product, apiOrMethod)
//			}
//		} else {
//			//
//			// Restful Call
//			// aliyun cs GET /clusters
//			// aliyun cs /clusters --roa GET
//			ok, method, path, err := checkRestfulMethod(ctx, apiOrMethod, path)
//			if !ok {
//				if err != nil {
//					ctx.Command().PrintFailed(err, "")
//				} else {
//					ctx.Command().PrintFailed(fmt.Errorf("product %s need restful call", product.Code), "")
//				}
//				return nil
//			}
//			c.InvokeRestful(ctx, &product, method, path)
//			if err != nil {
//				ctx.Command().PrintFailed(fmt.Errorf("call restful %s%s.%s faild %v", product.Code, path, method, err), "")
//				return nil
//			}
//		}
//	}
//	return nil
//}
//
//func (c *Caller) InitClient(ctx *cli.Context, product *meta.Product, isRpc bool) (*sdk.Client, *requests.CommonRequest, error) {
//	//
//	// call OpenApi
//	// return: if check failed return error, otherwise return nil
//	client, err := c.profile.GetClient()
//	if err != nil {
//		return nil, nil, fmt.Errorf("bad client %v", err)
//	}
//
//	request := requests.NewCommonRequest()
//	request.Headers["User-Agent"] = fmt.Sprintf("Aliyun-CLI-V%s", strings.Replace(cli.Version, " ", "-", -1))
//	request.RegionId = c.profile.RegionId
//	request.Product = product.Code
//	request.Version = product.Version
//
//	if v, ok := config.RegionFlag.GetValue(); ok {
//		request.RegionId = v
//	}
//	if v, ok := EndpointFlag.GetValue(); ok {
//		request.Domain = v
//	}
//	if v, ok := VersionFlag.GetValue(); ok {
//		request.Version = v
//	}
//
//	//if v, ok := ctx.Flags().GetValue("content-type"); ok {
//	//	request.SetContentType(v)
//	//} else if isRpc {
//	//	request.SetContentType("application/json")
//	//}
//
//	if v, ok := AcceptFlag.GetValue(); ok {
//		request.Headers["Accept"] = v
//		if strings.Contains(v, "xml") {
//			request.AcceptFormat = "XML"
//		} else if strings.Contains(v, "json") {
//			request.AcceptFormat = "JSON"
//		}
//	}
//
//	if request.Version == "" {
//		return nil, nil, fmt.Errorf("unknown version! Use flag --version 2016-07-09 to assign version")
//	}
//
//	if request.Domain == "" {
//		request.Domain, err = product.GetEndpoint(request.RegionId, client)
//		if err != nil {
//			return nil, nil, fmt.Errorf("unknown endpoint for %s/%s! Use flag --endpoint xxx.aliyuncs.com to assign endpoint"+
//				"\n  error: %s", product.Code, request.RegionId, err.Error())
//		}
//	}
//
//	return client, request, nil
//}
//
//func (c *Caller) UpdateRequest(ctx *cli.Context, request *requests.CommonRequest) error {
//	if _, ok := SecureFlag.GetValue(); ok {
//		request.Scheme = "https"
//	}
//
//	for _, v := range HeaderFlag.GetValues() {
//		if k2, v2, ok := cli.SplitStringWithPrefix(v, "="); ok {
//			request.Headers[k2] = v2
//		} else {
//			return fmt.Errorf("invaild flag --header `%s` use `--header HeaderName=Value`", v)
//		}
//	}
//
//	if accept, ok := AcceptFlag.GetValue(); ok {
//		accept = strings.ToLower(accept)
//		if strings.Contains(accept, "xml") {
//			request.AcceptFormat = "XML"
//		} else if strings.Contains(accept, "json") {
//			request.AcceptFormat = "JSON"
//		} else {
//			return fmt.Errorf("unsupported accept: %s", accept)
//		}
//	}
//
//	return nil
//}
