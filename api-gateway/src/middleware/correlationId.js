import crypto from "crypto"
const correlationId=(req,res,next)=>{
      let id = req.header("X-Correlation-ID");

  if (!id) {
    id = crypto.randomUUID();
  }

  req.headers["x-correlation-id"] = id;
  res.setHeader("X-Correlation-ID", id);

  next();
}

export default correlationId