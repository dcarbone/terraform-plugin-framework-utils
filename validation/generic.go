package validation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type GenericRequest struct {
	Path           path.Path
	PathExpression path.Expression
	Config         tfsdk.Config
	ConfigValue    attr.Value

	source interface{}
}

func (r GenericRequest) BoolRequest() (validator.BoolRequest, bool) {
	br, ok := r.source.(validator.BoolRequest)
	return br, ok
}

func (r GenericRequest) Float64Request() (validator.Float64Request, bool) {
	fr, ok := r.source.(validator.Float64Request)
	return fr, ok
}

func (r GenericRequest) Int64Request() (validator.Int64Request, bool) {
	ir, ok := r.source.(validator.Int64Request)
	return ir, ok
}

func (r GenericRequest) ListRequest() (validator.ListRequest, bool) {
	lr, ok := r.source.(validator.ListRequest)
	return lr, ok
}

func (r GenericRequest) MapRequest() (validator.MapRequest, bool) {
	mr, ok := r.source.(validator.MapRequest)
	return mr, ok
}

func (r GenericRequest) NumberRequest() (validator.NumberRequest, bool) {
	nr, ok := r.source.(validator.NumberRequest)
	return nr, ok
}

func (r GenericRequest) ObjectRequest() (validator.ObjectRequest, bool) {
	or, ok := r.source.(validator.ObjectRequest)
	return or, ok
}

func (r GenericRequest) SetRequest() (validator.SetRequest, bool) {
	sr, ok := r.source.(validator.SetRequest)
	return sr, ok
}

func (r GenericRequest) StringRequest() (validator.StringRequest, bool) {
	sr, ok := r.source.(validator.StringRequest)
	return sr, ok
}

func toGenericRequest(src interface{}) (GenericRequest, error) {
	var (
		err error

		req = GenericRequest{
			source: src,
		}
	)

	if br, ok := src.(validator.BoolRequest); ok {
		req.Path = br.Path
		req.PathExpression = br.PathExpression
		req.Config = br.Config
		req.ConfigValue = br.ConfigValue
	} else if fr, ok := src.(validator.Float64Request); ok {
		req.Path = fr.Path
		req.PathExpression = fr.PathExpression
		req.Config = fr.Config
		req.ConfigValue = fr.ConfigValue
	} else if ir, ok := src.(validator.Int64Request); ok {
		req.Path = ir.Path
		req.PathExpression = ir.PathExpression
		req.Config = ir.Config
		req.ConfigValue = ir.ConfigValue
	} else if lr, ok := src.(validator.ListRequest); ok {
		req.Path = lr.Path
		req.PathExpression = lr.PathExpression
		req.Config = lr.Config
		req.ConfigValue = lr.ConfigValue
	} else if mr, ok := src.(validator.MapRequest); ok {
		req.Path = mr.Path
		req.PathExpression = mr.PathExpression
		req.Config = mr.Config
		req.ConfigValue = mr.ConfigValue
	} else if nr, ok := src.(validator.NumberRequest); ok {
		req.Path = nr.Path
		req.PathExpression = nr.PathExpression
		req.Config = nr.Config
		req.ConfigValue = nr.ConfigValue
	} else if or, ok := src.(validator.ObjectRequest); ok {
		req.Path = or.Path
		req.PathExpression = or.PathExpression
		req.Config = or.Config
		req.ConfigValue = or.ConfigValue
	} else if sr, ok := src.(validator.SetRequest); ok {
		req.Path = sr.Path
		req.PathExpression = sr.PathExpression
		req.Config = sr.Config
		req.ConfigValue = sr.ConfigValue
	} else if sr, ok := src.(validator.StringRequest); ok {
		req.Path = sr.Path
		req.PathExpression = sr.PathExpression
		req.Config = sr.Config
		req.ConfigValue = sr.ConfigValue
	} else {
		err = fmt.Errorf("unknown validator request type %T seen", src)
	}

	return req, err
}

type GenericResponse struct {
	Diagnostics diag.Diagnostics

	nil    bool
	source interface{}
}

// Nil will return true if the source validator response instance was nil
func (r *GenericResponse) Nil() bool {
	return r.nil
}

func (r *GenericResponse) BoolResponse() (*validator.BoolResponse, bool) {
	br, ok := r.source.(*validator.BoolResponse)
	return br, ok
}

func (r *GenericResponse) Float64Response() (*validator.Float64Response, bool) {
	fr, ok := r.source.(*validator.Float64Response)
	return fr, ok
}

func (r *GenericResponse) Int64Response() (*validator.Int64Response, bool) {
	ir, ok := r.source.(*validator.Int64Response)
	return ir, ok
}

func (r *GenericResponse) ListResponse() (*validator.ListResponse, bool) {
	lr, ok := r.source.(*validator.ListResponse)
	return lr, ok
}

func (r *GenericResponse) MapResponse() (*validator.MapResponse, bool) {
	mr, ok := r.source.(*validator.MapResponse)
	return mr, ok
}

func (r *GenericResponse) NumberResponse() (*validator.NumberResponse, bool) {
	nr, ok := r.source.(*validator.NumberResponse)
	return nr, ok
}

func (r *GenericResponse) ObjectResponse() (*validator.ObjectResponse, bool) {
	or, ok := r.source.(*validator.ObjectResponse)
	return or, ok
}

func (r *GenericResponse) SetResponse() (*validator.SetResponse, bool) {
	sr, ok := r.source.(*validator.SetResponse)
	return sr, ok
}

func (r *GenericResponse) StringResponse() (*validator.StringResponse, bool) {
	sr, ok := r.source.(*validator.StringResponse)
	return sr, ok
}

func toGenericResponse(src interface{}) (*GenericResponse, error) {
	var (
		err error

		resp = &GenericResponse{
			source: src,
		}
	)

	if br, ok := src.(*validator.BoolResponse); ok {
		resp.nil = br == nil
		if br != nil {
			resp.Diagnostics = br.Diagnostics
		}
	} else if fr, ok := src.(*validator.Float64Response); ok {
		resp.nil = fr == nil
		if fr != nil {
			resp.Diagnostics = fr.Diagnostics
		}
	} else if ir, ok := src.(*validator.Int64Response); ok {
		resp.nil = ir == nil
		if ir != nil {
			resp.Diagnostics = ir.Diagnostics
		}
	} else if lr, ok := src.(*validator.ListResponse); ok {
		resp.nil = lr == nil
		if lr != nil {
			resp.Diagnostics = lr.Diagnostics
		}
	} else if mr, ok := src.(*validator.MapResponse); ok {
		resp.nil = mr == nil
		if mr != nil {
			resp.Diagnostics = mr.Diagnostics
		}
	} else if nr, ok := src.(*validator.NumberResponse); ok {
		resp.nil = nr == nil
		if nr != nil {
			resp.Diagnostics = nr.Diagnostics
		}
	} else if or, ok := src.(*validator.ObjectResponse); ok {
		resp.nil = or == nil
		if or != nil {
			resp.Diagnostics = or.Diagnostics
		}
	} else if sr, ok := src.(*validator.SetResponse); ok {
		resp.nil = sr == nil
		if sr != nil {
			resp.Diagnostics = sr.Diagnostics
		}
	} else if sr, ok := src.(*validator.StringResponse); ok {
		resp.nil = sr == nil
		if sr != nil {
			resp.Diagnostics = sr.Diagnostics
		}
	} else {
		err = fmt.Errorf("unknown validator response type: %T", src)
	}

	return resp, err
}

func toGenericTypes(srcReq interface{}, srcResp interface{}) (GenericRequest, *GenericResponse, error) {
	req, err := toGenericRequest(srcReq)
	if err != nil {
		return GenericRequest{}, nil, err
	}

	resp, err := toGenericResponse(srcResp)
	if err != nil {
		return GenericRequest{}, nil, err
	}

	return req, resp, nil
}