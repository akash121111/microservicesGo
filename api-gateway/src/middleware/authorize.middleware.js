import ApiError from "../utils/apiError.js"

const authorize=(...allowedRoles)=>{
    return (req,res,next)=>{
        if(!req.user){
            throw new ApiError(401,"authentication is required")
        }
        if(!allowedRoles.includes(req.user.role)){
            throw new ApiError(403,"you do not permision for this action")
        }
        next();    
    }
    
}
export default authorize;