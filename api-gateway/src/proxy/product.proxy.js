import { createProxyMiddleware } from "http-proxy-middleware";
import "dotenv/config"
export const productProxy=createProxyMiddleware({
    target:process.env.PRODUCT_SERVICE_URL,
    changeOrigin:true,
    proxyTimeout:5000,
    timeout:5000,
     pathRewrite:(path)=>{
        if (path === "/") {
            return "/product";
        }

        return `/product${path}`;
    },
    on:{
        error:(error,req,res)=>{
             console.error("product server error",error.message)
            if(!res.headersSent){
                res.status(503).json({
                    "sucess":false,
                    "message":"product service unavailable"
                })
            }
        }
    }
})