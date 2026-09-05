import { createProxyMiddleware } from "http-proxy-middleware";
import "dotenv/config"

export const orderProxy=createProxyMiddleware({
    target:process.env.ORDER_SERVICE_URL,
    changeOrigin:true,
    proxyTimeout:5000,
    timeout:5000,
   pathRewrite:(path)=>{
        if (path === "/") {
            return "/order";
        }

        return `/order${path}`;
    },
    on:{
        proxyReq:(proxyReq,req)=>{
            if(req.user){
                const idempotencyKey=req.headers["idempotency-key"]
                proxyReq.setHeader("X-USER-ID",req.user.id)
                proxyReq.setHeader("X-USER-EMAIL",req.user.email)
                proxyReq.setHeader("X-USER-ROLE",req.user.role)
                if (idempotencyKey) {
                    proxyReq.setHeader("Idempotency-Key", idempotencyKey);
                }
            }
        },
        error:(err,req,res)=>{
            console.error("order server error",err.message)
            if(!res.headersSent){
                res.status(503).json({
                    "sucess":false,
                    "message":"order service unavailable"
                })
            }
        }
    }
})