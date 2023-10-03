package object

import (
	"math"
)

func sin(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: taylorSin(float64(arg.Value))}
	case *Float:
		return &Float{Value: taylorSin(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func asin(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arcsin(float64(arg.Value))}
	case *Float:
		return &Float{Value: arcsin(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func cos(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: taylorCos(float64(arg.Value))}
	case *Float:
		return &Float{Value: taylorCos(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func acos(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arccos(float64(arg.Value))}
	case *Float:
		return &Float{Value: arccos(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func tan(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: taylorSin(float64(arg.Value)) / taylorCos(float64(arg.Value))}
	case *Float:
		return &Float{Value: taylorSin(arg.Value) / taylorCos(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func atan(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arctan(float64(arg.Value))}
	case *Float:
		return &Float{Value: arctan(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func cot(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: taylorCos(float64(arg.Value)) / taylorSin(float64(arg.Value))}
	case *Float:
		return &Float{Value: taylorCos(arg.Value) / taylorSin(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func acot(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}
	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arccot(float64(arg.Value))}
	case *Float:
		return &Float{Value: arccot(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}
}

func sec(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: 1.0 / taylorCos(float64(arg.Value))}
	case *Float:
		return &Float{Value: 1.0 / taylorCos(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func asec(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}
	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arcsec(float64(arg.Value))}
	case *Float:
		return &Float{Value: arcsec(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}
}

func csc(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: 1.0 / taylorSin(float64(arg.Value))}
	case *Float:
		return &Float{Value: 1.0 / taylorSin(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func acsc(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}
	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: arccsc(float64(arg.Value))}
	case *Float:
		return &Float{Value: arccsc(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}
}

func nlog(args ...Object) Object {

	if len(args) != 1 {
		return argumentAmountError(1, len(args))
	}

	switch arg := args[0].(type) {
	case *Integer:
		return &Float{Value: ln(float64(arg.Value))}
	case *Float:
		return &Float{Value: ln(arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func logBase(args ...Object) Object {

	if len(args) != 2 {
		return argumentAmountError(1, len(args))
	}
	var v float64
	switch arg := args[0].(type) {
	case *Integer:
		v = float64(arg.Value)
	case *Float:
		v = arg.Value
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

	switch arg := args[1].(type) {
	case *Integer:
		return &Float{Value: log(v, float64(arg.Value))}
	case *Float:
		return &Float{Value: log(v, arg.Value)}
	default:
		return newError("Object of Type %s Has wrong Type", args[0].Type())
	}

}

func taylorSin(num float64) float64 {
	for num > math.Pi {
		num = -math.Pi + (num - math.Pi)
	}

	for num < -math.Pi {
		num = math.Pi + (num + math.Pi)
	}
	return (num - (math.Pow(num, 3.0))/fac(3) + (math.Pow(num, 5.0))/fac(5) - (math.Pow(num, 7.0))/fac(7) + (math.Pow(num, 9.0))/fac(9))
}

func taylorCos(num float64) float64 {

	for num > math.Pi {
		num = -math.Pi + (num - math.Pi)
	}

	for num < -math.Pi {
		num = math.Pi + (num + math.Pi)
	}

	return (1 - (math.Pow(num, 2.0))/fac(2) + (math.Pow(num, 4.0))/fac(4) - (math.Pow(num, 6.0))/fac(6) + (math.Pow(num, 8.0))/fac(8))
}

func arcsin(num float64) float64 {
	if math.Abs(num) == 1 {
		return num * (math.Pi / 2.0)
	}
	return arctan(num / math.Sqrt(1-math.Pow(num, 2)))
}

func arccos(num float64) float64 {
	return (math.Pi / 2.0) - arcsin(num)
}

func arcsec(num float64) float64 {
	return arccos(1 / num)
}

func arccsc(num float64) float64 {
	return arcsin(1 / num)

}

func arccot(num float64) float64 {
	return (math.Pi / 2.0) - arctan(num)

}

func arctan(num float64) float64 {
	result := 0.0

	if num >= 1 {
		return (math.Pi / 2.0) - arctan(1/num)
	}

	sign := 1.0

	i := 1.0

	term := num

	for math.Abs(term) > 1e-5 {
		term = sign * (math.Pow(num, i) / i)
		result += term
		sign = -sign
		i += 2.0
	}

	return result
}

func fac(n int) float64 {
	f := 1
	for i := 1; i <= n; i++ {
		f *= i
	}
	return float64(f)
}

func Power(n float64, exponent float64) float64 {
	if exponent < 0 {
		n = 1 / n
		exponent = -exponent
	}
	if math.Trunc(exponent) == exponent {
		res := 1.0
		for i := 0; i < int(exponent); i++ {
			res *= n
		}
		return res
	} else {
		return math.Exp(exponent * ln(n))
	}
}

func log(n float64, base float64) float64 {
	return ln(n) / ln(base)
}

func ln(n float64) float64 {
	return math.Log(n)
}
