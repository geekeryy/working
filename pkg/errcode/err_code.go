package errcode

type Code uint32

const (
	ParamErr      Code = iota + 1000 // 参数错误
	SourceNotFind                    // 资源不存在
	SystemErr                        // 系统错误
)

func (i Code) GetCode() uint32 {
	return uint32(i)
}
