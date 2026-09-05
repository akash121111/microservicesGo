import { createProxyMiddleware } from "http-proxy-middleware";
import "dotenv/config";
export const userProxy=createProxyMiddleware({
    target:process.env.USER_SERVICE_URL,
    changeOrigin:true,
    proxyTimeout:5000,
    timeout:5000,
     pathRewrite:(path)=>{
        return `/api/v1/user${path}`
    },
    on:{
        error:(error,req,res)=>{
             console.error("user server error",error.message)
            if(!res.headersSent){
                res.status(503).json({
                    "sucess":false,
                    "message":"user service unavailable"
                })
            }
        }
    }
})