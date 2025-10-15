# Forum Application

A web-based forum application built with Go, featuring user authentication, post creation, commenting, and reaction systems. Users can create posts across different programming language categories, interact through comments, and express their opinions with likes and dislikes.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Database Schema](#database-schema)
- [API Routes](#api-routes)
- [Security](#security)
- [Known Issues & TODO](#known-issues--todo)
- [Contributing](#contributing)

## ✨ Features

### Current Features
- **User Authentication**
  - User registration with email, username, and password
  - Secure login system with bcrypt password hashing
  - Session-based authentication using cookies
  - Single active session per user
  - Logout functionality

- **Post Management**
  - Create posts with title and content
  - View all posts on the home page
  - View individual post details
  - Filter posts by programming language categories
  - Syntax-highlighted code display (using Chroma)

- **Categories**
  - Pre-defined programming language categories:
    - Golang, JavaScript, Rust, Ruby, Python, Java
  - Filter posts by category
  - View category-specific pages

- **Comments**
  - Add comments to posts
  - View all comments on a post
  - User attribution for comments

- **Reactions System**
  - Like/dislike posts
  - Like/dislike comments
  - Reaction counts displayed
  - Toggle reactions (cannot like and dislike simultaneously)

- **User Profiles**
  - View user profile pages
  - Display user information
  - Profile pictures support

### Guest Features
- Browse all posts and comments
- View user profiles
- Read content without authentication

## 🛠 Tech Stack

- **Backend**: Go 1.23.1
- **Database**: SQLite3
- **Authentication**: bcrypt password hashing, UUID-based sessions
- **Syntax Highlighting**: Chroma v2.20.0
- **Template Engine**: Go html/template
- **Frontend**: HTML, CSS, JavaScript

### Dependencies
```go
require (
    github.com/alecthomas/chroma/v2 v2.20.0
    github.com/google/uuid v1.6.0
    github.com/mattn/go-sqlite3 v1.14.32
    golang.org/x/crypto v0.41.0
)
```

## 📁 Project Structure

```
forum/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── forum.db               # SQLite database (generated)
├── AUDIT_REPORT.md        # Project audit and assessment
│
├── database/
│   ├── db.go              # Database initialization
│   └── db.sql             # Database schema
│
├── Internal/
│   ├── api/               # HTTP handlers
│   │   ├── server.go      # Server setup and routing
│   │   ├── login.go       # Login handler
│   │   ├── register.go    # Registration handler
│   │   ├── logout.go      # Logout handler
│   │   ├── createPost.go  # Post creation handler
│   │   ├── viewPost.go    # Post viewing handler
│   │   ├── discover_posts.go  # Posts discovery
│   │   ├── postInteractions.go  # Post reactions
│   │   ├── commentInteractions.go  # Comment reactions
│   │   ├── profile.go     # User profile handler
│   │   ├── categoryhand.go  # Category filtering
│   │   └── root.go        # Root/home handler
│   │
│   ├── model/
│   │   └── models.go      # Data models
│   │
│   ├── query/             # Database queries
│   │   ├── user.go        # User queries
│   │   ├── post.go        # Post queries
│   │   ├── comment.go     # Comment queries
│   │   ├── session.go     # Session queries
│   │   └── categories.go  # Category queries
│   │
│   └── service/           # Business logic
│       ├── service.go     # Service initialization
│       ├── Authentication.go  # Auth service
│       ├── user.go        # User service
│       ├── post.go        # Post service
│       ├── comment.go     # Comment service
│       ├── session.go     # Session service
│       ├── categories.go  # Category service
│       ├── profile.go     # Profile service
│       ├── error.go       # Error handling
│       └── golang.go      # Language-specific features
│
└── web/
    ├── static/
    │   ├── css/           # Stylesheets
    │   ├── icons/         # Icon assets
    │   ├── img/           # Images
    │   └── login-img/     # Login page images
    │
    └── templates/         # HTML templates
        ├── root.html      # Home page
        ├── login.html     # Login page
        ├── register.html  # Registration page
        ├── create-post.html  # Post creation
        ├── view-post.html    # Post detail view
        ├── profile.html   # User profile
        ├── category.html  # Category view
        ├── DiscoverPosts.html  # Posts discovery
        └── error.html     # Error page
```

## 📦 Prerequisites

- Go 1.23.1 or higher
- SQLite3
- Git (for cloning the repository)

## 🚀 Installation

1. **Clone the repository** (if using version control):
   ```bash
   git clone <repository-url>
   cd forum
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Initialize the database**:
   The database will be automatically created when you first run the application. The schema is defined in `database/db.sql`.

4. **Run the application**:
   ```bash
   go run .
   ```

5. **Access the application**:
   Open your browser and navigate to:
   ```
   http://localhost:7777
   ```

## 💻 Usage

### Registration
1. Navigate to `/register`
2. Fill in your full name, username, email, and password
3. Click "Register" to create your account

### Login
1. Navigate to `/login`
2. Enter your email and password
3. Click "Login" to access your account

### Creating a Post
1. Log in to your account
2. Navigate to `/create-post`
3. Enter a title and content
4. Select one or more categories
5. Click "Create Post"

### Interacting with Posts
- **View posts**: Click on any post from the home page
- **Add comments**: Scroll to the comment section on a post page
- **React**: Click the like or dislike buttons on posts and comments

### Filtering
- Use the category filter on the home page to view posts from specific programming languages

## 🗄 Database Schema

### Users Table
```sql
users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    photo TEXT DEFAULT 'default.png',
    fullname TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
)
```

### Posts Table
```sql
posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
```

### Comments Table
```sql
comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
)
```

### Categories Table
```sql
categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
)
```

### Post Categories (Many-to-Many)
```sql
post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
)
```

### Reactions Tables
```sql
post_reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    reaction_type TEXT CHECK(reaction_type IN ('like','dislike')),
    PRIMARY KEY (user_id, post_id)
)

comments_reactions (
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    reaction_type TEXT CHECK(reaction_type IN ('like','dislike')),
    PRIMARY KEY (user_id, comment_id)
)
```

### Sessions Table
```sql
sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
```

## 🔗 API Routes

| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| GET | `/` | Home page with all posts | No |
| GET | `/register` | Registration page | No |
| POST | `/register` | Create new user account | No |
| GET | `/login` | Login page | No |
| POST | `/login` | Authenticate user | No |
| GET | `/logout` | Logout user | Yes |
| GET | `/create-post` | Post creation page | Yes |
| POST | `/create-post` | Submit new post | Yes |
| GET | `/post/{id}` | View specific post | No |
| POST | `/create-comment` | Add comment to post | Yes |
| POST | `/post-reaction` | Like/dislike post | Yes |
| POST | `/comment-reaction` | Like/dislike comment | Yes |
| GET | `/profile/{username}` | View user profile | No |
| GET | `/category/{name}` | Filter posts by category | No |

## 🔒 Security

### Current Security Measures
- **Password Hashing**: Passwords are hashed using bcrypt
- **Session Management**: UUID-based session tokens
- **SQL Injection Prevention**: Prepared statements with parameterized queries
- **Template Escaping**: Automatic HTML escaping via Go's `html/template`
- **Cookie Security**: HttpOnly flag set, SameSite=Lax policy
- **Foreign Keys**: Enforced in SQLite for data integrity

### Security Improvements Needed
⚠️ **High Priority**:
- [ ] Add CSRF protection tokens
- [ ] Enforce session expiration checks
- [ ] Add `Secure` flag to cookies (HTTPS)
- [ ] Implement rate limiting for login attempts

⚠️ **Medium Priority**:
- [ ] Add input sanitization for user content
- [ ] Implement content security policy (CSP)
- [ ] Add brute force protection
- [ ] Session cleanup for expired sessions

## 🐛 Known Issues & TODO

### High Priority
- [ ] **CSRF Protection**: No CSRF tokens implemented
- [ ] **Session Expiration**: `expires_at` field exists but not enforced
- [ ] **Cookie Security**: Missing `Secure` flag for production
- [ ] **Comment Reaction Toggle**: Update logic needs refinement
- [ ] **Error Handling**: Inconsistent error page usage

### Medium Priority
- [ ] **Post Filtering**: Missing "my posts" and "my liked posts" filters
- [ ] **Pagination**: No pagination for posts list
- [ ] **Edit/Delete**: No post/comment editing or deletion
- [ ] **Validation**: Improve inline error messages
- [ ] **Dislike Counts**: Not displayed in UI

### Low Priority
- [ ] **Docker**: No Dockerfile or docker-compose
- [ ] **Testing**: No unit or integration tests
- [ ] **Documentation**: API documentation needed
- [ ] **Logging**: Enhanced logging system
- [ ] **Performance**: Add caching layer
- [ ] **Markdown Support**: Rich text formatting for posts

## 📊 Project Status

**Completion**: ~41% (based on audit report)

### What's Working ✅
- Core authentication flow
- Post creation and viewing
- Category filtering
- Comment system
- Basic reaction system
- Session management

### What Needs Work ⚠️
- Security enhancements (CSRF, session validation)
- Additional filtering options
- Pagination
- Edit/delete functionality
- Deployment artifacts
- Testing suite
- Comprehensive documentation

## 🧪 Testing

Currently, no automated tests are implemented. To test manually:

```bash
# Run the application
go run .

# In another terminal, check database
sqlite3 forum.db ".schema"
sqlite3 forum.db "SELECT COUNT(*) FROM users;"
sqlite3 forum.db "SELECT COUNT(*) FROM posts;"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 Quick Commands

```bash
# Run the application
go run .

# Run tests (when implemented)
go test ./...

# Build binary
go build -o forum.exe

# Check database schema
sqlite3 forum.db ".schema"

# View all tables
sqlite3 forum.db ".tables"

# Query users
sqlite3 forum.db "SELECT * FROM users;"

# Query posts with user info
sqlite3 forum.db "SELECT p.*, u.username FROM posts p JOIN users u ON p.user_id = u.id;"
```

## 📄 License

This project is currently unlicensed. Please add an appropriate license file for your use case.

## 👥 Authors

- Your Name/Team - Initial work

## 🙏 Acknowledgments

- Chroma library for syntax highlighting
- Go community for excellent standard library
- SQLite for a reliable embedded database

---

**Note**: This is an educational/development project. Before deploying to production, please address all security concerns listed in the "Known Issues & TODO" section.

For a detailed audit report, see [AUDIT_REPORT.md](AUDIT_REPORT.md).
