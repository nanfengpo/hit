package hit

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

/*
这里为什么要使用反射来调用函数，而不是直接传函数作为形参然后执行？
- 首先参数是interface而不是某个具体类型的函数
- 其次如果形参是函数，就必须定义好这个函数的形参类型，这样就无法通用

*/

// callFn if args[i] == func, run it
func callFn(f interface{}) interface{} {
	if f != nil {
		t := reflect.TypeOf(f) // 先获取f的类型
		if t.Kind() == reflect.Func && t.NumIn() == 0 { // 如果类型是Func, 且没有形参
			function := reflect.ValueOf(f) // 获取函数的值
			in := make([]reflect.Value, 0)
			out := function.Call(in) // function是reflect.Value类型; in和out是[]reflect.Value类型
			if num := len(out); num > 0 {
				list := make([]interface{}, num)
				for i, value := range out {
					list[i] = value.Interface() // 调用reflect.Value的Interface()方法，转换为接口类型
				}
				if num == 1 {
					return list[0]
				}
				return list
			}
			return nil
		}
	}
	return f
}

func isZero(f interface{}) bool  {
	v := reflect.ValueOf(f)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		str := v.String()
		if str == "" {
			return true
		}
		zero, error := strconv.ParseFloat(str, 10)
		if zero == 0 && error == nil {
			return true
		}
		boolean, error := strconv.ParseBool(str)
		return boolean == false && error == nil
	default:
		return false
	}
}

// TestFnTime run func use time
func TestFnTime(f interface{}) string {
	start := time.Now()
	callFn(f)
	end := time.Now()
	vf := reflect.ValueOf(f)
	str := fmt.Sprintf("[%s] runtime: %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
	fmt.Println(str)
	return str
}

// If - (a ? b : c) Or (a && b)
func If(args ...interface{}) interface{} {
	// 几个变量：
	// 1）condition：需要判断的条件量
	// 2) trueVal: 条件量为true时的返回值
	// 3) falseVal: 条件量为false时的返回值
	var condition = callFn(args[0]) // 调用第1个参数，结果作为条件量
	// (1) 如果只有1个参数，则直接返回调用callFn的结果
	if len(args) == 1 {
		return condition
	}
	var trueVal = args[1]
	var falseVal interface{}
	if len(args) > 2 {
		// (3) 如果参数数量为3及以上，则第3个参数为falseVal
		falseVal = args[2]
	} else {
		// (2) 如果参数数量为2，则falseVal为nil
		falseVal = nil
	}
	if condition == nil {
		return callFn(falseVal)
	} else if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(falseVal)
		}
	} else if isZero(condition) {
		return callFn(falseVal)
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return callFn(trueVal)
}

// Or - (a || b)
func Or(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	if condition == nil {
		return callFn(args[1])
	}
	if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(args[1])
		}
	} else if isZero(condition) {
		return callFn(args[1])
	} else if v, ok := condition.(error); ok {
		if v != nil {
			fmt.Println(v)
			return condition
		}
	}
	return condition
}
