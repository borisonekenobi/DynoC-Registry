package models

var SuccessResponse = Response{Message: "success"}                   // 200
var CreatedResponse = Response{Message: "created"}                   // 201
var BadRequestError = Response{Message: "invalid request"}           // 400
var MissingFieldError = Response{Message: "missing required fields"} // 400
var UnauthorizedError = Response{Message: "unauthorized"}            // 401
var ForbiddenError = Response{Message: "forbidden"}                  // 403
var NotFoundError = Response{Message: "not found"}                   // 404
var ConflictError = Response{Message: "conflict"}                    // 409
var ImATeapotError = Response{Message: "i'm a teapot"}               // 418
var InternalServerError = Response{Message: "internal server error"} // 500
var NotImplementedError = Response{Message: "not implemented"}       // 501

type Response struct {
	Message string `json:"message"`
}
