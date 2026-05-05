# PS Portal - Full Stack Template

A full-stack web application template with React frontend and Go backend, featuring OAuth authentication (Google & Microsoft), role-based access control, and dynamic resource management.

## 🏗️ Tech Stack

### Frontend
- **React 18** - UI Framework
- **React Router v6** - Client-side routing
- **Tailwind CSS** - Styling
- **Material-UI (MUI)** - Component library
- **Axios** - HTTP client
- **OAuth** - Google & Microsoft authentication
- **Chart.js & ApexCharts** - Data visualization
- **Monaco Editor** - Code editor
- **Syncfusion RTE** - Rich text editor

### Backend
- **Go 1.25** - Server runtime
- **Gin** - Web framework
- **MySQL** - Database
- **JWT** - Token-based authentication
- **CORS** - Cross-origin resource sharing
- **Rate Limiting** - Request throttling

---

## 🚀 Getting Started

### Prerequisites
- **Node.js** (v16 or higher)
- **Go** (v1.25 or higher)
- **MySQL** database
- Google OAuth credentials
- Microsoft Azure OAuth credentials (optional)

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ps_template
```


### 2. Backend Setup

#### Create `.env` file in `server/` directory:
```env
# Database Configuration
DB_USER=your_db_user
DB_PASS=your_db_password
DB_HOST=localhost:3306
DB_NAME=ps_portal

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id

# Microsoft OAuth (optional)
MICROSOFT_CLIENT_ID=your_microsoft_client_id
MICROSOFT_TENANT_ID=common

# Application Configuration
APP_BASE_PATH=/api
APP_DOMAIN=localhost
Origin=http://localhost:3000
```

#### Install dependencies and run:
```bash
cd server
go mod download
go run main.go
```

Server will start on `http://localhost:8080`

### 3. Frontend Setup

#### Create `.env` file in `client/` directory:
```env
REACT_APP_API=http://localhost:8080/api
REACT_APP_BASE=
REACT_APP_IMAGE_BASE=http://localhost:8080/api
REACT_APP_PORTAL_NAME=PS Portal
REACT_APP_GOOGLE_CLIENT_ID=your_google_client_id
REACT_APP_MSAL_CLIENT_ID=your_microsoft_client_id
REACT_APP_MSAL_TENANT_ID=common
```

#### Install dependencies and run:
```bash
cd client
npm install
npm start
```

Client will start on `http://localhost:3000`

---

## 📄 How to Create a New Page

### Step 1: Create Page Component
Create a new file in `client/src/pages/`:

```javascript
// client/src/pages/users/index.js
import React, { useState, useEffect } from "react";
import { apiGetRequest, apiPostRequest } from "../../utils/api";

function Users() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    const response = await apiGetRequest("/users");
    if (response.success) {
      setUsers(response.data);
    }
    setLoading(false);
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Users Management</h1>
      {/* Your component UI */}
    </div>
  );
}

export default Users;
```

### Step 2: Register Route
Add the route to `client/src/layouts/routes.js`:

```javascript
import { lazy } from "react";

const routes = {
  AppDashboard: lazy(() => import("../pages/dashboard")),
  AppUsers: lazy(() => import("../pages/users")), // Add this line
};

export default routes;
```

### Step 3: Add Resource to Database
Insert the resource into the database:

```sql
-- First, create or use an existing resource group
INSERT INTO master_resource_group (name, status) VALUES ('User Management', '1');

-- Add the resource (assuming res_group id is 1)
INSERT INTO master_resource_v2 (name, path, icon, element, menu, res_group, status, api_for, sort_by)
VALUES ('Users', '/users', 'bx-user', 'AppUsers', TRUE, 1, '1', 'app,api', 2);
```

### Step 4: Assign to Role
Update the role's resources field:

```sql
-- Update role to include the new resource group
UPDATE master_roles 
SET resources = CONCAT(resources, ',1') 
WHERE id = 1; -- Admin role
```

The page will automatically appear in the sidebar menu for authorized users!

---

## 📦 How to Add Resources

Resources control both menu visibility and API access permissions.

### Step 1: Create Resource Group (if needed)
```sql
INSERT INTO master_resource_group (name, status) 
VALUES ('Reports', '1');
```

### Step 2: Add Resource
```sql
INSERT INTO master_resource_v2 (
    name, 
    path, 
    icon, 
    element, 
    menu, 
    res_group, 
    status, 
    api_for, 
    sort_by
) VALUES (
    'Student Report',           -- Display name
    '/reports/students',        -- URL path
    'bx-file',                  -- Boxicon class
    'AppStudentReport',         -- React component name
    TRUE,                       -- Show in menu (TRUE/FALSE)
    2,                          -- Resource group ID
    '1',                        -- Active status
    'app,api',                  -- 'app' for menu, 'api' for API access
    5                           -- Sort order
);
```

### Step 3: Assign to Roles
```sql
-- Add resource group to role's resources
UPDATE master_roles 
SET resources = '1,2,3'  -- Comma-separated resource group IDs
WHERE id = 1;
```

### Resource Configuration Fields:
- **name**: Display name in menu
- **path**: URL path (must match route)
- **icon**: Boxicons class (e.g., 'bx-home', 'bx-user')
- **element**: React component name from routes.js
- **menu**: Show in sidebar (TRUE/FALSE)
- **api_for**: 
  - `'app'` - Menu only
  - `'api'` - API access only
  - `'app,api'` - Both menu and API
- **sort_by**: Menu order (ascending)

---

## 🔌 How to Create API

> **⚠️ Important Backend Convention**: Always use `c.Query()` for URL parameters, **NOT** `c.Param()`. This project uses query parameters (`?id=123`) instead of path parameters (`/users/123`) for consistency and flexibility.

### Step 1: Create API Handler

Create a new file in `server/api/<module>/<handler>.go`:

```go
package user

import (
    "net/http"
    "ps_portal/db"
    "ps_portal/utils"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID     int    `json:"id"`
    UserID string `json:"user_id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Dept   string `json:"dept"`
}

// GET /api/users
func GetAllUsers(c *gin.Context) {
    users := []User{}

    rows, err := db.DB.Query(`
        SELECT id, user_id, name, email, dept 
        FROM master_user 
        WHERE status NOT IN ('0', '9')
        ORDER BY name
    `)
    
    if err != nil {
        utils.Logging(c, err, 500)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.UserID, &user.Name, &user.Email, &user.Dept)
        if err != nil {
            utils.Logging(c, err, 500)
            return
        }
        users = append(users, user)
    }

    c.JSON(http.StatusOK, gin.H{
        "users": users,
        "total": len(users),
    })
}

// POST /api/users
func CreateUser(c *gin.Context) {
    var user User
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    result, err := db.DB.Exec(
        "INSERT INTO master_user (user_id, name, email, dept, role, status) VALUES (?, ?, ?, ?, 1, '1')",
        user.UserID, user.Name, user.Email, user.Dept,
    )

    if err != nil {
        utils.Logging(c, err, 500)
        return
    }

    id, _ := result.LastInsertId()
    c.JSON(http.StatusOK, gin.H{
        "message": "User created",
        "id": id,
    })
}

// PUT /api/users/:id
func UpdateUser(c *gin.Context) {
    id := c.Query("id")
    var user User
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    _, err := db.DB.Exec(
        "UPDATE master_user SET name = ?, email = ?, dept = ? WHERE id = ?",
        user.Name, user.Email, user.Dept, id,
    )

    if err != nil {
        utils.Logging(c, err, 500)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DELETE /api/users/:id
func DeleteUser(c *gin.Context) {
    id := c.Query("id")
    
    _, err := db.DB.Exec("UPDATE master_user SET status = '0' WHERE id = ?", id)
    
    if err != nil {
        utils.Logging(c, err, 500)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
```

### Step 2: Register Routes

Add routes in `server/routes/routes.go`:

```go
import (
    "ps_portal/api/user"
    // ... other imports
)

func SetupRouter() *gin.Engine {
    // ... existing code ...

    // Protected routes (requires JWT)
    router.Use(utils.JWTAuthMiddleware())
    
    // Scope middleware routes (requires permission check)
    router.Use(handles.ScopeMiddleware())
    
    // User routes
    router.GET(appBasePath+"/users", user.UpdateUser)      // Use ?id=123
    router.DELETE(appBasePath+"/users", user.DeleteUser)   // Use ?id=123

    return router
}
```

> **Note**: Routes don't use `:id` parameters. Instead, use query parameters like `/users?id=123`

### Step 3: Add API Resource Permission
    return router
}
```

### Step 3: Add API Resource Permission

Add the API path to the database for permission control:

```sql
-- Add API resource (menu = FALSE for API-only endpoints)
INSERT INTO master_resource_v2 (name, path, icon, element, menu, res_group, status, api_for)
VALUES ('Get Users API', '/users', '', '', FALSE, 1, '1', 'api');
```

### Step 4: Frontend API Call

Use the API in your React component:

```javascript
import { apiGetRequest, apiPostRequest, apiPutRequest, apiDeleteRequest } from "../../utils/api";

// GET request
const fetchUsers = async () => {
    const response = await apiGetRequest("/users");
    if (response.success) {
        setUsers(response.data.users);
    }
};

// POST request
const createUser = async (userData) => {
    const response = await apiPostRequest("/users", userData);
    if (respon (use query parameters)
const updateUser = async (id, userData) => {
    const response = await apiPutRequest(`/users?id=${id}`, userData);
    if (response.success) {
        console.log("User updated");
    }
};

// DELETE request (use query parameters)
const deleteUser = async (id) => {
    const response = await apiDeleteRequest(`/users?id=
};

// DELETE request
const deleteUser = async (id) => {
    const response = await apiDeleteRequest(`/users/${id}`);
    if (response.success) {
        console.log("User deleted");
    }
};
```

---

## 🔐 Authentication Flow

1. **User Login**: User authenticates via Google/Microsoft OAuth
2. **Token Validation**: Backend validates OAuth token
3. **User Lookup**: System checks if user exists in database
4. **JWT Generation**: Server generates JWT token with user claims
5. **Cookie Storage**: JWT stored in secure HTTP-only cookie
6. **Resource Fetch**: Frontend fetches user's accessible resources
7. **Dynamic Routing**: Routes and menus render based on permissions

---

## 🛡️ Middleware & Security

### Available Middlewares

1. **CorsMiddleware()**: Cross-origin resource sharing
2. **StrictOriginMiddleware()**: Origin validation
3. **RateLimit()**: Request throttling (1000 req/sec)
4. **JWTAuthMiddleware()**: JWT token validation
5. **ScopeMiddleware()**: Role-based access control

### Middleware Order
```go
router.Use(handles.CorsMiddleware())           // CORS
router.Use(handles.StrictOriginMiddleware())   // Origin check
router.POST("/auth/GLogin", auth.GoogleLogin)  // Public routes

router.Use(utils.JWTAuthMiddleware())          // Auth required
router.GET("/resources", resource.GetResources)

router.Use(handles.ScopeMiddleware())          // Permission check
router.GET("/users", user.GetAllUsers)         // Protected routes
```

---

## 📁 Project Structure

```
ps_template/
├── client/                      # React Frontend
│   ├── public/
│   │   ├── index.html
│   │   └── manifest.json
│   ├── src/
│   │   ├── assets/             # Images, fonts, etc.
│   │   ├── auth/               # Authentication components
│   │   │   ├── AuthContext.js  # Auth state management
│   │   │   ├── login.js        # Login page
│   │   │   ├── logout.js       # Logout handler
│   │   │   └── PrivateRoute.js # Protected route wrapper
│   │   ├── components/         # Reusable components
│   │   │   ├── button.js
│   │   │   ├── input.js
│   │   │   └── progress.js
│   │   ├── layouts/            # Layout components
│   │   │   ├── index.js        # Main layout
│   │   │   ├── routes.js       # Route configuration
│   │   │   └── sidebar.js      # Sidebar component
│   │   ├── pages/              # Page components
│   │   │   └── dashboard/
│   │   ├── theme/              # Theme configuration
│   │   ├── utils/              # Utility functions
│   │   │   ├── api.js          # API client
│   │   │   ├── crypto.js       # Encryption utilities
│   │   │   ├── dateUtils.js    # Date formatting
│   │   │   └── settings.js     # App configuration
│   │   ├── index.js            # App entry point
│   │   └── index.css           # Global styles
│   ├── package.json
│   └── tailwind.config.js
│
├── server/                      # Go Backend
│   ├── api/                    # API handlers
│   │   ├── auth/
│   │   │   ├── google.go       # Google OAuth
│   │   │   └── microsoft.go    # Microsoft OAuth
│   │   └── resource/
│   │       ├── getResources.go
│   │       ├── getActivity.go
│   │       └── getPresentationView.go
│   ├── config/
│   │   └── config.go           # Configuration loader
│   ├── db/
│   │   ├── db.go               # Database connection
│   │   └── scheme/             # SQL schemas
│   ├── handles/
│   │   ├── middleware.go       # Middleware functions
│   │   ├── getCourseImages.go
│   │   ├── getDeptImages.go
│   │   └── getProfileImage.go
│   ├── images/                 # Static image storage
│   │   ├── courses/
│   │   ├── departments/
│   │   └── users/
│   ├── routes/
│   │   └── routes.go           # Route definitions
│   ├── service/
│   │   ├── auth.go             # Auth business logic
│   │   └── qrGeneration.go     # QR code generation
│   ├── utils/
│   │   ├── client_ip.go        # IP extraction
│   │   ├── crypto.go           # Encryption
│   │   ├── jwt.go              # JWT utilities
│   │   ├── logging.go          # Error logging
│   │   └── response.go         # Response helpers
│   ├── main.go                 # Server entry point
│   ├── go.mod
│   └── .env
│
└── readme.md
```

---

## 🎨 UI Components

The template includes pre-built components in `client/src/components/`:

- **button.js**: Styled button component
- **input.js**: Form input component
- **progress.js**: Loading progress bar

### Using Components
```javascript
import Button from "../components/button";
import Input from "../components/input";

<Input 
    placeholder="Enter name" 
    value={name} 
    onChange={(e) => setName(e.target.value)} 
/>

<Button onClick={handleSubmit}>
    Submit
</Button>
```

---

## 🔧 Utility Functions

### Frontend Utilities (`client/src/utils/`)

#### API Client (api.js)
```javascript
import { apiGetRequest, apiPostRequest, apiPutRequest, apiDeleteRequest } from "./utils/api";

// GET
const response = await apiGetRequest("/endpoint");

// POST
const response = await apiPostRequest("/endpoint", { data });

// PUT
const response = await apiPutRequest("/endpoint", { data });

// DELETE
const response = await apiDeleteRequest("/endpoint");
```

#### Crypto (crypto.js)
Encryption/decryption utilities for sensitive data

#### Date Utils (dateUtils.js)
Date formatting and manipulation functions

#### Settings (settings.js)
```javascript
import { appBase, imagesBase, portalName } from "./utils/settings";
```

### Backend Utilities (`server/utils/`)

#### JWT (jwt.go)
```go
// Generate JWT
token, err := utils.GenerateJWT(name, email, id, userId, dept, year, yearGroup, role)

// Parse JWT
token, claims, err := utils.ParseToken(tokenString)
```

#### Logging (logging.go)
```go
utils.Logging(c, err, statusCode)
```


---

## 🐛 Debugging

### Check Server Logs
```bash
cd server
go run main.go
```

### Check Frontend Console
Open browser DevTools → Console tab

### Common Issues

**Issue**: Login fails
- Check OAuth credentials in `.env`
- Verify user exists in database
- Check CORS origin settings

**Issue**: Page not showing in menu
- Verify resource status is '1'
- Check user's role has access to resource group
- Ensure resource group status is '1'

**Issue**: API returns 401 Unauthorized
- Check JWT token validity
- Verify user has permission (check master_resource_v2 with api_for='api')
- Ensure ScopeMiddleware is configured

**Issue**: CORS error
- Add frontend URL to `Origin` in server `.env`
- Restart backend server

---

## 🔐 Environment Variables Reference

### Backend (server/.env)
```env
# Database
DB_USER=root
DB_PASS=password
DB_HOST=localhost:3306
DB_NAME=ps_portal

# OAuth
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
MICROSOFT_CLIENT_ID=xxx
MICROSOFT_TENANT_ID=common

# App Configuration
APP_BASE_PATH=/api
APP_DOMAIN=localhost
Origin=http://localhost:3000,http://localhost:3001
```

### Frontend (client/.env)
```env
REACT_APP_API=http://localhost:8080/api
REACT_APP_BASE=
REACT_APP_IMAGE_BASE=http://localhost:8080/api
REACT_APP_PORTAL_NAME=PS Portal
REACT_APP_GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
REACT_APP_MSAL_CLIENT_ID=xxx
REACT_APP_MSAL_TENANT_ID=common
```

---

## 📦 Production Deployment

### Frontend Build
```bash
cd client
npm run build
```
Serve the `build/` directory using Nginx, Apache, or static hosting.

### Backend Build
```bash
cd server
go build -o ps_portal main.go
./ps_portal
```
