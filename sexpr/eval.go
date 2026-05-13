package sexpr
 
import (
	"fmt"
	"errors"
	"math/big"
)
 
// ErrEval is the error value returned by the Evaluator if the contains
// an invalid token.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrEval = errors.New("eval error")

// checks if an sexpr represents nil, either as the empty list or the NIL symbol
func (expr *SExpr) isNilValue() bool {
	if expr == nil {
		return false
	}
	// check for it's the empty list (&SExpr{})
	if expr.isNil() {
		return true
	}
	// check symbol "nil"
	if expr.isSymbol() && expr.atom.literal == "NIL" {
		return true
	}
	return false
}

 
func (expr *SExpr) Eval() (*SExpr, error) {
	// check that expr is valid
	if expr == nil {
		return nil, ErrEval
	}

	if expr.isNilValue() {
		return mkNil(), nil
	}
 
	if expr.isNumber() {
		return expr.Number()
	}
 
	if expr.isConsCell() {
		// check car is def
		car, err := expr.Car()
		if err != nil {
			return nil, ErrEval
		}
		// fmt.Println("Eval: cons cell, car=", func() string {
		// 	if car == nil {
		// 		return "<nil>"
		// 	}
		// 	if car.isSymbol() {
		// 		return car.atom.literal
		// 	}
		// 	return car.SExprString()
		// }())
 
		if car.isSymbol() {
			switch car.atom.literal {
			case "QUOTE":
				return expr.Quote()
			case "CAR":
				return expr.CarFunc()
			case "CDR":
				return expr.CdrFunc()
			case "CONS":
				return expr.Cons()
			case "LENGTH":
				return expr.Length()
			case "ATOM":
				return expr.Atom()
			case "LISTP":
				return expr.Listp()
			case "ZEROP":
				return expr.Zerop()
			case "+":
				return expr.Sum()
			case "*":
				return expr.Product()
			default:
				return nil, ErrEval
			}
		}
	}
 
	return nil, ErrEval
}
 
func (expr *SExpr) Quote() (*SExpr, error) {
	// get arg by taking the cdr (arg . NIL)
	args, err := expr.Cdr()
	if err != nil || args == nil || args.isNil(){
		return nil, ErrEval
	}
 
	// args must be a cons cell
	if !args.isConsCell() {
		return nil, ErrEval
	}
 
	// verify that only one arg is passed
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}
 
	// check args CDR is nil VAL (empty list or NIL sym)
	if !argsCdr.isNilValue() {
		return nil, ErrEval
	}
 
	// get the first element (the thing to quote)
	quoteExpr, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}
 
	return quoteExpr, nil
}
 
func (expr *SExpr) Number() (*SExpr, error) {
	// we previously checked that we have a number
	return expr, nil
}
 
func (expr *SExpr) CarFunc() (*SExpr, error) {
	// expr is (CAR <arg>)
	args, err := expr.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	// verify exactly one argument
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	if argsCdr == nil || !argsCdr.isNil() {
		return nil, ErrEval
	}


	// get the argument and evaluate it
	arg, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}

	argVal, err := arg.Eval()
	if err != nil {
		return nil, ErrEval
	}
	if argVal.isNil() {
		return mkNil(), nil
	}

	// return the car of the evaluated result
	return argVal.Car()
}

func (expr *SExpr) CdrFunc() (*SExpr, error) {
	// expr is (CDR <arg>)
	args, err := expr.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	// verify exactly one argument
	argsCdr, err := args.Cdr()
	if err != nil {
		return nil, ErrEval
	}

	if argsCdr == nil || !argsCdr.isNil() {
		return nil, ErrEval
	}

	// get the argument and evaluate it
	arg, err := args.Car()
	if err != nil {
		return nil, ErrEval
	}

	argVal, err := arg.Eval()
	if err != nil {
		return nil, ErrEval
	}
	if argVal.isNil() {
		return mkNil(), nil
	}

	// return the cdr of the evaluated result
	return argVal.Cdr()
}
 
func (expr *SExpr) Cdr() (*SExpr, error) {
	// return the CDR of this cell
	return expr.cdr, nil
}
 
func (expr *SExpr) Car() (*SExpr, error) {
	// return the CAR of this cell
	return expr.car, nil
}

func (expr *SExpr) Cons() (*SExpr, error) {
	// get the first arg by taking the cdr (arg1 arg2 . NIL)
	arg1, err := expr.Cdr()
	fmt.Println(arg1)
	if err != nil || arg1 == nil || arg1.isNil(){
		return nil, ErrEval
	}
	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil {
		return nil, ErrEval
	}
	// eval arg1
	arg1Eval, err := arg1Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	
	arg2, err := arg1.Cdr()
	if err != nil {
		fmt.Println(2)
		return nil, ErrEval
	}
	arg2Cell, err := arg2.Car()
	if err != nil || arg2Cell == nil || arg2Cell.isNil(){
		return nil, ErrEval
	}
	arg2Eval, err := arg2Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	// check args2 CDR is nil VAL (empty list or NIL sym)
	argsCdr, err := arg2.Cdr()
	if err != nil {
		fmt.Println(3)
		return nil, ErrEval
	}
	if !argsCdr.isNilValue() {
		fmt.Println(4)
		return nil, ErrEval
	}
	return mkConsCell(arg1Eval, arg2Eval), nil
}
func (expr *SExpr) lengthHelper() (int64, error) {
	if expr == nil {
		return 0, nil
	}
	if expr.isNilValue() {
		return 0, nil
	}
	// must be a cons cell for a proper list
	if !expr.isConsCell() {
		return 0, ErrEval
	}
	// try to grab the cdr of this cell 
	exprCdr, err := expr.Cdr()
	if err != nil {
		return 0, ErrEval
	}
	recursiveCall, err := exprCdr.lengthHelper()
	if err != nil {
		return 0, ErrEval
	}
	return 1 + recursiveCall, nil
}

func (expr *SExpr) Length() (*SExpr, error) {
	// eval CAR
	arg1, err := expr.Cdr()
	if err != nil || arg1 == nil || arg1.isNil(){
		return nil, ErrEval
	}
	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil || arg1Cell.isNil(){
		return nil, ErrEval
	}
	// eval arg1
	arg1Eval, err := arg1Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	length, err := arg1Eval.lengthHelper()
	if err != nil {
		return nil, ErrEval
	}

	bigLen := big.NewInt(length)
	lengthCell := mkNumber(bigLen)
	return lengthCell, nil
}

func (expr *SExpr) Atom() (*SExpr, error) {
	// get arg by taking the cdr (arg . NIL)
	arg1, err := expr.Cdr()
	if err != nil || arg1 == nil {
		return nil, ErrEval
	}
	// args must be a cons cell
	if !arg1.isConsCell() {
		return nil, ErrEval
	}

	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil {
		return nil, ErrEval
	}
	// verify that only one arg is passed
	argsCdr, err := arg1.Cdr()
	if err != nil {
		return nil, ErrEval
	}
	// check args CDR is nil VAL (empty list or NIL sym)
	if !argsCdr.isNilValue() {
		return nil, ErrEval
	}

	if arg1Cell.isNilValue() {
		return mkSymbolTrue(), nil
	}

	// eval arg1
	arg1Eval, err := arg1Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	// get the first element (the thing to quote)
	if arg1Eval.isAtom() {
		return mkSymbolTrue(), nil
	}
	return mkNil(), nil
	
}

// Anything that is not a number or is nil
func (expr *SExpr) Listp() (*SExpr, error) {
	// get arg by taking the cdr (arg . NIL)
	arg1, err := expr.Cdr()
	if err != nil || arg1 == nil {
		return nil, ErrEval
	}
	// args must be a cons cell
	if !arg1.isConsCell() {
		return nil, ErrEval
	}

	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil {
		return nil, ErrEval
	}
	// verify that only one arg is passed
	argsCdr, err := arg1.Cdr()
	if err != nil {
		return nil, ErrEval
	}
	// check args CDR is nil VAL (empty list or NIL sym)
	if !argsCdr.isNilValue() {
		return nil, ErrEval
	}
	if arg1Cell.isNilValue() {
		return mkSymbolTrue(), nil
	}
	// eval arg1
	arg1Eval, err := arg1Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	if arg1Eval.isNilValue() || !arg1Eval.isAtom() {
		return mkSymbolTrue(), nil
	}

	return mkNil(), nil
}

// Anything that is not a number or is nil
func (expr *SExpr) Zerop() (*SExpr, error) {
	// get arg by taking the cdr (arg . NIL)
	arg1, err := expr.Cdr()
	if err != nil || arg1 == nil || arg1.isNil(){
		return nil, ErrEval
	}

	// args must be a cons cell
	if !arg1.isConsCell() {
		return nil, ErrEval
	}

	// get the CAR of cell
	arg1Cell, err := arg1.Car()
	if err != nil || arg1Cell == nil {
		return nil, ErrEval
	}
	// verify that only one arg is passed
	argsCdr, err := arg1.Cdr()
	if err != nil {
		return nil, ErrEval
	}
	// check args CDR is nil VAL (empty list or NIL sym)
	if !argsCdr.isNilValue() {
		return nil, ErrEval
	}
	
	// eval arg1
	arg1Eval, err := arg1Cell.Eval()
	if err != nil {
		return nil, ErrEval
	}

	// ZEROP only accepts numbers
	if arg1Eval == nil || !arg1Eval.isNumber() {
		return nil, ErrEval
	}

	zeroValue := big.NewInt(0)
	if zeroValue.Cmp(arg1Eval.atom.num) == 0  {
		return mkSymbolTrue(), nil
	}

	return mkNil(), nil
}
func (expr *SExpr) arithmeticHelper(op string) (*big.Int, error) {
	// fmt.Println("EVALUATING: " + expr.SExprString())
	if expr == nil || expr.isNilValue() {
		if (op == "add") {
			return big.NewInt(0), nil
		}
		return big.NewInt(1), nil
		
	}
	// grab current expr car
	exprCar, err := expr.Car()
	if err != nil {
		return nil, ErrEval
	}
	// evaluate the current expr
	exprEval, err := exprCar.Eval()
	if err != nil || !exprEval.isNumber() {
		return nil, ErrEval
	}
	// form big int from bigEval value
	exprEvalInt := exprEval.atom.num

	// recursive step
	exprCdr, err := expr.Cdr()
	if err != nil {
		return nil, ErrEval
	}
	exprCdrEval, err := exprCdr.arithmeticHelper(op)
	if err != nil {
		return nil, ErrEval
	}
	if (op == "add") {
		return new(big.Int).Add(exprEvalInt, exprCdrEval), nil
	}
	return new(big.Int).Mul(exprEvalInt, exprCdrEval), nil
}

func (expr *SExpr) Sum() (*SExpr, error) {
	// grab the first arg SExpr
	// get arg by taking the cdr (arg . NIL)
	arg1, err := expr.Cdr()

	if err != nil || arg1 == nil{
		return nil, ErrEval
	}

	// get the big int
	sumInt, err := arg1.arithmeticHelper("add")
	if err != nil {
		return nil, ErrEval
	}

	return mkNumber(sumInt), nil
}
func (expr *SExpr) Product() (*SExpr, error) {
	// grab the first arg SExpr
	// get arg by taking the cdr (arg . NIL)
	arg1, err := expr.Cdr()

	if err != nil || arg1 == nil{
		return nil, ErrEval
	}

	// get the big int
	prodInt, err := arg1.arithmeticHelper("mul")
	if err != nil {
		return nil, ErrEval
	}

	return mkNumber(prodInt), nil
}
