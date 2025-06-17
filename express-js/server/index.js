const express = require("express");
const path = require("path");
const cookieParser = require("cookie-parser");
const axios = require("axios");
const crypto = require("crypto");
const dotenv = require("dotenv");
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3100;

// Config
const REAL_BACKEND_URL = process.env.REAL_BACKEND_URL;
const COOKIE_SECRET = process.env.COOKIE_SECRET || "super-secret";
const PROTECT_API = process.env.PROTECT_API === 'true';
const CLIENT_BUILDS_FOLDER = process.env.CLIENT_BUILDS_FOLDER || 'dist';

// Middleware
app.use(express.static(path.join(__dirname, `./web/${CLIENT_BUILDS_FOLDER}`)));
app.use(express.json());
app.use(cookieParser(COOKIE_SECRET));

// Util to create a simple signed token
function generateToken() {
  return crypto.randomBytes(32).toString("hex");
}

// Middleware to assign token cookie if not present
if (PROTECT_API) {
  app.use((req, res, next) => {
    if (!req.signedCookies.token) {
      const token = generateToken();
      res.cookie("token", token, {
        httpOnly: true,
        signed: true,
        maxAge: 15 * 60 * 1000, // 15 minutes
        sameSite: "Strict",
        secure: false // set to true if using HTTPS
      });
    }
    next();
  });
}

// Proxy endpoint for API requests
app.use("/api", async (req, res) => {
  if (PROTECT_API) {
    const clientToken = req.signedCookies.token;
    if (!clientToken) {
      return res.status(403).json({ error: "Invalid or missing token" });
    }
  }

  try {
    const headers = req.headers
    // add the custom headers if need 
    // or the screet headers for accessing the real backend
    // headers.some_secret_header = some_secret_value
    const response = await axios({
      method: req.method,
      url: `${REAL_BACKEND_URL}${req.originalUrl}`,
      data: req.body,
      headers: headers
    });
    res.status(response.status).json(response.data);
  } catch (err) {
    res.status(err.response?.status || 500).json({
      error: `Internal error while contacting backend\n${err.message}`
    });
  }
});

// Fallback for React frontend
app.get("/{*any}", (req, res) => {
  res.sendFile(path.join(__dirname, `./web/${CLIENT_BUILDS_FOLDER}/index.html`));
});

app.listen(PORT, () => {
  console.log(`Proxy server running at http://localhost:${PORT}`);
});
