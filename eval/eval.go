package eval

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)
var builtins map[string]*object.BuiltIn

func Init() {

	builtins = map[string]*object.BuiltIn{

		//=========================

		"len": {
			Fn: length,
		},
		"push": {
			Fn: push,
		},
		"print": {
			Fn: print,
		},
		"printf": {
			Fn: printf,
		},
		"sort": {
			Fn: sort,
		},
		"typeof": {
			Fn: typeof,
		},
		"map": {
			Fn: maping,
		},

		//=========================

		"toStr": {
			Fn: toStr,
		},
		"toInt": {
			Fn: toInt,
		},
		"toFloat": {
			Fn: toFloat,
		},

		//=========================
		"randInt": {
			Fn: randInt,
		},
		"randFloat": {
			Fn: randFloat,
		},

		"randIntArr": {
			Fn: randIntArray,
		},

		//=========================

		"newWindow": {
			Fn: newWindow,
		},
		"openWindow": {
			Fn: openWindow,
		},
		"closeWindow": {
			Fn: closeWindow,
		},

		"addButton": {
			Fn: addButton,
		},
		"addEdit": {
			Fn: addEdit,
		},
		"addLabel": {
			Fn: addLabel,
		},
		"addCloseButton": {
			Fn: closeButton,
		},

		//=========================

		"read": {
			Fn: read,
		},
		"write": {
			Fn: write,
		},

		//=========================

	}
}

func Eval(node ast.Node, env *object.Env) object.Object {

	switch node := node.(type) {

	case *ast.Program:
		Init()
		return evalProgram(node, env)

	case *ast.ExpresssionStatement:

		return Eval(node.Expression, env)

	case *ast.BlockStatement:

		return evalBlockStatement(node, env)

	case *ast.ForLoop:
		return evalForLoop(node, env)

	case *ast.WhileLoop:
		return evalWhileLoop(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {

			return val
		}
		env.Set(node.Name.Value, val)

		return val

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)

		if isError(val) {

			return val
		}

		return &object.ReturnValue{
			Value: val,
		}
	case *ast.ReasignExpression:
		return evalReasignExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {

			return function
		}

		args := evalExpressions(node.Args, env)
		if len(args) == 1 && isError(args[0]) {

			return args[0]
		}

		return applyFunction(function, args)

	case *ast.FunctionLiteral:
		params := node.Params
		body := node.Body

		return &object.Function{
			Env:    env,
			Params: params,
			Body:   body,
		}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.Ident:

		return evalIdent(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{
			Elements: elements,
		}

	case *ast.StrLiteral:
		return &object.String{
			Value: node.Value,
		}
	case *ast.FloatLiteral:
		return &object.Float{
			Value: node.Value,
		}
	case *ast.IntLiteral:

		return &object.Integer{
			Value: int(node.Value),
		}

	case *ast.Boolean:

		return boolToBoolObj(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {

			return right
		}

		return evalPrefixExpression(right, node.Operator)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {

			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {

			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	}

	return nil
}

func evalWhileLoop(node *ast.WhileLoop, env *object.Env) object.Object {
	cond, ok := Eval(node.LoopCond, env).(*object.Boolean)
	if !ok {
		return newError("")
	}

	for cond.Value {
		Eval(node.Body, env)
		cond = Eval(node.LoopCond, env).(*object.Boolean)
	}

	return NULL
}

func evalForLoop(node *ast.ForLoop, env *object.Env) object.Object {
	Eval(node.LoopVar, env)

	cond, ok := Eval(node.LoopCond, env).(*object.Boolean)
	if !ok {
		return newError("")
	}

	for cond.Value {
		Eval(node.Body, env)
		Eval(node.PostLoop, env)
		cond = Eval(node.LoopCond, env).(*object.Boolean)
	}

	return NULL
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARR_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrIndexExpression(left, index)
	default:
		return newError("Index Operator Not Supported %s", left.Type())
	}

}
func evalArrIndexExpression(arr, index object.Object) object.Object {
	arrObj := arr.(*object.Array)
	idx := index.(*object.Integer).Value
	max := len(arrObj.Elements) - 1

	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elements[idx]
}

func evalReasignExpression(node *ast.ReasignExpression, env *object.Env) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	curVal, ok := env.Get(node.Var.Value)

	if !ok {
		return newError("Unkown Identefier %s", node.Var.Value)
	}

	var newVal object.Object

	switch node.Operator {
	case "=":
		newVal = val
	case "+=":
		newVal = evalInfixExpression("+", curVal, val)

	case "-=":
		newVal = evalInfixExpression("-", curVal, val)

	case "*=":
		newVal = evalInfixExpression("*", curVal, val)

	case "/=":
		newVal = evalInfixExpression("/", curVal, val)
	}

	env.Set(node.Var.Value, newVal)
	return newVal
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extEnv)
		return unwrapReturnValue(evaluated)

	case *object.BuiltIn:
		return fn.Fn(args...)

	default:
		return newError("not a function %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {

		return returnValue.Value
	}

	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)

	for i, p := range fn.Params {
		env.Set(p.Value, args[i])
	}

	return env
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var res []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {

			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}

	return res
}

func evalIdent(ident *ast.Ident, env *object.Env) object.Object {

	if val, ok := env.Get(ident.Value); ok {
		return val
	}
	if builtin, ok := builtins[ident.Value]; ok {
		return builtin
	}

	return newError("Identifier Not Found: %s", ident.Value)
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var res object.Object

	for _, stmt := range block.Statements {
		res = Eval(stmt, env)

		if res != nil {
			resType := res.Type()
			if resType == object.RETURN_OBJ || resType == object.ERROR_OBJ {

				return res
			}
		}
	}

	return res
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {

		return condition
	}

	if isTruthy(condition) {

		return Eval(ie.If, env)
	} else if ie.Else != nil {

		return Eval(ie.Else, env)
	} else {

		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:

		return false
	case FALSE:

		return false
	default:

		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:

		return evalIntegerInfix(operator, left, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:

		return evalFloatInfix(operator, left, right)

	case (left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ):

		return evalIntFloatInfix(operator, left, right, 1)

	case (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ):

		return evalIntFloatInfix(operator, left, right, 2)

	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:

		return evalBoolInfix(operator, left, right)

	case left.Type() == object.STR_OBJ && right.Type() == object.STR_OBJ:

		return evalStrInfix(operator, left, right)

	case left.Type() != right.Type():

		return newError("Type Mismatch %s %s %s", left.Type(), operator, right.Type())
	default:

		return newError("Unkown Operator %s %s %s", left.Type(), operator, right.Type())
	}

}
func evalIntFloatInfix(operator string, left, right object.Object, pos int) object.Object {
	var leftVal, rightVal float64
	switch pos {
	case 1:
		leftVal = left.(*object.Float).Value
		rightVal = float64(right.(*object.Integer).Value)
	case 2:
		leftVal = left.(*object.Float).Value
		rightVal = float64(right.(*object.Integer).Value)

	default:
		return newError("")
	}

	if obj := addNumberHelperFloat(operator, leftVal, rightVal); obj != nil {
		return obj
	} else {
		return compBoolHelper[float64](operator, leftVal, rightVal)
	}

}

func addNumberHelperFloat(operator string, leftVal, rightVal float64) object.Object {

	switch operator {
	case token.PLUS:

		return &object.Float{
			Value: leftVal + rightVal,
		}
	case token.MINUS:

		return &object.Float{
			Value: leftVal - rightVal,
		}
	case token.STAR:

		return &object.Float{
			Value: leftVal * rightVal,
		}
	case token.SLASH:

		return &object.Float{
			Value: leftVal / rightVal,
		}
	}
	return nil
}

func compBoolHelper[T int | float64 | string](operator string, leftVal, rightVal T) object.Object {
	switch operator {

	case token.GT:

		return boolToBoolObj(leftVal > rightVal)
	case token.LT:

		return boolToBoolObj(leftVal < rightVal)
	case token.EQ:

		return boolToBoolObj(leftVal == rightVal)
	case token.NOT_EQ:

		return boolToBoolObj(leftVal != rightVal)
	case token.LT_EQ:

		return boolToBoolObj(leftVal <= rightVal)

	case token.GT_EQ:

		return boolToBoolObj(leftVal >= rightVal)
	}

	return newError("")

}

func evalFloatInfix(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	if obj := addNumberHelperFloat(operator, leftVal, rightVal); obj != nil {
		return obj
	} else {
		return compBoolHelper[float64](operator, leftVal, rightVal)
	}

}

func evalStrInfix(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case token.PLUS:
		return &object.String{
			Value: leftVal + rightVal,
		}

	default:
		return compBoolHelper[string](operator, leftVal, rightVal)
	}

}

func evalBoolInfix(operator string, left, right object.Object) object.Object {
	switch operator {
	case token.EQ:

		return boolToBoolObj(left == right)
	case token.NOT_EQ:

		return boolToBoolObj(left != right)
	default:

		return newError("Unkown Oparator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfix(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case token.PLUS:

		return &object.Integer{
			Value: leftVal + rightVal,
		}
	case token.MINUS:

		return &object.Integer{
			Value: leftVal - rightVal,
		}
	case token.STAR:

		return &object.Integer{
			Value: leftVal * rightVal,
		}
	case token.SLASH:

		return &object.Integer{
			Value: leftVal / rightVal,
		}

	default:

		return compBoolHelper[int](operator, leftVal, rightVal)
	}
}

func evalPrefixExpression(right object.Object, operator string) object.Object {
	switch operator {
	case "!":

		return evalBangOperator(right)
	case "-":

		return evalMinusOperator(right)
	default:

		return newError("Unkown Operator %s", operator)
	}
}

func evalMinusOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {

		return newError("Unkown Oparator -%s", right.Type())
	}

	val := right.(*object.Integer).Value

	return &object.Integer{
		Value: -val,
	}
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:

		return FALSE
	case FALSE:

		return TRUE
	case NULL:

		return TRUE
	default:

		return FALSE
	}
}

func boolToBoolObj(b bool) *object.Boolean {
	if b {

		return TRUE
	}

	return FALSE
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var res object.Object

	for _, stmt := range program.Statements {

		res = Eval(stmt, env)
		switch res := res.(type) {

		case *object.ReturnValue:

			return res.Value
		case *object.Error:
			fmt.Println(res.Message)
			return res
		}

	}
	return res
}

func newError(format string, a ...interface{}) *object.Error {

	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {

		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
