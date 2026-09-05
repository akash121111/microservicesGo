package apperror

import "errors"

var ErrUserNotFound = errors.New("user Not found")
var ErrUserAlreadyExist = errors.New("user already exist")
var ErrInvalidUserID = errors.New("Invalid id")
var ErrUnauthorized = errors.New("unauthorize")
var ErrInvalidToken = errors.New("Invalid Token")
var ErrTokenExpired = errors.New("Toke expired")
var ErrForbidden = errors.New("Not authorize for this action")
var ErrInvalidCredentials = errors.New("Invalid credentials")
var ErrOrderNotFound = errors.New("Order Not Found")
var ErrProductAlreadyExist = errors.New("Product Name already Exist")
var ErrProductNotFound = errors.New("Product doesn't exist")
var ErrInsuffiecientProduct = errors.New("Insufficient stock")
var ErrOrderAlreadyCancel = errors.New("Order Already cancel")
var ErrInvalidStokes = errors.New("Invalid stocks")
