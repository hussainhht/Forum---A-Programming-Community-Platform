
main.go
└── main()
    ├── database.InitDB()
    ├── api.NewService(db)
    ├── api.NewServer(port, service)
    └── server.Start()
        ├── server.Get_RootHandler()
        ├── server.Get_RegisterHandler()
        ├── server.Post_RegisterHandler()
        ├── server.Get_LoginHandler()
        ├── server.Post_LoginHandler()
        ├── server.Get_CreatePostHandler()
        ├── server.Post_CreatePostHandler()
        ├── server.Post_CreateCommentHandler()
        ├── server.Get_PostHandler()
        ├── server.Post_ReactionHandler()
        ├── server.CommentReactionHandler()
        └── server.LogoutHandler()

Internal/api/server.go
├── NewService(db)
├── NewServer(port, service)
└── Server.Start()
    └── (see above for handlers)

Internal/api/root.go
└── Server.Get_RootHandler()
    ├── server.Service.GetAllPosts()

Internal/api/register.go
├── Server.Get_RegisterHandler()
├── Server.Post_RegisterHandler()
    └── server.Service.RegisterUser()
        ├── service.IsValidEmail()
        ├── service.IsValidPassword()
        ├── query.SelectUserWhereUsername()
        ├── query.SelectUserWhereEmail()
        └── query.InsertUser()

Internal/api/login.go
├── Server.Get_LoginHandler()
├── Server.Post_LoginHandler()
    └── server.Service.LoginUser()
        ├── query.GetUserByUsernameOrEmail()
        ├── query.RemoveSession()
        └── query.CreateSession()

Internal/api/logout.go
└── Server.LogoutHandler()
    ├── server.Service.GetUserIDFromSessionID()
    └── query.RemoveSession()

Internal/api/createPost.go
├── Server.Get_CreatePostHandler()
    └── server.Service.GetCategories()
├── Server.Post_CreatePostHandler()
    ├── server.Service.IsValidSession()
    ├── server.Service.GetSessionIDFromCookie()
    └── server.Service.CreatePost()

Internal/api/viewPost.go
└── Server.Get_PostHandler()
    └── server.Service.GetPostByID()

Internal/api/postInteractions.go
└── Server.Post_ReactionHandler()
    ├── server.Service.GetSessionIDFromCookie()
    ├── server.Service.GetUserIDFromSessionID()
    └── server.Service.PostReaction()

Internal/api/commentInteractions.go
├── Server.CommentReactionHandler()
    ├── server.Service.GetSessionIDFromCookie()
    ├── server.Service.GetUserIDFromSessionID()
    └── server.Service.CommentReaction()
├── Server.Post_CreateCommentHandler()
    ├── server.Service.GetSessionIDFromCookie()
    ├── server.Service.GetUserIDFromSessionID()
    └── server.Service.CreateComment()

Internal/service/service.go
└── Service struct

Internal/service/register.go
├── Service.RegisterUser()
├── hashPassword()
├── Service.IsValidPassword()
└── Service.IsValidEmail()

Internal/service/login.go
└── Service.LoginUser()

Internal/service/session.go
├── Service.IsValidSession()
├── Service.GetSessionIDFromCookie()
└── Service.GetUserIDFromSessionID()

Internal/service/post.go
├── Service.CreatePost()
├── Service.GetAllPosts()
├── Service.GetPostByID()
└── Service.PostReaction()

Internal/service/comment.go
├── Service.CommentReaction()
├── Service.GetCommentsByPostID()
└── Service.CreateComment()

Internal/service/categories.go
└── Service.GetCategories()

Internal/query/user.go
├── InsertUser()
├── SelectUserWhereUsername()
├── SelectUserWhereEmail()
├── GetUserByUsernameOrEmail()
└── GetUsernameByUserID()

Internal/query/post.go
├── InsertPost()
├── GetAllPosts()
├── GetPostByID()
├── InsertPostReaction()
├── GetPostReaction()
├── UpdatePostReaction()
└── GetPostLikeCount()

Internal/query/comment.go
├── InsertComment()
├── InsertCommentReaction()
└── GetCommentsByPostID()

Internal/query/session.go
├── RemoveSession()
├── CreateSession()
├── SelectUserFromSession()
└── GetUserIDFromSession()

database/db.go
└── InitDB()