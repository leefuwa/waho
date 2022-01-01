package register

type RegisterConf interface {
	Init()
	Start()
	IsRoutine() bool
}

type initRegister struct {
	notRoutine []RegisterConf
	routine    []RegisterConf
}

func (r *initRegister) AllRegisters() []RegisterConf {
	registers := make([]RegisterConf, 0)
	registers = append(registers, r.routine...)
	registers = append(registers, r.notRoutine...)
	return registers
}
func (r *initRegister) Register(register RegisterConf) {
	if register.IsRoutine() {
		r.routine = append(r.routine, register)
	} else {
		r.notRoutine = append(r.notRoutine, register)
	}
}

var InitRegister *initRegister = &initRegister{}

func Register(starter RegisterConf) {
	InitRegister.Register(starter)
}
func GetRegister() []RegisterConf {
	return InitRegister.AllRegisters()
}

type BaseRegister struct{}

func (baseRegister *BaseRegister) Init()           {}
func (baseRegister *BaseRegister) Start()          {}
func (baseRegister *BaseRegister) IsRoutine() bool { return true }
