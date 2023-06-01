// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
package gsrv
import (
	"context"
	mrz "github.com/tinkler/mqttadmin/mrz/v1"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	pb_role_v1 "github.com/tinkler/mqttadmin/role/v1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)


type roleGsrv struct {
	pb_role_v1.UnimplementedRoleGsrvServer
}

func NewRoleGsrv() *roleGsrv {
	return &roleGsrv{}
}


func (u *roleGsrv) RoleSave(ctx context.Context, in *anypb.Any) (out *anypb.Any, err error) {
	gm := mrz.ToTypedModel[*pb_role_v1.Role, *structpb.Struct](in)
	m := mrz.GetData[*role.Role](gm)
	res := mrz.NewTypedRes[*pb_role_v1.Role, *structpb.Value]()
	err = m.Save(ctx)
	if err != nil {
		return nil, err
	}
	mrz.SetResData(res, m)
	mrz.SetResRespNil(res)
	return res.ToAny(), nil
}

