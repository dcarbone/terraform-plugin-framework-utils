package validation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type GenericValidatorRequest struct {
	Path           path.Path
	PathExpression path.Expression
}

type GenericValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func toGenericValidatorRequest(src interface{}) (GenericValidatorRequest, error) {
	var (
		req GenericValidatorRequest
		err error
	)

	switch src.(type) {
	case validator.BoolRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.BoolRequest).Path,
			PathExpression: src.(validator.BoolRequest).PathExpression,
		}

	case validator.Float64Request:
		req = GenericValidatorRequest{
			Path:           src.(validator.Float64Request).Path,
			PathExpression: src.(validator.Float64Request).PathExpression,
		}
	case validator.Int64Request:
		req = GenericValidatorRequest{
			Path:           src.(validator.Int64Request).Path,
			PathExpression: src.(validator.Int64Request).PathExpression,
		}
	case validator.ListRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.ListRequest).Path,
			PathExpression: src.(validator.ListRequest).PathExpression,
		}
	case validator.MapRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.MapRequest).Path,
			PathExpression: src.(validator.MapRequest).PathExpression,
		}
	case validator.NumberRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.NumberRequest).Path,
			PathExpression: src.(validator.NumberRequest).PathExpression,
		}
	case validator.ObjectRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.ObjectRequest).Path,
			PathExpression: src.(validator.ObjectRequest).PathExpression,
		}
	case validator.SetRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.SetRequest).Path,
			PathExpression: src.(validator.SetRequest).PathExpression,
		}
	case validator.StringRequest:
		req = GenericValidatorRequest{
			Path:           src.(validator.StringRequest).Path,
			PathExpression: src.(validator.StringRequest).PathExpression,
		}

	default:
		err = fmt.Errorf("unknown validator request %T seen", src)
	}

	return req, err
}

func toGenericValidatorResponseType(src interface{}) (*GenericValidatorResponse, error) {
	if src == nil {
		return nil, nil
	}

	var (
		resp *GenericValidatorResponse
		err  error
	)

	switch src.(type) {
	case *validator.BoolResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.BoolResponse).Diagnostics,
		}
	case *validator.Float64Response:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.Float64Response).Diagnostics,
		}
	case *validator.Int64Response:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.Int64Response).Diagnostics,
		}
	case *validator.ListResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.ListResponse).Diagnostics,
		}
	case *validator.MapResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.MapResponse).Diagnostics,
		}
	case *validator.NumberResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.NumberResponse).Diagnostics,
		}
	case *validator.ObjectResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.ObjectResponse).Diagnostics,
		}
	case *validator.SetResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.SetResponse).Diagnostics,
		}
	case *validator.StringResponse:
		resp = &GenericValidatorResponse{
			Diagnostics: src.(*validator.StringResponse).Diagnostics,
		}

	default:
		err = fmt.Errorf("unknown validator response type: %T", src)
	}

	return resp, err
}

func toGenericValidatorTypes(srcReq interface{}, srcResp interface{}) (GenericValidatorRequest, *GenericValidatorResponse, error) {
	req, err := toGenericValidatorRequest(srcReq)
	if err != nil {
		return GenericValidatorRequest{}, nil, err
	}

	resp, err := toGenericValidatorResponseType(srcResp)
	if err != nil {
		return GenericValidatorRequest{}, nil, err
	}

	return req, resp, nil
}
