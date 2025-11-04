# Forum - A Programming Community Platform

A modern, feature-rich forum application built with Go, designed for programmers to share knowledge, discuss topics, and collaborate on various programming languages and technologies.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Docker Deployment](#docker-deployment)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Usage Guide](#usage-guide)
- [Contributing](#contributing)

## ✨ Features

### Core Features

- **User Authentication**: Secure registration, login, and session management
- **Post Management**: Create, view, and interact with posts
- **Category System**: Organize posts by programming languages (Go, JavaScript, Rust, Ruby, Python, MATLAB, Brainfuck, Java)
- **Comment System**: Engage in discussions with nested comments
- **Reaction System**: Like/dislike posts and comments
- **User Profiles**: Customizable profiles with photo uploads
- **Code Playground**: Write and preview code with syntax highlighting
- **Discovery Feed**: Browse all posts across categories
- **Responsive Design**: Mobile-friendly interface

### Advanced Features

- **Session-based Authentication**: Secure cookie-based sessions
- **Syntax Highlighting**: Using Chroma for beautiful code display
- **Profile Customization**: Update password and profile pictures
- **Category Filtering**: View posts by specific programming language
- **Interactive Code Editor**: Built-in playground for testing code snippets
- **Privacy & Help Pages**: Comprehensive documentation

## 🛠 Tech Stack

### Backend

- **Language**: Go 1.23.1
- **Database**: SQLite3 with foreign key constraints
- **Router**: Standard `net/http` package with Go 1.22+ routing enhancements

### Frontend

- **Templates**: HTML templates with Go's `html/template`
- **Styling**: Custom CSS with category-specific themes
- **Assets**: Static file serving for images, icons, and styles

### Dependencies

```go
- github.com/mattn/go-sqlite3 v1.14.32       // SQLite driver
- github.com/google/uuid v1.6.0              // UUID generation for sessions
- golang.org/x/crypto v0.41.0                // Password hashing
- github.com/alecthomas/chroma/v2 v2.20.0    // Syntax highlighting
- github.com/dlclark/regexp2 v1.11.5         // Regular expressions (indirect)
```

## 📁 Project Structure

```
forum/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── Dockerfile                   # Docker configuration
├── README.md                    # This file
│
├── database/                    # Database layer
│   ├── db.go                   # Database initialization
│   ├── db.sql                  # Schema definitions
│   └── profile-img/            # User profile images storage
│       └── README.md
│
├── Internal/                    # Application core (private)
│   ├── api/                    # HTTP handlers & routing
│   │   ├── server.go           # Server initialization & routes
│   │   ├── root.go             # Home page handler
│   │   ├── register.go         # User registration
│   │   ├── login.go            # User login
│   │   ├── logout.go           # User logout
│   │   ├── profile.go          # User profile view
│   │   ├── editProfile.go      # Profile editing
│   │   ├── createPost.go       # Post creation
│   │   ├── viewPost.go         # Single post view
│   │   ├── discover_posts.go   # All posts feed
│   │   ├── categoryhand.go     # Category filtering
│   │   ├── postInteractions.go # Post reactions
│   │   ├── commentInteractions.go # Comment reactions
│   │   ├── playground.go       # Code playground
│   │   └── help.go             # Help & privacy pages
│   │
│   ├── model/                  # Data structures
│   │   └── models.go           # User, Post, Comment, Category models
│   │
│   ├── query/                  # Database queries
│   │   ├── user.go             # User operations
│   │   ├── post.go             # Post operations
│   │   ├── comment.go          # Comment operations
│   │   ├── session.go          # Session management
│   │   └── categories.go       # Category operations
│   │
│   └── service/                # Business logic
│       ├── service.go          # Service initialization
│       ├── Authentication.go   # Auth logic
│       ├── user.go             # User service
│       ├── post.go             # Post service
│       ├── comment.go          # Comment service
│       ├── profile.go          # Profile service
│       ├── session.go          # Session service
│       ├── categories.go       # Category service
│       ├── playground.go       # Code execution service
│       ├── golang.go           # Go-specific features
│       └── error.go            # Error handling
│
└── web/                        # Frontend assets
    ├── static/                 # Static files
    │   ├── css/               # Stylesheets
    │   │   ├── styles.css     # Global styles
    │   │   ├── layout.css     # Layout components
    │   │   ├── home.css       # Home page styles
    │   │   ├── login.css      # Login/register styles
    │   │   ├── profile.css    # Profile styles
    │   │   ├── edit-profile.css
    │   │   ├── discover.css   # Feed styles
    │   │   ├── view-post.css  # Post view styles
    │   │   ├── editor.css     # Code editor styles
    │   │   ├── category-base.css
    │   │   ├── coding.css     # Playground styles
    │   │   ├── help-privacy.css
    │   │   ├── root.css
    │   │   └── err.css        # Error page styles
    │   │
    │   ├── icons/             # UI icons
    │   │   └── categories/    # Category-specific icons
    │   │
    │   ├── img/               # General images
    │   └── login-img/         # Login page images
    │
    └── templates/             # HTML templates
        ├── root.html          # Home page
        ├── login.html         # Login/register page
        ├── profile.html       # User profile
        ├── edit-profile.html  # Profile editor
        ├── create-post.html   # Post creation form
        ├── view-post.html     # Single post view
        ├── DiscoverPosts.html # All posts feed
        ├── category.html      # Category-filtered posts
        ├── startcoding.html   # Code playground
        ├── help-privacy.html  # Help & privacy
        └── error.html         # Error pages
```

## 📋 Prerequisites

- **Go**: Version 1.23.1 or higher
- **GCC**: Required for SQLite CGO compilation
- **Docker** (optional): For containerized deployment

### Platform-Specific Requirements

#### Windows

- Install [MinGW-w64](https://www.mingw-w64.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- Add GCC to your PATH

#### Linux

```bash
sudo apt-get install gcc
sudo apt-get install sqlite3
```

#### macOS

```bash
xcode-select --install
brew install sqlite3
```

## 🚀 Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/forum.git
cd forum
```

2. **Install dependencies**

```bash
go mod download
```

3. **Initialize the database**

```bash
# The database will be automatically initialized on first run
# Or manually create it:
sqlite3 forum.db < database/db.sql
```

4. **Verify installation**

```bash
go build -o forum.exe
```

## 🏃 Running the Application

### Development Mode

**Windows (PowerShell)**

```powershell
go run main.go
```

**Linux/macOS**

```bash
go run main.go
```

The server will start on `http://localhost:7777`

### Production Build

**Windows**

```powershell
go build -o forum.exe
./forum.exe
```

**Linux/macOS**

```bash
go build -o forum
./forum
```

## 🐳 Docker Deployment

### Build the Docker image

```bash
docker build -t forum-app .
```

### Run the container

```bash
docker run -p 7777:7777 -v $(pwd)/database:/root/database forum-app
```

### Using Docker Compose (create `docker-compose.yml`)

```yaml
version: "3.8"
services:
  forum:
    build: .
    ports:
      - "7777:7777"
    volumes:
      - ./database:/root/database
    restart: unless-stopped
```

Run with:

```bash
docker-compose up -d
```

## 🗄 Database Schema

### Tables

#### `users`

- `id`: Primary key (auto-increment)
- `photo`: Profile image filename (default: 'default.png')
- `username`: Unique username
- `email`: Unique email address
- `password`: Bcrypt hashed password

#### `posts`

- `id`: Primary key
- `user_id`: Foreign key to users
- `title`: Post title
- `content`: Post content (supports markdown/code)
- `created_at`: Timestamp

#### `comments`

- `id`: Primary key
- `user_id`: Foreign key to users
- `post_id`: Foreign key to posts
- `content`: Comment text

#### `categories`

- `id`: Primary key
- `name`: Category name (unique)
- Predefined: golang, javascript, rust, ruby, python, matlab, brainfuck, java

#### `post_categories`

- Many-to-many relationship between posts and categories
- `post_id`, `category_id`: Composite primary key

#### `post_reactions`

- `user_id`, `post_id`: Composite primary key
- `reaction_type`: 'like' or 'dislike'

#### `comments_reactions`

- `user_id`, `comment_id`: Composite primary key
- `reaction_type`: 'like' or 'dislike'

#### `sessions`

- `id`: Primary key
- `user_id`: Foreign key to users
- `session_id`: UUID for session identification
- `expires_at`: Session expiration timestamp

## 🔌 API Endpoints

### Authentication

- `GET /register` - Display registration form
- `POST /register` - Create new user account
- `GET /login` - Display login form
- `POST /login` - Authenticate user
- `GET /logout` - End user session

### Posts

- `GET /` - Home page
- `GET /posts` - Discover all posts
- `GET /post/{id}` - View specific post
- `GET /create-post` - Display post creation form
- `POST /create-post` - Submit new post
- `POST /post-reaction` - Like/dislike a post

### Comments

- `POST /create-comment` - Add comment to post
- `POST /comment-reaction` - Like/dislike a comment

### Categories

- `GET /category/{slug}` - View posts in category (e.g., `/category/golang`)

### User Profile

- `GET /profile` - View user profile (query param: `?username=xyz`)
- `GET /edit-profile` - Display profile editor
- `POST /edit-profile/password` - Update password
- `POST /edit-profile/photo` - Upload new profile picture

### Playground

- `GET /playground` - Code playground interface
- `POST /playground/preview` - Preview code execution
- `POST /download` - Download code file

### Static Pages

- `GET /help` - Help & documentation
- `GET /privacy-terms` - Privacy policy & terms

### Assets

- `/web/*` - Static files (CSS, JS, images)
- `/profile-img/*` - User profile images
- `/assets/profile.css` - Dynamic profile CSS
- `/assets/edit-profile.css` - Dynamic edit profile CSS

## 📖 Usage Guide

### For Users

1. **Registration**

   - Navigate to `/register`
   - Provide username, email, and password
   - Upload a profile picture (optional)

2. **Creating Posts**

   - Login to your account
   - Click "Create Post"
   - Add title and content (supports code blocks)
   - Select relevant categories
   - Submit your post

3. **Engaging with Content**

   - Like/dislike posts and comments
   - Add comments to posts
   - Browse by category
   - View user profiles

4. **Code Playground**
   - Navigate to `/playground`
   - Write code with syntax highlighting
   - Preview execution results
   - Download your code

### For Developers

#### Adding New Categories

Edit `database/db.sql`:

```sql
INSERT OR IGNORE INTO categories (name)
VALUES ('new-category');
```

#### Custom Handlers

Add to `Internal/api/server.go`:

```go
router.HandleFunc("GET /your-route", server.YourHandler)
```

#### Database Queries

Add functions to respective files in `Internal/query/`

#### Business Logic

Implement in `Internal/service/`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code Style

- Follow Go conventions and best practices
- Run `go fmt` before committing
- Add comments for exported functions
- Write tests for new features

## 📝 License

This project is created for educational purposes.

## 👥 Authors

**Team GITEA**

- Batool Sayed (@basayed)
- HUSSAIN ALI (@hussainali7)
- Shaikh Alkaabi (@shaikaabi)
- Bader Alafoo (@balafoo)

## 🐛 Known Issues

- Session cleanup needs to be implemented for expired sessions
- Profile image validation could be more strict
- Consider adding pagination for large post lists

## 🔮 Future Enhancements

- [ ] Real-time notifications
- [ ] Direct messaging between users
- [ ] Post editing functionality
- [ ] Advanced search with filters
- [ ] User reputation system
- [ ] Post bookmarking
- [ ] Email verification
- [ ] OAuth integration (GitHub, Google)
- [ ] API rate limiting
- [ ] Admin dashboard
- [ ] Post tags in addition to categories
- [ ] Markdown editor with preview
- [ ] Image uploads in posts
- [ ] User following system

## 📞 Support

For issues and questions, please open an issue in the GitHub repository.

---

**Made with ❤️ using Go**
