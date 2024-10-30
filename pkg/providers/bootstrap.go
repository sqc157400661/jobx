package providers

import (
	"reflect"
	"sync"
)

type Inputer interface {
	GetRootID() int
	GetTaskID() int
	GetInput() map[string]interface{}
	GetEnv() map[string]interface{}
}

type DefaultInput struct {
	i      map[string]interface{}
	env    map[string]interface{}
	rootID int
	taskID int
}

func NewDefaultInput(rootID, taskID int, i, env map[string]interface{}) *DefaultInput {
	return &DefaultInput{
		i:      i,
		env:    env,
		rootID: rootID,
		taskID: taskID,
	}
}

func (d *DefaultInput) GetRootID() int {
	return d.rootID
}

func (d *DefaultInput) GetTaskID() int {
	return d.taskID
}

func (d *DefaultInput) GetInput() map[string]interface{} {
	return d.i
}

func (d *DefaultInput) GetEnv() map[string]interface{} {
	return d.env
}

type TaskProvider interface {
	Name() string
	Input(Inputer) (err error)
	Output() (ctx map[string]interface{}, res map[string]interface{}, err error)
	Run(int) (err error)
}

type providerPool struct {
	pool     *sync.Pool
	provider TaskProvider
}
type providers map[string]*providerPool

func NewProviderPool(provider TaskProvider) (pool *providerPool) {
	pool = &providerPool{
		provider: provider,
		pool: &sync.Pool{New: func() interface{} {
			return clone(provider)
		}},
	}
	return
}

var providersConfig providers
var registerOneProvidersConfig sync.Once

func init() {
	registerOneProvidersConfig.Do(func() {
		delay := &DelayTasker{}
		providersConfig = map[string]*providerPool{
			delay.Name(): NewProviderPool(delay),
		}
	})
}

func Has(name string) bool {
	_, has := providersConfig[name]
	return has
}

func Set(provider TaskProvider) {
	if provider == nil {
		return
	}
	if !Has(provider.Name()) {
		providersConfig[provider.Name()] = NewProviderPool(provider)
	}
}

func ReSet(provider TaskProvider) {
	if provider == nil {
		return
	}
	if Has(provider.Name()) {
		pool := providersConfig[provider.Name()]
		pool.pool.Put(provider)
	}
}

func Get(name string) TaskProvider {
	if Has(name) {
		pool := providersConfig[name]
		return pool.pool.Get().(TaskProvider)
	}
	return nil
}

func clone(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		data = data.Elem()             //源数据实际值（否则为指针）
		dst.Set(data)                  //设置数据
		dst = dst.Addr()               //创建对象的地址（否则返回值）
		return dst.Interface()         //返回地址
	} else {
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		dst.Set(data)                  //设置数据
		return dst.Interface()         //返回
	}
}

// Desc the interface used to define details for provider
type Desc interface {
	Desc() string
}

// Rollback the interface used to define rollback function for provider
type Rollback interface {
	Rollback(string)
}

var (
	tpDesc   = reflect.TypeOf((*Desc)(nil)).Elem()
	rollback = reflect.TypeOf((*Rollback)(nil)).Elem()
)

func GetDesc(action string, defaultDesc string) string {
	if !Has(action) {
		return defaultDesc
	}
	provider := Get(action)
	defer ReSet(provider)
	if reflect.TypeOf(provider).Implements(tpDesc) {
		return reflect.ValueOf(provider).Interface().(Desc).Desc()
	}
	return defaultDesc
}

func GetRollback(provider TaskProvider) Rollback {
	if reflect.TypeOf(provider).Implements(rollback) {
		return reflect.ValueOf(provider).Interface().(Rollback)
	}
	return nil
}
