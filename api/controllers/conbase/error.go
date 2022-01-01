package conbase

type ErrorController struct {
	BaseController
}

func (e *ErrorController) Error404()  {
	e.Return.Code, e.Return.Data = GetCodeMessage(Err404)
	e.handle.Msg()
}